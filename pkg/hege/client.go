// Copyright (c) 2018-2021 Contributors as noted in the AUTHORS file
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"context"
	"github.com/google/uuid"
	authclient "github.com/hegemonie-rpg/engine/pkg/auth/client"
	evtclient "github.com/hegemonie-rpg/engine/pkg/event/client"
	mapclient "github.com/hegemonie-rpg/engine/pkg/map/client"
	regclient "github.com/hegemonie-rpg/engine/pkg/region/client"
	"github.com/hegemonie-rpg/engine/pkg/utils"
	"github.com/juju/errors"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/metadata"
	"gopkg.in/yaml.v3"
	"os"
	"strconv"
	"time"
)

type endpointConfig struct {
	Addr    string        `yaml:"addr" json:"addr"`
	Timeout time.Duration `yaml:"timeout" json:"timeout"`
}

func (ec *endpointConfig) reset() {
	ec.Addr = ""
	ec.Timeout = 0
}

type clientConfig struct {
	Proxy   endpointConfig `yaml:"proxy" json:"proxy"`
	Maps    endpointConfig `yaml:"maps" json:"maps"`
	Events  endpointConfig `yaml:"events" json:"events"`
	Regions endpointConfig `yaml:"regions" json:"regions"`
}

func (cc *clientConfig) reset() {
	for _, cfg := range []endpointConfig{cc.Proxy, cc.Events, cc.Regions, cc.Maps} {
		cfg.reset()
	}
}

func clients() *cobra.Command {
	var pathConfig string

	cmd := &cobra.Command{
		Use:   "client",
		Short: "Client tool for various Hegemonie services",
		Args:  cobra.MinimumNArgs(1),
		RunE:  nonLeaf,
	}

	// Set a common reasonable timeout to all client RPC
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	sessionID := os.Getenv("HEGE_CLI_SESSIONID")
	if sessionID == "" {
		sessionID = "cli/" + uuid.New().String()
	}

	// Inherit a session-id from the env
	ctx = metadata.AppendToOutgoingContext(ctx, "session-id", sessionID)

	// Override the discovery if a proxy is configured
	cmd.PersistentFlags().StringVarP(
		&pathConfig, "config", "f", "",
		"IP:PORT endpoint for the gRPC proxy")

	cmd.PersistentPreRunE = func(_ *cobra.Command, _ []string) error {
		home, err := os.UserHomeDir()
		if err != nil {
			return errors.Annotate(err, "home directory error")
		}
		err = loadDiscovery("/etc/hegemonie/client.yml", false)
		if err != nil {
			return errors.Annotate(err, "system configuration")
		}
		err = loadDiscovery(home+"/.hegemonie/client.yml", false)
		if err != nil {
			return errors.Annotate(err, "home configuration")
		}
		err = loadDiscovery(pathConfig, true)
		if err != nil {
			return errors.Annotate(err, "explicit configuration")
		}
		return nil
	}
	//cmd.PersistentPostRun = func(_ *cobra.Command, _ []string) { cancel() }

	cmd.AddCommand(
		clientMap(ctx), clientEvent(ctx), clientAuth(ctx),
		clientRegion(ctx), clientCity(ctx), clientTemplates(ctx))
	return cmd
}

func loadDiscovery(path string, must bool) error {
	var config clientConfig

	if path == "" {
		return nil
	}

	fin, err := os.Open(path)
	if err != nil {
		if must {
			return errors.Annotate(err, "invalid configuration file")
		}
		utils.Logger.Debug().Str("path", path).Msg("Not Found")
		return nil
	}
	decoder := yaml.NewDecoder(fin)
	if err = decoder.Decode(&config); err != nil {
		return errors.Annotate(err, "malformed configuration")
	}

	// Override the configuration if any value is specified
	if config.Proxy.Addr != "" {
		utils.DefaultDiscovery = utils.SingleEndpoint(config.Proxy.Addr)
	} else {
		sc := utils.NewStaticConfig().(*utils.StaticConfig)
		if config.Maps.Addr != "" {
			sc.SetMap(config.Maps.Addr)
		}
		if config.Events.Addr != "" {
			sc.SetEvent(config.Events.Addr)
		}
		if config.Regions.Addr != "" {
			sc.SetRegion(config.Regions.Addr)
		}
		utils.DefaultDiscovery = sc
	}

	utils.Logger.Info().Str("path", path).RawJSON("cfg", utils.JSON2Buf(config)).Msg("Loaded")
	return nil
}

func clientMap(ctx context.Context) *cobra.Command {
	var cfg mapclient.ClientCLI
	var pathArgs mapclient.PathArgs

	hook := func(action func() error) func(cmd *cobra.Command, args []string) error {
		return func(cmd *cobra.Command, args []string) error {
			if err := pathArgs.Parse(args); err != nil {
				return errors.Trace(err)
			}
			return action()
		}
	}

	cmd := &cobra.Command{
		Use:   "maps",
		Short: "Client of a Maps API service",
		Args:  cobra.MinimumNArgs(1),
		RunE:  nonLeaf,
	}

	list := &cobra.Command{
		Use:     "list",
		Short:   "List all the maps registered",
		Example: "map list [$MAPID_MARKER]",
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			pathArgs.MapName = first(args)
			return cfg.ListMaps(ctx, pathArgs)
		},
	}

	path := &cobra.Command{
		Use:     "path",
		Short:   "Compute the path between two nodes",
		Example: "map path $MAPID $SRC $DST",
		Args:    cobra.ExactArgs(3),
		RunE:    hook(func() error { return cfg.GetPath(ctx, pathArgs) }),
	}
	path.Flags().Uint32VarP(&pathArgs.Max, "max", "m", 0, "Max path length")

	step := &cobra.Command{
		Use:     "step",
		Short:   "Get the next step of the path between two nodes",
		Example: "map step $REGION $SRC $DST",
		Args:    cobra.ExactArgs(3),
		RunE:    hook(func() error { return cfg.GetStep(ctx, pathArgs) }),
	}

	cities := &cobra.Command{
		Use:     "cities",
		Short:   "List the Cities when the map is instantiated",
		Example: "map cities $REGION [$MARKER]",
		Args:    cobra.RangeArgs(1, 2),
		RunE:    hook(func() error { return cfg.GetCities(ctx, pathArgs) }),
	}
	cities.Flags().Uint32VarP(&pathArgs.Max, "max", "m", 0, "List max N cities")

	roads := &cobra.Command{
		Use:     "roads",
		Short:   "List of the roads of the map",
		Example: "map roads $REGION [$MARKER_SRC [$MARKER_DST]]",
		Args:    cobra.RangeArgs(1, 3),
		RunE:    hook(func() error { return cfg.GetRoads(ctx, pathArgs) }),
	}
	roads.Flags().Uint32VarP(&pathArgs.Max, "max", "m", 0, "List max N roads")

	positions := &cobra.Command{
		Use:     "positions",
		Short:   "List the positions of the map",
		Example: "map positions $REGION [$MARKER]",
		Args:    cobra.RangeArgs(1, 2),
		RunE:    hook(func() error { return cfg.GetPositions(ctx, pathArgs) }),
	}
	positions.Flags().Uint32VarP(&pathArgs.Max, "max", "m", 0, "List max N positions")

	cmd.AddCommand(list, path, step, cities, roads, positions)
	return cmd
}

func clientEvent(ctx context.Context) *cobra.Command {
	var max uint32
	var cfg evtclient.ClientCLI

	cmd := &cobra.Command{
		Use:   "events",
		Short: "Client of an Events API service",
		Args:  cobra.MinimumNArgs(1),
		RunE:  nonLeaf,
	}

	push := &cobra.Command{
		Use:     "push",
		Short:   "Push events in the Character's log",
		Example: `server event push "${CHARACTER}" "${MSG0}" "${MSG1}"`,
		Args:    cobra.MinimumNArgs(2),
		RunE:    func(cmd *cobra.Command, args []string) error { return cfg.DoPush(ctx, args[0], args[1:]...) },
	}

	list := &cobra.Command{
		Use:     "list",
		Short:   "List the events",
		Example: `server event list "${CHARACTER}" "${EVENT_TIMESTAMP}" [${EVENT_MARKER}]`,
		Args:    cobra.RangeArgs(1, 3),
		RunE: func(cmd *cobra.Command, args []string) error {
			var when uint64
			var marker string
			if len(args) > 1 {
				var err error
				when, err = strconv.ParseUint(args[1], 10, 63)
				if err != nil {
					return errors.Trace(err)
				}
				if len(args) > 2 {
					marker = args[2]
				}
			}
			return cfg.DoList(ctx, args[0], when, marker, max)
		},
	}
	list.Flags().Uint32VarP(&max, "max", "m", 0, "List at most N events")

	ack := &cobra.Command{
		Use:     "ack",
		Short:   "Acknowledge an event",
		Example: `server event ack "${CHARACTER}" "${EVENT_UUID}" "${EVENT_TIMESTAMP}"`,
		Args:    cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			var when uint64
			if len(args) > 2 {
				var err error
				when, err = strconv.ParseUint(args[1], 10, 63)
				if err != nil {
					return errors.Trace(err)
				}
			}
			return cfg.DoAck(ctx, args[0], args[1], when)
		},
	}

	cmd.AddCommand(push, ack, list)
	return cmd

}

func clientAuth(ctx context.Context) *cobra.Command {
	var cfg authclient.ClientCLI

	cmd := &cobra.Command{
		Use:     "auth",
		Short:   "Authorization and Authentication client",
		Example: "auth (users|details|create|invite|affect) ...",
		Args:    cobra.MinimumNArgs(1),
		RunE:    nonLeaf,
	}

	users := &cobra.Command{
		Use:     "users",
		Short:   "List the registered USERS",
		Example: "auth list",
		Args:    cobra.NoArgs,
		RunE:    func(cmd *cobra.Command, args []string) error { return cfg.DoList(ctx, args) },
	}

	details := &cobra.Command{
		Use:     "detail",
		Short:   "Show the details of specific users",
		Long:    "Print a detailed JSON representation of the information and permissions for each user specified as a positional argument",
		Example: "show a4ddeee6-b72a-4a27-8e2d-35c3cc62c7d3 ab2bca77-efdb-4dc2-b80a-fc03e0fc5226 ...",
		Args:    cobra.MinimumNArgs(1),
		RunE:    func(cmd *cobra.Command, args []string) error { return cfg.DoShow(ctx, args) },
	}

	create := &cobra.Command{
		Use:     "create",
		Short:   "Create a User",
		Example: "auth create forced.user@example.com",
		Args:    cobra.MinimumNArgs(1),
		RunE:    func(cmd *cobra.Command, args []string) error { return cfg.DoCreate(ctx, args) },
	}

	invite := &cobra.Command{
		Use:     "invite",
		Short:   "Invite a user identified by its email",
		Example: "auth invite invited.user@example.com",
		Args:    cobra.ExactArgs(1),
		RunE:    func(cmd *cobra.Command, args []string) error { return cfg.DoInvite(ctx, args) },
	}

	affect := &cobra.Command{
		Use:   "affect",
		Short: "Invite a user identified by its email",
		RunE:  func(cmd *cobra.Command, args []string) error { return cfg.DoInvite(ctx, args) },
	}

	cmd.AddCommand(users, details, create, invite, affect)
	return cmd
}

func clientRegion(ctx context.Context) *cobra.Command {
	var cfg regclient.ClientCLI

	cmd := &cobra.Command{
		Use:     "regions",
		Aliases: []string{"reg", "region"},
		Short:   "Client of a Regions API service",
		Example: "hege client regions (create|list) ...",
		Args:    cobra.MinimumNArgs(1),
		RunE:    nonLeaf,
	}

	createRegion := &cobra.Command{
		Use:     "create",
		Aliases: []string{"new"},
		Short:   "Create a new region",
		Example: "hege client regions create $REGION_ID $MAP_ID",
		Args:    cobra.ExactArgs(2),
		RunE:    func(cmd *cobra.Command, args []string) error { return cfg.DoCreateRegion(ctx, args[0], args[1]) },
	}

	listRegions := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls", "all"},
		Short:   "List the existing regions",
		Example: "hege client regions list [$REGION_ID_MARKER]",
		Args:    cobra.MaximumNArgs(1),
		RunE:    func(cmd *cobra.Command, args []string) error { return cfg.DoListRegions(ctx, first(args)) },
	}

	roundMovement := &cobra.Command{
		Use:     "move",
		Aliases: []string{"movement", "mov"},
		Short:   "Execute a movement round on the region",
		Example: "hege client regions move $REGION_ID",
		Args:    cobra.ExactArgs(1),
		RunE:    func(cmd *cobra.Command, args []string) error { return cfg.DoRegionMovement(ctx, args[0]) },
	}

	roundProduction := &cobra.Command{
		Use:     "produce",
		Aliases: []string{"production", "prod"},
		Short:   "Execute a movement round on the region",
		Example: "hege client regions move $REGION_ID",
		Args:    cobra.ExactArgs(1),
		RunE:    func(cmd *cobra.Command, args []string) error { return cfg.DoRegionProduction(ctx, args[0]) },
	}

	pushStats := &cobra.Command{
		Use:     "stats_refresh",
		Aliases: []string{"refresh", "stx"},
		Short:   "Trigger a stats refresh by the region service, for the given region",
		Example: "hege client regions refresh $REGION_ID",
		Args:    cobra.ExactArgs(1),
		RunE:    func(cmd *cobra.Command, args []string) error { return cfg.DoRegionPushStats(ctx, args[0]) },
	}

	getStats := &cobra.Command{
		Use:     "stats",
		Aliases: []string{"stat", "st"},
		Short:   "Get the stats board of the region",
		Example: "hege client regions stats $REGION_ID",
		Args:    cobra.ExactArgs(1),
		RunE:    func(cmd *cobra.Command, args []string) error { return cfg.DoRegionGetStats(ctx, args[0]) },
	}

	getScore := &cobra.Command{
		Use:     "score",
		Aliases: []string{"scores"},
		Short:   "Get the score board of the region",
		Example: "hege client regions score $REGION_ID",
		Args:    cobra.ExactArgs(1),
		RunE:    func(cmd *cobra.Command, args []string) error { return cfg.DoRegionGetScores(ctx, args[0]) },
	}

	buildings := &cobra.Command{
		Use:     "buildings",
		Aliases: []string{"b", "building"},
		Short:   "Get the available buildings",
		Example: "hege client regions buildings $REGION_ID",
		Args:    cobra.ExactArgs(1),
		RunE:    func(cmd *cobra.Command, args []string) error { return cfg.DoRegionGetBuildings(ctx, args[0]) },
	}

	skills := &cobra.Command{
		Use:     "skills",
		Aliases: []string{"s", "skill", "k", "knowledge"},
		Short:   "Get the available buildings",
		Example: "hege client regions buildings $REGION_ID",
		Args:    cobra.ExactArgs(1),
		RunE:    func(cmd *cobra.Command, args []string) error { return cfg.DoRegionGetSkills(ctx, args[0]) },
	}

	units := &cobra.Command{
		Use:     "units",
		Aliases: []string{"u", "unit", "troops", "troop"},
		Short:   "Get the available buildings",
		Example: "hege client regions buildings $REGION_ID",
		Args:    cobra.ExactArgs(1),
		RunE:    func(cmd *cobra.Command, args []string) error { return cfg.DoRegionGetUnits(ctx, args[0]) },
	}

	cmd.AddCommand(
		createRegion, listRegions,
		roundMovement, roundProduction,
		pushStats, getStats,
		getScore,
		buildings, skills, units)
	return cmd
}

func clientCity(ctx context.Context) *cobra.Command {
	var cfg regclient.ClientCLI

	cmd := &cobra.Command{
		Use:     "city",
		Short:   "Client of a Regions API service",
		Example: "hege client regions (create|list) ...",
		Args:    cobra.MinimumNArgs(1),
		RunE:    nonLeaf,
	}

	show := &cobra.Command{
		Use:     "show",
		Aliases: []string{"all", "detail", "details"},
		Short:   "Show the City details",
		Example: "hege client city show $REGION_ID $CITY_ID",
		Args:    cobra.ExactArgs(2),
		RunE:    func(cmd *cobra.Command, args []string) error { return cfg.DoCityShow(ctx, args[0], args[1]) },
	}

	build := &cobra.Command{
		Use:     "build",
		Aliases: []string{"construct"},
		Short:   "Start the construction of a building",
		Example: "hege client city build $REGION_ID $CITY_ID $BUILDING_TYPE__ID",
		Args:    cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			return cfg.DoCityBuild(ctx, args[0], args[1], args[2])
		},
	}

	learn := &cobra.Command{
		Use:     "study",
		Aliases: []string{"learn"},
		Short:   "Start the study of a skill",
		Example: "hege client city study $REGION_ID $CITY_ID $SKILL_TYPE_ID",
		Args:    cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			return cfg.DoCityStudy(ctx, args[0], args[1], args[2])
		},
	}

	train := &cobra.Command{
		Use:     "train",
		Short:   "Start the training of a new unit",
		Example: "hege client city train $REGION_ID $CITY_ID $UNIT_TYPE_ID",
		Args:    cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			return cfg.DoCityTrain(ctx, args[0], args[1], args[2])
		},
	}

	cmd.AddCommand(show, build, learn, train)
	return cmd
}

func clientTemplates(ctx context.Context) *cobra.Command {
	var cfg regclient.ClientCLI

	cmd := &cobra.Command{
		Use:     "templates",
		Aliases: []string{"template", "tpl"},
		Short:   "Client of a Regions API service",
		Example: "hege client templates (list|create|update|delete) ...",
		Args:    cobra.MinimumNArgs(1),
		RunE:    nonLeaf,
	}

	list := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls", "get", "show"},
		Short:   "Show the Templates details",
		Example: "hege client templates list $REGION_ID [$TPL_NAME]",
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return cfg.DoTemplateList(ctx, args[0], second(args))
		},
	}

	create := &cobra.Command{
		Use:     "create",
		Aliases: []string{"new", "add", "put"},
		Short:   "Create a Template whose details are read from stdin",
		Example: "hege client templates create $REGION_ID",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return cfg.DoTemplateCreate(ctx, args[0])
		},
	}

	update := &cobra.Command{
		Use:     "update",
		Aliases: []string{"set", "configure", "post"},
		Short:   "Update the Templates details",
		Example: "hege client templates create $REGION_ID $TPL_NAME",
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return cfg.DoTemplateUpdate(ctx, args[0], args[1])
		},
	}

	del := &cobra.Command{
		Use:     "delete",
		Aliases: []string{"del", "erase", "remove", "rem"},
		Short:   "Delete te template with the given ID",
		Example: "hege client templates delete $REGION_ID $TPL_NAME",
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return cfg.DoTemplateDelete(ctx, args[0], args[1])
		},
	}

	cmd.AddCommand(list, create, update, del)
	return cmd
}

func first(args []string) string {
	if len(args) <= 0 {
		return ""
	}
	return args[0]
}

func second(args []string) string {
	if len(args) <= 1 {
		return ""
	}
	return args[1]
}
