// Copyright (c) 2018-2021 Contributors as noted in the AUTHORS file
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package regagent

import (
	"context"
	"github.com/hegemonie-rpg/engine/pkg/region/model"
	"github.com/hegemonie-rpg/engine/pkg/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"

	proto "github.com/hegemonie-rpg/engine/pkg/region/proto"
)

type cityApp struct {
	proto.UnimplementedCityServer

	app *regionApp
}

func (s *cityApp) List(req *proto.CitiesByCharReq, stream proto.City_ListServer) error {
	return s.app._regLock('r', req.Region, func(r *region.Region) error {
		last := req.Marker
		for {
			tab := r.Cities.Slice(last, 100)
			if len(tab) <= 0 {
				return nil
			}
			for _, c := range tab {
				last = c.ID
				if c.Owner != req.Character {
					continue
				}
				err := stream.Send(showCityPublic(s.app.w, c, false))
				if err == io.EOF {
					return nil
				}
				if err != nil {
					return err
				}
			}
		}
	})
}

func (s *cityApp) AllCities(req *proto.PaginatedU64Query, stream proto.City_AllCitiesServer) error {
	return s.app._regLock('r', req.Region, func(r *region.Region) error {
		last := req.Marker
		for {
			tab := r.Cities.Slice(last, 100)
			if len(tab) <= 0 {
				return nil
			}
			for _, c := range tab {
				last = c.ID
				err := stream.Send(showCityPublic(s.app.w, c, false))
				if err == io.EOF {
					return nil
				}
				if err != nil {
					return err
				}
			}
		}
	})
}

func (s *cityApp) ShowAll(ctx context.Context, req *proto.CityId) (reply *proto.CityView, err error) {
	err = s.app.cityLock('r', req, func(r *region.Region, c *region.City) error {
		view := showCity(s.app.w, c)
		utils.Logger.Debug().
			Int("#a", len(view.Assets.Armies)).
			Int("#k", len(view.Assets.Knowledges)).
			Int("#b", len(view.Assets.Buildings)).
			Int("#u", len(view.Assets.Units)).
			Interface("prod", view.Production).
			Interface("stock", view.Stock).
			Msg("city")
		reply = view
		return nil
	})
	return reply, err
}

func (s *cityApp) Study(ctx context.Context, req *proto.StudyReq) (*proto.None, error) {
	return none, s.app.cityLock('w', req.City, func(r *region.Region, c *region.City) error {
		_, e := c.Study(r, req.KnowledgeType)
		return e
	})
}

func (s *cityApp) Build(ctx context.Context, req *proto.BuildReq) (*proto.None, error) {
	return none, s.app.cityLock('w', req.City, func(r *region.Region, c *region.City) error {
		_, e := c.Build(r, req.BuildingType)
		return e
	})
}

func (s *cityApp) Train(ctx context.Context, req *proto.TrainReq) (*proto.None, error) {
	return none, s.app.cityLock('w', req.City, func(r *region.Region, c *region.City) error {
		_, e := c.Train(r, req.UnitType)
		return e
	})
}

func (s *cityApp) ListArmies(req *proto.CityId, stream proto.City_ListArmiesServer) error {
	return s.app.cityLock('r', req, func(r *region.Region, c *region.City) error {
		var last string
		for {
			tab := c.Armies.Slice(last, 100)
			if len(tab) <= 0 {
				return nil
			}
			for _, a := range c.Armies {
				last = a.ID
				err := stream.Send(&proto.ArmyName{Id: a.ID, Name: a.Name})
				if err == io.EOF {
					return nil
				}
				if err != nil {
					return err
				}
			}
		}
	})
}

// Create an army made of only Units (no Resources carried)
func (s *cityApp) CreateArmy(ctx context.Context, req *proto.CreateArmyReq) (*proto.None, error) {
	return none, s.app.cityLock('w', req.City, func(r *region.Region, c *region.City) error {
		_, e := c.CreateArmyFromIds(r, req.Unit...)
		return e
	})
}

// Create an army made of only Resources (no Units)
func (s *cityApp) CreateTransport(ctx context.Context, req *proto.CreateTransportReq) (*proto.None, error) {
	return none, s.app.cityLock('w', req.City, func(r *region.Region, c *region.City) error {
		_, e := c.CreateTransport(r, resAbsP2M(req.Stock))
		return e
	})
}

func (s *cityApp) TransferUnit(ctx context.Context, req *proto.TransferUnitReq) (*proto.None, error) {
	return none, s.app.cityLock('w', req.City, func(r *region.Region, c *region.City) error {
		army := c.Armies.Get(req.Army)
		if army == nil {
			return status.Error(codes.NotFound, "no such army")
		}
		return c.TransferOwnUnit(army, req.Unit...)
	})
}

func (s *cityApp) TransferResources(ctx context.Context, req *proto.TransferResourcesReq) (*proto.None, error) {
	return none, s.app.cityLock('w', req.City, func(r *region.Region, c *region.City) error {
		army := c.Armies.Get(req.Army)
		if army == nil {
			return status.Error(codes.NotFound, "no such army")
		}
		return c.TransferOwnResources(army, resAbsP2M(req.Stock))
	})
}
