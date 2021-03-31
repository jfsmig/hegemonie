// Copyright (c) 2018-2021 Contributors as noted in the AUTHORS file
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package regagent

import (
	"github.com/hegemonie-rpg/engine/pkg/region/model"
	proto "github.com/hegemonie-rpg/engine/pkg/region/proto"
)

func showEvolution(w *region.World, c *region.City) *proto.CityEvolution {
	cv := &proto.CityEvolution{}

	for _, kt := range c.KnowledgeFrontier(w) {
		cv.KFrontier = append(cv.KFrontier, &proto.KnowledgeTypeView{
			Id: kt.ID, Name: kt.Name,
		})
	}
	for _, bt := range c.BuildingFrontier(w) {
		cv.BFrontier = append(cv.BFrontier, &proto.BuildingTypeRef{
			Id: bt.ID, Name: bt.Name,
		})
	}
	for _, ut := range c.UnitFrontier(w) {
		cv.UFrontier = append(cv.UFrontier, &proto.UnitTypeView{
			Id: ut.ID, Name: ut.Name,
		})
	}

	return cv
}

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

func showProduction(w *region.World, c *region.City) *proto.ProductionView {
	v := &proto.ProductionView{}
	prod := c.GetProduction(w)
	v.Base = resAbsM2P(prod.Base)
	v.Buildings = resModM2P(prod.Buildings)
	v.Knowledge = resModM2P(prod.Knowledge)
	v.Actual = resAbsM2P(prod.Actual)
	return v
}

func showStock(w *region.World, c *region.City) *proto.StockView {
	v := &proto.StockView{}
	stock := c.GetStock(w)
	v.Base = resAbsM2P(stock.Base)
	v.Buildings = resModM2P(stock.Buildings)
	v.Knowledge = resModM2P(stock.Knowledge)
	v.Actual = resAbsM2P(stock.Actual)
	v.Usage = resAbsM2P(stock.Usage)
	return v
}

func showAssets(w *region.World, c *region.City) *proto.CityAssets {
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

func showCity(w *region.World, c *region.City) *proto.CityView {
	cv := &proto.CityView{
		Public: &proto.PublicCity{
			Id:   c.ID,
			Name: c.Name,

			Cult:      c.Cult,
			Chaos:     c.Chaotic,
			Alignment: c.Alignment,
			Ethny:     c.EthnicGroup,
			Politics:  c.PoliticalGroup,
		},

		Owner: c.Owner,

		TickMassacres: c.TicksMassacres,

		Politics: &proto.CityPolitics{
			Overlord: c.Overlord,
			Lieges:   []uint64{},
		},
	}

	for _, c := range c.GetLieges() {
		cv.Politics.Lieges = append(cv.Politics.Lieges, c.ID)
	}

	cv.Evol = showEvolution(w, c)
	cv.Production = showProduction(w, c)
	cv.Stock = showStock(w, c)
	cv.Assets = showAssets(w, c)
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

func showArmy(w *region.World, a *region.Army) *proto.ArmyView {
	view := &proto.ArmyView{
		Id:       a.ID,
		Name:     a.Name,
		Location: a.Cell,
		Stock:    resAbsM2P(a.Stock),
	}
	for _, u := range a.Units {
		view.Units = append(view.Units, showUnit(w, u))
	}
	for _, c := range a.Targets {
		view.Commands = append(view.Commands, showArmyCommand(&c))
	}
	return view
}

func showUnit(w *region.World, u *region.Unit) *proto.UnitView {
	return &proto.UnitView{
		Id:     u.ID,
		IdType: u.Type,
		Name:   "",
		Ticks:  u.Ticks,
		Health: u.Health,
	}
}

func showCityPublic(w *region.World, c *region.City, scored bool) *proto.PublicCity {
	var score int64
	if scored {
		score = c.GetActualPopularity(w)
	}
	return &proto.PublicCity{
		Id:        c.ID,
		Name:      c.Name,
		Score:     score,
		Alignment: c.Alignment,
		Chaos:     c.Chaotic,
		Cult:      c.Cult,
		Politics:  c.PoliticalGroup,
		Ethny:     c.EthnicGroup,
	}
}

func showCityStats(r *region.Region, c *region.City) *proto.CityStats {
	stats := c.GetStats(r)
	return &proto.CityStats{
		// Identifiers
		Id:   c.ID,
		Name: c.Name,
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
	rc.Public = showCityPublic(r.GetWorld(), c, false)
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
