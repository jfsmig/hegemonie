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
)

func (cli *ClientCLI) DoRegionGetBuildings(ctx context.Context, reg string) error {
	return cli.connect(ctx, func(ctx context.Context, cnx *grpc.ClientConn) error {
		out, err := proto.NewDefinitionsClient(cnx).ListBuildings(ctx, &proto.PaginatedU64Query{})
		if err != nil {
			return errors.Trace(err)
		}
		return utils.EncodeStream(func() (interface{}, error) { return out.Recv() })
	})
}

func (cli *ClientCLI) DoRegionGetSkills(ctx context.Context, reg string) error {
	return cli.connect(ctx, func(ctx context.Context, cnx *grpc.ClientConn) error {
		out, err := proto.NewDefinitionsClient(cnx).ListKnowledges(ctx, &proto.PaginatedU64Query{})
		if err != nil {
			return errors.Trace(err)
		}
		return utils.EncodeStream(func() (interface{}, error) { return out.Recv() })
	})
}

func (cli *ClientCLI) DoRegionGetUnits(ctx context.Context, reg string) error {
	return cli.connect(ctx, func(ctx context.Context, cnx *grpc.ClientConn) error {
		out, err := proto.NewDefinitionsClient(cnx).ListUnits(ctx, &proto.PaginatedU64Query{})
		if err != nil {
			return errors.Trace(err)
		}
		return utils.EncodeStream(func() (interface{}, error) { return out.Recv() })
	})
}
