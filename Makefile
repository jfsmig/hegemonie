# Copyright (c) 2018-2021 Contributors as noted in the AUTHORS file
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at http://mozilla.org/MPL/2.0/.

GO=go
PROTOC=protoc
COV_OUT=coverage.txt

AUTO=
# gen-set
AUTO+= pkg/gen-set/genset_auto_test.go
AUTO+= pkg/map/graph/map_auto.go
AUTO+= pkg/region/model/world_auto.go
# grpc
AUTO+= pkg/map/proto/map_grpc.pb.go
AUTO+= pkg/map/proto/map.pb.go
AUTO+= pkg/event/proto/event_grpc.pb.go
AUTO+= pkg/event/proto/event.pb.go
AUTO+= pkg/region/proto/region_grpc.pb.go
AUTO+= pkg/region/proto/region.pb.go
AUTO+= pkg/healthcheck/healthcheck_grpc.pb.go
AUTO+= pkg/healthcheck/healthcheck.pb.go

default: hege

all: prepare hege

gen-set: pkg/gen-set/gen-set.go
	cd pkg/gen-set ; $(GO) install

hege: gen-set
	cd pkg/hege ; $(GO) install 

.PHONY: all default prepare clean test benchmark fmt hege gen-set

prepare:
	-rm $(AUTO)
	set -ex ; go list ./... | while read D ; do $(GO) generate $$D ; done

fmt:
	./bin/hege-ci-fmt

clean:
	-rm $(AUTO)

test: all
	-rm $(COV_OUT)
	./bin/hege-ci-test $(COV_OUT)

benchmark: all
	./bin/hege-ci-benchmark

