// Copyright (c) 2018-2021 Contributors as noted in the AUTHORS file
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package regagent

import (
	"github.com/hegemonie-rpg/engine/pkg/region/model"
	proto "github.com/hegemonie-rpg/engine/pkg/region/proto"
)

// M2P -> Model to Proto
func resMultM2P(r region.ResourcesMultiplier) *proto.ResourcesMult {
	rm := proto.ResourcesMult{}
	// Fuck, protobuf has no array of fixed size
	rm.R0 = r[0]
	rm.R1 = r[1]
	rm.R2 = r[2]
	rm.R3 = r[3]
	rm.R4 = r[4]
	rm.R5 = r[5]
	return &rm
}

// M2P -> Model to Proto
func resPlusM2P(r region.ResourcesIncrement) *proto.ResourcesPlus {
	rm := proto.ResourcesPlus{}
	// Fuck, protobuf has no array of fixed size
	rm.R0 = r[0]
	rm.R1 = r[1]
	rm.R2 = r[2]
	rm.R3 = r[3]
	rm.R4 = r[4]
	rm.R5 = r[5]
	return &rm
}

// M2P -> Model to Proto
func resAbsM2P(r region.Resources) *proto.ResourcesAbs {
	rm := proto.ResourcesAbs{}
	// Fuck, protobuf has no array of fixed size
	rm.R0 = r[0]
	rm.R1 = r[1]
	rm.R2 = r[2]
	rm.R3 = r[3]
	rm.R4 = r[4]
	rm.R5 = r[5]
	return &rm
}

func resAbsP2M(rm *proto.ResourcesAbs) region.Resources {
	r := region.Resources{}
	r[0] = rm.R0
	r[1] = rm.R1
	r[2] = rm.R2
	r[3] = rm.R3
	r[4] = rm.R4
	r[5] = rm.R5
	return r
}

// M2P -> Model to Proto
func resModM2P(r region.ResourceModifiers) *proto.ResourcesMod {
	rm := proto.ResourcesMod{}
	rm.Mult = resMultM2P(r.Mult)
	rm.Plus = resPlusM2P(r.Plus)
	return &rm
}

func showProduction(r *region.Region, c *region.City) *proto.ProductionView {
	v := &proto.ProductionView{}
	prod := c.ComputeProduction(r)
	v.Base = resAbsM2P(prod.Base)
	v.Buildings = resModM2P(prod.Buildings)
	v.Knowledge = resModM2P(prod.Knowledge)
	v.Actual = resAbsM2P(prod.Actual)
	return v
}

func showStock(r *region.Region, c *region.City) *proto.StockView {
	v := &proto.StockView{}
	stock := c.ComputeStock(r)
	v.Base = resAbsM2P(stock.Base)
	v.Buildings = resModM2P(stock.Buildings)
	v.Knowledge = resModM2P(stock.Knowledge)
	v.Actual = resAbsM2P(stock.Actual)
	v.Usage = resAbsM2P(stock.Usage)
	return v
}

func showEvolution(r *region.Region, c *region.City) *proto.CityEvolution {
	cv := &proto.CityEvolution{}
	for _, kt := range c.KnowledgeFrontier(r) {
		cv.KFrontier = append(cv.KFrontier, &proto.KnowledgeTypeView{
			Id: kt.ID, Name: kt.Name,
		})
	}
	for _, bt := range c.BuildingFrontier(r) {
		cv.BFrontier = append(cv.BFrontier, &proto.BuildingTypeRef{
			Id: bt.ID, Name: bt.Name,
		})
	}
	for _, ut := range c.UnitFrontier(r) {
		cv.UFrontier = append(cv.UFrontier, &proto.UnitTypeView{
			Id: ut.ID, Name: ut.Name,
		})
	}
	return cv
}

func showAssets(r *region.Region, c *region.City) *proto.CityAssets {
	v := &proto.CityAssets{}

	for _, k := range c.Knowledges {
		v.Knowledges = append(v.Knowledges, &proto.KnowledgeView{
			Id: k.ID, IdType: k.Type, Ticks: uint32(k.Ticks),
		})
	}
	for _, b := range c.Buildings {
		v.Buildings = append(v.Buildings, &proto.BuildingView{
			Id: b.ID, IdType: b.Type, Ticks: uint32(b.Ticks),
		})
	}
	for _, u := range c.Units {
		v.Units = append(v.Units, &proto.UnitView{
			Id: u.ID, IdType: u.Type, Ticks: uint32(u.Ticks), Health: u.Health,
		})
	}

	for _, a := range c.Armies {
		v.Armies = append(v.Armies, &proto.ArmyView{
			Id: a.ID, Name: a.Name, Location: a.Cell,
			Stock: resAbsM2P(a.Stock),
		})
	}

	return v
}

func showCityConstants(r *region.Region, c *region.City) *proto.CityConstants {
	return &proto.CityConstants{
		Cult:      c.Cult,
		Chaos:     c.Chaotic,
		Alignment: c.Alignment,
		Ethny:     c.EthnicGroup,
		Politics:  c.PoliticalGroup,
	}
}

func showCity(r *region.Region, c *region.City) *proto.CityView {
	cv := &proto.CityView{
		Key:       showCityKey(r, c),
		Constants: showCityConstants(r, c),
		Variables: showCityVars(r, c),
		Politics: &proto.CityPolitics{
			Overlord: c.Overlord,
			Lieges:   []uint64{},
		},
	}

	for _, c := range c.GetLieges() {
		cv.Politics.Lieges = append(cv.Politics.Lieges, c.ID)
	}

	cv.Evol = showEvolution(r, c)
	cv.Production = showProduction(r, c)
	cv.Stock = showStock(r, c)
	cv.Assets = showAssets(r, c)
	return cv
}

func showArmyCommand(c *region.Command) *proto.ArmyCommand {
	cmd := proto.ArmyCommand{Type: proto.ArmyCommandType_Unknown, Target: c.Cell}
	switch c.Action {
	case region.CmdMove:
		cmd.Type = proto.ArmyCommandType_Move
		cmd.Move = &proto.ArmyMoveArgs{}
	case region.CmdWait:
		cmd.Type = proto.ArmyCommandType_Wait
	case region.CmdCityAttack:
		cmd.Type = proto.ArmyCommandType_Attack
		cmd.Attack = &proto.ArmyAssaultArgs{}
	case region.CmdCityDefend:
		cmd.Type = proto.ArmyCommandType_Defend
	case region.CmdCityDisband:
		cmd.Type = proto.ArmyCommandType_Disband
	}
	return &cmd
}

func showArmy(r *region.Region, a *region.Army) *proto.ArmyView {
	view := &proto.ArmyView{
		Id:       a.ID,
		Name:     a.Name,
		Location: a.Cell,
		Stock:    resAbsM2P(a.Stock),
	}
	for _, u := range a.Units {
		view.Units = append(view.Units, showUnit(r, u))
	}
	for _, c := range a.Targets {
		view.Commands = append(view.Commands, showArmyCommand(&c))
	}
	return view
}

func showUnit(r *region.Region, u *region.Unit) *proto.UnitView {
	return &proto.UnitView{
		Id:     u.ID,
		IdType: u.Type,
		Ticks:  u.Ticks,
		Health: u.Health,
	}
}

func showCityKey(r *region.Region, c *region.City) *proto.CityKey {
	return &proto.CityKey{
		City:   c.ID,
		Name:   c.Name,
		Region: r.Name,
	}
}

func showCityPublic(r *region.Region, c *region.City) *proto.CityView {
	return &proto.CityView{
		Key:   showCityKey(r, c),
		Constants: &proto.CityConstants{
			Alignment: c.Alignment,
			Chaos:     c.Chaotic,
			Cult:      c.Cult,
			Politics:  c.PoliticalGroup,
			Ethny:     c.EthnicGroup,
		},
	}
}

func showCityVars(r *region.Region, c *region.City) *proto.CityVars {
	return &proto.CityVars{
		TickMassacres:    c.TicksMassacres,
		PermIntelligence: c.PermanentIntelligence,
		PermHealth:       c.PermanentHealth,
		PermPop:          c.PermanentPopularity,
	}
}

func showCityStats(r *region.Region, c *region.City) *proto.CityStats {
	stats := c.ComputeStats(r)
	return &proto.CityStats{
		// Gauges
		StockUsage:     resAbsM2P(stats.StockUsage),
		StockCapacity:  resAbsM2P(stats.StockCapacity),
		ScoreArmy:      stats.ScoreMilitary,
		ScoreBuilding:  stats.ScoreBuildings,
		ScoreKnowledge: stats.ScoreKnowledge,
		// Counters
		ResourceProduced: resAbsM2P(stats.Activity.ResourceProduced),
		ResourceReceived: resAbsM2P(stats.Activity.ResourceReceived),
		ResourceSent:     resAbsM2P(stats.Activity.ResourceSent),
		FightJoined:      stats.Activity.FightsJoined,
		FightLeft:        stats.Activity.FightsLeft,
		FightLost:        stats.Activity.FightsLost,
		FightWon:         stats.Activity.FightsWon,
		Moves:            stats.Activity.Moves,
		TaxReceived:      resAbsM2P(stats.Activity.TaxReceived),
		TaxSent:          resAbsM2P(stats.Activity.TaxSent),
		UnitLost:         stats.Activity.UnitsLost,
		UnitRaised:       stats.Activity.UnitsRaised,
	}
}

func showCityTemplate(r *region.Region, c *region.City) *proto.CityTemplate {
	rc := &proto.CityTemplate{}
	rc.Constants = showCityConstants(r, c)
	rc.Stock = resAbsM2P(c.Stock)
	rc.StockCapacity = resAbsM2P(c.StockCapacity)
	rc.Production = resAbsM2P(c.Production)
	rc.BuildingTypes = make([]uint64, 0)
	rc.SkillTypes = make([]uint64, 0)
	rc.UnitTypes = make([]uint64, 0)
	for _, x := range c.Buildings {
		rc.BuildingTypes = append(rc.BuildingTypes, x.Type)
	}
	for _, x := range c.Knowledges {
		rc.SkillTypes = append(rc.SkillTypes, x.Type)
	}
	for _, x := range c.Units {
		rc.UnitTypes = append(rc.UnitTypes, x.Type)
	}
	return rc
}

func showCityStatsRecord(r *region.Region, c *region.City) *proto.CityStatsRecord {
	return &proto.CityStatsRecord{
		City:  showCityKey(r, c),
		Stats: showCityStats(r, c),
	}
}
