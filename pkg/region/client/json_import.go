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
	StockCapacity  _resourcesAbs `json:"stockCapacity"`
	StockUsage     _resourcesAbs `json:"stockUsage"`
	ScoreBuilding  uint64        `json:"scoreBuilding"`
	ScoreKnowledge uint64        `json:"scoreKnowledge"`
	ScoreArmy      uint64        `json:"scoreArmy"`
	Popularity     int64         `json:"popularity"`
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
	BuildingTypes []uint64                   `json:"buildings"`
	SkillTypes    []uint64                   `json:"skills"`
	UnitTypes     []uint64                   `json:"units"`
}

type _cityStatsRecord struct {
	ID    _cityID    `json:"id"`
	Stats _cityStats `json:"stats"`
}

type _resourcesAbs struct {
	R0 uint64 `json:"r0"`
	R1 uint64 `json:"r1"`
	R2 uint64 `json:"r2"`
	R3 uint64 `json:"r3"`
	R4 uint64 `json:"r4"`
	R5 uint64 `json:"r5"`
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
		SkillTypes:    make([]uint64, 0),
		BuildingTypes: make([]uint64, 0),
		UnitTypes:     make([]uint64, 0),
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

func (t *_cityConstants) importFrom(tpl *proto.CityConstants) *_cityConstants {
	t.Ethny = tpl.Ethny
	t.Cult = tpl.Cult
	t.Chaos = tpl.Chaos
	t.Politics = tpl.Politics
	t.Alignment = tpl.Alignment
	return t
}

func (t *_cityID) importFrom(k *proto.CityKey) *_cityID {
	t.ID = k.City
	t.Name = k.Name
	return t
}

func (t *_cityTemplate) importFrom(tpl *proto.CityTemplate) {
	t.Name = tpl.Name
	t.Constants.importFrom(tpl.Constants)

	t.Stock[0] = tpl.Stock.R0
	t.Stock[1] = tpl.Stock.R1
	t.Stock[2] = tpl.Stock.R2
	t.Stock[3] = tpl.Stock.R3
	t.Stock[4] = tpl.Stock.R4
	t.Stock[5] = tpl.Stock.R5

	t.StockCapacity[0] = tpl.StockCapacity.R0
	t.StockCapacity[1] = tpl.StockCapacity.R1
	t.StockCapacity[2] = tpl.StockCapacity.R2
	t.StockCapacity[3] = tpl.StockCapacity.R3
	t.StockCapacity[4] = tpl.StockCapacity.R4
	t.StockCapacity[5] = tpl.StockCapacity.R5

	t.Production[0] = tpl.Production.R0
	t.Production[1] = tpl.Production.R1
	t.Production[2] = tpl.Production.R2
	t.Production[3] = tpl.Production.R3
	t.Production[4] = tpl.Production.R4
	t.Production[5] = tpl.Production.R5

	empty := make([]uint64, 0)
	tpl.BuildingTypes = empty
	if tpl.BuildingTypes != nil {
		t.BuildingTypes = tpl.BuildingTypes
	}
	tpl.SkillTypes = empty
	if tpl.SkillTypes != nil {
		t.SkillTypes = tpl.SkillTypes
	}
	tpl.UnitTypes = empty
	if tpl.UnitTypes != nil {
		t.UnitTypes = tpl.UnitTypes
	}
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
	u64, err := strconv.ParseUint(vs, 10, 63)
	if err != nil {
		return errors.NewNotValid(err, "Invalid value")
	}
	switch k {
	case "b", "build", "building":
		t.BuildingTypes = append(t.BuildingTypes, u64)
	case "s", "skill", "skills":
		t.SkillTypes = append(t.SkillTypes, u64)
	case "u", "unit", "units":
		t.UnitTypes = append(t.UnitTypes, u64)
	}
	return nil
}

func (t *_cityTemplate) del(k, vs string) error {
	filter := func(tab []uint64, bad uint64) []uint64 {
		out := make([]uint64, 0)
		for _, v := range tab {
			if v != bad {
				out = append(out, v)
			}
		}
		return out
	}
	u64, err := strconv.ParseUint(vs, 10, 63)
	if err != nil {
		return errors.NewNotValid(err, "Invalid value")
	}
	switch k {
	case "b", "build", "building":
		t.BuildingTypes = filter(t.BuildingTypes, u64)
	case "s", "skill", "skills":
		t.SkillTypes = filter(t.SkillTypes, u64)
	case "u", "unit", "units":
		t.UnitTypes = filter(t.UnitTypes, u64)
	}
	return nil
}

func (t *_resourcesAbs) importFrom(in *proto.ResourcesAbs) *_resourcesAbs {
	t.R0 = in.R0
	t.R1 = in.R1
	t.R2 = in.R2
	t.R3 = in.R3
	t.R4 = in.R4
	t.R5 = in.R5
	return t
}

func (t *_cityStats) importFrom(in *proto.CityStats) *_cityStats {
	t.StockCapacity.importFrom(in.StockCapacity)
	t.StockUsage.importFrom(in.StockUsage)

	t.ScoreBuilding = in.ScoreBuilding
	t.ScoreKnowledge = in.ScoreKnowledge
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
	return t
}

func (t *_cityStatsRecord) importFrom(in *proto.CityStatsRecord) *_cityStatsRecord {
	t.ID.importFrom(in.City)
	t.Stats.importFrom(in.Stats)
	return t
}
