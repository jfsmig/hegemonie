// Copyright (c) 2018-2021 Contributors as noted in the AUTHORS file
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package regagent

import (
	"context"
	"github.com/go-openapi/errors"
	"github.com/hegemonie-rpg/engine/pkg/region/proto"
)

type playerApp struct {
	proto.UnimplementedPlayerServer

	app *regionBackend
}

func (app *playerApp) LifecycleConfigure(ctx context.Context, req *proto.LifecycleConfigureReq) (*proto.None, error) {
	return none, errors.NotImplemented("NYI")
}

func (app *playerApp) LifecycleAcquire(ctx context.Context, req *proto.CityId) (*proto.None, error) {
	return none, errors.NotImplemented("NYI")
}

func (app *playerApp) LifecycleLeave(ctx context.Context, req *proto.CityId) (*proto.None, error) {
	return none, errors.NotImplemented("NYI")
}

func (app *playerApp) LifecycleAuto(ctx context.Context, req *proto.CityId) (*proto.None, error) {
	return none, errors.NotImplemented("NYI")
}
