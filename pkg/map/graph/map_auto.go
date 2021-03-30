// Code generated : DO NOT EDIT.

// Copyright (c) 2018-2021 Contributors as noted in the AUTHORS file
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package mapgraph

import (
	"github.com/juju/errors"
	"math/rand"
	"sort"
)

type SetOfVertices []*Vertex

func (s SetOfVertices) CheckThenFail() {
	if err := s.Check(); err != nil {
		panic(err.Error())
	}
}

func (s SetOfVertices) Len() int {
	return len(s)
}

func (s SetOfVertices) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s *SetOfVertices) Add(a *Vertex) {
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

func (s SetOfVertices) Check() error {
	if !sort.IsSorted(s) {
		return errors.NotValidf("sorting (%v) %v", s.Len(), s)
	}
	if !s.areItemsUnique() {
		return errors.NotValidf("unicity")
	}
	return nil
}

func (s *SetOfVertices) testRandomVacuum() {
	for s.Len() > 0 {
		idx := rand.Intn(s.Len())
		s.Remove((*s)[idx])
		s.CheckThenFail()
	}
}

func (s SetOfVertices) Less(i, j int) bool {
	return s[i].ID < s[j].ID
}

func (s SetOfVertices) areItemsUnique() bool {
	var lastId uint64
	for _, a := range s {
		if lastId == a.ID {
			return false
		}
		lastId = a.ID
	}
	return true
}

func (s SetOfVertices) Slice(marker uint64, max uint32) []*Vertex {
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

func (s SetOfVertices) getIndex(id uint64) int {
	i := sort.Search(len(s), func(i int) bool {
		return s[i].ID >= id
	})
	if i < len(s) && s[i].ID == id {
		return i
	}
	return -1
}

func (s SetOfVertices) Get(id uint64) *Vertex {
	var out *Vertex
	idx := s.getIndex(id)
	if idx >= 0 {
		out = s[idx]
	}
	return out
}

func (s SetOfVertices) Has(id uint64) bool {
	return s.getIndex(id) >= 0
}

func (s *SetOfVertices) Remove(a *Vertex) {
	s.RemovePK(a.ID)
}

func (s *SetOfVertices) RemovePK(pk uint64) {
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

type SetOfEdges []*Edge

func (s SetOfEdges) CheckThenFail() {
	if err := s.Check(); err != nil {
		panic(err.Error())
	}
}

func (s SetOfEdges) Len() int {
	return len(s)
}

func (s SetOfEdges) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s *SetOfEdges) Add(a *Edge) {
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

func (s SetOfEdges) Check() error {
	if !sort.IsSorted(s) {
		return errors.NotValidf("sorting (%v) %v", s.Len(), s)
	}
	if !s.areItemsUnique() {
		return errors.NotValidf("unicity")
	}
	return nil
}

func (s *SetOfEdges) testRandomVacuum() {
	for s.Len() > 0 {
		idx := rand.Intn(s.Len())
		s.Remove((*s)[idx])
		s.CheckThenFail()
	}
}

func (s SetOfEdges) Less(i, j int) bool {
	p0, p1 := s[i], s[j]
	return p0.S < p1.S || (p0.S == p1.S && p0.D < p1.D)
}

func (s SetOfEdges) First(at uint64) int {
	return sort.Search(len(s), func(i int) bool { return s[i].S >= at })
}

func (s SetOfEdges) areItemsUnique() bool {
	var l0 uint64
	var l1 uint64
	for _, a := range s {
		if l0 == a.S && l1 == a.D {
			return false
		}
		l0 = a.S
	}
	return true
}

func (s SetOfEdges) Slice(m0 uint64, m1 uint64, max uint32) []*Edge {
	if max == 0 {
		max = 1000
	} else if max > 100000 {
		max = 100000
	}

	iMax := s.Len()
	start := s.First(m0)
	for start < iMax && s[start].S == m0 && s[start].D <= m1 {
		start++
	}

	remaining := uint32(iMax - start)
	if remaining > max {
		remaining = max
	}
	return s[start : uint32(start)+remaining]
}

func (s SetOfEdges) getIndex(f0 uint64, f1 uint64) int {
	i := sort.Search(len(s), func(i int) bool {
		return s[i].S >= f0 || (s[i].S == f0 && s[i].D >= f1)
	})
	if i < len(s) && s[i].S == f0 && s[i].D == f1 {
		return i
	}
	return -1
}

func (s SetOfEdges) Get(f0 uint64, f1 uint64) *Edge {
	var out *Edge
	idx := s.getIndex(f0, f1)
	if idx >= 0 {
		out = s[idx]
	}
	return out
}

func (s SetOfEdges) Has(f0 uint64, f1 uint64) bool {
	return s.getIndex(f0, f1) >= 0
}

func (s *SetOfEdges) Remove(a *Edge) {
	s.RemovePK(a.S, a.D)
}

func (s *SetOfEdges) RemovePK(f0 uint64, f1 uint64) {
	idx := s.getIndex(f0, f1)
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

type SetOfMaps []*Map

func (s SetOfMaps) CheckThenFail() {
	if err := s.Check(); err != nil {
		panic(err.Error())
	}
}

func (s SetOfMaps) Len() int {
	return len(s)
}

func (s SetOfMaps) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s *SetOfMaps) Add(a *Map) {
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

func (s SetOfMaps) Check() error {
	if !sort.IsSorted(s) {
		return errors.NotValidf("sorting (%v) %v", s.Len(), s)
	}
	if !s.areItemsUnique() {
		return errors.NotValidf("unicity")
	}
	return nil
}

func (s *SetOfMaps) testRandomVacuum() {
	for s.Len() > 0 {
		idx := rand.Intn(s.Len())
		s.Remove((*s)[idx])
		s.CheckThenFail()
	}
}

func (s SetOfMaps) Less(i, j int) bool {
	return s[i].ID < s[j].ID
}

func (s SetOfMaps) areItemsUnique() bool {
	var lastId string
	for _, a := range s {
		if lastId == a.ID {
			return false
		}
		lastId = a.ID
	}
	return true
}

func (s SetOfMaps) Slice(marker string, max uint32) []*Map {
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

func (s SetOfMaps) getIndex(id string) int {
	i := sort.Search(len(s), func(i int) bool {
		return s[i].ID >= id
	})
	if i < len(s) && s[i].ID == id {
		return i
	}
	return -1
}

func (s SetOfMaps) Get(id string) *Map {
	var out *Map
	idx := s.getIndex(id)
	if idx >= 0 {
		out = s[idx]
	}
	return out
}

func (s SetOfMaps) Has(id string) bool {
	return s.getIndex(id) >= 0
}

func (s *SetOfMaps) Remove(a *Map) {
	s.RemovePK(a.ID)
}

func (s *SetOfMaps) RemovePK(pk string) {
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
