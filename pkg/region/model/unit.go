// Copyright (c) 2018-2021 Contributors as noted in the AUTHORS file
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package region

func (reg *Region) UnitGet(city uint64, id string) *Unit {
	c := reg.CityGet(city)
	if c != nil {
		return c.Unit(id)
	}
	return nil
}

func (s SetOfUnitTypes) Frontier(owned []*Building) []*UnitType {
	bIndex := make(map[uint64]bool)
	for _, b := range owned {
		bIndex[b.Type] = true
	}
	result := make([]*UnitType, 0)
	for _, ut := range s {
		if ut.RequiredBuilding == 0 || bIndex[ut.RequiredBuilding] {
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
