// Copyright (c) 2018-2021 Contributors as noted in the AUTHORS file
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package evtagent

import (
	"context"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	back "github.com/jfsmig/hegemonie/pkg/event/backend-local"
	"github.com/jfsmig/hegemonie/pkg/event/proto"
	"github.com/jfsmig/hegemonie/pkg/utils"
	"github.com/juju/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"math"
)

// Config gathers the configuration fields required to start a gRPC Event API service.
type Config struct {
	PathBase string `yaml:"base" json:"base"`
}

type eventService struct {
	proto.UnimplementedConsumerServer
	proto.UnimplementedProducerServer

	cfg     Config
	backend *back.Backend
}

// Application implements the expectations of the application backend
func (cfg Config) Application(_ context.Context) (utils.RegisterableMonitorable, error) {
	if cfg.PathBase == "" {
		return nil, errors.New("missing path to the event data directory")
	}

	var err error
	app := eventService{cfg: cfg}
	app.backend, err = back.Open(app.cfg.PathBase)
	if err != nil {
		return nil, errors.NewNotValid(err, "backend error")
	}

	return &app, nil
}

// Register pugs the internal gRPC routes into the given server
func (es *eventService) Register(grpcSrv *grpc.Server) error {
	proto.RegisterProducerServer(grpcSrv, es)
	proto.RegisterConsumerServer(grpcSrv, es)
	grpc_prometheus.Register(grpcSrv)
	utils.Logger.Info().
		Str("base", es.cfg.PathBase).
		Msg("ready")
	return nil
}

// Ack1 marks an event as read so that it won't be listed again.
func (es *eventService) Ack1(ctx context.Context, req *proto.Ack1Req) (*proto.None, error) {
	err := es.backend.Ack1(req.CharId, req.When, req.EvtId)
	return &proto.None{}, err
}

// List streams event objects belonging to the user with the given ID. The objects are sorted by
// decreasing timestamp then by increasing UUID. The events are served as they are stored, the
// messages are not rendered.
func (es *eventService) List(ctx context.Context, req *proto.ListReq) (*proto.ListRep, error) {
	items, err := es.backend.List(req.CharId, req.Marker, req.Max)
	if err != nil {
		return nil, err
	}

	rep := proto.ListRep{}
	for _, x := range items {
		rep.Items = append(rep.Items, &proto.ListItem{
			CharId:  x.CharID,
			When:    math.MaxUint64 - x.When,
			EvtId:   x.ID,
			Payload: x.Payload,
		})
	}
	return &rep, nil
}

// Push1 inserts an event in the log of the Character with the given ID.
// The current timestamp will be used. An UUID will be generated.
func (es *eventService) Push1(ctx context.Context, req *proto.Push1Req) (*proto.None, error) {
	err := es.backend.Push1(req.CharId, req.EvtId, req.Payload)
	return &proto.None{}, err
}

// Check implements the one-shot health-check of the gRPC service
func (es *eventService) Check(ctx context.Context) grpc_health_v1.HealthCheckResponse_ServingStatus {
	return grpc_health_v1.HealthCheckResponse_SERVING
}
