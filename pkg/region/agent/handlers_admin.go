// Copyright (c) 2018-2021 Contributors as noted in the AUTHORS file
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package regagent

import (
	"context"
	mproto "github.com/hegemonie-rpg/engine/pkg/map/proto"
	"github.com/hegemonie-rpg/engine/pkg/region/model"
	"github.com/hegemonie-rpg/engine/pkg/region/proto"
	"github.com/hegemonie-rpg/engine/pkg/utils"
	"github.com/influxdata/influxdb-client-go/v2"
	"github.com/juju/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"strconv"
	"time"
)

type adminApp struct {
	proto.UnimplementedAdminServer

	app *regionBackend
}

func (app *adminApp) CreateRegion(ctx context.Context, req *proto.RegionCreateReq) (*proto.Created, error) {
	//  first, load the cities from the maps repository
	endpoint, err := utils.DefaultDiscovery.Map()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "configuration error: %v", err)
	}

	out := make([]region.NamedCity, 0)

	err = utils.Connect(ctx, endpoint, func(ctx context.Context, cli *grpc.ClientConn) error {
		client := mproto.NewMapClient(cli)
		marker := uint64(0)
		rep, err := client.Cities(ctx, &mproto.ListCitiesReq{
			MapName: req.MapName,
			Marker:  marker,
		})
		if err != nil {
			return errors.Trace(err)
		}

		for {
			x, err := rep.Recv()
			if err != nil {
				if err == io.EOF {
					break
				}
				return errors.Trace(err)
			}
			marker = x.GetId()
			out = append(out, region.NamedCity{Name: x.GetName(), ID: x.GetId()})
		}
		return nil
	})
	if err != nil {
		return nil, errors.Trace(err)
	}

	var id string
	err = app.app._worldLock('w', func() error {
		r, err := app.app.w.CreateRegion(req.Name, req.MapName, out)
		id = r.Name
		return err
	})
	return &proto.Created{Id: id}, err
}

func (app *adminApp) PushStats(ctx context.Context, req *proto.RegionId) (*proto.None, error) {
	return none, app.app._regLock('r', req.Region, func(r *region.Region) error {
		when := time.Now().Truncate(time.Minute)
		client := influxdb2.NewClientWithOptions(
			"http://localhost:8086",
			"1987b7a4-2bd4-4fdd-a701-5bc9ced89c94",
			influxdb2.DefaultOptions().SetPrecision(time.Second).SetBatchSize(100))
		defer client.Close()
		writeAPI := client.WriteAPI(
			"hegemonie",
			"hege_region_0")
		defer writeAPI.Flush()
		for _, c := range r.Cities {
			stats := showCityStats(r, c)
			p := influxdb2.NewPointWithMeasurement("stat").
				AddTag("region", r.Name).
				AddTag("city", strconv.FormatUint(c.ID, 10)).
				AddField("r_sent_0", s(stats.ResourceSent.R0)).
				AddField("r_sent_1", s(stats.ResourceSent.R1)).
				AddField("r_sent_2", s(stats.ResourceSent.R2)).
				AddField("r_sent_3", s(stats.ResourceSent.R3)).
				AddField("r_sent_4", s(stats.ResourceSent.R4)).
				AddField("r_sent_5", s(stats.ResourceSent.R5)).
				AddField("r_recv_0", s(stats.ResourceReceived.R0)).
				AddField("r_recv_1", s(stats.ResourceReceived.R1)).
				AddField("r_recv_2", s(stats.ResourceReceived.R2)).
				AddField("r_recv_3", s(stats.ResourceReceived.R3)).
				AddField("r_recv_4", s(stats.ResourceReceived.R4)).
				AddField("r_recv_5", s(stats.ResourceReceived.R5)).
				AddField("r_prod_0", s(stats.ResourceProduced.R0)).
				AddField("r_prod_1", s(stats.ResourceProduced.R1)).
				AddField("r_prod_2", s(stats.ResourceProduced.R2)).
				AddField("r_prod_3", s(stats.ResourceProduced.R3)).
				AddField("r_prod_4", s(stats.ResourceProduced.R4)).
				AddField("r_prod_5", s(stats.ResourceProduced.R5)).
				AddField("t_sent_0", s(stats.TaxSent.R0)).
				AddField("t_sent_1", s(stats.TaxSent.R1)).
				AddField("t_sent_2", s(stats.TaxSent.R2)).
				AddField("t_sent_3", s(stats.TaxSent.R3)).
				AddField("t_sent_4", s(stats.TaxSent.R4)).
				AddField("t_sent_5", s(stats.TaxSent.R5)).
				AddField("t_recv_0", s(stats.TaxReceived.R0)).
				AddField("t_recv_1", s(stats.TaxReceived.R1)).
				AddField("t_recv_2", s(stats.TaxReceived.R2)).
				AddField("t_recv_3", s(stats.TaxReceived.R3)).
				AddField("t_recv_4", s(stats.TaxReceived.R4)).
				AddField("t_recv_5", s(stats.TaxReceived.R5)).
				AddField("s_used_0", s(stats.StockUsage.R0)).
				AddField("s_used_1", s(stats.StockUsage.R1)).
				AddField("s_used_2", s(stats.StockUsage.R2)).
				AddField("s_used_3", s(stats.StockUsage.R3)).
				AddField("s_used_4", s(stats.StockUsage.R4)).
				AddField("s_used_5", s(stats.StockUsage.R5)).
				AddField("s_max_0", s(stats.StockCapacity.R0)).
				AddField("s_max_1", s(stats.StockCapacity.R1)).
				AddField("s_max_2", s(stats.StockCapacity.R2)).
				AddField("s_max_3", s(stats.StockCapacity.R3)).
				AddField("s_max_4", s(stats.StockCapacity.R4)).
				AddField("s_max_5", s(stats.StockCapacity.R5)).
				AddField("u_raised", s(stats.UnitRaised)).
				AddField("u_lost", s(stats.UnitLost)).
				AddField("a_score", s(stats.ScoreArmy)).
				AddField("k_score", s(stats.ScoreKnowledge)).
				AddField("b_score", s(stats.ScoreBuilding)).
				AddField("f_lost", s(stats.FightLost)).
				AddField("f_won", s(stats.FightWon)).
				AddField("f_joined", s(stats.FightJoined)).
				AddField("f_left", s(stats.FightLeft)).
				AddField("moves", s(stats.Moves)).
				SetTime(when)
			writeAPI.WritePoint(p)
		}
		return nil
	})
}
