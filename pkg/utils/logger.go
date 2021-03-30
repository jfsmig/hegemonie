// Copyright (c) 2018-2021 Contributors as noted in the AUTHORS file
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/juju/errors"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"os"
	"strings"
	"time"
)

var (
	// LoggerContext is the builder of a zerolog.Logger that is exposed to the application so that
	// options at the CLI might alter the formatting and the output of the logs.
	LoggerContext = zerolog.
			New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).
			With().Timestamp()

	// Logger is a zerolog logger, that can be safely used from any part of the application.
	// It gathers the format and the output.
	Logger = LoggerContext.Logger()
)

var (
	flagVerbose = 0
	flagQuiet   = false

	// TODO(jfs): implement a syslog output
	flagSyslog = false

	// ServiceID is used in log traces to identify the service emitting the trace.
	// The value can be safely altered before emitting the first trace.
	syslogID = "hege"
)

func OverrideLogID(ID string) {
	syslogID = ID
}

// AddLogFlagsToCommand add to cmd a set of persistent flags that will alter the logging behavior
// of the current process.
func AddLogFlagsToCommand(cmd *cobra.Command) {
	cmd.PersistentFlags().CountVarP(&flagVerbose, "verbose", "v", "Increase the verbosity level")
	cmd.PersistentFlags().BoolVarP(&flagQuiet, "quiet", "q", flagQuiet, "Shut the logs")
	cmd.PersistentFlags().BoolVarP(&flagQuiet, "syslog", "s", flagQuiet, "Log in syslog")
	cmd.PersistentFlags().StringVarP(&syslogID, "syslog-id", "i", syslogID, "Use that service ID in the logs")
}

func ApplyLogModifiers() {
	Logger = LoggerContext.Str("id", syslogID).Logger()
	if flagQuiet {
		zerolog.SetGlobalLevel(zerolog.Disabled)
	} else {
		switch flagVerbose {
		case 0:
			zerolog.SetGlobalLevel(zerolog.InfoLevel)
		case 1:
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
		case 2:
			zerolog.SetGlobalLevel(zerolog.TraceLevel)
		}
	}
}

type logEvt struct {
	z     *zerolog.Event
	start time.Time
}

func newEvent(method string) *logEvt {
	return &logEvt{z: Logger.Info().Str("uri", method), start: time.Now()}
}

func (evt *logEvt) send() { evt.z.Send() }

func (evt *logEvt) setResult(err error) *logEvt {
	evt.z = evt.z.TimeDiff("t", time.Now(), evt.start)
	if err != nil {
		evt.z.Int("rc", 500)
		evt.z.Err(err)
	} else {
		evt.z.Int("rc", 200)
	}
	return evt
}

func (evt *logEvt) patchWithRequest(ctx context.Context) *logEvt {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		auth := md.Get(":authority")
		if len(auth) > 0 {
			evt.z.Str("local", auth[0])
		}
		sessionID := md.Get("session-id")
		if len(sessionID) > 0 {
			evt.z.Str("session", sessionID[0])
		}
	}
	return evt
}

func (evt *logEvt) pathWithReply(ctx context.Context) *logEvt {
	if peer, ok := peer.FromContext(ctx); ok {
		addr := peer.Addr.String()
		if i := strings.LastIndex(addr, ":"); i > -1 {
			addr = addr[:i]
		}
		evt.z.Str("peer", addr)
	}
	return evt
}

func NewStreamServerInterceptorZerolog() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		evt := newEvent(info.FullMethod)
		err := handler(srv, ss)
		ctx := ss.Context()
		evt.setResult(err).patchWithRequest(ctx).pathWithReply(ctx).send()
		return errors.Trace(err)
	}
}

func NewUnaryServerInterceptorZerolog() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		evt := newEvent(info.FullMethod)
		resp, err := handler(ctx, req)
		evt.setResult(err).patchWithRequest(ctx).pathWithReply(ctx).send()
		return resp, err
	}
}

// DumpJSON encodes the argument in JSON and writes the output on os.Stdout
func DumpJSON(x interface{}) error {
	encoder := json.NewEncoder(os.Stdout)
	// encoder.SetIndent("", " ")
	return encoder.Encode(x)
}

// JSON2Str generates a string with the JSON representation of the given object.
// No indentation, everything on one line. Errors ignored.
func JSON2Str(x interface{}) string {
	sb := strings.Builder{}
	encoder := json.NewEncoder(&sb)
	_ = encoder.Encode(x)
	return sb.String()
}

// JSON2Buf generates the JSON-encoded form of the given object, stored in the returned byte array.
// No indentation, no CRLF
func JSON2Buf(obj interface{}) []byte {
	buf := bytes.Buffer{}
	json.NewEncoder(&buf).Encode(obj)
	return buf.Bytes()
}

// StatusJSON encodes a standard error structure and then forwards it to DumpJSON
func StatusJSON(code int, ID, msg string) error {
	type status struct {
		Msg  string
		Code int
		ID   string
	}
	return DumpJSON(status{Msg: "Create", Code: code, ID: ID})
}
