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

// DoListRegions dumps to os.Stdout a JSON stream of the known regions, sorted by name
func (cli *ClientCLI) DoListRegions(ctx context.Context, marker string) error {
	return cli.connect(ctx, func(ctx context.Context, cnx *grpc.ClientConn) error {
		req := proto.RegionListReq{NameMarker: marker}
		rep, err := proto.NewPublicClient(cnx).ListRegions(ctx, &req)
		if err != nil {
			return errors.Trace(err)
		}
		return utils.EncodeStream(func() (interface{}, error) { return rep.Recv() })
	})
}

func (cli *ClientCLI) DoListCities(ctx context.Context, regionID string, marker string) error {
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
