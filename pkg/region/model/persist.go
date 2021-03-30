// Copyright (c) 2018-2021 Contributors as noted in the AUTHORS file
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package region

import (
	"encoding/json"
	"github.com/jfsmig/hegemonie/pkg/utils"
	"github.com/juju/errors"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func (defs *DefinitionsBase) Check() error {
	if !sort.IsSorted(&defs.Knowledges) {
		return errors.NotValidf("knowledge types unsorted")
	}
	if !sort.IsSorted(&defs.Buildings) {
		return errors.NotValidf("building types unsorted")
	}
	if !sort.IsSorted(&defs.Units) {
		return errors.NotValidf("unit types unsorted")
	}

	return nil
}

func walkJSON(path string, hook func(path string, decoder *json.Decoder) error) error {
	return filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if !strings.HasSuffix(path, ".json") {
			return nil
		}
		if !info.Mode().IsRegular() {
			return nil
		}
		if info.Size() <= 2 {
			return nil
		}
		f, err := os.Open(path)
		return hook(path, json.NewDecoder(f))
	})

}

func (defs *DefinitionsBase) loadUnits(basedir string) error {
	return walkJSON(basedir, func(path string, decoder *json.Decoder) error {
		utils.Logger.Info().Str("kind", "unit").Str("path", path).Msg("loading")
		tmp := make([]*UnitType, 0)
		if err := decoder.Decode(&tmp); err != nil {
			return errors.NewNotValid(err, "invalid json")
		}
		defs.Units = append(defs.Units, tmp...)
		return nil
	})
}

func (defs *DefinitionsBase) loadKnowledge(basedir string) (err error) {
	return walkJSON(basedir, func(path string, decoder *json.Decoder) error {
		utils.Logger.Info().Str("kind", "skill").Str("path", path).Msg("loading")
		tmp := make([]*KnowledgeType, 0)
		if err = decoder.Decode(&tmp); err != nil {
			return errors.NewNotValid(err, "invalid json")
		}
		defs.Knowledges = append(defs.Knowledges, tmp...)
		return nil
	})
}

func (defs *DefinitionsBase) loadBuildings(basedir string) (err error) {
	return walkJSON(basedir, func(path string, decoder *json.Decoder) error {
		utils.Logger.Info().Str("kind", "building").Str("path", path).Msg("loading")
		tmp := make([]*BuildingType, 0)
		if err = decoder.Decode(&tmp); err != nil {
			return errors.NewNotValid(err, "invalid json")
		}
		defs.Buildings = append(defs.Buildings, tmp...)
		return nil
	})
}

func (defs *DefinitionsBase) load(path string) (err error) {
	err = defs.loadUnits(path + "/units")
	if err == nil {
		err = defs.loadKnowledge(path + "/knowledge")
	}
	if err == nil {
		err = defs.loadBuildings(path + "/buildings")
	}

	if err == nil {
		sort.Sort(&defs.Knowledges)
		sort.Sort(&defs.Buildings)
		sort.Sort(&defs.Units)
	}
	return err
}

func (w *World) Check() error {
	if w.notifier == nil {
		return errors.NotValidf("missing an event notifier")
	}
	if w.mapView == nil {
		return errors.NotValidf("missing a path resolution object")
	}
	if err := w.Definitions.Check(); err != nil {
		return errors.Annotate(err, "bad definitions")
	}
	for _, r := range w.Regions {
		if err := r.Check(); err != nil {
			return errors.Annotatef(err, "bad region '%s'", r.Name)
		}
	}
	return nil
}

func (reg *Region) Check() error {
	if !sort.IsSorted(&reg.Cities) {
		return errors.NotValidf("cities unsorted")
	}
	if !sort.IsSorted(&reg.Fights) {
		return errors.NotValidf("fights unsorted")
	}

	for _, a := range reg.Fights {
		if !sort.IsSorted(&a.Attack) {
			return errors.NotValidf("fight attack unsorted")
		}
		if !sort.IsSorted(&a.Defense) {
			return errors.NotValidf("fight defense unsorted")
		}
	}
	for _, a := range reg.Cities {
		if !sort.IsSorted(&a.Knowledges) {
			return errors.NotValidf("knowledge unsorted")
		}
		if !sort.IsSorted(&a.Buildings) {
			return errors.NotValidf("building unsorted")
		}
		if !sort.IsSorted(&a.Units) {
			return errors.NotValidf("unit sequence: unsorted")
		}
		if !sort.IsSorted(&a.lieges) {
			return errors.NotValidf("city lieges unsorted")
		}
		if !sort.IsSorted(&a.Armies) {
			return errors.NotValidf("city armies unsorted")
		}
		for _, a := range a.Armies {
			if !sort.IsSorted(&a.Units) {
				return errors.NotValidf("units unsorted")
			}
		}
	}

	return nil
}

func (reg *Region) PostLoad() error {
	// Sort all the lookup arrays
	sort.Sort(&reg.Cities)
	sort.Sort(&reg.Fights)

	for _, c := range reg.Cities {
		sort.Sort(&c.Knowledges)
		sort.Sort(&c.Buildings)
		sort.Sort(&c.Units)
		if c.Armies == nil {
			c.Armies = make(SetOfArmies, 0)
		} else {
			sort.Sort(&c.Armies)
		}
		if c.lieges == nil {
			c.lieges = make(SetOfCities, 0)
		} else {
			sort.Sort(&c.lieges)
		}

		for _, a := range c.Armies {
			// Link Armies to their City
			a.City = c
			// FIXME: Link each Army to its Fight
		}
	}

	return nil
}

func (w *World) PostLoad() error {
	sort.Sort(&w.Regions)
	for _, r := range w.Regions {
		if err := r.PostLoad(); err != nil {
			return errors.Annotate(err, "postload/link error")
		}
	}
	return nil
}

func (w *World) LoadDefinitions(basedir string) (err error) {
	err = w.Definitions.load(basedir)
	if err != nil {
		return errors.Annotatef(err, "invalid world from [%s]", basedir)
	}

	err = w.Definitions.Check()
	if err != nil {
		return errors.Annotatef(err, "inconsistent world from [%s]", basedir)
	}

	return nil
}

func (w *World) LoadRegions(basedir string) error {
	return walkJSON(basedir, func(path string, decoder *json.Decoder) error {
		reg := &Region{}
		err := decoder.Decode(&reg)
		if err != nil {
			return errors.Annotatef(err, "region decoding error [%s]", path)
		}
		w.Regions.Add(reg)
		return nil
	})
}
