// Copyright (c) 2018-2021 Contributors as noted in the AUTHORS file
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package region

import (
	"testing"
)

func TestUnitFrontier(t *testing.T) {
	ut := SetOfUnitTypes{}
	ut.Add(&UnitType{ID: 1})
	ut.Add(&UnitType{ID: 2, RequiredBuilding: 2})
	ut.Add(&UnitType{ID: 3, RequiredBuilding: 2})
	ut.Add(&UnitType{ID: 4, RequiredBuilding: 3})

	b := SetOfBuildings{}
	b.Add(&Building{ID: "1", Type: 1})
	b.Add(&Building{ID: "2", Type: 2})
	b.Add(&Building{ID: "3", Type: 3})

	var f []*UnitType

	// Troops without requirement
	f = ut.Frontier(SetOfBuildings{})
	if len(f) != 1 {
		t.Fatal()
	}
	f = ut.Frontier(SetOfBuildings{&Building{ID: "1", Type: 1}})
	if len(f) != 1 {
		t.Fatal()
	}

	// Troops with requirements
	f = ut.Frontier(SetOfBuildings{&Building{ID: "1", Type: 1}, &Building{ID: "3", Type: 3}})
	if len(f) != 2 {
		t.Fatal()
	}
	f = ut.Frontier(SetOfBuildings{&Building{ID: "1", Type: 1}, &Building{ID: "3", Type: 3}, &Building{ID: "2", Type: 2}})
	if len(f) != 4 {
		t.Fatal()
	}
}
