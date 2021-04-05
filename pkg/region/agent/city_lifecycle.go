// Copyright (c) 2018-2021 Contributors as noted in the AUTHORS file
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package regagent

type UserCityController interface {
	// Config applies to the current City a model identified by its ID
	Config(modelID string) error

	// Acquire takes the ownership on a City and assigns it to userID
	Acquire(userID string) error

	// Leave voluntarily abandon the ownership on the current City
	Leave() error

	// Auto toggles the automatic mode on the current City
	Auto() error
}

type GMCityController interface {
	// Config applies a model City on the current City
	Config(modelID string) error

	// Assign sets the ownership on the current City
	Assign(userID string) error

	// Resume allows the user to control the City
	Resume() error

	// Dismiss removes the ownership on the current City
	Dismiss() error

	// Suspend prevents the user from controlling the City
	Suspend() error

	// Reset removes any activity footprint from the current City
	Reset() error
}
