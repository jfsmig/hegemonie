// Copyright (c) 2018-2021 Contributors as noted in the AUTHORS file
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package mapclient

import (
	"context"
	"github.com/jfsmig/hegemonie/pkg/map/proto"
	"github.com/jfsmig/hegemonie/pkg/utils"
	"github.com/juju/errors"
	"google.golang.org/grpc"
	"strconv"
)

// ClientCLI gathers the command line interface actions related to the Map API service.
// It wraps the main gRPC client calls and dumps all the output as JSON on os.Stdout.
type ClientCLI struct{}

// ListMaps produces to os.Stdout a JSON stream of objects, each object describing
// a Map registered in the repository.
func (c *ClientCLI) ListMaps(ctx context.Context, args PathArgs) error {
	return c.connect(ctx, func(ctx context.Context, cnx *grpc.ClientConn) error {
		rep, err := proto.NewMapClient(cnx).Maps(ctx, &proto.ListMapsReq{
			Marker: args.MapName,
		})
		if err != nil {
			return errors.Trace(err)
		}
		return utils.EncodeStream(func() (interface{}, error) { return rep.Recv() })
	})
}

// GetCities produces to os.Stdout a JSON array of cities, where each City is an <id,name> tuple.
func (c *ClientCLI) GetCities(ctx context.Context, args PathArgs) error {
	return c.connect(ctx, func(ctx context.Context, cnx *grpc.ClientConn) error {
		rep, err := proto.NewMapClient(cnx).Cities(ctx, &proto.ListCitiesReq{
			MapName: args.MapName,
			Marker:  args.Src,
		})
		if err != nil {
			return errors.Trace(err)
		}
		return utils.EncodeStream(func() (interface{}, error) { return rep.Recv() })
	})
}

// GetRoads produces to os.Stdout a JSON stream of <int,int> pairs objects, with one pair
// for a road in place on the given map.
func (c *ClientCLI) GetRoads(ctx context.Context, args PathArgs) error {
	return c.connect(ctx, func(ctx context.Context, cnx *grpc.ClientConn) error {
		rep, err := proto.NewMapClient(cnx).Edges(ctx, &proto.ListEdgesReq{
			MapName:   args.MapName,
			MarkerSrc: args.Src,
			MarkerDst: args.Dst,
		})
		if err != nil {
			return errors.Trace(err)
		}
		return utils.EncodeStream(func() (interface{}, error) { return rep.Recv() })
	})
}

// GetPositions produces to os.Stdout a JSON array of <id,x,y> tuples, all fields being integers,
// with one tuple for each position on the map, i.e. one tuple per vertex in the graph.
func (c *ClientCLI) GetPositions(ctx context.Context, args PathArgs) error {
	return c.connect(ctx, func(ctx context.Context, cnx *grpc.ClientConn) error {
		rep, err := proto.NewMapClient(cnx).Vertices(ctx, &proto.ListVerticesReq{
			MapName: args.MapName,
			Marker:  args.Src,
		})
		if err != nil {
			return errors.Trace(err)
		}
		return utils.EncodeWhole(func() (interface{}, error) { return rep.Recv() })
	})
}

// GetPath produces to os.Stdout a JSON array of integers, all fields being the unique ID
// of a position involved in the path from the given source to the given destination.
func (c *ClientCLI) GetPath(ctx context.Context, args PathArgs) error {
	return c.getPath(ctx, args)
}

// GetStep produces to os.Stdout a JSON singleton array of integers, with the unique value
// equals to the ID of the position that is the next step in the path from the given source
// to the given destination.
func (c *ClientCLI) GetStep(ctx context.Context, args PathArgs) error {
	args.Max = 1
	return c.getPath(ctx, args)
}

func (c *ClientCLI) getPath(ctx context.Context, args PathArgs) error {
	return c.connect(ctx, func(ctx context.Context, cnx *grpc.ClientConn) error {
		rep, err := proto.NewMapClient(cnx).GetPath(ctx, &proto.PathRequest{MapName: args.MapName, Src: args.Src, Dst: args.Dst})
		if err != nil {
			return errors.Trace(err)
		}
		return utils.EncodeWhole(func() (interface{}, error) { return rep.Recv() })
	})
}

func (c *ClientCLI) connect(ctx context.Context, action utils.ActionFunc) error {
	endpoint, err := utils.DefaultDiscovery.Map()
	if err != nil {
		return errors.Trace(err)
	}
	return utils.Connect(ctx, endpoint, action)
}

// PathArgs gathers the possible arguments for the map-related calls.
// Hopefully they are always the same.
type PathArgs struct {
	MapName  string
	Src, Dst uint64
	Max      uint32
}

// Parse extracts the elements of the PathArgs from the array of command line
// positional arguments.
func (pa *PathArgs) Parse(args []string) (err error) {
	if len(args) >= 1 {
		pa.MapName = args[0]
		if len(args) >= 2 {
			pa.Src, err = strconv.ParseUint(args[1], 10, 63)
			if err != nil {
				return errors.Trace(err)
			}
			if len(args) >= 3 {
				pa.Dst, err = strconv.ParseUint(args[2], 10, 63)
				if err != nil {
					return errors.Trace(err)
				}
				if len(args) >= 4 {
					return errors.BadRequestf("max 3 arguments expected: MAPNAME [INT [INT]]")
				}
			}
		}
	} else {
		return errors.BadRequestf("min 1 argument expected: MAPNAME [INT [INT]]")
	}
	return nil
}
