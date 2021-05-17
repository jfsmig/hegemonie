// Copyright (c) 2018-2021 Contributors as noted in the AUTHORS file
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package region

import (
	"github.com/google/uuid"
	"github.com/juju/errors"
)

func (s SetOfBuildingTypes) Frontier(pop int64, built []*Building, owned []*Skill) []*BuildingType {
	bmap := make(map[string]bool)
	pending := make(map[string]bool)
	finished := make(map[string]bool)
	for _, k := range owned {
		if k.Ticks == 0 {
			finished[k.Type] = true
		} else {
			pending[k.Type] = true
		}
	}
	for _, b := range built {
		bmap[b.Type] = true
	}

	valid := func(bt *BuildingType) bool {
		if bt.PopRequired > pop {
			return false
		}
		if !bt.MultipleAllowed && bmap[bt.ID] {
			return false
		}
		for _, c := range bt.Conflicts {
			if finished[c] || pending[c] {
				return false
			}
		}
		for _, c := range bt.Requires {
			if !finished[c] {
				return false
			}
		}
		return true
	}

	result := make([]*BuildingType, 0)
	for _, bt := range s {
		if valid(bt) {
			result = append(result, bt)
		}
	}
	return result
}

// Abruptly terminate the construction of the Building.
// The number of building ticks suddenly drop to 0, whatever its prior value.
func (b *Building) Finish() *Building {
	b.Ticks = 0
	return b
}

// StartBuilding declares a new Building with the given type and the ticks gauge at its max
// No check is performed to verify the City has all the requirements.
func (c *City) StartBuilding(t *BuildingType) *Building {
	id := uuid.New().String()
	b := &Building{ID: id, Type: t.ID, Ticks: t.Ticks}
	c.Buildings.Add(b)
	return b
}

// Build forwards the call to StartBuilding if all the conditions are met, after the initial fee
// given in the BuildingType
func (c *City) Build(w *Region, bID string) (string, error) {
	t := w.world.BuildingTypeGet(bID)
	if t == nil {
		return "", errors.NotFoundf("Building Type not found")
	}
	if !t.MultipleAllowed {
		for _, b := range c.Buildings {
			if b.Type == bID {
				return "", errors.AlreadyExistsf("building already present")
			}
		}
	}
	if !CheckSkillDependencies(c.ownedSkillTypes(w), t.Requires, t.Conflicts) {
		return "", errors.Forbiddenf("dependencies unmet")
	}
	if !c.Stock.GreaterOrEqualTo(t.Cost0) {
		return "", errors.Forbiddenf("insufficient resources")
	}

	c.Stock.Remove(t.Cost0)
	return c.StartBuilding(t).ID, nil
}

// Build forwards the call to StartBuilding if all the conditions are met, after the initial fee
// given in the BuildingType
func (c *City) InstantBuild(w *Region, bID string) error {
	t := w.world.BuildingTypeGet(bID)
	if t == nil {
		return errors.NotFoundf("Building Type not found")
	}
	if !t.MultipleAllowed {
		for _, b := range c.Buildings {
			if b.Type == bID {
				return nil
			}
		}
	}

	c.StartBuilding(t).Ticks = 0
	return nil
}
