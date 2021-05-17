// Copyright (c) 2018-2021 Contributors as noted in the AUTHORS file
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package region

import (
	"github.com/google/uuid"
	"github.com/juju/errors"
)

func (reg *Region) UnitGet(city uint64, id string) *Unit {
	c := reg.CityGet(city)
	if c != nil {
		return c.GetUnit(id)
	}
	return nil
}

func (s SetOfUnitTypes) Frontier(owned []*Building) []*UnitType {
	bIndex := make(map[string]bool)
	for _, b := range owned {
		bIndex[b.Type] = true
	}
	result := make([]*UnitType, 0)
	for _, ut := range s {
		if ut.RequiredBuilding == "" || bIndex[ut.RequiredBuilding] {
			result = append(result, ut)
		}
	}
	return result
}

// Abruptly terminate the training of the Unit.
// The number of training ticks suddenly drop to 0, whatever its prior value.
func (u *Unit) Finish() *Unit {
	u.Ticks = 0
	return u
}

// Create a Unit of the given UnitType.
// No check is performed to verify the City has all the requirements.
func (c *City) StartUnit(w *Region, pType *UnitType) *Unit {
	id := uuid.New().String()
	u := &Unit{ID: id, Type: pType.ID, Ticks: pType.Ticks, Health: pType.Health}
	c.Units.Add(u)
	return u
}

// Start the training of a Unit of the given UnitType (id).
// The whole chain of requirements will be checked.
func (c *City) Train(w *Region, typeID string) (string, error) {
	t := w.world.UnitTypeGet(typeID)
	if t == nil {
		return "", errors.NotFoundf("unit type not found")
	}
	if !c.UnitAllowed(t) {
		return "", errors.Forbiddenf("no suitable building")
	}

	u := c.StartUnit(w, t)
	// TODO(jfs): emit a notification
	return u.ID, nil
}

func (c *City) InstantTrain(w *Region, typeID string) error {
	t := w.world.UnitTypeGet(typeID)
	if t == nil {
		return errors.NotFoundf("unit type not found")
	}
	c.StartUnit(w, t).Finish()
	return nil
}
