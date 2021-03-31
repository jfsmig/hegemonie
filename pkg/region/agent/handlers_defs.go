// Copyright (c) 2018-2021 Contributors as noted in the AUTHORS file
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package regagent

import (
	proto "github.com/hegemonie-rpg/engine/pkg/region/proto"
	"io"
)

type defsApp struct {
	proto.UnimplementedDefinitionsServer

	app *regionApp
}

func (app *defsApp) ListUnits(req *proto.PaginatedU64Query, stream proto.Definitions_ListUnitsServer) error {
	return app.app._worldLock('r', func() error {
		for last := req.GetMarker(); ; {
			tab := app.app.w.Definitions.Units.Slice(last, 100)
			if len(tab) <= 0 {
				return nil
			}
			for _, i := range tab {
				last = i.ID
				err := stream.Send(&proto.UnitTypeView{
					Id: i.ID, Name: i.Name, Ticks: i.Ticks, Health: i.Health})
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

func (app *defsApp) ListBuildings(req *proto.PaginatedU64Query, stream proto.Definitions_ListBuildingsServer) error {
	return app.app._worldLock('r', func() error {
		for last := req.GetMarker(); ; {
			tab := app.app.w.Definitions.Buildings.Slice(last, 100)
			if len(tab) <= 0 {
				return nil
			}
			for _, i := range tab {
				last = i.ID
				v := &proto.BuildingTypeView{
					Ref: &proto.BuildingTypeRef{Id: i.ID, Name: i.Name, Ticks: i.Ticks},
					Public: &proto.BuildingTypePublic{
						Cost0:    resAbsM2P(i.Cost0),
						Cost:     resAbsM2P(i.Cost),
						Stock:    resModM2P(i.Stock),
						Prod:     resModM2P(i.Prod),
						Multiple: i.MultipleAllowed,
					},
					Private: &proto.BuildingTypePrivate{
						PopBuild:     i.PopBonusBuild,
						PopFall:      i.PopBonusFall,
						PopDismantle: i.PopBonusDismantle,
						PopDestroy:   i.PopBonusDestroy,
						Requires:     i.Requires,
						Conflicts:    i.Conflicts,
					},
				}
				err := stream.Send(v)
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

func (app *defsApp) ListKnowledges(req *proto.PaginatedU64Query, stream proto.Definitions_ListKnowledgesServer) error {
	return app.app._worldLock('r', func() error {
		for last := req.GetMarker(); ; {
			tab := app.app.w.Definitions.Knowledges.Slice(last, 100)
			if len(tab) <= 0 {
				return nil
			}
			for _, i := range tab {
				last = i.ID
				err := stream.Send(&proto.KnowledgeTypeView{
					Id: i.ID, Name: i.Name, Ticks: i.Ticks})
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
