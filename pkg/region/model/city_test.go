// Copyright (c) 2018-2021 Contributors as noted in the AUTHORS file
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package region

import (
	"context"
	"github.com/hegemonie-rpg/engine/pkg/utils"
	"testing"
)

func TestCityCreateDuplicate(t *testing.T) {
	fixtureRegion(t, func(ctx context.Context, t *testing.T, r *Region) {
		for _, c := range r.Cities {
			_, err := r.CityCreate(c.ID)
			if err == nil {
				t.Fatal()
			}
		}
	})
}

func TestCityProductionNoModifier(t *testing.T) {
	fixtureRegion(t, func(ctx context.Context, t *testing.T, r *Region) {
		city := r.Cities[0]
		city.Stock.Zero()
		city.StockCapacity.SetValue(2)
		city.Production.SetValue(1)

		expectedProd := ResourcesUniform(1)
		prod := city.GetProduction(r.world)
		if !prod.Actual.Equals(expectedProd) {
			t.Fatal("unexpected stock, got", utils.JSON2Str(prod), "expected", utils.JSON2Str(expectedProd))
		}

		step := prod.Actual
		for _, expectation := range []Resources{step, city.StockCapacity, city.StockCapacity} {
			r.Produce(ctx)
			if !city.Stock.Equals(expectation) {
				t.Fatal("unexpected production, got", city.Stock, "expected", expectation, "step", step)
			}
		}
	})
}

func TestCityProductionWithModifier(t *testing.T) {
	fixtureRegion(t, func(ctx context.Context, t *testing.T, r *Region) {
		city := r.Cities[0]
		city.Stock.Zero()
		city.StockCapacity.SetValue(2)
		city.Production.SetValue(1)

		bt := r.world.Definitions.Buildings[0]
		b := city.StartBuilding(bt)
		b.Ticks = 0

		expectedProd := ResourcesUniform(2)
		prod := city.GetProduction(r.world)
		if !prod.Actual.Equals(expectedProd) {
			t.Fatal("unexpected production, got", utils.JSON2Str(prod), "expected", utils.JSON2Str(expectedProd))
		}

		step := prod.Actual
		for _, expectation := range []Resources{step, city.StockCapacity, city.StockCapacity} {
			r.Produce(ctx)
			if !city.Stock.Equals(expectation) {
				t.Fatal("unexpected stock, got", city.Stock, "expected", expectation, "step", step)
			}
		}
	})
}

func TestCityDefenceCreation(t *testing.T) {
	fixtureRegion(t, func(ctx context.Context, t *testing.T, r *Region) {
		//t.Fail()
	})
}
