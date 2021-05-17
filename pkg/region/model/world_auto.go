// Code generated : DO NOT EDIT.

// Copyright (c) 2018-2021 Contributors as noted in the AUTHORS file
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package region

import (
	"github.com/juju/errors"
	"math/rand"
	"sort"
)

type SetOfArtifacts []*Artifact

func (s SetOfArtifacts) CheckThenFail() {
	if err := s.Check(); err != nil {
		panic(err.Error())
	}
}

func (s SetOfArtifacts) Len() int {
	return len(s)
}

func (s SetOfArtifacts) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s *SetOfArtifacts) Add(a *Artifact) {
	*s = append(*s, a)
	switch nb := len(*s); nb {
	case 0:
		panic("yet another attack of a solar eruption")
	case 1:
		return
	case 2:
		sort.Sort(s)
	default:
		if !sort.IsSorted((*s)[nb-2:]) {
			sort.Sort(s)
		}
	}
}

func (s SetOfArtifacts) Check() error {
	if !sort.IsSorted(s) {
		return errors.NotValidf("sorting (%v) %v", s.Len(), s)
	}
	if !s.areItemsUnique() {
		return errors.NotValidf("unicity")
	}
	return nil
}

func (s *SetOfArtifacts) testRandomVacuum() {
	for s.Len() > 0 {
		idx := rand.Intn(s.Len())
		s.Remove((*s)[idx])
		s.CheckThenFail()
	}
}

func (s SetOfArtifacts) Less(i, j int) bool {
	return s[i].ID < s[j].ID
}

func (s SetOfArtifacts) areItemsUnique() bool {
	var lastId string
	for _, a := range s {
		if lastId == a.ID {
			return false
		}
		lastId = a.ID
	}
	return true
}

func (s SetOfArtifacts) Slice(marker string, max uint32) []*Artifact {
	if max == 0 {
		max = 1000
	} else if max > 100000 {
		max = 100000
	}
	start := sort.Search(len(s), func(i int) bool {
		return s[i].ID > marker
	})
	if start < 0 || start >= s.Len() {
		return s[:0]
	}
	remaining := uint32(s.Len() - start)
	if remaining > max {
		remaining = max
	}
	return s[start : uint32(start)+remaining]
}

func (s SetOfArtifacts) getIndex(id string) int {
	i := sort.Search(len(s), func(i int) bool {
		return s[i].ID >= id
	})
	if i < len(s) && s[i].ID == id {
		return i
	}
	return -1
}

func (s SetOfArtifacts) Get(id string) *Artifact {
	var out *Artifact
	idx := s.getIndex(id)
	if idx >= 0 {
		out = s[idx]
	}
	return out
}

func (s SetOfArtifacts) Has(id string) bool {
	return s.getIndex(id) >= 0
}

func (s *SetOfArtifacts) Remove(a *Artifact) {
	s.RemovePK(a.ID)
}

func (s *SetOfArtifacts) RemovePK(pk string) {
	idx := s.getIndex(pk)
	if idx >= 0 && idx < len(*s) {
		if len(*s) == 1 {
			*s = (*s)[:0]
		} else {
			s.Swap(idx, s.Len()-1)
			*s = (*s)[:s.Len()-1]
			sort.Sort(*s)
		}
	}
}

type SetOfArmies []*Army

func (s SetOfArmies) CheckThenFail() {
	if err := s.Check(); err != nil {
		panic(err.Error())
	}
}

func (s SetOfArmies) Len() int {
	return len(s)
}

func (s SetOfArmies) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s *SetOfArmies) Add(a *Army) {
	*s = append(*s, a)
	switch nb := len(*s); nb {
	case 0:
		panic("yet another attack of a solar eruption")
	case 1:
		return
	case 2:
		sort.Sort(s)
	default:
		if !sort.IsSorted((*s)[nb-2:]) {
			sort.Sort(s)
		}
	}
}

func (s SetOfArmies) Check() error {
	if !sort.IsSorted(s) {
		return errors.NotValidf("sorting (%v) %v", s.Len(), s)
	}
	if !s.areItemsUnique() {
		return errors.NotValidf("unicity")
	}
	return nil
}

func (s *SetOfArmies) testRandomVacuum() {
	for s.Len() > 0 {
		idx := rand.Intn(s.Len())
		s.Remove((*s)[idx])
		s.CheckThenFail()
	}
}

func (s SetOfArmies) Less(i, j int) bool {
	return s[i].ID < s[j].ID
}

func (s SetOfArmies) areItemsUnique() bool {
	var lastId string
	for _, a := range s {
		if lastId == a.ID {
			return false
		}
		lastId = a.ID
	}
	return true
}

func (s SetOfArmies) Slice(marker string, max uint32) []*Army {
	if max == 0 {
		max = 1000
	} else if max > 100000 {
		max = 100000
	}
	start := sort.Search(len(s), func(i int) bool {
		return s[i].ID > marker
	})
	if start < 0 || start >= s.Len() {
		return s[:0]
	}
	remaining := uint32(s.Len() - start)
	if remaining > max {
		remaining = max
	}
	return s[start : uint32(start)+remaining]
}

func (s SetOfArmies) getIndex(id string) int {
	i := sort.Search(len(s), func(i int) bool {
		return s[i].ID >= id
	})
	if i < len(s) && s[i].ID == id {
		return i
	}
	return -1
}

func (s SetOfArmies) Get(id string) *Army {
	var out *Army
	idx := s.getIndex(id)
	if idx >= 0 {
		out = s[idx]
	}
	return out
}

func (s SetOfArmies) Has(id string) bool {
	return s.getIndex(id) >= 0
}

func (s *SetOfArmies) Remove(a *Army) {
	s.RemovePK(a.ID)
}

func (s *SetOfArmies) RemovePK(pk string) {
	idx := s.getIndex(pk)
	if idx >= 0 && idx < len(*s) {
		if len(*s) == 1 {
			*s = (*s)[:0]
		} else {
			s.Swap(idx, s.Len()-1)
			*s = (*s)[:s.Len()-1]
			sort.Sort(*s)
		}
	}
}

type SetOfBuildings []*Building

func (s SetOfBuildings) CheckThenFail() {
	if err := s.Check(); err != nil {
		panic(err.Error())
	}
}

func (s SetOfBuildings) Len() int {
	return len(s)
}

func (s SetOfBuildings) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s *SetOfBuildings) Add(a *Building) {
	*s = append(*s, a)
	switch nb := len(*s); nb {
	case 0:
		panic("yet another attack of a solar eruption")
	case 1:
		return
	case 2:
		sort.Sort(s)
	default:
		if !sort.IsSorted((*s)[nb-2:]) {
			sort.Sort(s)
		}
	}
}

func (s SetOfBuildings) Check() error {
	if !sort.IsSorted(s) {
		return errors.NotValidf("sorting (%v) %v", s.Len(), s)
	}
	if !s.areItemsUnique() {
		return errors.NotValidf("unicity")
	}
	return nil
}

func (s *SetOfBuildings) testRandomVacuum() {
	for s.Len() > 0 {
		idx := rand.Intn(s.Len())
		s.Remove((*s)[idx])
		s.CheckThenFail()
	}
}

func (s SetOfBuildings) Less(i, j int) bool {
	return s[i].ID < s[j].ID
}

func (s SetOfBuildings) areItemsUnique() bool {
	var lastId string
	for _, a := range s {
		if lastId == a.ID {
			return false
		}
		lastId = a.ID
	}
	return true
}

func (s SetOfBuildings) Slice(marker string, max uint32) []*Building {
	if max == 0 {
		max = 1000
	} else if max > 100000 {
		max = 100000
	}
	start := sort.Search(len(s), func(i int) bool {
		return s[i].ID > marker
	})
	if start < 0 || start >= s.Len() {
		return s[:0]
	}
	remaining := uint32(s.Len() - start)
	if remaining > max {
		remaining = max
	}
	return s[start : uint32(start)+remaining]
}

func (s SetOfBuildings) getIndex(id string) int {
	i := sort.Search(len(s), func(i int) bool {
		return s[i].ID >= id
	})
	if i < len(s) && s[i].ID == id {
		return i
	}
	return -1
}

func (s SetOfBuildings) Get(id string) *Building {
	var out *Building
	idx := s.getIndex(id)
	if idx >= 0 {
		out = s[idx]
	}
	return out
}

func (s SetOfBuildings) Has(id string) bool {
	return s.getIndex(id) >= 0
}

func (s *SetOfBuildings) Remove(a *Building) {
	s.RemovePK(a.ID)
}

func (s *SetOfBuildings) RemovePK(pk string) {
	idx := s.getIndex(pk)
	if idx >= 0 && idx < len(*s) {
		if len(*s) == 1 {
			*s = (*s)[:0]
		} else {
			s.Swap(idx, s.Len()-1)
			*s = (*s)[:s.Len()-1]
			sort.Sort(*s)
		}
	}
}

type SetOfBuildingTypes []*BuildingType

func (s SetOfBuildingTypes) CheckThenFail() {
	if err := s.Check(); err != nil {
		panic(err.Error())
	}
}

func (s SetOfBuildingTypes) Len() int {
	return len(s)
}

func (s SetOfBuildingTypes) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s *SetOfBuildingTypes) Add(a *BuildingType) {
	*s = append(*s, a)
	switch nb := len(*s); nb {
	case 0:
		panic("yet another attack of a solar eruption")
	case 1:
		return
	case 2:
		sort.Sort(s)
	default:
		if !sort.IsSorted((*s)[nb-2:]) {
			sort.Sort(s)
		}
	}
}

func (s SetOfBuildingTypes) Check() error {
	if !sort.IsSorted(s) {
		return errors.NotValidf("sorting (%v) %v", s.Len(), s)
	}
	if !s.areItemsUnique() {
		return errors.NotValidf("unicity")
	}
	return nil
}

func (s *SetOfBuildingTypes) testRandomVacuum() {
	for s.Len() > 0 {
		idx := rand.Intn(s.Len())
		s.Remove((*s)[idx])
		s.CheckThenFail()
	}
}

func (s SetOfBuildingTypes) Less(i, j int) bool {
	return s[i].ID < s[j].ID
}

func (s SetOfBuildingTypes) areItemsUnique() bool {
	var lastId string
	for _, a := range s {
		if lastId == a.ID {
			return false
		}
		lastId = a.ID
	}
	return true
}

func (s SetOfBuildingTypes) Slice(marker string, max uint32) []*BuildingType {
	if max == 0 {
		max = 1000
	} else if max > 100000 {
		max = 100000
	}
	start := sort.Search(len(s), func(i int) bool {
		return s[i].ID > marker
	})
	if start < 0 || start >= s.Len() {
		return s[:0]
	}
	remaining := uint32(s.Len() - start)
	if remaining > max {
		remaining = max
	}
	return s[start : uint32(start)+remaining]
}

func (s SetOfBuildingTypes) getIndex(id string) int {
	i := sort.Search(len(s), func(i int) bool {
		return s[i].ID >= id
	})
	if i < len(s) && s[i].ID == id {
		return i
	}
	return -1
}

func (s SetOfBuildingTypes) Get(id string) *BuildingType {
	var out *BuildingType
	idx := s.getIndex(id)
	if idx >= 0 {
		out = s[idx]
	}
	return out
}

func (s SetOfBuildingTypes) Has(id string) bool {
	return s.getIndex(id) >= 0
}

func (s *SetOfBuildingTypes) Remove(a *BuildingType) {
	s.RemovePK(a.ID)
}

func (s *SetOfBuildingTypes) RemovePK(pk string) {
	idx := s.getIndex(pk)
	if idx >= 0 && idx < len(*s) {
		if len(*s) == 1 {
			*s = (*s)[:0]
		} else {
			s.Swap(idx, s.Len()-1)
			*s = (*s)[:s.Len()-1]
			sort.Sort(*s)
		}
	}
}

type SetOfCities []*City

func (s SetOfCities) CheckThenFail() {
	if err := s.Check(); err != nil {
		panic(err.Error())
	}
}

func (s SetOfCities) Len() int {
	return len(s)
}

func (s SetOfCities) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s *SetOfCities) Add(a *City) {
	*s = append(*s, a)
	switch nb := len(*s); nb {
	case 0:
		panic("yet another attack of a solar eruption")
	case 1:
		return
	case 2:
		sort.Sort(s)
	default:
		if !sort.IsSorted((*s)[nb-2:]) {
			sort.Sort(s)
		}
	}
}

func (s SetOfCities) Check() error {
	if !sort.IsSorted(s) {
		return errors.NotValidf("sorting (%v) %v", s.Len(), s)
	}
	if !s.areItemsUnique() {
		return errors.NotValidf("unicity")
	}
	return nil
}

func (s *SetOfCities) testRandomVacuum() {
	for s.Len() > 0 {
		idx := rand.Intn(s.Len())
		s.Remove((*s)[idx])
		s.CheckThenFail()
	}
}

func (s SetOfCities) Less(i, j int) bool {
	return s[i].ID < s[j].ID
}

func (s SetOfCities) areItemsUnique() bool {
	var lastId uint64
	for _, a := range s {
		if lastId == a.ID {
			return false
		}
		lastId = a.ID
	}
	return true
}

func (s SetOfCities) Slice(marker uint64, max uint32) []*City {
	if max == 0 {
		max = 1000
	} else if max > 100000 {
		max = 100000
	}
	start := sort.Search(len(s), func(i int) bool {
		return s[i].ID > marker
	})
	if start < 0 || start >= s.Len() {
		return s[:0]
	}
	remaining := uint32(s.Len() - start)
	if remaining > max {
		remaining = max
	}
	return s[start : uint32(start)+remaining]
}

func (s SetOfCities) getIndex(id uint64) int {
	i := sort.Search(len(s), func(i int) bool {
		return s[i].ID >= id
	})
	if i < len(s) && s[i].ID == id {
		return i
	}
	return -1
}

func (s SetOfCities) Get(id uint64) *City {
	var out *City
	idx := s.getIndex(id)
	if idx >= 0 {
		out = s[idx]
	}
	return out
}

func (s SetOfCities) Has(id uint64) bool {
	return s.getIndex(id) >= 0
}

func (s *SetOfCities) Remove(a *City) {
	s.RemovePK(a.ID)
}

func (s *SetOfCities) RemovePK(pk uint64) {
	idx := s.getIndex(pk)
	if idx >= 0 && idx < len(*s) {
		if len(*s) == 1 {
			*s = (*s)[:0]
		} else {
			s.Swap(idx, s.Len()-1)
			*s = (*s)[:s.Len()-1]
			sort.Sort(*s)
		}
	}
}

type SetOfTemplates []*City

func (s SetOfTemplates) CheckThenFail() {
	if err := s.Check(); err != nil {
		panic(err.Error())
	}
}

func (s SetOfTemplates) Len() int {
	return len(s)
}

func (s SetOfTemplates) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s *SetOfTemplates) Add(a *City) {
	*s = append(*s, a)
	switch nb := len(*s); nb {
	case 0:
		panic("yet another attack of a solar eruption")
	case 1:
		return
	case 2:
		sort.Sort(s)
	default:
		if !sort.IsSorted((*s)[nb-2:]) {
			sort.Sort(s)
		}
	}
}

func (s SetOfTemplates) Check() error {
	if !sort.IsSorted(s) {
		return errors.NotValidf("sorting (%v) %v", s.Len(), s)
	}
	if !s.areItemsUnique() {
		return errors.NotValidf("unicity")
	}
	return nil
}

func (s *SetOfTemplates) testRandomVacuum() {
	for s.Len() > 0 {
		idx := rand.Intn(s.Len())
		s.Remove((*s)[idx])
		s.CheckThenFail()
	}
}

func (s SetOfTemplates) Less(i, j int) bool {
	return s[i].Name < s[j].Name
}

func (s SetOfTemplates) areItemsUnique() bool {
	var lastId string
	for _, a := range s {
		if lastId == a.Name {
			return false
		}
		lastId = a.Name
	}
	return true
}

func (s SetOfTemplates) Slice(marker string, max uint32) []*City {
	if max == 0 {
		max = 1000
	} else if max > 100000 {
		max = 100000
	}
	start := sort.Search(len(s), func(i int) bool {
		return s[i].Name > marker
	})
	if start < 0 || start >= s.Len() {
		return s[:0]
	}
	remaining := uint32(s.Len() - start)
	if remaining > max {
		remaining = max
	}
	return s[start : uint32(start)+remaining]
}

func (s SetOfTemplates) getIndex(id string) int {
	i := sort.Search(len(s), func(i int) bool {
		return s[i].Name >= id
	})
	if i < len(s) && s[i].Name == id {
		return i
	}
	return -1
}

func (s SetOfTemplates) Get(id string) *City {
	var out *City
	idx := s.getIndex(id)
	if idx >= 0 {
		out = s[idx]
	}
	return out
}

func (s SetOfTemplates) Has(id string) bool {
	return s.getIndex(id) >= 0
}

func (s *SetOfTemplates) Remove(a *City) {
	s.RemovePK(a.Name)
}

func (s *SetOfTemplates) RemovePK(pk string) {
	idx := s.getIndex(pk)
	if idx >= 0 && idx < len(*s) {
		if len(*s) == 1 {
			*s = (*s)[:0]
		} else {
			s.Swap(idx, s.Len()-1)
			*s = (*s)[:s.Len()-1]
			sort.Sort(*s)
		}
	}
}

type SetOfId []uint64

func (s SetOfId) CheckThenFail() {
	if err := s.Check(); err != nil {
		panic(err.Error())
	}
}

func (s SetOfId) Len() int {
	return len(s)
}

func (s SetOfId) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s *SetOfId) Add(a uint64) {
	*s = append(*s, a)
	switch nb := len(*s); nb {
	case 0:
		panic("yet another attack of a solar eruption")
	case 1:
		return
	case 2:
		sort.Sort(s)
	default:
		if !sort.IsSorted((*s)[nb-2:]) {
			sort.Sort(s)
		}
	}
}

func (s SetOfId) Check() error {
	if !sort.IsSorted(s) {
		return errors.NotValidf("sorting (%v) %v", s.Len(), s)
	}
	if !s.areItemsUnique() {
		return errors.NotValidf("unicity")
	}
	return nil
}

func (s *SetOfId) testRandomVacuum() {
	for s.Len() > 0 {
		idx := rand.Intn(s.Len())
		s.Remove((*s)[idx])
		s.CheckThenFail()
	}
}

func (s SetOfId) Less(i, j int) bool {
	return s[i] < s[j]
}

func (s SetOfId) areItemsUnique() bool {
	var lastId uint64
	for _, a := range s {
		if lastId == a {
			return false
		}
		lastId = a
	}
	return true
}

func (s SetOfId) Slice(marker uint64, max uint32) []uint64 {
	if max == 0 {
		max = 1000
	} else if max > 100000 {
		max = 100000
	}
	start := sort.Search(len(s), func(i int) bool {
		return s[i] > marker
	})
	if start < 0 || start >= s.Len() {
		return s[:0]
	}
	remaining := uint32(s.Len() - start)
	if remaining > max {
		remaining = max
	}
	return s[start : uint32(start)+remaining]
}

func (s SetOfId) getIndex(id uint64) int {
	i := sort.Search(len(s), func(i int) bool {
		return s[i] >= id
	})
	if i < len(s) && s[i] == id {
		return i
	}
	return -1
}

func (s SetOfId) Get(id uint64) uint64 {
	var out uint64
	idx := s.getIndex(id)
	if idx >= 0 {
		out = s[idx]
	}
	return out
}

func (s SetOfId) Has(id uint64) bool {
	return s.getIndex(id) >= 0
}

func (s *SetOfId) Remove(a uint64) {
	s.RemovePK(a)
}

func (s *SetOfId) RemovePK(pk uint64) {
	idx := s.getIndex(pk)
	if idx >= 0 && idx < len(*s) {
		if len(*s) == 1 {
			*s = (*s)[:0]
		} else {
			s.Swap(idx, s.Len()-1)
			*s = (*s)[:s.Len()-1]
			sort.Sort(*s)
		}
	}
}

type SetOfKnowledges []*Knowledge

func (s SetOfKnowledges) CheckThenFail() {
	if err := s.Check(); err != nil {
		panic(err.Error())
	}
}

func (s SetOfKnowledges) Len() int {
	return len(s)
}

func (s SetOfKnowledges) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s *SetOfKnowledges) Add(a *Knowledge) {
	*s = append(*s, a)
	switch nb := len(*s); nb {
	case 0:
		panic("yet another attack of a solar eruption")
	case 1:
		return
	case 2:
		sort.Sort(s)
	default:
		if !sort.IsSorted((*s)[nb-2:]) {
			sort.Sort(s)
		}
	}
}

func (s SetOfKnowledges) Check() error {
	if !sort.IsSorted(s) {
		return errors.NotValidf("sorting (%v) %v", s.Len(), s)
	}
	if !s.areItemsUnique() {
		return errors.NotValidf("unicity")
	}
	return nil
}

func (s *SetOfKnowledges) testRandomVacuum() {
	for s.Len() > 0 {
		idx := rand.Intn(s.Len())
		s.Remove((*s)[idx])
		s.CheckThenFail()
	}
}

func (s SetOfKnowledges) Less(i, j int) bool {
	return s[i].ID < s[j].ID
}

func (s SetOfKnowledges) areItemsUnique() bool {
	var lastId string
	for _, a := range s {
		if lastId == a.ID {
			return false
		}
		lastId = a.ID
	}
	return true
}

func (s SetOfKnowledges) Slice(marker string, max uint32) []*Knowledge {
	if max == 0 {
		max = 1000
	} else if max > 100000 {
		max = 100000
	}
	start := sort.Search(len(s), func(i int) bool {
		return s[i].ID > marker
	})
	if start < 0 || start >= s.Len() {
		return s[:0]
	}
	remaining := uint32(s.Len() - start)
	if remaining > max {
		remaining = max
	}
	return s[start : uint32(start)+remaining]
}

func (s SetOfKnowledges) getIndex(id string) int {
	i := sort.Search(len(s), func(i int) bool {
		return s[i].ID >= id
	})
	if i < len(s) && s[i].ID == id {
		return i
	}
	return -1
}

func (s SetOfKnowledges) Get(id string) *Knowledge {
	var out *Knowledge
	idx := s.getIndex(id)
	if idx >= 0 {
		out = s[idx]
	}
	return out
}

func (s SetOfKnowledges) Has(id string) bool {
	return s.getIndex(id) >= 0
}

func (s *SetOfKnowledges) Remove(a *Knowledge) {
	s.RemovePK(a.ID)
}

func (s *SetOfKnowledges) RemovePK(pk string) {
	idx := s.getIndex(pk)
	if idx >= 0 && idx < len(*s) {
		if len(*s) == 1 {
			*s = (*s)[:0]
		} else {
			s.Swap(idx, s.Len()-1)
			*s = (*s)[:s.Len()-1]
			sort.Sort(*s)
		}
	}
}

type SetOfKnowledgeTypes []*KnowledgeType

func (s SetOfKnowledgeTypes) CheckThenFail() {
	if err := s.Check(); err != nil {
		panic(err.Error())
	}
}

func (s SetOfKnowledgeTypes) Len() int {
	return len(s)
}

func (s SetOfKnowledgeTypes) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s *SetOfKnowledgeTypes) Add(a *KnowledgeType) {
	*s = append(*s, a)
	switch nb := len(*s); nb {
	case 0:
		panic("yet another attack of a solar eruption")
	case 1:
		return
	case 2:
		sort.Sort(s)
	default:
		if !sort.IsSorted((*s)[nb-2:]) {
			sort.Sort(s)
		}
	}
}

func (s SetOfKnowledgeTypes) Check() error {
	if !sort.IsSorted(s) {
		return errors.NotValidf("sorting (%v) %v", s.Len(), s)
	}
	if !s.areItemsUnique() {
		return errors.NotValidf("unicity")
	}
	return nil
}

func (s *SetOfKnowledgeTypes) testRandomVacuum() {
	for s.Len() > 0 {
		idx := rand.Intn(s.Len())
		s.Remove((*s)[idx])
		s.CheckThenFail()
	}
}

func (s SetOfKnowledgeTypes) Less(i, j int) bool {
	return s[i].ID < s[j].ID
}

func (s SetOfKnowledgeTypes) areItemsUnique() bool {
	var lastId string
	for _, a := range s {
		if lastId == a.ID {
			return false
		}
		lastId = a.ID
	}
	return true
}

func (s SetOfKnowledgeTypes) Slice(marker string, max uint32) []*KnowledgeType {
	if max == 0 {
		max = 1000
	} else if max > 100000 {
		max = 100000
	}
	start := sort.Search(len(s), func(i int) bool {
		return s[i].ID > marker
	})
	if start < 0 || start >= s.Len() {
		return s[:0]
	}
	remaining := uint32(s.Len() - start)
	if remaining > max {
		remaining = max
	}
	return s[start : uint32(start)+remaining]
}

func (s SetOfKnowledgeTypes) getIndex(id string) int {
	i := sort.Search(len(s), func(i int) bool {
		return s[i].ID >= id
	})
	if i < len(s) && s[i].ID == id {
		return i
	}
	return -1
}

func (s SetOfKnowledgeTypes) Get(id string) *KnowledgeType {
	var out *KnowledgeType
	idx := s.getIndex(id)
	if idx >= 0 {
		out = s[idx]
	}
	return out
}

func (s SetOfKnowledgeTypes) Has(id string) bool {
	return s.getIndex(id) >= 0
}

func (s *SetOfKnowledgeTypes) Remove(a *KnowledgeType) {
	s.RemovePK(a.ID)
}

func (s *SetOfKnowledgeTypes) RemovePK(pk string) {
	idx := s.getIndex(pk)
	if idx >= 0 && idx < len(*s) {
		if len(*s) == 1 {
			*s = (*s)[:0]
		} else {
			s.Swap(idx, s.Len()-1)
			*s = (*s)[:s.Len()-1]
			sort.Sort(*s)
		}
	}
}

type SetOfUnits []*Unit

func (s SetOfUnits) CheckThenFail() {
	if err := s.Check(); err != nil {
		panic(err.Error())
	}
}

func (s SetOfUnits) Len() int {
	return len(s)
}

func (s SetOfUnits) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s *SetOfUnits) Add(a *Unit) {
	*s = append(*s, a)
	switch nb := len(*s); nb {
	case 0:
		panic("yet another attack of a solar eruption")
	case 1:
		return
	case 2:
		sort.Sort(s)
	default:
		if !sort.IsSorted((*s)[nb-2:]) {
			sort.Sort(s)
		}
	}
}

func (s SetOfUnits) Check() error {
	if !sort.IsSorted(s) {
		return errors.NotValidf("sorting (%v) %v", s.Len(), s)
	}
	if !s.areItemsUnique() {
		return errors.NotValidf("unicity")
	}
	return nil
}

func (s *SetOfUnits) testRandomVacuum() {
	for s.Len() > 0 {
		idx := rand.Intn(s.Len())
		s.Remove((*s)[idx])
		s.CheckThenFail()
	}
}

func (s SetOfUnits) Less(i, j int) bool {
	return s[i].ID < s[j].ID
}

func (s SetOfUnits) areItemsUnique() bool {
	var lastId string
	for _, a := range s {
		if lastId == a.ID {
			return false
		}
		lastId = a.ID
	}
	return true
}

func (s SetOfUnits) Slice(marker string, max uint32) []*Unit {
	if max == 0 {
		max = 1000
	} else if max > 100000 {
		max = 100000
	}
	start := sort.Search(len(s), func(i int) bool {
		return s[i].ID > marker
	})
	if start < 0 || start >= s.Len() {
		return s[:0]
	}
	remaining := uint32(s.Len() - start)
	if remaining > max {
		remaining = max
	}
	return s[start : uint32(start)+remaining]
}

func (s SetOfUnits) getIndex(id string) int {
	i := sort.Search(len(s), func(i int) bool {
		return s[i].ID >= id
	})
	if i < len(s) && s[i].ID == id {
		return i
	}
	return -1
}

func (s SetOfUnits) Get(id string) *Unit {
	var out *Unit
	idx := s.getIndex(id)
	if idx >= 0 {
		out = s[idx]
	}
	return out
}

func (s SetOfUnits) Has(id string) bool {
	return s.getIndex(id) >= 0
}

func (s *SetOfUnits) Remove(a *Unit) {
	s.RemovePK(a.ID)
}

func (s *SetOfUnits) RemovePK(pk string) {
	idx := s.getIndex(pk)
	if idx >= 0 && idx < len(*s) {
		if len(*s) == 1 {
			*s = (*s)[:0]
		} else {
			s.Swap(idx, s.Len()-1)
			*s = (*s)[:s.Len()-1]
			sort.Sort(*s)
		}
	}
}

type SetOfUnitTypes []*UnitType

func (s SetOfUnitTypes) CheckThenFail() {
	if err := s.Check(); err != nil {
		panic(err.Error())
	}
}

func (s SetOfUnitTypes) Len() int {
	return len(s)
}

func (s SetOfUnitTypes) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s *SetOfUnitTypes) Add(a *UnitType) {
	*s = append(*s, a)
	switch nb := len(*s); nb {
	case 0:
		panic("yet another attack of a solar eruption")
	case 1:
		return
	case 2:
		sort.Sort(s)
	default:
		if !sort.IsSorted((*s)[nb-2:]) {
			sort.Sort(s)
		}
	}
}

func (s SetOfUnitTypes) Check() error {
	if !sort.IsSorted(s) {
		return errors.NotValidf("sorting (%v) %v", s.Len(), s)
	}
	if !s.areItemsUnique() {
		return errors.NotValidf("unicity")
	}
	return nil
}

func (s *SetOfUnitTypes) testRandomVacuum() {
	for s.Len() > 0 {
		idx := rand.Intn(s.Len())
		s.Remove((*s)[idx])
		s.CheckThenFail()
	}
}

func (s SetOfUnitTypes) Less(i, j int) bool {
	return s[i].ID < s[j].ID
}

func (s SetOfUnitTypes) areItemsUnique() bool {
	var lastId string
	for _, a := range s {
		if lastId == a.ID {
			return false
		}
		lastId = a.ID
	}
	return true
}

func (s SetOfUnitTypes) Slice(marker string, max uint32) []*UnitType {
	if max == 0 {
		max = 1000
	} else if max > 100000 {
		max = 100000
	}
	start := sort.Search(len(s), func(i int) bool {
		return s[i].ID > marker
	})
	if start < 0 || start >= s.Len() {
		return s[:0]
	}
	remaining := uint32(s.Len() - start)
	if remaining > max {
		remaining = max
	}
	return s[start : uint32(start)+remaining]
}

func (s SetOfUnitTypes) getIndex(id string) int {
	i := sort.Search(len(s), func(i int) bool {
		return s[i].ID >= id
	})
	if i < len(s) && s[i].ID == id {
		return i
	}
	return -1
}

func (s SetOfUnitTypes) Get(id string) *UnitType {
	var out *UnitType
	idx := s.getIndex(id)
	if idx >= 0 {
		out = s[idx]
	}
	return out
}

func (s SetOfUnitTypes) Has(id string) bool {
	return s.getIndex(id) >= 0
}

func (s *SetOfUnitTypes) Remove(a *UnitType) {
	s.RemovePK(a.ID)
}

func (s *SetOfUnitTypes) RemovePK(pk string) {
	idx := s.getIndex(pk)
	if idx >= 0 && idx < len(*s) {
		if len(*s) == 1 {
			*s = (*s)[:0]
		} else {
			s.Swap(idx, s.Len()-1)
			*s = (*s)[:s.Len()-1]
			sort.Sort(*s)
		}
	}
}

type SetOfRegions []*Region

func (s SetOfRegions) CheckThenFail() {
	if err := s.Check(); err != nil {
		panic(err.Error())
	}
}

func (s SetOfRegions) Len() int {
	return len(s)
}

func (s SetOfRegions) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s *SetOfRegions) Add(a *Region) {
	*s = append(*s, a)
	switch nb := len(*s); nb {
	case 0:
		panic("yet another attack of a solar eruption")
	case 1:
		return
	case 2:
		sort.Sort(s)
	default:
		if !sort.IsSorted((*s)[nb-2:]) {
			sort.Sort(s)
		}
	}
}

func (s SetOfRegions) Check() error {
	if !sort.IsSorted(s) {
		return errors.NotValidf("sorting (%v) %v", s.Len(), s)
	}
	if !s.areItemsUnique() {
		return errors.NotValidf("unicity")
	}
	return nil
}

func (s *SetOfRegions) testRandomVacuum() {
	for s.Len() > 0 {
		idx := rand.Intn(s.Len())
		s.Remove((*s)[idx])
		s.CheckThenFail()
	}
}

func (s SetOfRegions) Less(i, j int) bool {
	return s[i].Name < s[j].Name
}

func (s SetOfRegions) areItemsUnique() bool {
	var lastId string
	for _, a := range s {
		if lastId == a.Name {
			return false
		}
		lastId = a.Name
	}
	return true
}

func (s SetOfRegions) Slice(marker string, max uint32) []*Region {
	if max == 0 {
		max = 1000
	} else if max > 100000 {
		max = 100000
	}
	start := sort.Search(len(s), func(i int) bool {
		return s[i].Name > marker
	})
	if start < 0 || start >= s.Len() {
		return s[:0]
	}
	remaining := uint32(s.Len() - start)
	if remaining > max {
		remaining = max
	}
	return s[start : uint32(start)+remaining]
}

func (s SetOfRegions) getIndex(id string) int {
	i := sort.Search(len(s), func(i int) bool {
		return s[i].Name >= id
	})
	if i < len(s) && s[i].Name == id {
		return i
	}
	return -1
}

func (s SetOfRegions) Get(id string) *Region {
	var out *Region
	idx := s.getIndex(id)
	if idx >= 0 {
		out = s[idx]
	}
	return out
}

func (s SetOfRegions) Has(id string) bool {
	return s.getIndex(id) >= 0
}

func (s *SetOfRegions) Remove(a *Region) {
	s.RemovePK(a.Name)
}

func (s *SetOfRegions) RemovePK(pk string) {
	idx := s.getIndex(pk)
	if idx >= 0 && idx < len(*s) {
		if len(*s) == 1 {
			*s = (*s)[:0]
		} else {
			s.Swap(idx, s.Len()-1)
			*s = (*s)[:s.Len()-1]
			sort.Sort(*s)
		}
	}
}
