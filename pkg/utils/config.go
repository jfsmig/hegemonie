// Copyright (c) 2018-2021 Contributors as noted in the AUTHORS file
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package utils

import (
	"fmt"
	"github.com/juju/errors"
	"gopkg.in/yaml.v3"
	"os"
	"time"
)

const (
	// DefaultPortKratos is the default port of the Kratos authentication service by ORY
	DefaultPortKratos = 4434
	// DefaultPortKeto is the default port of the Keto autorisation service by ORY
	DefaultPortKeto = 4466

	// DefaultPortCommon is the default port used to hast the service.
	// This is also the configuration expected by the docker image.
	DefaultPortCommon = 6000

	// DefaultPortMonitoring is the default port used for clear-text HTTP
	// serving the /metrics route dedicated to Prometheus exporters.
	DefaultPortMonitoring = 6001
)

type TargetConfig struct {
	Address string        `yaml:"addr" json:"addr"`
	Timeout time.Duration `yaml:"timeout" json:"timeout"`
	PathKey string        `yaml:"key" json:"key"`
	PathCrt string        `yaml:"cert" json:"cert"`
}

type ClientConfig struct {
	Default TargetConfig `yaml:"default" json:"default"`
	Proxy   TargetConfig `yaml:"proxy" json:"proxy"`
	Maps    TargetConfig `yaml:"maps" json:"maps"`
	Events  TargetConfig `yaml:"events" json:"events"`
	Regions TargetConfig `yaml:"regions" json:"regions"`
}

type ServerConfig struct {
	EndpointService string `yaml:"bind" json:"bind"`
	EndpointMonitor string `yaml:"monitor" json:"monitor"`
	ServiceType     string `yaml:"type" json:"type"`
	PathKey         string `yaml:"key" json:"key"`
	PathCrt         string `yaml:"cert" json:"cert"`

	MapConfig MapServiceConfig    `yaml:"map" json:"map"`
	EvtConfig EventServiceConfig  `yaml:"evt" json:"evt"`
	RegConfig RegionServiceConfig `yaml:"reg" json:"reg"`
}

type EventServiceConfig struct {
	PathBase string `yaml:"base" json:"base"`
}

type MapServiceConfig struct {
	PathRepository string `yaml:"repository" json:"repository"`
}

type RegionServiceConfig struct {
	PathDefs string `yaml:"definitions" json:"definitions"`
	PathLive string `yaml:"live" json:"live"`
}

type MainConfig struct {
	Client ClientConfig `yaml:"client" json:"client"`
	Server ServerConfig `yaml:"server" json:"server"`
}

// DefaultTargetConfig
func DefaultTargetConfig() TargetConfig {
	return TargetConfig{
		Address: fmt.Sprintf("127.0.0.1:%v", DefaultPortCommon),
	}
}

// DefaultClientConfig
func DefaultClientConfig() ClientConfig {
	return ClientConfig{
		Default: TargetConfig{
			Timeout: time.Second,
			PathCrt: "TLS-certificate-NOT-SET",
			PathKey: "TLS-key-NOT-SET",
		},
		Proxy:   TargetConfig{},
		Maps:    DefaultTargetConfig(),
		Events:  DefaultTargetConfig(),
		Regions: DefaultTargetConfig(),
	}
}

// DefaultConfig
func DefaultConfig() MainConfig {
	return MainConfig{
		Client: DefaultClientConfig(),
		Server: ServerConfig{
			EndpointService: fmt.Sprintf("0.0.0.0:%v", DefaultPortCommon),
			EndpointMonitor: fmt.Sprintf("0.0.0.0:%v", DefaultPortMonitoring),
			ServiceType:     "type-NOT-SET",
			PathCrt:         "TLS-certificate-NOT-SET",
			PathKey:         "TLS-key-NOT-SET",
		},
	}
}

// LoadFile
func (cfg *MainConfig) LoadFile(path string, must bool) error {
	if path == "" {
		return nil
	}

	fin, err := os.Open(path)
	if err != nil {
		if must {
			return errors.Annotate(err, "invalid configuration file")
		}
		Logger.Debug().Str("path", path).Msg("Not Found")
		return nil
	}
	decoder := yaml.NewDecoder(fin)
	if err = decoder.Decode(cfg); err != nil {
		return errors.Annotate(err, "malformed configuration")
	}

	return cfg.Client.ApplyToDiscovery()
}

// ApplyToDiscovery
func (cfg *ClientConfig) ApplyToDiscovery() error {
	sc := new(StaticConfig)

	// set the addresses
	if cfg.Maps.Address != "" {
		sc.SetMap(cfg.Maps.Address)
	}
	if cfg.Events.Address != "" {
		sc.SetEvent(cfg.Events.Address)
	}
	if cfg.Regions.Address != "" {
		sc.SetRegion(cfg.Regions.Address)
	}
	if cfg.Proxy.Address != "" {
		sc.SetEvent(cfg.Proxy.Address)
		sc.SetMap(cfg.Proxy.Address)
		sc.SetRegion(cfg.Proxy.Address)
	}

	DefaultDiscovery = sc
	return nil
}
