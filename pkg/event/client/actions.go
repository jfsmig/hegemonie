// Copyright (c) 2018-2021 Contributors as noted in the AUTHORS file
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package evtclient

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jfsmig/hegemonie/pkg/event/proto"
	"github.com/jfsmig/hegemonie/pkg/utils"
	"github.com/juju/errors"
	"google.golang.org/grpc"
)

// ClientCLI gathers the event-related client actions available at the command line.
type ClientCLI struct{}

func (cfg *ClientCLI) connect(ctx context.Context, action utils.ActionFunc) error {
	endpoint, err := utils.DefaultDiscovery.Event()
	if err != nil {
		return errors.Trace(err)
	}
	return utils.Connect(ctx, endpoint, action)
}

// DoPush insert an event whose content and target are described on the command line.
// A descriptive error is returned in case of failure.
// FIXME(jfsmig): no retry is performed upon error
func (cfg *ClientCLI) DoPush(ctx context.Context, UserID string, msg ...string) error {
	return cfg.connect(ctx, func(ctx context.Context, cnx *grpc.ClientConn) error {
		var err error
		client := proto.NewProducerClient(cnx)
		for _, a := range msg {
			id := uuid.New().String()
			_, e := client.Push1(ctx, &proto.Push1Req{UserId: UserID, EvtId: id, Payload: []byte(a)})
			if e != nil {
				if err == nil {
					err = errors.Trace(e)
				} else {
					oldE := err
					err = errors.New("errors occured")
					err = errors.Annotate(err, oldE.Error())
					err = errors.Annotate(err, e.Error())
				}
				utils.Logger.Error().Str("char", UserID).Str("msg", a).Str("uuid", id).Err(err).Msg("PUSH")
			} else {
				utils.Logger.Info().Str("char", UserID).Str("msg", a).Str("uuid", id).Msg("PUSH")
			}
		}
		return err
	})
}

// DoAck consumes a message whose owner, timestamp and ID are described on the command line.
// FIXME(jfsmig): no retry is performed upon error
func (cfg *ClientCLI) DoAck(ctx context.Context, UserID, evtID string, when uint64) error {
	return cfg.connect(ctx, func(ctx context.Context, cnx *grpc.ClientConn) error {
		client := proto.NewConsumerClient(cnx)
		_, err := client.Ack1(ctx, &proto.Ack1Req{UserId: UserID, When: when, EvtId: evtID})
		if err != nil {
			return errors.Trace(err)
		}
		utils.Logger.Info().
			Str("char", UserID).
			Uint64("when", when).
			Str("uuid", evtID).
			Msg("ACK")
		return nil
	})
}

// DoList dumps to os.Stdout the Event objects streamed by the contacted service. The output consists
// in a JSON stream of objects separated by a CRLF (i.e. one object per line)
// FIXME(jfsmig): no retry is performed upon error
func (cfg *ClientCLI) DoList(ctx context.Context, UserID string, when uint64, marker string, max uint32) error {
	return cfg.connect(ctx, func(ctx context.Context, cnx *grpc.ClientConn) error {
		client := proto.NewConsumerClient(cnx)
		rep, err := client.List(ctx, &proto.ListReq{UserId: UserID, Marker: when, Max: 100})
		if err != nil {
			return errors.Trace(err)
		}
		anyError := false
		for _, x := range rep.Items {
			fmt.Printf("%s %d %s %s\n", x.UserId, x.When, x.EvtId, x.Payload)
		}
		if anyError {
			return errors.New("Invalid events matched")
		}
		return nil
	})
}
