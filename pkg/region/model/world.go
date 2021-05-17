// Copyright (c) 2018-2021 Contributors as noted in the AUTHORS file
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package region

import (
	"github.com/juju/errors"
)

// WLock acquires an exclusive ("writer") lock on the current world
// The request for a writer's lock has the priority on any equest for a reader's lock.
// It also forbid any concurrent writer or reader lock to be acquired.
func (w *World) WLock() { w.rw.Lock() }

// WUnlock releases an exclusive ("writer") lock on the current world
func (w *World) WUnlock() { w.rw.Unlock() }

// RLock acquires a shared("reader") lock on the current world
// Several requests for a reader's lock may be granted simultaneously but they are
// incompatible with any writer's locks. Their usage is reserved for read-nonly operations.
func (w *World) RLock() { w.rw.RLock() }

// RUnlock releases a shared("reader") lock on the current world
func (w *World) RUnlock() { w.rw.RUnlock() }

// NewWorld instantiates an empty World, free of any definition
func NewWorld() (*World, error) {
	return &World{
		notifier: LogEvent(&noEvt{}),
		mapView:  &noopMapClient{},
		Regions:  make(SetOfRegions, 0),
		Definitions: DefinitionsBase{
			Units:     make(SetOfUnitTypes, 0),
			Buildings: make(SetOfBuildingTypes, 0),
			Skills:    make(SetOfSkillTypes, 0),
		},
	}, nil
}

// SetNotifier changes the event Notifier associated with the current World.
// That Notifier is used for in-game notifications, to collect a per-Character log.
func (w *World) SetNotifier(n Notifier) {
	w.notifier = LogEvent(n)
}

// SetMapClient changes the MapClient associated to the current World.
// that MapClient will further be used for paths resolution before armies movements
func (w *World) SetMapClient(c MapClient) {
	// TODO(jfs): maybe wrap with a cache if perf is too poor
	w.mapView = c
}

// NamedCity is the information that a City with named Name should exist at the position ID on the graph map.
type NamedCity struct {
	Name string
	ID   uint64
}

// CreateRegion instantiates and registers a Region into the current World.
// There is no check that the map exists!
// There is no check that the set of City exist on the map!
func (w *World) CreateRegion(name, mapName string, cities []NamedCity) (*Region, error) {
	if w.Regions.Has(name) {
		return nil, errors.AlreadyExistsf("region found with id [%s]", name)
	}
	r := &Region{
		Name:    name,
		MapName: mapName,
		Cities:  make(SetOfCities, 0),
		Fights:  make(SetOfFights, 0),
		world:   w,
	}
	for _, x := range cities {
		city, err := r.CityCreate(x.ID)
		if err != nil {
			return nil, err
		}
		city.Name = x.Name
	}
	w.Regions.Add(r)
	return r, nil
}

func (w *World) UnitTypeGet(id string) *UnitType {
	return w.Definitions.Units.Get(id)
}

func (w *World) UnitGetFrontier(owned []*Building) []*UnitType {
	return w.Definitions.Units.Frontier(owned)
}

func (w *World) BuildingTypeGet(id string) *BuildingType {
	return w.Definitions.Buildings.Get(id)
}

func (w *World) BuildingGetFrontier(pop int64, built []*Building, owned []*Skill) []*BuildingType {
	// TODO(jfs): Maybe speed the execution with a reverse index of Requires
	return w.Definitions.Buildings.Frontier(pop, built, owned)
}

func (w *World) SkillTypeGet(id string) *SkillType {
	return w.Definitions.Skills.Get(id)
}

func (w *World) SkillGetFrontier(owned []*Skill) []*SkillType {
	return w.Definitions.Skills.Frontier(owned)
}
