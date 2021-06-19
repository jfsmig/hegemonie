// Copyright (c) 2018-2021 Contributors as noted in the AUTHORS file
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package utils

import (
	"fmt"
	"github.com/go-openapi/errors"
)

// StatelessDiscovery is the simplest form of Discovery API providing one call per
// type of service. Each call returns either a usable endpoint string or the error
// that occurred during the discovery process.
// The implementation of the StatelessDiscovery interface is responsible for the
// management of its concurrent accesses.
type StatelessDiscovery interface {
	// Kratos locates an ORY kratos service (Authentication)
	Kratos() (string, error)

	// Keto locates an ORY keto service (Authorisation)
	Keto() (string, error)

	// Map locates map services in Hegemonie
	Map() (string, error)

	// Region locates an hegemonie's region service
	// Please note that those services are typically sharded. Stateless weighted polling
	// is only meaningful when it is necessary to instantiate a new Region.
	Region() (string, error)

	// Event locates an hegemonie's event services
	Event() (string, error)
}

// DefaultDiscovery is the default implementation of a discovery.
// Valued by default to the discovery of test services, all located on
// localhost and serving default ports.
var DefaultDiscovery = NewStaticConfig()

// StaticConfig is a StatelessDiscovery implementation with a different
// endpoint for each kind of service, configured once in the application.
type StaticConfig struct {
	kratos  string
	keto    string
	maps    string
	regions string
	events  string
}

func f1(p int) string { return fmt.Sprintf("127.0.0.1:%v", p) }

// NewStaticConfig instantiates a StaticConfig with the default endpoint value
// for each service type.
func NewStaticConfig() StatelessDiscovery {
	return &StaticConfig{
		maps:    f1(DefaultPortCommon),
		events:  f1(DefaultPortCommon),
		regions: f1(DefaultPortCommon),
		keto:    f1(DefaultPortKeto),
		kratos:  f1(DefaultPortKratos),
	}
}

func (d *StaticConfig) nyi() (string, error) { return "", errors.NotImplemented("NYI") }

// Kratos ... see StatelessDiscovery.Kratos
func (d *StaticConfig) Kratos() (string, error) { return d.nyi() }

// Keto ... see StatelessDiscovery.Keto
func (d *StaticConfig) Keto() (string, error) { return d.nyi() }

// Map ... see StatelessDiscovery.Map
func (d *StaticConfig) Map() (string, error) { return d.maps, nil }

// Region ... see StatelessDiscovery.Region
func (d *StaticConfig) Region() (string, error) { return d.regions, nil }

// Event ... see StatelessDiscovery.Event
func (d *StaticConfig) Event() (string, error) { return d.events, nil }

// SetMap updates the endpoint of the maps service
func (d *StaticConfig) SetMap(ep string) { d.maps = ep }

// SetRegion updates the endpoint of the region service
func (d *StaticConfig) SetRegion(ep string) { d.regions = ep }

// SetEvent updates the endpoint of the event service
func (d *StaticConfig) SetEvent(ep string) { d.events = ep }
