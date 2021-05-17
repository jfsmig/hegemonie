// Copyright (c) 2018-2021 Contributors as noted in the AUTHORS file
// This Source Code Form is subject to the terms of the Mozilla Constants
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package regclient

import (
	region "github.com/hegemonie-rpg/engine/pkg/region/model"
	proto "github.com/hegemonie-rpg/engine/pkg/region/proto"
	"github.com/juju/errors"
	"strconv"
)

type _cityID struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

type _cityKey struct {
	ID     uint64 `json:"id"`
	Name   string `json:"name"`
	Region string `json:"region"`
}

type _cityConstants struct {
	Alignment int32  `json:"alignment"`
	Chaos     int32  `json:"chaos"`
	Politics  uint32 `json:"politics"`
	Cult      uint32 `json:"cult"`
	Ethny     uint32 `json:"ethny"`
	Score     int64  `json:"score"`
}

type _cityStats struct {
	// Gauges
	StockCapacity _resourcesAbs `json:"stockCapacity"`
	StockUsage    _resourcesAbs `json:"stockUsage"`
	ScoreBuilding uint64        `json:"scoreBuilding"`
	ScoreSkill    uint64        `json:"scoreSkill"`
	ScoreArmy     uint64        `json:"scoreArmy"`
	Popularity    int64         `json:"popularity"`
	// Counters
	ResourceProduced _resourcesAbs `json:"resourceProduced"`
	ResourceSent     _resourcesAbs `json:"resourceSent"`
	ResourceReceived _resourcesAbs `json:"resourceReceived"`
	TaxSent          _resourcesAbs `json:"taxSent"`
	TaxReceived      _resourcesAbs `json:"taxReceived"`
	Moves            uint64        `json:"moves"`
	UnitRaised       uint64        `json:"unitRaised"`
	UnitLost         uint64        `json:"unitLost"`
	FightJoined      uint64        `json:"fightJoined"`
	FightLeft        uint64        `json:"fightLeft"`
	FightWon         uint64        `json:"fightWon"`
	FightLost        uint64        `json:"fightLost"`
}

type _cityTemplate struct {
	Name          string                     `json:"name"`
	Constants     _cityConstants             `json:"constants"`
	Stock         [region.ResourceMax]uint64 `json:"stock"`
	StockCapacity [region.ResourceMax]uint64 `json:"capacity"`
	Production    [region.ResourceMax]uint64 `json:"production"`
	BuildingTypes []string                   `json:"buildings"`
	SkillTypes    []string                   `json:"skills"`
	UnitTypes     []string                   `json:"troops"`
}

type _cityStatsRecord struct {
	ID    _cityID    `json:"id"`
	Stats _cityStats `json:"stats"`
}

type _cityPolitics struct {
	Overlord uint64   `json:"overlord"`
	Lieges   []uint64 `json:"lieges"`
}

type _cityStock struct {
	Base      _resourcesAbs       `json:"base"`
	Skills    _resourcesModifiers `json:"skills"`
	Buildings _resourcesModifiers `json:"buildings"`
	Troops    _resourcesModifiers `json:"troops"`
	Actual    _resourcesAbs       `json:"actual"`
	Usage     _resourcesAbs       `json:"usage"`
}

type _cityProduction struct {
	Base      _resourcesAbs       `json:"base"`
	Skills    _resourcesModifiers `json:"skills"`
	Buildings _resourcesModifiers `json:"buildings"`
	Troops    _resourcesModifiers `json:"troops"`
	Actual    _resourcesAbs       `json:"actual"`
}

type _cityLightView struct {
	ID         _cityKey      `json:"id"`
	Politics   _cityPolitics `json:"politics"`
	Stock      _resourcesAbs `json:"stock"`
	Production _resourcesAbs `json:"production"`
}

type _cityView struct {
	ID         _cityKey        `json:"id"`
	Politics   _cityPolitics   `json:"politics"`
	Stock      _cityStock      `json:"stock"`
	Production _cityProduction `json:"production"`
	Assets     _cityAssets     `json:"assets"`
	Evolution  _cityEvolution  `json:"evolution"`
}

type _cityAssets struct {
	Troops    []_troopView    `json:"troops"`
	Buildings []_buildingView `json:"buildings"`
	Skills    []_skillView    `json:"skills"`
	Armies    []_armyView     `json:"armies"`
}

type _cityEvolution struct {
	Troops    []_troopTypeView    `json:"troops"`
	Buildings []_buildingTypeView `json:"buildings"`
	Skills    []_skillTypeView    `json:"skills"`
}

type _troopView struct {
	ID     string `json:"id"`
	TypeID string `json:"type"`
	Ticks  uint32 `json:"ticks"`
	Health uint32 `json:"health"`
}

type _skillView struct {
	ID     string `json:"id"`
	TypeID string `json:"type"`
	Ticks  uint32 `json:"ticks"`
}

type _buildingView struct {
	ID     string `json:"id"`
	TypeID string `json:"type"`
	Ticks  uint32 `json:"ticks"`
}

type _armyView struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Location uint64 `json:"location"`
}

type _troopTypeView struct {
	ID     string `json:"id"`
	Name   string `json:"type"`
	Ticks  uint32 `json:"ticks"`
	Health uint32 `json:"health"`
}

type _skillTypeView struct {
	ID    string `json:"id"`
	Name  string `json:"type"`
	Ticks uint32 `json:"ticks"`
}

type _buildingTypeView struct {
	ID    string `json:"id"`
	Name  string `json:"type"`
	Ticks uint32 `json:"ticks"`
}

type _resourcesModifiers struct {
	Plus _resourcesPlus `json:"plus"`
	Mult _resourcesMult `json:"mult"`
}

type _resourcesAbs struct {
	R0 uint64 `json:"r0"`
	R1 uint64 `json:"r1"`
	R2 uint64 `json:"r2"`
	R3 uint64 `json:"r3"`
	R4 uint64 `json:"r4"`
	R5 uint64 `json:"r5"`
}

type _resourcesPlus struct {
	R0 int64 `json:"r0"`
	R1 int64 `json:"r1"`
	R2 int64 `json:"r2"`
	R3 int64 `json:"r3"`
	R4 int64 `json:"r4"`
	R5 int64 `json:"r5"`
}

type _resourcesMult struct {
	R0 float64 `json:"r0"`
	R1 float64 `json:"r1"`
	R2 float64 `json:"r2"`
	R3 float64 `json:"r3"`
	R4 float64 `json:"r4"`
	R5 float64 `json:"r5"`
}

func (t _cityTemplate) exportTo() *proto.CityTemplate {
	s := func(i int) uint64 { return t.Stock[i] }
	c := func(i int) uint64 { return t.StockCapacity[i] }
	p := func(i int) uint64 { return t.Production[i] }

	pub := t.Constants.exportTo()

	tpl := proto.CityTemplate{Constants: pub}
	tpl.Stock = &proto.ResourcesAbs{R0: s(0), R1: s(1), R2: s(2), R3: s(3), R4: s(4), R5: s(5)}
	tpl.StockCapacity = &proto.ResourcesAbs{R0: c(0), R1: c(1), R2: c(2), R3: c(3), R4: c(4), R5: c(5)}
	tpl.Production = &proto.ResourcesAbs{R0: p(0), R1: p(1), R2: p(2), R3: p(3), R4: p(4), R5: p(5)}
	tpl.BuildingTypes = t.BuildingTypes
	tpl.SkillTypes = t.SkillTypes
	tpl.UnitTypes = t.UnitTypes
	return &tpl
}

func emptyCityTemplate() _cityTemplate {
	return _cityTemplate{
		SkillTypes:    make([]string, 0),
		BuildingTypes: make([]string, 0),
		UnitTypes:     make([]string, 0),
	}
}

func (t _cityConstants) exportTo() *proto.CityConstants {
	return &proto.CityConstants{
		Ethny:     t.Ethny,
		Cult:      t.Cult,
		Politics:  t.Politics,
		Chaos:     t.Chaos,
		Alignment: t.Alignment,
	}
}

func (t *_troopView) importFrom(in *proto.UnitView) *_troopView {
	if in != nil {
		t.ID = in.Id
		t.TypeID = in.IdType
		t.Ticks = in.Ticks
		t.Health = in.Health
	}
	return t
}

func (t *_buildingView) importFrom(in *proto.BuildingView) *_buildingView {
	if in != nil {
		t.ID = in.Id
		t.TypeID = in.IdType
		t.Ticks = in.Ticks
	}
	return t
}

func (t *_skillView) importFrom(in *proto.SkillView) *_skillView {
	if in != nil {
		t.ID = in.Id
		t.TypeID = in.IdType
		t.Ticks = in.Ticks
	}
	return t
}

func (t *_armyView) importFrom(in *proto.ArmyView) *_armyView {
	if in != nil {
		t.ID = in.Id
		t.Name = in.Name
		t.Location = in.Location
	}
	return t
}

func (t *_troopTypeView) importFrom(in *proto.UnitTypeView) *_troopTypeView {
	if in != nil {
		t.ID = in.Id
		t.Name = in.Name
		t.Ticks = in.Ticks
		t.Health = in.Health
	}
	return t
}

func (t *_buildingTypeView) importFrom(in *proto.BuildingTypeRef) *_buildingTypeView {
	if in != nil {
		t.ID = in.Id
		t.Name = in.Name
		t.Ticks = in.Ticks
		// TODO(jfs)
	}
	return t
}

func (t *_skillTypeView) importFrom(in *proto.SkillTypeView) *_skillTypeView {
	if in != nil {
		t.ID = in.Id
		t.Name = in.Name
		t.Ticks = in.Ticks
	}
	return t
}

func (t *_cityConstants) importFrom(in *proto.CityConstants) *_cityConstants {
	if in != nil {
		t.Ethny = in.Ethny
		t.Cult = in.Cult
		t.Chaos = in.Chaos
		t.Politics = in.Politics
		t.Alignment = in.Alignment
	}
	return t
}

func (t *_cityID) importFrom(in *proto.CityKey) *_cityID {
	if in != nil {
		t.ID = in.City
		t.Name = in.Name
	}
	return t
}

func (t *_cityKey) importFrom(in *proto.CityKey) *_cityKey {
	if in != nil {
		t.ID = in.City
		t.Name = in.Name
		t.Region = in.Region
	}
	return t
}

func (t *_resourcesAbs) importFrom(in *proto.ResourcesAbs) *_resourcesAbs {
	if in != nil {
		t.R0 = in.R0
		t.R1 = in.R1
		t.R2 = in.R2
		t.R3 = in.R3
		t.R4 = in.R4
		t.R5 = in.R5
	}
	return t
}

func (t *_resourcesPlus) importFrom(in *proto.ResourcesPlus) *_resourcesPlus {
	if in != nil {
		t.R0 = in.R0
		t.R1 = in.R1
		t.R2 = in.R2
		t.R3 = in.R3
		t.R4 = in.R4
		t.R5 = in.R5
	}
	return t
}

func (t *_resourcesMult) importFrom(in *proto.ResourcesMult) *_resourcesMult {
	if in != nil {
		t.R0 = in.R0
		t.R1 = in.R1
		t.R2 = in.R2
		t.R3 = in.R3
		t.R4 = in.R4
		t.R5 = in.R5
	}
	return t
}

func (t *_resourcesModifiers) importFrom(in *proto.ResourcesMod) *_resourcesModifiers {
	if in != nil {
		t.Plus.importFrom(in.Plus)
		t.Mult.importFrom(in.Mult)
	}
	return t
}

func (t *_cityStock) importFrom(in *proto.StockView) *_cityStock {
	if in != nil {
		t.Base.importFrom(in.Base)
		t.Skills.importFrom(in.Skill)
		t.Buildings.importFrom(in.Buildings)
		t.Troops.importFrom(in.Troops)
		t.Actual.importFrom(in.Actual)
		t.Usage.importFrom(in.Usage)
	}
	return t
}

func (t *_cityProduction) importFrom(in *proto.ProductionView) *_cityProduction {
	if in != nil {
		t.Base.importFrom(in.Base)
		t.Skills.importFrom(in.Skill)
		t.Buildings.importFrom(in.Buildings)
		t.Troops.importFrom(in.Troops)
		t.Actual.importFrom(in.Actual)
	}
	return t
}

func (t *_cityAssets) importFrom(in *proto.CityAssets) *_cityAssets {
	if in != nil {
		t.Troops = make([]_troopView, 0)
		for _, x := range in.Units {
			var v _troopView
			v.importFrom(x)
			t.Troops = append(t.Troops, v)
		}

		t.Skills = make([]_skillView, 0)
		for _, x := range in.Skills {
			var v _skillView
			v.importFrom(x)
			t.Skills = append(t.Skills, v)
		}

		t.Buildings = make([]_buildingView, 0)
		for _, x := range in.Buildings {
			var v _buildingView
			v.importFrom(x)
			t.Buildings = append(t.Buildings, v)
		}

		t.Armies = make([]_armyView, 0)
		for _, x := range in.Armies {
			var v _armyView
			v.importFrom(x)
			t.Armies = append(t.Armies, v)
		}
	}
	return t
}

func (t *_cityEvolution) importFrom(in *proto.CityEvolution) *_cityEvolution {
	if in != nil {
		t.Troops = make([]_troopTypeView, 0)
		for _, x := range in.UFrontier {
			var v _troopTypeView
			v.importFrom(x)
			t.Troops = append(t.Troops, v)
		}

		t.Skills = make([]_skillTypeView, 0)
		for _, x := range in.KFrontier {
			var v _skillTypeView
			v.importFrom(x)
			t.Skills = append(t.Skills, v)
		}

		t.Buildings = make([]_buildingTypeView, 0)
		for _, x := range in.BFrontier {
			var v _buildingTypeView
			v.importFrom(x)
			t.Buildings = append(t.Buildings, v)
		}
	}
	return t
}

func (t *_cityStats) importFrom(in *proto.CityStats) *_cityStats {
	if in != nil {
		t.StockCapacity.importFrom(in.StockCapacity)
		t.StockUsage.importFrom(in.StockUsage)

		t.ScoreBuilding = in.ScoreBuilding
		t.ScoreSkill = in.ScoreSkill
		t.ScoreArmy = in.ScoreArmy
		t.Popularity = in.Popularity

		t.ResourceProduced.importFrom(in.ResourceProduced)
		t.ResourceSent.importFrom(in.ResourceSent)
		t.ResourceReceived.importFrom(in.ResourceReceived)
		t.TaxSent.importFrom(in.TaxSent)
		t.TaxReceived.importFrom(in.TaxReceived)

		t.Moves = in.Moves
		t.UnitRaised = in.UnitRaised
		t.UnitLost = in.UnitLost
		t.FightJoined = in.FightJoined
		t.FightLeft = in.FightLeft
		t.FightWon = in.FightWon
		t.FightLost = in.FightLost
	}
	return t
}

func (t *_cityView) importFrom(in *proto.CityView) *_cityView {
	if in != nil {
		t.ID.importFrom(in.Key)
		t.Politics.importFrom(in.Politics)
		t.Stock.importFrom(in.Stock)
		t.Production.importFrom(in.Production)
		t.Assets.importFrom(in.Assets)
		t.Evolution.importFrom(in.Evol)
	}
	return t
}

func (t *_cityLightView) importFrom(in *proto.CityView) *_cityLightView {
	if in != nil {
		t.ID.importFrom(in.Key)
		t.Politics.importFrom(in.Politics)
		t.Stock.importFrom(in.Stock.Actual)
		t.Production.importFrom(in.Production.Actual)
	}
	return t
}

func (t *_cityPolitics) importFrom(in *proto.CityPolitics) *_cityPolitics {
	if in != nil {
		t.Overlord = in.Overlord
		t.Lieges = make([]uint64, len(in.Lieges), len(in.Lieges))
		copy(t.Lieges, in.Lieges)
	}
	return t
}

func (t *_cityStatsRecord) importFrom(in *proto.CityStatsRecord) *_cityStatsRecord {
	if in != nil {
		t.ID.importFrom(in.City)
		t.Stats.importFrom(in.Stats)
	}
	return t
}

func (t *_cityTemplate) importFrom(in *proto.CityTemplate) *_cityTemplate {
	if in != nil {
		t.Name = in.Name
		t.Constants.importFrom(in.Constants)

		t.Stock[0] = in.Stock.R0
		t.Stock[1] = in.Stock.R1
		t.Stock[2] = in.Stock.R2
		t.Stock[3] = in.Stock.R3
		t.Stock[4] = in.Stock.R4
		t.Stock[5] = in.Stock.R5

		t.StockCapacity[0] = in.StockCapacity.R0
		t.StockCapacity[1] = in.StockCapacity.R1
		t.StockCapacity[2] = in.StockCapacity.R2
		t.StockCapacity[3] = in.StockCapacity.R3
		t.StockCapacity[4] = in.StockCapacity.R4
		t.StockCapacity[5] = in.StockCapacity.R5

		t.Production[0] = in.Production.R0
		t.Production[1] = in.Production.R1
		t.Production[2] = in.Production.R2
		t.Production[3] = in.Production.R3
		t.Production[4] = in.Production.R4
		t.Production[5] = in.Production.R5

		in.BuildingTypes = make([]string, 0)
		if in.BuildingTypes != nil {
			t.BuildingTypes = in.BuildingTypes
		}
		in.SkillTypes = make([]string, 0)
		if in.SkillTypes != nil {
			t.SkillTypes = in.SkillTypes
		}
		in.UnitTypes = make([]string, 0)
		if in.UnitTypes != nil {
			t.UnitTypes = in.UnitTypes
		}
	}
	return t
}

func (t *_cityTemplate) set(k, vs string) error {
	i64, err := strconv.ParseInt(vs, 10, 63)
	if err != nil {
		return errors.NewNotValid(err, "invalid value")
	}
	switch k {
	case "align":
		t.Constants.Alignment = int32(i64)
	case "chaos":
		t.Constants.Chaos = int32(i64)
	case "politics":
		t.Constants.Politics = uint32(i64)
	case "cult":
		t.Constants.Cult = uint32(i64)
	case "ethny":
		t.Constants.Ethny = uint32(i64)

	case "stock.0", "s0":
		t.Stock[0] = uint64(i64)
	case "stock.1", "s1":
		t.Stock[1] = uint64(i64)
	case "stock.2", "s2":
		t.Stock[2] = uint64(i64)
	case "stock.3", "s3":
		t.Stock[3] = uint64(i64)
	case "stock.4", "s4":
		t.Stock[4] = uint64(i64)
	case "stock.5", "s5":
		t.Stock[5] = uint64(i64)

	case "capa.0", "c0":
		t.StockCapacity[0] = uint64(i64)
	case "capa.1", "c1":
		t.StockCapacity[1] = uint64(i64)
	case "capa.2", "c2":
		t.StockCapacity[2] = uint64(i64)
	case "capa.3", "c3":
		t.StockCapacity[3] = uint64(i64)
	case "capa.4", "c4":
		t.StockCapacity[4] = uint64(i64)
	case "capa.5", "c5":
		t.StockCapacity[5] = uint64(i64)

	case "prod.0", "p0":
		t.Production[0] = uint64(i64)
	case "prod.1", "p1":
		t.Production[1] = uint64(i64)
	case "prod.2", "p2":
		t.Production[2] = uint64(i64)
	case "prod.3", "p3":
		t.Production[3] = uint64(i64)
	case "prod.4", "p4":
		t.Production[4] = uint64(i64)
	case "prod.5", "p5":
		t.Production[5] = uint64(i64)
	}
	return nil
}

func (t *_cityTemplate) add(k, vs string) error {
	switch k {
	case "b", "build", "building":
		t.BuildingTypes = append(t.BuildingTypes, vs)
	case "s", "skill", "skills":
		t.SkillTypes = append(t.SkillTypes, vs)
	case "u", "unit", "units":
		t.UnitTypes = append(t.UnitTypes, vs)
	}
	return nil
}

func (t *_cityTemplate) del(k, vs string) error {
	filter := func(tab []string, bad string) []string {
		out := make([]string, 0)
		for _, v := range tab {
			if v != bad {
				out = append(out, v)
			}
		}
		return out
	}
	switch k {
	case "b", "build", "building":
		t.BuildingTypes = filter(t.BuildingTypes, vs)
	case "s", "skill", "skills":
		t.SkillTypes = filter(t.SkillTypes, vs)
	case "u", "unit", "units":
		t.UnitTypes = filter(t.UnitTypes, vs)
	}
	return nil
}
