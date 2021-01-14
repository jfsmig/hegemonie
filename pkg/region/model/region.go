// Copyright (c) 2018-2021 Contributors as noted in the AUTHORS file
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package region

import (
	"context"
	"github.com/juju/errors"
)

// Produce performs a production round that involves all the cities on the map of the region.
// The round action might take long. But there is no notion of a transaction.
// As a consequence, the action will ignore the cancellation signal brought by the context.Context.
func (reg *Region) Produce(ctx context.Context) {
	for _, c := range reg.Cities {
		c.Produce(ctx, reg)
	}
}

// Move performs a movement round that involves all the armies of all the cities on the map of the region.
// The round action might take long. But there is no notion of a transaction.
// As a consequence, the action will ignore the cancellation signal brought by the context.Context.
func (reg *Region) Move(ctx context.Context) {
	for _, c := range reg.Cities {
		for _, a := range c.Armies {
			a.Move(ctx, reg)
		}
	}
}

func (reg *Region) CityGet(id uint64) *City {
	return reg.Cities.Get(id)
}

func (reg *Region) CityGetAt(loc uint64) *City {
	return reg.CityGet(loc)
}

func (reg *Region) CityCheck(id uint64) bool {
	return reg.CityGet(id) != nil
}

func (reg *Region) CityCreateModel(loc uint64, model *City) (*City, error) {
	if reg.Cities.Has(loc) {
		return nil, errors.AlreadyExistsf("city found at [%v]", loc)
	}
	city := CopyCity(model)
	city.ID = loc
	city.Name = "NOT-SET"
	reg.Cities.Add(city)
	return city, nil
}

func (reg *Region) CityCreate(loc uint64) (*City, error) {
	return reg.CityCreateModel(loc, nil)
}

func (reg *Region) CitiesList(idChar string) []*City {
	rep := make([]*City, 0)
	for _, c := range reg.Cities {
		if c.Owner == idChar || c.Deputy == idChar {
			rep = append(rep, c)
		}
	}
	return rep
}
