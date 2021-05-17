// Copyright (c) 2018-2021 Contributors as noted in the AUTHORS file
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package region

import (
	"testing"
)

func testFrontier(t *testing.T, f []*BuildingType, nb int) {
	if len(f) != nb {
		t.Log("Expected:", nb, "Got:", len(f))
		for _, bt := range f {
			t.Log("->", *bt)
		}
		t.Fatal()
	}
}

func TestBuildingFrontier(t *testing.T) {
	k := SetOfSkills{}
	k.Add(&Skill{ID: "1", Type: "1"})
	k.Add(&Skill{ID: "2", Type: "2"})
	k.Add(&Skill{ID: "3", Type: "3"})

	bt := SetOfBuildingTypes{}
	bt.Add(&BuildingType{ID: "1"})
	bt.Add(&BuildingType{ID: "2", PopRequired: 1})
	bt.Add(&BuildingType{ID: "3", Requires: []string{"3"}, MultipleAllowed: false})

	b := SetOfBuildings{}
	b.Add(&Building{ID: "1", Type: "1"})

	var f []*BuildingType

	// Pop & Req not matched
	f = bt.Frontier(0, []*Building{}, []*Skill{})
	testFrontier(t, f, 1)

	// Pop matched, not Req
	f = bt.Frontier(1, []*Building{}, []*Skill{})
	testFrontier(t, f, 2)

	// Pop & Req matched
	f = bt.Frontier(1, []*Building{}, []*Skill{{ID: "3", Type: "3"}})
	testFrontier(t, f, 3)

	// Pop & Req matched + Unicity
	f = bt.Frontier(1, []*Building{{ID: "1", Type: "3"}}, []*Skill{{ID: "3", Type: "3"}})
	testFrontier(t, f, 2)
}
