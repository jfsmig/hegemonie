// Copyright (c) 2018-2021 Contributors as noted in the AUTHORS file
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package region

import (
	"testing"
)

func TestResources(t *testing.T) {
	var zero, one, r0 Resources
	if !zero.IsZero() || !r0.IsZero() || !one.IsZero() {
		t.Fatal()
	}
	if !r0.GreaterOrEqualTo(r0) {
		t.Fatal()
	}
	one[1] = 1
	if !one.GreaterOrEqualTo(one) {
		t.Fatal()
	}
	if one.IsZero() {
		t.Fatal()
	}
	if !one.GreaterOrEqualTo(r0) || r0.GreaterOrEqualTo(one) {
		t.Fatal()
	}

	r0.Add(one)
	if !r0.Equals(one) {
		t.Fatal()
	}
	r0.Add(one)
	if r0.Equals(one) {
		t.Fatal()
	}
	r0.TrimTo(zero)
	if !r0.IsZero() {
		t.Fatal()
	}
	if !r0.Equals(zero) {
		t.Fatal()
	}

	one.Remove(one)
	if !one.IsZero() {
		t.Fatal()
	}

	one[1] = 1
	one.Zero()
	if !one.IsZero() {
		t.Fatal()
	}

	allOneAbs := ResourcesUniform(1)
	for i := 0; i < ResourceMax; i++ {
		if allOneAbs[i] != 1 {
			t.Fatal()
		}
	}
}

func TestModifiers(t *testing.T) {
	inc := IncrementUniform(1)
	for i := 0; i < ResourceMax; i++ {
		if inc[i] != 1 {
			t.Fatal()
		}
	}

	mult := MultiplierUniform(2.5)
	for i := 0; i < ResourceMax; i++ {
		if mult[i] != 2.5 {
			t.Fatal()
		}
	}

	abs := ResourcesUniform(2)
	abs.Multiply(mult)
	for i := 0; i < ResourceMax; i++ {
		if abs[i] != 5 {
			t.Fatal()
		}
	}

	mod := ResourceModifierNoop()
	mod1 := ResourceModifiers{Mult: mult, Plus: inc}
	mod.ComposeWith(mod1)
	mod.ComposeWith(mod1)
	for i := 0; i < ResourceMax; i++ {
		if mod.Plus[i] != mod1.Plus[i]*2 {
			t.Fatal()
		}
		if mod.Mult[i] != mod1.Mult[i]*mod1.Mult[i] {
			t.Fatal()
		}
	}

	for _, op := range []ResourceModifiers{
		ResourceModifierUniform(5.0, 0.0),
		ResourceModifierUniform(1.0, 7.0),
	} {
		if ResourceModifierNoop().Equals(op) {
			t.Fatal()
		}
	}
}

func TestCompose(t *testing.T) {
	mod := ResourceModifierNoop()
	mod.ComposeWith(ResourceModifierNoop())
	if !mod.Equals(ResourceModifierNoop()) {
		t.Fatal()
	}

	mod = ResourceModifierUniform(1.0, 1.0)
	mod.ComposeWith(mod)
	if !mod.Equals(ResourceModifierUniform(1.0, 2.0)) {
		t.Fatal()
	}
}
