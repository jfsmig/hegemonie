// Copyright (c) 2018-2021 Contributors as noted in the AUTHORS file
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package regagent

import (
	"context"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/hegemonie-rpg/engine/pkg/region/model"
	"github.com/hegemonie-rpg/engine/pkg/region/proto"
	"github.com/hegemonie-rpg/engine/pkg/utils"
	"github.com/juju/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
	"strconv"
)

// Config gathers the configuration fields required to start a gRPC region API service.
type Config struct {
	PathDefs string `yaml:"definitions" json:"definitions"`
	PathLive string `yaml:"live" json:"live"`
}

type regionBackend struct {
	cfg Config
	w   *region.World
}

var none = &proto.None{}

// Application implements the expectations of the application backend
func (cfg Config) Application(ctx context.Context) (utils.RegisterableMonitorable, error) {
	w, err := region.NewWorld()
	if err != nil {
		return nil, errors.Annotate(err, "")
	}

	if cfg.PathDefs == "" {
		return nil, errors.NotValidf("Missing path for definition data directory")
	}

	if cfg.PathLive == "" {
		return nil, errors.NotValidf("Missing path for live data directory")
	}

	err = w.LoadDefinitions(cfg.PathDefs)
	if err != nil {
		return nil, errors.Annotate(err, "pathDefs error")
	}

	err = w.LoadRegions(cfg.PathLive)
	if err != nil {
		return nil, errors.Annotate(err, "pathLive error")
	}

	err = w.Check()
	if err != nil {
		return nil, errors.Annotate(err, "inconsistent world")
	}

	var notifier region.Notifier
	notifier, err = NewEventStoreClient(ctx)
	if err != nil {
		return nil, errors.Annotate(err, "event store client")
	}
	w.SetNotifier(notifier)

	var mc region.MapClient
	mc, err = region.NewDirectMapClient(ctx)
	if err != nil {
		return nil, errors.Annotate(err, "map client")
	}
	w.SetMapClient(mc)

	return &regionBackend{w: w, cfg: cfg}, nil
}

// Register pugs the internal gRPC routes into the given server
func (app *regionBackend) Register(grpcSrv *grpc.Server) error {
	proto.RegisterAdminServer(grpcSrv, &adminApp{app: app})
	proto.RegisterGameMasterServer(grpcSrv, &gmApp{app: app})
	proto.RegisterTemplatesServer(grpcSrv, &templatesApp{app: app})
	proto.RegisterPlayerServer(grpcSrv, &playerApp{app: app})
	proto.RegisterCityServer(grpcSrv, &cityApp{app: app})
	proto.RegisterArmyServer(grpcSrv, &armyApp{app: app})
	proto.RegisterDefinitionsServer(grpcSrv, &defsApp{app: app})
	proto.RegisterPublicServer(grpcSrv, &publicApp{app: app})
	grpc_prometheus.Register(grpcSrv)

	utils.Logger.Info().
		Str("defs", app.cfg.PathDefs).
		Str("live", app.cfg.PathLive).
		Msg("starting")

	return nil
}

// Make the RegionApp monitorable by the server stub
func (app *regionBackend) Check(ctx context.Context) grpc_health_v1.HealthCheckResponse_ServingStatus {
	return grpc_health_v1.HealthCheckResponse_SERVING
}

func (app *regionBackend) _worldLock(mode rune, action func() error) error {
	switch mode {
	case 'r':
		app.w.RLock()
		defer app.w.RUnlock()
	case 'w':
		app.w.WLock()
		defer app.w.WUnlock()
	default:
		return status.Error(codes.Internal, "invalid lock type")
	}
	return action()
}

func (app *regionBackend) _regLock(mode rune, regID string, action func(*region.Region) error) error {
	switch mode {
	case 'r':
		app.w.RLock()
		defer app.w.RUnlock()
	case 'w':
		app.w.WLock()
		defer app.w.WUnlock()
	default:
		return status.Error(codes.Internal, "invalid lock type")
	}
	r := app.w.Regions.Get(regID)
	if r == nil {
		return status.Error(codes.NotFound, "no such region")
	}
	return action(r)
}

func (app *regionBackend) cityLock(mode rune, req *proto.CityId, action func(*region.Region, *region.City) error) error {
	return app._regLock('r', req.Region, func(r *region.Region) error {
		switch mode {
		case 'r':
			// TODO(jfs) NYI
		case 'w':
			// TODO(jfs) NYI
		default:
			return status.Error(codes.Internal, "invalid lock type")
		}

		c := r.CityGet(req.City)
		if c == nil {
			return status.Error(codes.NotFound, "no such city")
		}
		// TODO FIXME(jfs): check the permissions
		return action(r, c)
	})
}

func (app *regionBackend) armyLock(mode rune, req *proto.ArmyId, action func(*region.Region, *region.City, *region.Army) error) error {
	cID := proto.CityId{Region: req.Region, City: req.City}
	return app.cityLock(mode, &cID, func(r *region.Region, c *region.City) error {
		a := c.Armies.Get(req.Army)
		if a == nil {
			return status.Error(codes.NotFound, "no such army")
		}
		return action(r, c, a)
	})
}

func s(u uint64) string { return strconv.FormatUint(u, 10) }
