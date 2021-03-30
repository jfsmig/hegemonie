// Copyright (c) 2018-2021 Contributors as noted in the AUTHORS file
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package region

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/juju/errors"
)

func MakeCity() *City {
	return &City{
		ID:         0,
		State:      CityStateIdle,
		Units:      make(SetOfUnits, 0),
		Buildings:  make(SetOfBuildings, 0),
		Knowledges: make(SetOfKnowledges, 0),
		Armies:     make(SetOfArmies, 0),
		lieges:     make(SetOfCities, 0),
	}
}

func (c *City) Apply(model *City) {
	c.Owner = ""
	c.Overlord = 0
	c.TaxRate.SetValue(0.0)
	c.Assault = nil

	c.Chaotic = model.Chaotic
	c.Alignment = model.Alignment
	c.EthnicGroup = model.EthnicGroup
	c.PoliticalGroup = model.PoliticalGroup
	c.Cult = model.Cult

	c.PermanentPopularity = 0
	c.PermanentHealth = model.PermanentHealth
	c.PermanentIntelligence = model.PermanentIntelligence
	c.TicksMassacres = 0

	c.Stock.Set(model.Stock)
	c.Production.Set(model.Production)
	c.StockCapacity.Set(model.StockCapacity)

	c.Buildings = make(SetOfBuildings, 0)
	for _, x0 := range model.Buildings {
		x := new(Building)
		*x = *x0
		x.ID = uuid.New().String()
		x.Ticks = 0
		c.Buildings.Add(x)
	}

	c.Knowledges = make(SetOfKnowledges, 0)
	for _, x0 := range model.Knowledges {
		x := new(Knowledge)
		*x = *x0
		x.ID = uuid.New().String()
		x.Ticks = 0
		c.Knowledges.Add(x)
	}

	c.Units = make(SetOfUnits, 0)
	for _, x0 := range model.Units {
		x := new(Unit)
		*x = *x0
		x.ID = uuid.New().String()
		x.Ticks = 0
		c.Units.Add(x)
	}

	c.Armies = make(SetOfArmies, 0)
	c.Artifacts = make(SetOfArtifacts, 0)

	c.Counters.Reset()
}

// GetUnit returns a Unit owned by the current City, given the Unit ID
func (c *City) GetUnit(id string) *Unit {
	return c.Units.Get(id)
}

// GetBuilding returns a Building owned by the current City, given the Building
// ID
func (c *City) GetBuilding(id string) *Building {
	return c.Buildings.Get(id)
}

// GetKnowledge returns a Knowledge owned by the current City, given the
// Knowledge ID
func (c *City) GetKnowledge(id string) *Knowledge {
	return c.Knowledges.Get(id)
}

// GetActualPopularity returns the total Popularity of the current City
// (permanent + transient)
func (c *City) GetActualPopularity(w *World) int64 {
	var pop int64 = c.PermanentPopularity

	// Add Transient values for Units in the Armies
	for _, a := range c.Armies {
		for _, u := range a.Units {
			ut := w.UnitTypeGet(u.Type)
			pop += ut.PopBonus
		}
		pop += w.Config.PopBonusArmyAlive
	}

	// Add Transient values for Units in the City
	for _, u := range c.Units {
		ut := w.UnitTypeGet(u.Type)
		pop += ut.PopBonus
	}

	// Add Transient values for Buildings
	for _, b := range c.Buildings {
		bt := w.BuildingTypeGet(b.Type)
		pop += bt.PopBonus
	}

	// Add Transient values for Knowledges
	for _, k := range c.Knowledges {
		kt := w.KnowledgeTypeGet(k.Type)
		pop += kt.PopBonus
	}

	return pop
}

// GetProduction computes the actual production of the local City,
// and a summary of the main steps leading to the result. In other words,
// a summary of all the City assets that have an impact.
func (c *City) GetProduction(w *World) *CityProduction {
	p := &CityProduction{
		Buildings: ResourceModifierNoop(),
		Knowledge: ResourceModifierNoop(),
	}

	for _, b := range c.Buildings {
		t := w.BuildingTypeGet(b.Type)
		p.Buildings.ComposeWith(t.Prod)
	}
	for _, u := range c.Knowledges {
		t := w.KnowledgeTypeGet(u.Type)
		p.Knowledge.ComposeWith(t.Prod)
	}

	p.Base = c.Production.Copy()
	p.Actual = c.Production.Copy()
	p.Actual.Apply(p.Buildings, p.Knowledge)
	return p
}

// GetStock computes the actual stock of the local City, and a summary of the
// main steps leading to the result. In other words, a summary of all the City
// assets that have an impact.
func (c *City) GetStock(w *World) *CityStock {
	p := &CityStock{
		Buildings: ResourceModifierNoop(),
		Knowledge: ResourceModifierNoop(),
	}

	for _, b := range c.Buildings {
		t := w.BuildingTypeGet(b.Type)
		p.Buildings.ComposeWith(t.Stock)
	}
	for _, b := range c.Knowledges {
		t := w.BuildingTypeGet(b.Type)
		p.Buildings.ComposeWith(t.Stock)
	}

	p.Base = c.StockCapacity.Copy()
	p.Actual = c.StockCapacity.Copy()
	p.Actual.Apply(p.Buildings, p.Knowledge)
	p.Usage = c.Stock.Copy()
	return p
}

// GetLieges returns a list of all the Lieges of the current City.
func (c *City) GetLieges() []*City {
	return c.lieges[:]
}

// GetStats computes the gauges and extract the counters to build a CityStats
// about the current City.
func (c *City) GetStats(w *Region) CityStats {
	stock := c.GetStock(w.world)
	return CityStats{
		Activity:       c.Counters,
		StockCapacity:  stock.Actual,
		StockUsage:     stock.Usage,
		ScoreBuildings: uint64(c.Buildings.Len()),
		ScoreKnowledge: uint64(c.Knowledges.Len()),
		ScoreMilitary:  uint64(c.Armies.Len()),
	}
}

func (c *City) CreateEmptyArmy(w *Region) *Army {
	aid := uuid.New().String()
	a := &Army{
		ID:       aid,
		City:     c,
		Cell:     c.ID,
		Name:     fmt.Sprintf("A-%v", aid),
		Units:    make(SetOfUnits, 0),
		Postures: []int64{int64(c.ID)},
		Targets:  make([]Command, 0),
	}
	c.Armies.Add(a)
	return a
}

func unitsToIDs(uv []*Unit) (out []string) {
	for _, u := range uv {
		out = append(out, u.ID)
	}
	return out
}

func unitsFilterIdle(uv []*Unit) (out []*Unit) {
	for _, u := range uv {
		if u.Health > 0 && u.Ticks <= 0 {
			out = append(out, u)
		}
	}
	return out
}

// Create an Army made of some Unit of the City
func (c *City) CreateArmyFromUnit(w *Region, units ...*Unit) (*Army, error) {
	return c.CreateArmyFromIds(w, unitsToIDs(unitsFilterIdle(units))...)
}

// Create an Army made of some Unit of the City
func (c *City) CreateArmyFromIds(w *Region, ids ...string) (*Army, error) {
	a := c.CreateEmptyArmy(w)
	err := c.TransferOwnUnit(a, ids...)
	if err != nil { // Rollback
		a.Disband(w, c, false)
		return nil, errors.Annotate(err, "transfer error")
	}
	return a, nil
}

// Create an Army made of all the Units defending the City
func (c *City) CreateArmyDefence(w *Region) (*Army, error) {
	ids := unitsToIDs(unitsFilterIdle(c.Units))
	if len(ids) <= 0 {
		return nil, errors.NotFoundf("unit not found")
	}
	return c.CreateArmyFromIds(w, ids...)
}

// Create an Army carrying resources you own
func (c *City) CreateTransport(w *Region, r Resources) (*Army, error) {
	if !c.Stock.GreaterOrEqualTo(r) {
		return nil, errors.Forbiddenf("insufficient resources")
	}

	a := c.CreateEmptyArmy(w)
	c.Stock.Remove(r)
	a.Stock.Add(r)
	return a, nil
}

// Play one round of local production and return the
func (c *City) ProduceLocally(w *Region, p *CityProduction) Resources {
	var prod Resources = p.Actual
	if c.TicksMassacres > 0 {
		mult := MultiplierUniform(w.world.Config.MassacreImpact)
		for i := uint32(0); i < c.TicksMassacres; i++ {
			prod.Multiply(mult)
		}
		c.TicksMassacres--
	}
	return prod
}

func (c *City) Produce(_ context.Context, w *Region) {
	// Pre-compute the modified values of Stock and Production.
	// We just reuse a function that already does it (despite it does more)
	prod0 := c.GetProduction(w.world)
	stock := c.GetStock(w.world)

	// Make the local City generate resources (and recover the massacres)
	prod := c.ProduceLocally(w, prod0)
	c.Stock.Add(prod)

	if c.Overlord != 0 && c.pOverlord != nil {
		// Compute the expected Tax based on the local production
		var tax Resources = prod
		tax.Multiply(c.TaxRate)
		// Ensure the tax isn't superior to the actual production (to cope with
		// invalid tax rates)
		tax.TrimTo(c.Stock)
		// Then preempt the tax from the stock
		c.Stock.Remove(tax)

		// TODO(jfs): check for potential shortage
		//  shortage := c.Tax.GreaterThan(tax)

		if w.world.Config.InstantTransfers {
			c.pOverlord.Stock.Add(tax)
		} else {
			c.SendResourcesTo(w, c.pOverlord, tax)
		}

		// FIXME(jfs): notify overlord
		// FIXME(jfs): notify c
	}

	// ATM the stock maybe still stores resources. We use them to make the assets evolve.
	// We arbitrarily give the preference to Units, then Buildings and eventually the
	// Knowledge.

	for _, u := range c.Units {
		if u.Ticks > 0 {
			ut := w.world.UnitTypeGet(u.Type)
			if c.Stock.GreaterOrEqualTo(ut.Cost) {
				c.Stock.Remove(ut.Cost)
				u.Ticks--
				if u.Ticks <= 0 {
					// FIXME(jfs): Notify the City that a Unit is OK
				}
			} else {
				// FIXME(jfs): Notify the maintenance fault to the City
			}
		}
	}

	for _, b := range c.Buildings {
		if b.Ticks > 0 {
			bt := w.world.BuildingTypeGet(b.Type)
			if c.Stock.GreaterOrEqualTo(bt.Cost) {
				c.Stock.Remove(bt.Cost)
				b.Ticks--
				if b.Ticks <= 0 {
					// FIXME(jfs): Notify the City
				}
			}
		}
	}

	for _, k := range c.Knowledges {
		if k.Ticks > 0 {
			bt := w.world.KnowledgeTypeGet(k.Type)
			if c.Stock.GreaterOrEqualTo(bt.Cost) {
				c.Stock.Remove(bt.Cost)
				k.Ticks--
			}
			if k.Ticks <= 0 {
				// FIXME(jfs): Notify the City
			}
		}
	}

	// At the end of the turn, ensure we do not hold more resources than the actual
	// stock capacity (with the effect of all the multipliers)
	c.Stock.TrimTo(stock.Actual)
}

// Set a tax rate on the current City, with the same ratio on every Resource.
func (c *City) SetUniformTaxRate(nb float64) {
	c.TaxRate = MultiplierUniform(nb)
}

// Set the given tax rate to the current City.
func (c *City) SetTaxRate(m ResourcesMultiplier) {
	c.TaxRate = m
}

func (c *City) LiberateCity(w *World, other *City) {
	pre := other.pOverlord
	if pre == nil {
		return
	}

	other.Overlord = 0
	other.pOverlord = nil

	// FIXME(jfs): Notify 'pre'
	// FIXME(jfs): Notify 'c'
	// FIXME(jfs): Notify 'other'
}

func (c *City) GainFreedom(w *World) {
	pre := c.pOverlord
	if pre == nil {
		return
	}

	c.Overlord = 0
	c.pOverlord = nil

	// FIXME(jfs): Notify 'pre'
	// FIXME(jfs): Notify 'c'
}

func (c *City) ConquerCity(w *World, other *City) {
	if other.pOverlord == c {
		c.pOverlord = nil
		c.Overlord = 0
		c.TaxRate = MultiplierUniform(0)
		return
	}

	//pre := other.pOverlord
	other.pOverlord = c
	other.Overlord = c.ID
	other.TaxRate = MultiplierUniform(w.Config.RateOverlord)

	// FIXME(jfs): Notify 'pre'
	// FIXME(jfs): Notify 'c'
	// FIXME(jfs): Notify 'other'
}

func (c *City) SendResourcesTo(w *Region, overlord *City, amount Resources) error {
	// FIXME(jfs): NYI
	return errors.New("SendResourcesTo() not implemented")
}

func (c *City) TransferOwnResources(a *Army, r Resources) error {
	if a.City != c {
		return errors.Forbiddenf("army not controlled by the city")
	}
	if !c.Stock.GreaterOrEqualTo(r) {
		return errors.Forbiddenf("insufficient resources")
	}

	c.Stock.Remove(r)
	a.Stock.Add(r)
	return nil
}

func (c *City) TransferOwnUnit(a *Army, units ...string) error {
	if len(units) <= 0 || a == nil {
		panic("EINVAL")
	}

	if a.City != c {
		return errors.Forbiddenf("army not controlled by the city")
	}

	allUnits := make(map[string]*Unit)
	for _, uid := range units {
		if _, ok := allUnits[uid]; ok {
			continue
		}
		if u := c.Units.Get(uid); u == nil {
			return errors.NotFoundf("unit not found")
		} else if u.Ticks > 0 || u.Health <= 0 {
			continue
		} else {
			allUnits[uid] = u
		}
	}

	for _, u := range allUnits {
		c.Units.Remove(u)
		a.Units.Add(u)
	}
	return nil
}

func (c *City) KnowledgeFrontier(w *World) []*KnowledgeType {
	return w.KnowledgeGetFrontier(c.Knowledges)
}

func (c *City) BuildingFrontier(w *World) []*BuildingType {
	return w.BuildingGetFrontier(c.GetActualPopularity(w), c.Buildings, c.Knowledges)
}

// Return a collection of UnitType that may be trained by the current City
// because all the requirements are met.
// Each UnitType 'p' returned validates 'c.UnitAllowed(p)'.
func (c *City) UnitFrontier(w *World) []*UnitType {
	return w.UnitGetFrontier(c.Buildings)
}

// check the current City has all the requirements to train a Unti of the
// given UnitType.
func (c *City) UnitAllowed(t *UnitType) bool {
	if t.RequiredBuilding == 0 {
		return true
	}
	for _, b := range c.Buildings {
		if b.Type == t.RequiredBuilding {
			return true
		}
	}
	return false
}

func (c *City) ownedKnowledgeTypes(reg *Region) SetOfKnowledgeTypes {
	out := make(SetOfKnowledgeTypes, 0)
	for _, k := range c.Knowledges {
		out.Add(reg.world.Definitions.Knowledges.Get(k.Type))
	}
	return out
}

func (cnt *CityActivityCounters) Reset() {
	cnt.ResourceProduced.Zero()
	cnt.ResourceSent.Zero()
	cnt.ResourceReceived.Zero()
	cnt.TaxSent.Zero()
	cnt.TaxReceived.Zero()
	cnt.Moves = 0
	cnt.FightsJoined = 0
	cnt.FightsLeft = 0
	cnt.FightsWon = 0
	cnt.FightsLeft = 0
	cnt.UnitsRaised = 0
	cnt.UnitsLost = 0
}
