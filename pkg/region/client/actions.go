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

func (cli *ClientCLI) DoCityBuild(ctx context.Context, regID, cityStrID, bStrID string) error {
	return cli.connect(ctx, func(ctx context.Context, cnx *grpc.ClientConn) error {
		var err error
		var cityID, bID uint64

		cityID, err = strconv.ParseUint(cityStrID, 10, 64)
		if err != nil {
			return errors.Annotate(err, "Invalid city ID")
		}
		bID, err = strconv.ParseUint(bStrID, 10, 64)
		if err != nil {
			return errors.Annotate(err, "Invalid city ID")
		}

		inArgs := proto.BuildReq{
			City: &proto.CityId{
				Region: regID, City: cityID,
			},
			BuildingType: bID,
		}
		_, err = proto.NewCityClient(cnx).Build(ctx, &inArgs)
		if err != nil {
			return errors.Trace(err)
		}
		return utils.StatusJSON(200, "-", "Started")
	})
}

func (cli *ClientCLI) DoCityStudy(ctx context.Context, regID, cityStrID, sStrID string) error {
	return cli.connect(ctx, func(ctx context.Context, cnx *grpc.ClientConn) error {
		var err error
		var cityID, skillID uint64

		cityID, err = strconv.ParseUint(cityStrID, 10, 64)
		if err != nil {
			return errors.Annotate(err, "Invalid city ID")
		}
		skillID, err = strconv.ParseUint(sStrID, 10, 64)
		if err != nil {
			return errors.Annotate(err, "Invalid city ID")
		}

		inArgs := proto.StudyReq{
			City: &proto.CityId{
				Region: regID, City: cityID,
			},
			KnowledgeType: skillID,
		}
		_, err = proto.NewCityClient(cnx).Study(ctx, &inArgs)
		if err != nil {
			return errors.Trace(err)
		}
		return utils.StatusJSON(200, "-", "Started")
	})
}

func (cli *ClientCLI) DoCityTrain(ctx context.Context, regID, cityStrID, uStrID string) error {
	return cli.connect(ctx, func(ctx context.Context, cnx *grpc.ClientConn) error {
		var err error
		var cityID, unitID uint64

		cityID, err = strconv.ParseUint(cityStrID, 10, 64)
		if err != nil {
			return errors.Annotate(err, "Invalid city ID")
		}
		unitID, err = strconv.ParseUint(uStrID, 10, 64)
		if err != nil {
			return errors.Annotate(err, "Invalid city ID")
		}

		inArgs := proto.TrainReq{
			City: &proto.CityId{
				Region: regID, City: cityID,
			},
			UnitType: unitID,
		}
		_, err = proto.NewCityClient(cnx).Train(ctx, &inArgs)
		if err != nil {
			return errors.Trace(err)
		}
		return utils.StatusJSON(200, "-", "Started")
	})
}

func (cli *ClientCLI) DoCityShow(ctx context.Context, regID, cityStrID string) error {
	return cli.connect(ctx, func(ctx context.Context, cnx *grpc.ClientConn) error {
		var err error
		var cityID uint64
		var cityView *proto.CityView

		cityID, err = strconv.ParseUint(cityStrID, 10, 64)
		if err != nil {
			return errors.Annotate(err, "Invalid city ID")
		}
		id := proto.CityId{
			Region: regID, City: cityID,
		}

		cityView, err = proto.NewCityClient(cnx).ShowAll(ctx, &id)
		if err != nil {
			return errors.Trace(err)
		}
		return utils.DumpJSON(cityView)
	})
}

func (cli *ClientCLI) connect(ctx context.Context, action utils.ActionFunc) error {
	endpoint, err := utils.DefaultDiscovery.Region()
	if err != nil {
		return errors.Trace(err)
	}
	return utils.Connect(ctx, endpoint, action)
}
