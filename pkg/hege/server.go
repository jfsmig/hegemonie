// Copyright (c) 2018-2021 Contributors as noted in the AUTHORS file
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	evtagent "github.com/hegemonie-rpg/engine/pkg/event/agent"
	mapagent "github.com/hegemonie-rpg/engine/pkg/map/agent"
	regagent "github.com/hegemonie-rpg/engine/pkg/region/agent"
	"github.com/hegemonie-rpg/engine/pkg/utils"
	"github.com/juju/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health/grpc_health_v1"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

type srvCommons struct {
	pathConfig string

	EndpointService string `yaml:"bind" json:"bind"`
	EndpointMonitor string `yaml:"monitor" json:"monitor"`
	ServiceType     string `yaml:"type" json:"type"`
	PathKey         string `yaml:"key" json:"key"`
	PathCrt         string `yaml:"cert" json:"cert"`

	MapConfig mapagent.Config `yaml:"map" json:"map"`
	EvtConfig evtagent.Config `yaml:"evt" json:"evt"`
	RegConfig regagent.Config `yaml:"reg" json:"reg"`

	apps []utils.RegisterableMonitorable
}

// appGenerator is implemented by all the "Config" structs in the accepted applications.
// U know what's the funiest? Those "Config" structs are not even aware of it.
type appGenerator interface {
	Application(ctx context.Context) (utils.RegisterableMonitorable, error)
}

const (
	defaultKeyPath = "/etc/hegemonie/pki/<SRVTYPE>.key"
	defaultCrtPath = "/etc/hegemonie/pki/<SRVTYPE>.crt"
)

func servers() *cobra.Command {
	srv := srvCommons{
		EndpointService: fmt.Sprintf("0.0.0.0:%v", utils.DefaultPortCommon),
		EndpointMonitor: fmt.Sprintf("0.0.0.0:%v", utils.DefaultPortMonitoring),
		ServiceType:     "type-NOT-SET",
		PathCrt:         defaultCrtPath,
		PathKey:         defaultKeyPath,
	}
	cmd := &cobra.Command{
		Use:              "server",
		Short:            "Run Hegemonie services",
		Args:             cobra.MinimumNArgs(1),
		RunE:             nonLeaf,
		TraverseChildren: true,
	}
	cmd.PersistentFlags().StringVarP(
		&srv.pathConfig, "config", "f", "/etc/hegemonie/server.yml",
		"Path to the configuration file")
	cmd.AddCommand(srv.maps(), srv.events(), srv.regions(), srv.bundle())
	return cmd
}

func (srv *srvCommons) events() *cobra.Command {
	return &cobra.Command{
		Use:               "events",
		Short:             "Event Log Service",
		Example:           "hege server events -f /path/to/config",
		Args:              cobra.NoArgs,
		PersistentPreRunE: srv.wrapPreRun("events"),
		RunE:              func(cmd *cobra.Command, args []string) error { return srv.runServer(srv.EvtConfig) },
	}
}

func (srv *srvCommons) maps() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "maps",
		Short:             "Map Service",
		Example:           "hege server maps -f /path/to/config",
		Args:              cobra.NoArgs,
		PersistentPreRunE: srv.wrapPreRun("maps"),
		RunE:              func(cmd *cobra.Command, args []string) error { return srv.runServer(srv.MapConfig) },
	}
	return cmd
}

func (srv *srvCommons) regions() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "regions",
		Short:             "Region Service",
		Example:           "hege server regions -f /path/to/config",
		Args:              cobra.NoArgs,
		PersistentPreRunE: srv.wrapPreRun("regions"),
		RunE:              func(cmd *cobra.Command, args []string) error { return srv.runServer(srv.RegConfig) },
	}
	return cmd
}

func (srv *srvCommons) bundle() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "all",
		Short:             "All services at once",
		Example:           "hege server all -f /path/to/config",
		Args:              cobra.NoArgs,
		PersistentPreRunE: srv.wrapPreRun("all"),
		RunE: func(cmd *cobra.Command, args []string) error {
			return srv.runServer(srv.RegConfig, srv.EvtConfig, srv.MapConfig)
		},
	}
	return cmd
}

func (srv *srvCommons) replaceTag(ps *string) {
	*ps = strings.Replace(*ps, "<SRVTYPE>", srv.ServiceType, 1)
}

// ServerTLS automates the creation of a grpc.Server over a TLS connection
// with the proper interceptors.
func (srv *srvCommons) serverTLS() (*grpc.Server, error) {
	if len(srv.PathCrt) <= 0 {
		return nil, errors.NotValidf("invalid TLS/x509 certificate path [%s]", srv.PathCrt)
	}
	if len(srv.PathKey) <= 0 {
		return nil, errors.NotValidf("invalid TLS/x509 key path [%s]", srv.PathKey)
	}
	var certBytes, keyBytes []byte
	var err error

	utils.Logger.Info().Str("key", srv.PathKey).Str("crt", srv.PathCrt).Msg("TLS config")

	if certBytes, err = ioutil.ReadFile(srv.PathCrt); err != nil {
		return nil, errors.Annotate(err, "certificate file error")
	}
	if keyBytes, err = ioutil.ReadFile(srv.PathKey); err != nil {
		return nil, errors.Annotate(err, "key file error")
	}

	certPool := x509.NewCertPool()
	ok := certPool.AppendCertsFromPEM(certBytes)
	if !ok {
		return nil, errors.New("invalid certificates")
	}

	cert, err := tls.X509KeyPair(certBytes, keyBytes)
	if err != nil {
		return nil, errors.Annotate(err, "x509 key pair error")
	}

	return grpc.NewServer(
		grpc.Creds(credentials.NewServerTLSFromCert(&cert)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_prometheus.UnaryServerInterceptor,
			utils.NewUnaryServerInterceptorZerolog())),
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_prometheus.StreamServerInterceptor,
			utils.NewStreamServerInterceptorZerolog()))), nil
}

func (srv *srvCommons) runServer(registrators ...appGenerator) error {
	var listenerSrv, listenerMon net.Listener
	var grpcSrv *grpc.Server
	var prometheusExporter *http.Server
	var err error

	ctx, cancel := context.WithCancel(context.Background())

	grpcSrv, err = srv.serverTLS()
	if err != nil {
		return errors.Annotate(err, "TLS server error")
	}

	for _, reg := range registrators {
		app, err := reg.Application(ctx)
		if err != nil {
			return errors.Annotate(err, "App startup error")
		}
		err = app.Register(grpcSrv)
		if err != nil {
			return errors.Annotate(err, "App config error")
		}
	}

	grpc_health_v1.RegisterHealthServer(grpcSrv, srv)

	listenerSrv, err = net.Listen("tcp", srv.EndpointService)
	if err != nil {
		return errors.NewNotValid(err, "listen error")
	}

	if srv.EndpointMonitor != "" {
		listenerMon, err = net.Listen("tcp", srv.EndpointMonitor)
		if err != nil {
			cancel()
			return errors.NewNotValid(err, "listen error")
		}

		prometheusExporter = &http.Server{Handler: promhttp.Handler()}
	}

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGABRT)
	defer signal.Stop(stopChan)

	var barrier sync.WaitGroup
	runner := func(wg *sync.WaitGroup, tag string, cb func() error) {
		defer wg.Done()
		if err := cb(); err != nil && err != http.ErrServerClosed {
			utils.Logger.Error().Err(err).Str("itf", tag).Msg("failed")
		} else {
			utils.Logger.Info().Str("itf", tag).Msg("exiting")
		}
		cancel()
	}

	barrier.Add(1)
	go runner(&barrier, "service", func() error { return grpcSrv.Serve(listenerSrv) })

	if prometheusExporter != nil {
		barrier.Add(1)
		go runner(&barrier, "monitor", func() error { return prometheusExporter.Serve(listenerMon) })
	}

	select {
	case <-stopChan:
		break
	case <-ctx.Done():
		break
	}
	cancel()

	grpcSrv.GracefulStop()

	if prometheusExporter != nil {
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()
		if err := prometheusExporter.Shutdown(shutdownCtx); err != nil {
			utils.Logger.Warn().Err(err).Msg("shutdown error")
		}
	}

	barrier.Wait()
	return nil
}

func (srv *srvCommons) wrapPreRun(srvtype string) func(*cobra.Command, []string) error {
	return func(*cobra.Command, []string) (err error) {
		srv.ServiceType = srvtype
		utils.OverrideLogID("hege," + srv.ServiceType)
		utils.ApplyLogModifiers()

		fin, err := os.Open(srv.pathConfig)
		if err != nil {
			return errors.Annotate(err, "configuration path error")
		}
		err = yaml.NewDecoder(fin).Decode(srv)
		if err != nil {
			return errors.Annotate(err, "invalid configuration")
		}
		utils.Logger.Info().RawJSON("srv", utils.JSON2Buf(srv)).Msg("configuration")

		utils.OverrideLogID("hege," + srv.ServiceType)
		utils.ApplyLogModifiers()
		srv.replaceTag(&srv.PathKey)
		srv.replaceTag(&srv.PathCrt)
		return nil
	}
}

func (srv *srvCommons) getCheckResponse(ctx context.Context, service string) *grpc_health_v1.HealthCheckResponse {
	// FIXME(jfs): check the service ID
	status := grpc_health_v1.HealthCheckResponse_UNKNOWN
	for _, app := range srv.apps {
		if s := app.Check(ctx); s > status {
			status = s
		}
	}
	return &grpc_health_v1.HealthCheckResponse{Status: status}
}

// Check implements the one-shot healthcheck of the gRPC service
func (srv *srvCommons) Check(ctx context.Context, req *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	return srv.getCheckResponse(ctx, req.Service), nil
}

// Watch implements the long polling healthcheck of the gRPC service
func (srv *srvCommons) Watch(req *grpc_health_v1.HealthCheckRequest, s grpc_health_v1.Health_WatchServer) error {
	for {
		err := s.Send(srv.getCheckResponse(s.Context(), req.Service))
		if err != nil {
			return errors.Trace(err)
		}
	}
}
