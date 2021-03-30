// Copyright (c) 2018-2021 Contributors as noted in the AUTHORS file
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package region

import (
	"github.com/google/uuid"
	"github.com/juju/errors"
)

// TODO(jfs): Maybe speed the execution with a reverse index of Requires
func (s SetOfKnowledgeTypes) Frontier(owned []*Knowledge) []*KnowledgeType {
	pending := make(map[uint64]bool)
	finished := make(map[uint64]bool)
	for _, k := range owned {
		if k.Ticks == 0 {
			finished[k.Type] = true
		} else {
			pending[k.Type] = true
		}
	}

	valid := func(kt *KnowledgeType) bool {
		if finished[kt.ID] || pending[kt.ID] {
			return false
		}
		for _, c := range kt.Conflicts {
			if finished[c] || pending[c] {
				return false
			}
		}
		for _, c := range kt.Requires {
			if !finished[c] {
				return false
			}
		}
		return true
	}

	result := make([]*KnowledgeType, 0)
	for _, kt := range s {
		if valid(kt) {
			result = append(result, kt)
		}
	}
	return result
}

func CheckKnowledgeDependencies(owned SetOfKnowledgeTypes, required, forbidden []uint64) bool {
	for _, k := range forbidden {
		if owned.Has(k) {
			return false
		}
	}
	for _, k := range required {
		if !owned.Has(k) {
			return false
		}
	}
	return true
}

// Finish abruptly terminates the study of the Skill.
// The number of study ticks suddenly drop to 0, whatever its prior value.
func (k *Knowledge) Finish() *Knowledge {
	k.Ticks = 0
	return k
}

// StartStudy declares a new Skill with the given type and the tick  sgauge at its max
func (c *City) StartStudy(pType *KnowledgeType) *Knowledge {
	id := uuid.New().String()
	k := &Knowledge{ID: id, Type: pType.ID, Ticks: pType.Ticks}
	c.Knowledges.Add(k)
	return k
}

func (c *City) Study(w *Region, typeID uint64) (string, error) {
	t := w.world.KnowledgeTypeGet(typeID)
	if t == nil {
		return "", errors.NotFoundf("knowledge type not found")
	}
	for _, k := range c.Knowledges {
		if typeID == k.Type {
			return "", errors.AlreadyExistsf("already started")
		}
	}
	if !CheckKnowledgeDependencies(c.ownedKnowledgeTypes(w), t.Requires, t.Conflicts) {
		return "", errors.Forbiddenf("dependencies unmet")
	}

	k := c.StartStudy(t)
	// TODO(jfs): emit a notification
	return k.ID, nil
}

func (c *City) InstantStudy(w *Region, typeID uint64) error {
	t := w.world.KnowledgeTypeGet(typeID)
	if t == nil {
		return errors.NotFoundf("knowledge type not found")
	}
	for _, k := range c.Knowledges {
		if typeID == k.Type {
			return nil
		}
	}

	c.StartStudy(t).Ticks = 0
	return nil
}
