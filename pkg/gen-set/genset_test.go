// Copyright (c) 2018-2021 Contributors as noted in the AUTHORS file
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"github.com/google/uuid"
	"math/rand"
	"testing"
)

type singleStringField struct {
	f0 string
}

type singleUint64Field struct {
	f0 uint64
}

type twoFieldsSU struct {
	f0 string
	f1 uint64
}

type twoFieldsUS struct {
	f0 uint64
	f1 string
}

func rstr() string { return uuid.New().String() }
func rint() uint64 { return rand.Uint64() }
func rounds() int  { return 256 + rand.Intn(256) }

func TestSetOfString(t *testing.T) {
	s := make(setOfString, 0)
	for i, max := 0, rounds(); i < max; i++ {
		s.Add(rstr())
		s.CheckThenFail()
	}
	s.testRandomVacuum()
}

func TestSetOfUint64(t *testing.T) {
	s := make(setOfUint64, 0)
	for i, max := 0, rounds(); i < max; i++ {
		s.Add(rint())
		s.CheckThenFail()
	}
	s.testRandomVacuum()
}

func TestSetOfSingleString(t *testing.T) {
	s := make(setOfSingleString, 0)
	for i, max := 0, rounds(); i < max; i++ {
		s.Add(&singleStringField{rstr()})
		s.CheckThenFail()
	}
	s.testRandomVacuum()
}

func TestSetOfSingleUint64(t *testing.T) {
	s := make(setOfSingleUint64, 0)
	for i, max := 0, rounds(); i < max; i++ {
		s.Add(&singleUint64Field{rint()})
		s.CheckThenFail()
	}
	s.testRandomVacuum()
}

func TestSetOfTwoFieldSU(t *testing.T) {
	s := make(setOfTwoFieldsSU, 0)
	for i, max := 0, rounds(); i < max; i++ {
		s.Add(&twoFieldsSU{rstr(), rint()})
		s.CheckThenFail()
	}
	s.testRandomVacuum()
}

func TestSetOfTwoFieldUS(t *testing.T) {
	s := make(setOfTwoFieldsSU, 0)
	for i, max := 0, rounds(); i < max; i++ {
		s.Add(&twoFieldsSU{rstr(), rint()})
		s.CheckThenFail()
	}
	s.testRandomVacuum()
}

// Only the tests around the functional cases are tested here.
//go:generate go run github.com/hegemonie-rpg/engine/pkg/gen-set ./genset_auto_test.go main:setOfUint64:uint64 :uint64
//go:generate go run github.com/hegemonie-rpg/engine/pkg/gen-set ./genset_auto_test.go main:setOfString:string :string
//go:generate go run github.com/hegemonie-rpg/engine/pkg/gen-set ./genset_auto_test.go main:setOfSingleUint64:*singleUint64Field f0:uint64
//go:generate go run github.com/hegemonie-rpg/engine/pkg/gen-set ./genset_auto_test.go main:setOfSingleString:*singleStringField f0:string
//go:generate go run github.com/hegemonie-rpg/engine/pkg/gen-set ./genset_auto_test.go main:setOfTwoFieldsSU:*twoFieldsSU f0:string f1:uint64
//go:generate go run github.com/hegemonie-rpg/engine/pkg/gen-set ./genset_auto_test.go main:setOfTwoFieldsUS:*twoFieldsUS f0:uint64 f1:string
