// Copyright (c) 2018-2021 Contributors as noted in the AUTHORS file
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package regagent

import (
	"context"
	"github.com/google/uuid"
	region "github.com/hegemonie-rpg/engine/pkg/region/model"
	proto "github.com/hegemonie-rpg/engine/pkg/region/proto"
	"github.com/juju/errors"
	"io"
)

type templatesApp struct {
	proto.UnimplementedTemplatesServer

	app *regionBackend
}

func (app *templatesApp) ListTemplates(req *proto.PaginatedStrQuery, stream proto.Templates_ListTemplatesServer) error {
	return app.app._regLock('r', req.Region, func(r *region.Region) error {
		for last := req.GetMarker(); ; {
			tab := r.Models.Slice(last, 100)
			if len(tab) <= 0 {
				return nil
			}
			for _, i := range tab {
				last = i.Name
				err := stream.Send(showCityTemplate(r, i))
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

func templateProtoToModel(r *region.Region, reqTpl *proto.CityTemplate) (*region.City, error) {
	tpl := region.City{
		ID:            0,
		Name:          "",
		Stock:         resAbsP2M(reqTpl.Stock),
		StockCapacity: resAbsP2M(reqTpl.StockCapacity),
		Production:    resAbsP2M(reqTpl.Production),
	}
	if tpl.Name == "" {
		tpl.Name = uuid.New().String()
	}
	for _, t := range reqTpl.BuildingTypes {
		err := tpl.InstantBuild(r, t)
		if err != nil {
			return nil, errors.Trace(err)
		}
	}
	for _, t := range reqTpl.SkillTypes {
		err := tpl.InstantStudy(r, t)
		if err != nil {
			return nil, err
		}
	}
	for _, t := range reqTpl.UnitTypes {
		err := tpl.InstantTrain(r, t)
		if err != nil {
			return nil, errors.Trace(err)
		}
	}
	return &tpl, nil
}

func (app *templatesApp) CreateTemplate(ctx context.Context, req *proto.CityTemplateReq) (*proto.Created, error) {
	created := proto.Created{}
	err := app.app._regLock('w', req.Region, func(r *region.Region) error {
		tpl, err := templateProtoToModel(r, req.Tpl)
		if err != nil {
			return errors.Trace(err)
		}
		tpl.Name = req.Id
		if tpl.Name == "" {
			tpl.Name = uuid.New().String()
		}
		if r.Models.Has(tpl.Name) {
			return errors.AlreadyExistsf("template '%s'", tpl.Name)
		}
		r.Models.Add(tpl)
		created.Id = tpl.Name
		return nil
	})
	return &created, err
}

func (app *templatesApp) UpdateTemplate(ctx context.Context, req *proto.CityTemplateReq) (*proto.None, error) {
	err := app.app._regLock('w', req.Region, func(r *region.Region) error {
		tpl, err := templateProtoToModel(r, req.Tpl)
		if err != nil {
			return errors.Trace(err)
		}
		tpl.Name = req.Id
		if !r.Models.Has(tpl.Name) {
			return errors.NotFoundf("template '%s'", tpl.Name)
		}
		tpl0 := r.Models.Get(tpl.Name)
		*tpl0 = *tpl
		return nil
	})
	return none, err
}

func (app *templatesApp) DeleteTemplate(ctx context.Context, req *proto.TemplateId) (*proto.None, error) {
	err := app.app._regLock('w', req.Region, func(r *region.Region) error {
		if !r.Models.Has(req.Id) {
			return errors.NotFoundf("template '%s'", req.Id)
		}
		r.Models.RemovePK(req.Id)
		return nil
	})
	return none, err
}
