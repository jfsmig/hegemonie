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

// DoRegionList dumps to os.Stdout a JSON stream of the known regions, sorted by name
func (cli *ClientCLI) DoRegionList(ctx context.Context, marker string) error {
	return cli.connect(ctx, func(ctx context.Context, cnx *grpc.ClientConn) error {
		req := proto.RegionListReq{NameMarker: marker}
		rep, err := proto.NewPublicClient(cnx).ListRegions(ctx, &req)
		if err != nil {
			return errors.Trace(err)
		}
		return utils.EncodeStream(func() (interface{}, error) { return rep.Recv() })
	})
}

// DoRegionCreate triggers the synchronous creation of a region with the given name, modeled on the named map.
func (cli *ClientCLI) DoRegionCreate(ctx context.Context, regID, mapID string) error {
	return cli.connect(ctx, func(ctx context.Context, cnx *grpc.ClientConn) error {
		_, err := proto.NewAdminClient(cnx).CreateRegion(ctx, &proto.RegionCreateReq{MapName: mapID, Name: regID})
		if err != nil {
			return errors.Trace(err)
		}
		return utils.StatusJSON(200, regID, "created")
	})
}

// DoRegionCities
func (cli *ClientCLI) DoRegionCities(ctx context.Context, regionID string, marker string) error {
	return cli.connect(ctx, func(ctx context.Context, cnx *grpc.ClientConn) error {
		var cityID uint64
		var err error
		if marker != "" {
			cityID, err = strconv.ParseUint(marker, 10, 64)
			if err != nil {
				return errors.Annotate(err, "Invalid city ID")
			}
		}
		req := proto.PaginatedU64Query{Region: regionID, Marker: cityID}
		rep, err := proto.NewPublicClient(cnx).AllCities(ctx, &req)
		if err != nil {
			return errors.Trace(err)
		}
		return utils.EncodeStream(func() (interface{}, error) { return rep.Recv() })
	})
}

// DoRegionStatsPush triggers a refresh of the stats (of the Region with the
// given ID) by the pointed region service
func (cli *ClientCLI) DoRegionStatsPush(ctx context.Context, reg string) error {
	return cli.connect(ctx, func(ctx context.Context, cnx *grpc.ClientConn) error {
		_, err := proto.NewAdminClient(cnx).PushStats(ctx, &proto.RegionId{Region: reg})
		if err != nil {
			return errors.Trace(err)
		}
		return utils.StatusJSON(200, reg, "Done")
	})
}

// DoRegionStatsGet triggers a refresh of the stats (of the Region with the
// given ID) by the pointed region service
func (cli *ClientCLI) DoRegionStatsGet(ctx context.Context, reg string) error {
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

// DoRegionRoundProduction triggers the production of resources on all the cities of the named region
func (cli *ClientCLI) DoRegionRoundProduction(ctx context.Context, reg string) error {
	return cli.connect(ctx, func(ctx context.Context, cnx *grpc.ClientConn) error {
		_, err := proto.NewGameMasterClient(cnx).Produce(ctx, &proto.RegionId{Region: reg})
		if err != nil {
			return errors.Trace(err)
		}
		return utils.StatusJSON(200, reg, "Produced")
	})
}

// DoRegionRoundMovement triggers one round of armies movement on all the cities of the named region
func (cli *ClientCLI) DoRegionRoundMovement(ctx context.Context, reg string) error {
	return cli.connect(ctx, func(ctx context.Context, cnx *grpc.ClientConn) error {
		_, err := proto.NewGameMasterClient(cnx).Move(ctx, &proto.RegionId{Region: reg})
		if err != nil {
			return errors.Trace(err)
		}
		return utils.StatusJSON(200, reg, "Moved")
	})
}
