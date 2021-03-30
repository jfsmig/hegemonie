// Code generated : DO NOT EDIT.

// Copyright (c) 2018-2021 Contributors as noted in the AUTHORS file
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"github.com/juju/errors"
	"math/rand"
	"sort"
)

type setOfUint64 []uint64

func (s setOfUint64) CheckThenFail() {
	if err := s.Check(); err != nil {
		panic(err.Error())
	}
}

func (s setOfUint64) Len() int {
	return len(s)
}

func (s setOfUint64) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s *setOfUint64) Add(a uint64) {
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

func (s setOfUint64) Check() error {
	if !sort.IsSorted(s) {
		return errors.NotValidf("sorting (%v) %v", s.Len(), s)
	}
	if !s.areItemsUnique() {
		return errors.NotValidf("unicity")
	}
	return nil
}

func (s *setOfUint64) testRandomVacuum() {
	for s.Len() > 0 {
		idx := rand.Intn(s.Len())
		s.Remove((*s)[idx])
		s.CheckThenFail()
	}
}

func (s setOfUint64) Less(i, j int) bool {
	return s[i] < s[j]
}

func (s setOfUint64) areItemsUnique() bool {
	var lastId uint64
	for _, a := range s {
		if lastId == a {
			return false
		}
		lastId = a
	}
	return true
}

func (s setOfUint64) Slice(marker uint64, max uint32) []uint64 {
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

func (s setOfUint64) getIndex(id uint64) int {
	i := sort.Search(len(s), func(i int) bool {
		return s[i] >= id
	})
	if i < len(s) && s[i] == id {
		return i
	}
	return -1
}

func (s setOfUint64) Get(id uint64) uint64 {
	var out uint64
	idx := s.getIndex(id)
	if idx >= 0 {
		out = s[idx]
	}
	return out
}

func (s setOfUint64) Has(id uint64) bool {
	return s.getIndex(id) >= 0
}

func (s *setOfUint64) Remove(a uint64) {
	s.RemovePK(a)
}

func (s *setOfUint64) RemovePK(pk uint64) {
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

type setOfString []string

func (s setOfString) CheckThenFail() {
	if err := s.Check(); err != nil {
		panic(err.Error())
	}
}

func (s setOfString) Len() int {
	return len(s)
}

func (s setOfString) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s *setOfString) Add(a string) {
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

func (s setOfString) Check() error {
	if !sort.IsSorted(s) {
		return errors.NotValidf("sorting (%v) %v", s.Len(), s)
	}
	if !s.areItemsUnique() {
		return errors.NotValidf("unicity")
	}
	return nil
}

func (s *setOfString) testRandomVacuum() {
	for s.Len() > 0 {
		idx := rand.Intn(s.Len())
		s.Remove((*s)[idx])
		s.CheckThenFail()
	}
}

func (s setOfString) Less(i, j int) bool {
	return s[i] < s[j]
}

func (s setOfString) areItemsUnique() bool {
	var lastId string
	for _, a := range s {
		if lastId == a {
			return false
		}
		lastId = a
	}
	return true
}

func (s setOfString) Slice(marker string, max uint32) []string {
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

func (s setOfString) getIndex(id string) int {
	i := sort.Search(len(s), func(i int) bool {
		return s[i] >= id
	})
	if i < len(s) && s[i] == id {
		return i
	}
	return -1
}

func (s setOfString) Get(id string) string {
	var out string
	idx := s.getIndex(id)
	if idx >= 0 {
		out = s[idx]
	}
	return out
}

func (s setOfString) Has(id string) bool {
	return s.getIndex(id) >= 0
}

func (s *setOfString) Remove(a string) {
	s.RemovePK(a)
}

func (s *setOfString) RemovePK(pk string) {
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

type setOfSingleUint64 []*singleUint64Field

func (s setOfSingleUint64) CheckThenFail() {
	if err := s.Check(); err != nil {
		panic(err.Error())
	}
}

func (s setOfSingleUint64) Len() int {
	return len(s)
}

func (s setOfSingleUint64) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s *setOfSingleUint64) Add(a *singleUint64Field) {
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

func (s setOfSingleUint64) Check() error {
	if !sort.IsSorted(s) {
		return errors.NotValidf("sorting (%v) %v", s.Len(), s)
	}
	if !s.areItemsUnique() {
		return errors.NotValidf("unicity")
	}
	return nil
}

func (s *setOfSingleUint64) testRandomVacuum() {
	for s.Len() > 0 {
		idx := rand.Intn(s.Len())
		s.Remove((*s)[idx])
		s.CheckThenFail()
	}
}

func (s setOfSingleUint64) Less(i, j int) bool {
	return s[i].f0 < s[j].f0
}

func (s setOfSingleUint64) areItemsUnique() bool {
	var lastId uint64
	for _, a := range s {
		if lastId == a.f0 {
			return false
		}
		lastId = a.f0
	}
	return true
}

func (s setOfSingleUint64) Slice(marker uint64, max uint32) []*singleUint64Field {
	if max == 0 {
		max = 1000
	} else if max > 100000 {
		max = 100000
	}
	start := sort.Search(len(s), func(i int) bool {
		return s[i].f0 > marker
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

func (s setOfSingleUint64) getIndex(id uint64) int {
	i := sort.Search(len(s), func(i int) bool {
		return s[i].f0 >= id
	})
	if i < len(s) && s[i].f0 == id {
		return i
	}
	return -1
}

func (s setOfSingleUint64) Get(id uint64) *singleUint64Field {
	var out *singleUint64Field
	idx := s.getIndex(id)
	if idx >= 0 {
		out = s[idx]
	}
	return out
}

func (s setOfSingleUint64) Has(id uint64) bool {
	return s.getIndex(id) >= 0
}

func (s *setOfSingleUint64) Remove(a *singleUint64Field) {
	s.RemovePK(a.f0)
}

func (s *setOfSingleUint64) RemovePK(pk uint64) {
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

type setOfSingleString []*singleStringField

func (s setOfSingleString) CheckThenFail() {
	if err := s.Check(); err != nil {
		panic(err.Error())
	}
}

func (s setOfSingleString) Len() int {
	return len(s)
}

func (s setOfSingleString) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s *setOfSingleString) Add(a *singleStringField) {
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

func (s setOfSingleString) Check() error {
	if !sort.IsSorted(s) {
		return errors.NotValidf("sorting (%v) %v", s.Len(), s)
	}
	if !s.areItemsUnique() {
		return errors.NotValidf("unicity")
	}
	return nil
}

func (s *setOfSingleString) testRandomVacuum() {
	for s.Len() > 0 {
		idx := rand.Intn(s.Len())
		s.Remove((*s)[idx])
		s.CheckThenFail()
	}
}

func (s setOfSingleString) Less(i, j int) bool {
	return s[i].f0 < s[j].f0
}

func (s setOfSingleString) areItemsUnique() bool {
	var lastId string
	for _, a := range s {
		if lastId == a.f0 {
			return false
		}
		lastId = a.f0
	}
	return true
}

func (s setOfSingleString) Slice(marker string, max uint32) []*singleStringField {
	if max == 0 {
		max = 1000
	} else if max > 100000 {
		max = 100000
	}
	start := sort.Search(len(s), func(i int) bool {
		return s[i].f0 > marker
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

func (s setOfSingleString) getIndex(id string) int {
	i := sort.Search(len(s), func(i int) bool {
		return s[i].f0 >= id
	})
	if i < len(s) && s[i].f0 == id {
		return i
	}
	return -1
}

func (s setOfSingleString) Get(id string) *singleStringField {
	var out *singleStringField
	idx := s.getIndex(id)
	if idx >= 0 {
		out = s[idx]
	}
	return out
}

func (s setOfSingleString) Has(id string) bool {
	return s.getIndex(id) >= 0
}

func (s *setOfSingleString) Remove(a *singleStringField) {
	s.RemovePK(a.f0)
}

func (s *setOfSingleString) RemovePK(pk string) {
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

type setOfTwoFieldsSU []*twoFieldsSU

func (s setOfTwoFieldsSU) CheckThenFail() {
	if err := s.Check(); err != nil {
		panic(err.Error())
	}
}

func (s setOfTwoFieldsSU) Len() int {
	return len(s)
}

func (s setOfTwoFieldsSU) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s *setOfTwoFieldsSU) Add(a *twoFieldsSU) {
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

func (s setOfTwoFieldsSU) Check() error {
	if !sort.IsSorted(s) {
		return errors.NotValidf("sorting (%v) %v", s.Len(), s)
	}
	if !s.areItemsUnique() {
		return errors.NotValidf("unicity")
	}
	return nil
}

func (s *setOfTwoFieldsSU) testRandomVacuum() {
	for s.Len() > 0 {
		idx := rand.Intn(s.Len())
		s.Remove((*s)[idx])
		s.CheckThenFail()
	}
}

func (s setOfTwoFieldsSU) Less(i, j int) bool {
	p0, p1 := s[i], s[j]
	return p0.f0 < p1.f0 || (p0.f0 == p1.f0 && p0.f1 < p1.f1)
}

func (s setOfTwoFieldsSU) First(at string) int {
	return sort.Search(len(s), func(i int) bool { return s[i].f0 >= at })
}

func (s setOfTwoFieldsSU) areItemsUnique() bool {
	var l0 string
	var l1 uint64
	for _, a := range s {
		if l0 == a.f0 && l1 == a.f1 {
			return false
		}
		l0 = a.f0
	}
	return true
}

func (s setOfTwoFieldsSU) Slice(m0 string, m1 uint64, max uint32) []*twoFieldsSU {
	if max == 0 {
		max = 1000
	} else if max > 100000 {
		max = 100000
	}

	iMax := s.Len()
	start := s.First(m0)
	for start < iMax && s[start].f0 == m0 && s[start].f1 <= m1 {
		start++
	}

	remaining := uint32(iMax - start)
	if remaining > max {
		remaining = max
	}
	return s[start : uint32(start)+remaining]
}

func (s setOfTwoFieldsSU) getIndex(f0 string, f1 uint64) int {
	i := sort.Search(len(s), func(i int) bool {
		return s[i].f0 >= f0 || (s[i].f0 == f0 && s[i].f1 >= f1)
	})
	if i < len(s) && s[i].f0 == f0 && s[i].f1 == f1 {
		return i
	}
	return -1
}

func (s setOfTwoFieldsSU) Get(f0 string, f1 uint64) *twoFieldsSU {
	var out *twoFieldsSU
	idx := s.getIndex(f0, f1)
	if idx >= 0 {
		out = s[idx]
	}
	return out
}

func (s setOfTwoFieldsSU) Has(f0 string, f1 uint64) bool {
	return s.getIndex(f0, f1) >= 0
}

func (s *setOfTwoFieldsSU) Remove(a *twoFieldsSU) {
	s.RemovePK(a.f0, a.f1)
}

func (s *setOfTwoFieldsSU) RemovePK(f0 string, f1 uint64) {
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

type setOfTwoFieldsUS []*twoFieldsUS

func (s setOfTwoFieldsUS) CheckThenFail() {
	if err := s.Check(); err != nil {
		panic(err.Error())
	}
}

func (s setOfTwoFieldsUS) Len() int {
	return len(s)
}

func (s setOfTwoFieldsUS) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s *setOfTwoFieldsUS) Add(a *twoFieldsUS) {
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

func (s setOfTwoFieldsUS) Check() error {
	if !sort.IsSorted(s) {
		return errors.NotValidf("sorting (%v) %v", s.Len(), s)
	}
	if !s.areItemsUnique() {
		return errors.NotValidf("unicity")
	}
	return nil
}

func (s *setOfTwoFieldsUS) testRandomVacuum() {
	for s.Len() > 0 {
		idx := rand.Intn(s.Len())
		s.Remove((*s)[idx])
		s.CheckThenFail()
	}
}

func (s setOfTwoFieldsUS) Less(i, j int) bool {
	p0, p1 := s[i], s[j]
	return p0.f0 < p1.f0 || (p0.f0 == p1.f0 && p0.f1 < p1.f1)
}

func (s setOfTwoFieldsUS) First(at uint64) int {
	return sort.Search(len(s), func(i int) bool { return s[i].f0 >= at })
}

func (s setOfTwoFieldsUS) areItemsUnique() bool {
	var l0 uint64
	var l1 string
	for _, a := range s {
		if l0 == a.f0 && l1 == a.f1 {
			return false
		}
		l0 = a.f0
	}
	return true
}

func (s setOfTwoFieldsUS) Slice(m0 uint64, m1 string, max uint32) []*twoFieldsUS {
	if max == 0 {
		max = 1000
	} else if max > 100000 {
		max = 100000
	}

	iMax := s.Len()
	start := s.First(m0)
	for start < iMax && s[start].f0 == m0 && s[start].f1 <= m1 {
		start++
	}

	remaining := uint32(iMax - start)
	if remaining > max {
		remaining = max
	}
	return s[start : uint32(start)+remaining]
}

func (s setOfTwoFieldsUS) getIndex(f0 uint64, f1 string) int {
	i := sort.Search(len(s), func(i int) bool {
		return s[i].f0 >= f0 || (s[i].f0 == f0 && s[i].f1 >= f1)
	})
	if i < len(s) && s[i].f0 == f0 && s[i].f1 == f1 {
		return i
	}
	return -1
}

func (s setOfTwoFieldsUS) Get(f0 uint64, f1 string) *twoFieldsUS {
	var out *twoFieldsUS
	idx := s.getIndex(f0, f1)
	if idx >= 0 {
		out = s[idx]
	}
	return out
}

func (s setOfTwoFieldsUS) Has(f0 uint64, f1 string) bool {
	return s.getIndex(f0, f1) >= 0
}

func (s *setOfTwoFieldsUS) Remove(a *twoFieldsUS) {
	s.RemovePK(a.f0, a.f1)
}

func (s *setOfTwoFieldsUS) RemovePK(f0 uint64, f1 string) {
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
