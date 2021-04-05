// Copyright (c) 2018-2021 Contributors as noted in the AUTHORS file
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package regclient

import (
	"context"
	"github.com/hegemonie-rpg/engine/pkg/region/proto"
	"github.com/hegemonie-rpg/engine/pkg/utils"
	"github.com/juju/errors"
	"google.golang.org/grpc"
	"strconv"
)

// DoRegionProduction triggers the production of resources on all the cities of the named region
func (cli *ClientCLI) DoRegionProduction(ctx context.Context, reg string) error {
	return cli.connect(ctx, func(ctx context.Context, cnx *grpc.ClientConn) error {
		_, err := proto.NewGameMasterClient(cnx).Produce(ctx, &proto.RegionId{Region: reg})
		if err != nil {
			return errors.Trace(err)
		}
		return utils.StatusJSON(200, reg, "Produced")
	})
}

// DoRegionMovement triggers one round of armies movement on all the cities of the named region
func (cli *ClientCLI) DoRegionMovement(ctx context.Context, reg string) error {
	return cli.connect(ctx, func(ctx context.Context, cnx *grpc.ClientConn) error {
		_, err := proto.NewGameMasterClient(cnx).Move(ctx, &proto.RegionId{Region: reg})
		if err != nil {
			return errors.Trace(err)
		}
		return utils.StatusJSON(200, reg, "Moved")
	})
}

// DoRegionGetStats triggers a refresh of the stats (of the Region with the
// given ID) by the pointed region service
func (cli *ClientCLI) DoRegionGetStats(ctx context.Context, reg string) error {
	return cli.connect(ctx, func(ctx context.Context, cnx *grpc.ClientConn) error {
		rep, err := proto.NewGameMasterClient(cnx).GetStats(ctx, &proto.RegionId{Region: reg})
		if err != nil {
			return errors.Trace(err)
		}
		return utils.EncodeStream(func() (interface{}, error) {
			itf, err := rep.Recv()
			if err != nil {
				return nil, err
			}
			return (&_cityStatsRecord{}).importFrom(itf), nil
		})
	})
}

type lifecycleFunc func(ctx context.Context, cnx proto.GameMasterClient, cityID *proto.CityId) (string, error)

func (cli *ClientCLI) DoRegionGetCityDetail(ctx context.Context, reg, cityStrID string) error {
	cityID, err := strconv.ParseUint(cityStrID, 10, 64)
	if err != nil {
		return errors.Annotate(err, "Invalid city ID")
	}
	return cli.connect(ctx, func(ctx context.Context, cnx *grpc.ClientConn) error {
		out, err := proto.NewGameMasterClient(cnx).GetDetail(ctx, &proto.CityId{Region: reg, City: cityID})
		if err != nil {
			return errors.Trace(err)
		}
		return utils.DumpJSON(out)
	})
}

func (cli *ClientCLI) doLifecyleAction(ctx context.Context, reg, cityStrID string, fn lifecycleFunc) error {
	cityID, err := strconv.ParseUint(cityStrID, 10, 64)
	if err != nil {
		return errors.Annotate(err, "Invalid city ID")
	}
	return cli.connect(ctx, func(ctx context.Context, cnx *grpc.ClientConn) error {
		msg, err := fn(ctx, proto.NewGameMasterClient(cnx), &proto.CityId{Region: reg, City: cityID})
		if err != nil {
			return errors.Trace(err)
		}
		return utils.StatusJSON(200, reg, msg)
	})
}

func (cli *ClientCLI) DoLifecyleConfigure(ctx context.Context, reg, cityStrID, model string) error {
	return cli.doLifecyleAction(ctx, reg, cityStrID,
		func(ctx context.Context, cnx proto.GameMasterClient, cityID *proto.CityId) (string, error) {
			req := proto.LifecycleConfigureReq{Region: cityID.Region, City: cityID.City, Model: model}
			_, err := cnx.LifecycleConfigure(ctx, &req)
			return "configured", err
		})
}

func (cli *ClientCLI) DoLifecyleAssign(ctx context.Context, reg, cityStrID, user string) error {
	return cli.doLifecyleAction(ctx, reg, cityStrID,
		func(ctx context.Context, cnx proto.GameMasterClient, cityID *proto.CityId) (string, error) {
			req := proto.LifecycleAssignReq{
				User: user,
				City: &proto.CityId{Region: cityID.Region, City: cityID.City},
			}
			_, err := cnx.LifecycleAssign(ctx, &req)
			return "assigned", err
		})
}

func (cli *ClientCLI) DoLifecyleResume(ctx context.Context, reg, cityStrID string) error {
	return cli.doLifecyleAction(ctx, reg, cityStrID,
		func(ctx context.Context, cnx proto.GameMasterClient, cityID *proto.CityId) (string, error) {
			_, err := cnx.LifecycleResume(ctx, &proto.LifecycleAbstractReq{City: cityID})
			return "resumed", err
		})
}

func (cli *ClientCLI) DoLifecyleDismiss(ctx context.Context, reg, cityStrID string) error {
	return cli.doLifecyleAction(ctx, reg, cityStrID,
		func(ctx context.Context, cnx proto.GameMasterClient, cityID *proto.CityId) (string, error) {
			_, err := cnx.LifecycleResume(ctx, &proto.LifecycleAbstractReq{City: cityID})
			return "dismissed", err
		})
}

func (cli *ClientCLI) DoLifecyleSuspend(ctx context.Context, reg, cityStrID string) error {
	return cli.doLifecyleAction(ctx, reg, cityStrID,
		func(ctx context.Context, cnx proto.GameMasterClient, cityID *proto.CityId) (string, error) {
			_, err := cnx.LifecycleSuspend(ctx, &proto.LifecycleAbstractReq{City: cityID})
			return "suspended", err
		})
}

func (cli *ClientCLI) DoLifecyleReset(ctx context.Context, reg, cityStrID string) error {
	return cli.doLifecyleAction(ctx, reg, cityStrID,
		func(ctx context.Context, cnx proto.GameMasterClient, cityID *proto.CityId) (string, error) {
			_, err := cnx.LifecycleReset(ctx, &proto.LifecycleAbstractReq{City: cityID})
			return "reset", err
		})
}
