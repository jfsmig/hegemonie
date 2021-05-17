// Copyright (c) 2018-2021 Contributors as noted in the AUTHORS file
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package region

import (
	"context"
	"github.com/google/uuid"
	"testing"
)

type fixtureWorldFunc func(ctx context.Context, t *testing.T, w *World)

type fixtureRegionFunc func(ctx context.Context, t *testing.T, r *Region)

type localFullMeshMap struct{}

// Step always returns the destination because, by definition, the map is a full-mesh and any location
// is connected to any other location.
func (r *localFullMeshMap) Step(ctx context.Context, mapName string, src, dst uint64) (uint64, error) {
	return dst, nil
}

func definitionSandbox() DefinitionsBase {
	db := DefinitionsBase{}
	db.Skills.Add(&SkillType{
		ID: "1", Name: uuid.New().String(), Ticks: 1,
		Cost0: ResourcesUniform(0),
		Cost:  ResourcesUniform(0),
		Prod:  ResourceModifierNoop(),
		Stock: ResourceModifierNoop(),
	})
	db.Buildings.Add(&BuildingType{
		ID: "1", Name: uuid.New().String(), Ticks: 1,
		Cost0:    ResourcesUniform(0),
		Cost:     ResourcesUniform(0),
		Prod:     ResourceModifierUniform(1.0, 1.0),
		Stock:    ResourceModifierNoop(),
		Requires: []string{"1"},
	})
	db.Units.Add(&UnitType{
		ID: "1", Name: uuid.New().String(), Ticks: 1,
		Cost0:            ResourcesUniform(0),
		Cost:             ResourcesUniform(0),
		Prod:             ResourceModifierNoop(),
		Health:           100,
		HealthFactor:     1.0,
		RequiredBuilding: "1",
	})
	return db
}

func fixtureWorld(t *testing.T, hook fixtureWorldFunc) {
	ctx := context.Background()
	if dl, ok := t.Deadline(); ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithDeadline(ctx, dl)
		defer cancel()
	}

	w, err := NewWorld()
	if err != nil {
		t.Fatalf("World fixture creation error: %v", err)
	}
	w.SetMapClient(&localFullMeshMap{})
	w.Definitions = definitionSandbox()
	hook(ctx, t, w)
}

func fixtureRegion(t *testing.T, hook fixtureRegionFunc) {
	cities := []NamedCity{
		{Name: uuid.New().String(), ID: 1},
		{Name: uuid.New().String(), ID: 2},
		{Name: uuid.New().String(), ID: 3},
	}
	fixtureWorld(t, func(ctx context.Context, t *testing.T, w *World) {
		r, err := w.CreateRegion(uuid.New().String(), uuid.New().String(), cities)
		if err != nil {
			t.Fatal(err)
		}
		hook(ctx, t, r)
	})
}

func TestWorldCreateRegion(t *testing.T) {
	fixtureWorld(t, func(_ context.Context, t *testing.T, w *World) {
		cities := make([]NamedCity, 0)
		for i := uint64(0); i < 64; i++ {
			cities = append(cities, NamedCity{Name: uuid.New().String(), ID: i})
		}
		for i := 0; i < len(cities); i++ {
			name := "region-" + uuid.New().String()
			mapName := "map-" + uuid.New().String()
			region, err := w.CreateRegion(name, mapName, cities[:i])
			if err != nil {
				t.Fatal(err)
			}
			_, err = w.CreateRegion(name, mapName, cities[:i])
			if err == nil {
				t.Fatal("duplicated region creation succeeded")
			}
			if region.Name != name || region.MapName != mapName {
				t.Fatal("Invalid region parameters")
			}
			if region.Cities.Len() != i {
				t.Fatal("Wrong number of cities")
			}
		}
	})
}
