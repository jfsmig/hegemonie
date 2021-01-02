// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package region

import (
	"sort"
	"testing"
)

func TestSetOfBuildingType(t *testing.T) {
	u2 := &BuildingType{ID: 2}
	s := SetOfBuildingTypes{}
	s.Add(&BuildingType{ID: 1})
	s.Add(&BuildingType{ID: 3})
	s.Add(u2)
	s.Add(&BuildingType{ID: 4})
	if len(s) != 4 {
		t.Fatal()
	}
	if !sort.IsSorted(s) {
		t.Fatal()
	}
	if u2 != s.Get(2) {
		t.Fatal()
	}
}

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
	k := SetOfKnowledges{}
	k.Add(&Knowledge{ID: "1", Type: 1})
	k.Add(&Knowledge{ID: "2", Type: 2})
	k.Add(&Knowledge{ID: "3", Type: 3})

	bt := SetOfBuildingTypes{}
	bt.Add(&BuildingType{ID: 1})
	bt.Add(&BuildingType{ID: 2, PopRequired: 1})
	bt.Add(&BuildingType{ID: 3, Requires: []uint64{3}, MultipleAllowed: false})

	b := SetOfBuildings{}
	b.Add(&Building{ID: "1", Type: 1})

	var f []*BuildingType

	// Pop & Req not matched
	f = bt.Frontier(0, []*Building{}, []*Knowledge{})
	testFrontier(t, f, 1)

	// Pop matched, not Req
	f = bt.Frontier(1, []*Building{}, []*Knowledge{})
	testFrontier(t, f, 2)

	// Pop & Req matched
	f = bt.Frontier(1, []*Building{}, []*Knowledge{{ID: "3", Type: 3}})
	testFrontier(t, f, 3)

	// Pop & Req matched + Unicity
	f = bt.Frontier(1, []*Building{{ID: "1", Type: 3}}, []*Knowledge{{ID: "3", Type: 3}})
	testFrontier(t, f, 2)
}
