// Copyright (c) 2018-2021 Contributors as noted in the AUTHORS file
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package region

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/jfsmig/hegemonie/pkg/utils"
	"github.com/juju/errors"
	"math/rand"
	"sort"
	"strings"
)

func (a *Army) PopCommand() {
	a.Targets = a.Targets[1:]
}

func (a *Army) ApplyAgressivity(w *Region) {
	// FIXME(jfs): NYI
}

func (a *Army) Move(ctx context.Context, r *Region) {
	w := r.world

	if a.FightID != "" {
		return
	}

	if len(a.Targets) <= 0 {
		// The Army has no command pending. It just stays.
	} else {
		cmd := a.Targets[0]
		src := a.Cell
		dst := cmd.Cell

		pLocalCity := r.CityGet(a.Cell)

		nxt, err := w.mapView.Step(ctx, r.MapName, src, dst)
		if err != nil || nxt == 0 {
			if err != nil {
				utils.Logger.Warn().Err(err).Uint64("src", src).Uint64("dst", dst).Send()
			}
			w.notifier.Army(a.City).Item(a).NoRoute(src, dst).Send()
		} else {
			a.Cell = nxt
			w.notifier.Army(a.City).Item(a).Move(src, dst).Send()
			if pLocalCity != nil && a.City.ID != pLocalCity.ID {
				w.notifier.Army(pLocalCity).Item(a).Move(src, dst).Send()
			}
		}

		if nxt == dst {
			var preventPopping bool
			switch cmd.Action {
			case CmdMove:
				// Just a stop on the way
			case CmdCityAttack:
				a.JoinCityAttack(r, pLocalCity)
			case CmdCityDefend:
				if a.JoinCityDefence(r, pLocalCity) {
					preventPopping = true
				}
			case CmdCityDisband:
				a.Disband(r, pLocalCity, true)
			}
			if !preventPopping {
				a.PopCommand()
			}
		}
	}
	a.ApplyAgressivity(r)
}

func (a *Army) Deposit(w *Region, pCity *City) {
	if pCity == nil {
		panic("Impossible action: nil city")
	}

	pCity.Stock.Add(a.Stock)
	a.Stock.Zero()

	// FIXME(jfs): Popularities
	// FIXME(jfs): Notify pLocalCity
	// FIXME(jfs): Notify a.City
}

func (a *Army) Massacre(w *Region, pCity *City) {
	if pCity == nil {
		panic("Impossible action: nil city")
	}

	pCity.TicksMassacres++

	// FIXME(jfs): Popularities
	// FIXME(jfs): Notify pLocalCity
	// FIXME(jfs): Notify a.City
}

func (a *Army) Disband(w *Region, pCity *City, shouldNotify bool) {
	if pCity == nil {
		panic("Impossible action: nil city")
	}

	nb := len(a.Units)
	if nb > 0 {
		for _, u := range a.Units {
			pCity.Units = append(pCity.Units, u)
		}
		sort.Sort(pCity.Units)
		a.Units = a.Units[:0]

		if shouldNotify {
			// FIXME(jfs): Notify pCity the arrival of 'nb' units
			// FIXME(jfs): Notify a.City the transfer of 'nb' units
		}
	}
}

func (a *Army) BreakBuilding(w *Region, pCity *City) {
	if pCity == nil {
		panic("Impossible action: nil city")
	}

	idx := rand.Intn(len(pCity.Buildings))
	pCity.Buildings.Remove(pCity.Buildings[idx])

	// FIXME(jfs): Popularities
	// FIXME(jfs): Notify pLocalCity
	// FIXME(jfs): Notify a.City
}

func (a *Army) Conquer(w *World, pCity *City) {
	if pCity == nil {
		panic("Impossible action: nil city")
	}
	a.City.ConquerCity(w, pCity)
}

func (a *Army) JoinCityDefence(w *Region, pCity *City) bool {
	if pCity.Assault == nil {
		return false
	}
	if pCity.Assault.Cell != a.Cell {
		panic("inconsistency")
	}
	if pCity.Assault.Cell != pCity.ID {
		panic("inconsistency")
	}

	a.FightID = pCity.Assault.ID
	pCity.Assault.Defense.Add(a)

	return true
}

func (a *Army) JoinCityAttack(w *Region, pCity *City) {
	if pCity.Assault == nil {
		pCity.Assault = &Fight{
			ID:      uuid.New().String(),
			Defense: make(SetOfArmies, 0),
			Attack:  make(SetOfArmies, 0)}
		if def, _ := pCity.CreateArmyDefence(w); def != nil {
			def.FightID = pCity.Assault.ID
			pCity.Assault.Defense.Add(def)
		}
	}

	if pCity.Assault.Cell != a.Cell {
		panic("inconsistency")
	}
	if pCity.Assault.Cell != pCity.ID {
		panic("inconsistency")
	}

	a.FightID = pCity.Assault.ID
	pCity.Assault.Attack.Add(a)
}

// Leave the Fight as a loser
func (a *Army) Flea(w *Region) error {
	return errors.New("Flea NYI")
}

// Change the side in the Fight.
// If the Army was defending, it becomes an attacker, if it was an attacker
// it becomes a defender.
func (a *Army) Flip(w *Region) error {
	return errors.NotImplementedf("NYI")
}

func (a *Army) Cancel(w *Region) error {
	return errors.NotImplementedf("NYI")
}

func (a *Army) DeferAttack(w *Region, loc uint64, args ActionArgAssault) error {
	var sb strings.Builder
	err := json.NewEncoder(&sb).Encode(&args)
	if err != nil {
		return errors.NewBadRequest(err, "invalid action argument")
	}
	a.Targets = append(a.Targets, Command{Action: CmdCityAttack, Cell: loc, Args: sb.String()})
	return nil
}

func (a *Army) DeferDefend(w *Region, loc uint64) error {
	a.Targets = append(a.Targets, Command{Action: CmdCityDefend, Cell: loc})
	return nil
}

func (a *Army) DeferDisband(w *Region, loc uint64) error {
	a.Targets = append(a.Targets, Command{Action: CmdCityDisband, Cell: loc})
	return nil
}

func (a *Army) DeferMove(w *Region, loc uint64, args ActionArgMove) error {
	var sb strings.Builder
	err := json.NewEncoder(&sb).Encode(&args)
	if err != nil {
		return errors.NewBadRequest(err, "invalid action argument")
	}
	a.Targets = append(a.Targets, Command{Action: CmdMove, Cell: loc, Args: sb.String()})
	return nil
}

func (a *Army) DeferWait(w *Region, loc uint64) error {
	a.Targets = append(a.Targets, Command{Action: CmdWait, Cell: loc})
	return nil
}
