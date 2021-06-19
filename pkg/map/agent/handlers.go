// Copyright (c) 2018-2021 Contributors as noted in the AUTHORS file
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package mapagent

import (
	"context"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/hegemonie-rpg/engine/pkg/map/graph"
	"github.com/hegemonie-rpg/engine/pkg/map/proto"
	"github.com/hegemonie-rpg/engine/pkg/utils"
	"github.com/juju/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type srvMap struct {
	proto.UnimplementedMapServer

	config utils.MapServiceConfig
	maps   mapgraph.SetOfMaps
	rw     sync.RWMutex
}

type AppGenerator struct{}

// Application implements the expectations of the application backend
func (gen *AppGenerator) Application(ctx context.Context, config utils.MainConfig) (utils.RegisterableMonitorable, error) {
	cfg := config.Server.MapConfig
	app := &srvMap{config: cfg, maps: make(mapgraph.SetOfMaps, 0)}
	if err := app.LoadDirectory(cfg.PathRepository); err != nil {
		return nil, errors.Trace(err)
	}
	return app, nil
}

// Register starts an Map API service bond to Endpoint
// ctx is used for a clean stop of the service.
func (s *srvMap) Register(grpcSrv *grpc.Server) error {
	proto.RegisterMapServer(grpcSrv, s)
	grpc_prometheus.Register(grpcSrv)
	utils.Logger.Info().
		Int("maps", s.maps.Len()).
		Msg("Ready")
	for _, m := range s.maps {
		utils.Logger.Debug().
			Str("name", m.ID).
			Int("sites", m.Cells.Len()).
			Int("roads", m.Roads.Len()).
			Msg("map>")
	}
	return nil
}

func (s *srvMap) Check(ctx context.Context) grpc_health_v1.HealthCheckResponse_ServingStatus {
	return grpc_health_v1.HealthCheckResponse_SERVING
}

// Vertices streams Vertice objects, sorted by ID.
func (s *srvMap) Vertices(req *proto.ListVerticesReq, stream proto.Map_VerticesServer) error {
	return s._get('r', req.MapName, func(m *mapgraph.Map) error {
		next := req.Marker
		for {
			vertices := m.Cells.Slice(next, 100)
			if len(vertices) <= 0 {
				return nil
			}
			for _, x := range vertices {
				err := stream.Send(&proto.Vertex{Id: x.ID, X: x.X, Y: x.Y})
				if err != nil {
					return errors.Trace(err)
				}
				next = x.ID
			}
		}
	})
}

// Edges streams Edge objects, sorted by source then by destination.
func (s *srvMap) Edges(req *proto.ListEdgesReq, stream proto.Map_EdgesServer) error {
	return s._get('r', req.MapName, func(m *mapgraph.Map) error {
		src, dst := req.MarkerSrc, req.MarkerDst
		for {
			edges := m.Roads.Slice(src, dst, 100)
			if len(edges) <= 0 {
				return nil
			}
			for _, x := range edges {
				err := stream.Send(&proto.Edge{Src: x.S, Dst: x.D})
				if err != nil {
					return errors.Trace(err)
				}
				src, dst = x.S, x.D
			}
		}
	})
}

// GetPath streams the Vertice elements of the path from the source to the destination.
func (s *srvMap) GetPath(req *proto.PathRequest, stream proto.Map_GetPathServer) error {
	return s._get('r', req.MapName, func(m *mapgraph.Map) error {
		src := req.Src
		for {
			next, err := m.PathNextStep(src, req.Dst)
			if err != nil {
				return errors.Trace(err)
			}
			err = stream.Send(&proto.PathElement{Id: src})
			if err != nil {
				return errors.Trace(err)
			}
			if next == req.Dst {
				return nil
			}
			src = next
		}
	})
}

// Cities streams City <ID,name> pair objects
func (s *srvMap) Cities(req *proto.ListCitiesReq, stream proto.Map_CitiesServer) error {
	return s._get('r', req.MapName, func(m *mapgraph.Map) error {
		next := req.Marker
		for {
			cities := m.Cells.Slice(next, 100)
			if len(cities) <= 0 {
				return nil
			}
			for _, v := range cities {
				if v.City != "" {
					err := stream.Send(&proto.CityLocation{Id: v.ID, Name: v.City})
					if err == io.EOF {
						return nil
					}
					if err != nil {
						return errors.Trace(err)
					}
				}
				next = v.ID
			}
		}
	})
}

// Maps streams the name of the maps registered in the current service
func (s *srvMap) Maps(req *proto.ListMapsReq, stream proto.Map_MapsServer) error {
	// Extract stats on a slice of the array of cities, under the umbrella of 	a read-lock
	slice := func(marker string) []proto.MapName {
		out := make([]proto.MapName, 0)
		s.rw.RLock()
		defer s.rw.RUnlock()
		for _, m := range s.maps.Slice(marker, 100) {
			out = append(out, proto.MapName{
				Name:          m.ID,
				CountEdges:    uint32(len(m.Roads)),
				CountVertices: uint32(len(m.Cells)),
				CountCities: func() (total uint32) {
					for _, c := range m.Cells {
						if c.City != "" {
							total++
						}
					}
					return total
				}(),
			})
		}
		return out
	}

	next := req.Marker
	for {
		names := slice(next)
		if len(names) <= 0 {
			return nil
		}
		for _, v := range names {
			err := stream.Send(&v)
			if err == io.EOF {
				return nil
			}
			if err != nil {
				return errors.Trace(err)
			}
			next = v.Name
		}
	}
}

// LoadDirectory loads all the maps stored as files, containing JSON objects desribing maps.
// Only the filenames with a .final.json suffix are considered.
func (s *srvMap) LoadDirectory(path string) error {
	return filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return errors.Trace(err)
		}

		// Only accept non-hidden JSON files
		_, fn := filepath.Split(path)
		if info.IsDir() || info.Size() <= 0 {
			return nil
		}
		if len(fn) < 2 || fn[0] == '.' {
			return nil
		}
		if !strings.HasSuffix(fn, ".final.json") {
			return nil
		}

		var f *os.File
		f, err = os.Open(path)
		if err != nil {
			return errors.NewNotValid(err, "fs error")
		}
		defer f.Close()

		m := mapgraph.NewMap()
		if err = m.Load(f); err != nil {
			return errors.NewNotValid(err, "format error")
		}
		s.maps.Add(m)
		return nil
	})
}

func (s *srvMap) _lock(mode rune, action func() error) error {
	switch mode {
	case 'r':
		s.rw.RLock()
		defer s.rw.RUnlock()
	case 'w':
		s.rw.Lock()
		defer s.rw.Unlock()
	default:
		panic("unexpected lock mode")
	}
	return action()
}

func (s *srvMap) _get(mode rune, name string, action func(*mapgraph.Map) error) error {
	return s._lock(mode, func() error {
		m := s.maps.Get(name)
		if m == nil {
			return errors.NotFoundf("no such map")
		}
		return action(m)
	})
}
