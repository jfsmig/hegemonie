// Copyright (c) 2018-2021 Contributors as noted in the AUTHORS file
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package regagent

import (
	region "github.com/hegemonie-rpg/engine/pkg/region/model"
	"github.com/hegemonie-rpg/engine/pkg/region/proto"
	"io"
)

type publicApp struct {
	proto.UnimplementedPublicServer

	app *regionBackend
}

func (app *publicApp) ListRegions(req *proto.RegionListReq, stream proto.Public_ListRegionsServer) error {
	return app.app._worldLock('r', func() error {
		for marker := req.NameMarker; ; {
			tab := app.app.w.Regions.Slice(marker, 100)
			if len(tab) <= 0 {
				return nil
			}
			for _, x := range tab {
				marker = x.Name
				summary := &proto.RegionSummary{
					Name:        x.Name,
					MapName:     x.MapName,
					CountCities: uint32(len(x.Cities)),
					CountFights: uint32(len(x.Fights)),
				}
				err := stream.Send(summary)
				if err != nil {
					if err == io.EOF {
						return nil
					}
					return err
				}
			}
		}
	})
}

func (s *publicApp) AllCities(req *proto.PaginatedU64Query, stream proto.Public_AllCitiesServer) error {
	return s.app._regLock('r', req.Region, func(r *region.Region) error {
		last := req.Marker
		for {
			tab := r.Cities.Slice(last, 100)
			if len(tab) <= 0 {
				return nil
			}
			for _, c := range tab {
				last = c.ID
				err := stream.Send(showCityKey(r, c))
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
