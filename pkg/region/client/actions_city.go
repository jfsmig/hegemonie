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

// ClientCLI gathers the actions destined to be exposed at the CLI, to manage a region service.
type ClientCLI struct{}

func (cli *ClientCLI) connect(ctx context.Context, action utils.ActionFunc) error {
	endpoint, err := utils.DefaultDiscovery.Region()
	if err != nil {
		return errors.Trace(err)
	}
	return utils.Connect(ctx, endpoint, action)
}

func (cli *ClientCLI) DoCityBuild(ctx context.Context, regID, cityStrID, bStrID string) error {
	cityID, err := strconv.ParseUint(cityStrID, 10, 64)
	if err != nil {
		return errors.Annotate(err, "Invalid city ID")
	}

	return cli.connect(ctx, func(ctx context.Context, cnx *grpc.ClientConn) error {
		inArgs := proto.BuildReq{
			City:         &proto.CityId{Region: regID, City: cityID},
			BuildingType: bStrID,
		}
		_, err = proto.NewCityClient(cnx).Build(ctx, &inArgs)
		if err != nil {
			return errors.Trace(err)
		}
		return utils.StatusJSON(200, "-", "Started")
	})
}

func (cli *ClientCLI) DoCityStudy(ctx context.Context, regID, cityStrID, sStrID string) error {
	cityID, err := strconv.ParseUint(cityStrID, 10, 64)
	if err != nil {
		return errors.Annotate(err, "Invalid city ID")
	}

	return cli.connect(ctx, func(ctx context.Context, cnx *grpc.ClientConn) error {
		inArgs := proto.StudyReq{
			City:          &proto.CityId{Region: regID, City: cityID},
			KnowledgeType: sStrID,
		}
		_, err = proto.NewCityClient(cnx).Study(ctx, &inArgs)
		if err != nil {
			return errors.Trace(err)
		}
		return utils.StatusJSON(200, "-", "Started")
	})
}

func (cli *ClientCLI) DoCityTrain(ctx context.Context, regID, cityStrID, uStrID string) error {
	cityID, err := strconv.ParseUint(cityStrID, 10, 64)
	if err != nil {
		return errors.Annotate(err, "Invalid city ID")
	}

	return cli.connect(ctx, func(ctx context.Context, cnx *grpc.ClientConn) error {
		inArgs := proto.TrainReq{
			City:     &proto.CityId{Region: regID, City: cityID},
			UnitType: uStrID,
		}
		_, err = proto.NewCityClient(cnx).Train(ctx, &inArgs)
		if err != nil {
			return errors.Trace(err)
		}
		return utils.StatusJSON(200, "-", "Started")
	})
}

func (cli *ClientCLI) DoCityShow(ctx context.Context, regID, cityStrID string) error {
	cityID, err := strconv.ParseUint(cityStrID, 10, 64)
	if err != nil {
		return errors.Annotate(err, "Invalid city ID")
	}
	return cli.connect(ctx, func(ctx context.Context, cnx *grpc.ClientConn) error {
		cityView, err := proto.NewCityClient(cnx).ShowAll(ctx,
			&proto.CityId{Region: regID, City: cityID})
		if err != nil {
			return errors.Trace(err)
		}
		return utils.DumpJSON((&_cityLightView{}).importFrom(cityView))
	})
}

func (cli *ClientCLI) DoCityShow2(ctx context.Context, regID, cityStrID string) error {
	cityID, err := strconv.ParseUint(cityStrID, 10, 64)
	if err != nil {
		return errors.Annotate(err, "Invalid city ID")
	}
	return cli.connect(ctx, func(ctx context.Context, cnx *grpc.ClientConn) error {
		cityView, err := proto.NewCityClient(cnx).ShowAll(ctx,
			&proto.CityId{Region: regID, City: cityID})
		if err != nil {
			return errors.Trace(err)
		}
		return utils.DumpJSON((&_cityView{}).importFrom(cityView))
	})
}

func (cli *ClientCLI) DoCityDetail(ctx context.Context, reg, cityStrID string) error {
	cityID, err := strconv.ParseUint(cityStrID, 10, 64)
	if err != nil {
		return errors.Annotate(err, "Invalid city ID")
	}
	return cli.connect(ctx, func(ctx context.Context, cnx *grpc.ClientConn) error {
		out, err := proto.NewGameMasterClient(cnx).GetDetail(ctx, &proto.CityId{Region: reg, City: cityID})
		if err != nil {
			return errors.Trace(err)
		}
		return utils.DumpJSON((&_cityLightView{}).importFrom(out))
	})
}

type lifecycleFunc func(ctx context.Context, cnx proto.GameMasterClient, cityID *proto.CityId) (string, error)

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
