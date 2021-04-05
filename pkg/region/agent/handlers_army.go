// Copyright (c) 2018-2021 Contributors as noted in the AUTHORS file
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package regagent

import (
	"context"
	"github.com/hegemonie-rpg/engine/pkg/region/model"
	"github.com/hegemonie-rpg/engine/pkg/region/proto"
)

type armyApp struct {
	proto.UnimplementedArmyServer

	app *regionBackend
}

func (app *armyApp) Show(ctx context.Context, req *proto.ArmyId) (reply *proto.ArmyView, err error) {
	err = app.app.armyLock('r', req, func(r *region.Region, _ *region.City, army *region.Army) error {
		reply = showArmy(r, army)
		return nil
	})
	return reply, err
}

func (app *armyApp) Flea(ctx context.Context, req *proto.ArmyId) (*proto.None, error) {
	return none, app.app.armyLock('w', req, func(r *region.Region, _ *region.City, a *region.Army) error {
		return a.Flea(r)
	})
}

func (app *armyApp) Flip(ctx context.Context, req *proto.ArmyId) (*proto.None, error) {
	return none, app.app.armyLock('w', req, func(r *region.Region, _ *region.City, a *region.Army) error {
		return a.Flip(r)
	})
}

func (app *armyApp) Move(ctx context.Context, req *proto.ArmyMoveReq) (*proto.None, error) {
	return none, app.app.armyLock('w', req.Id, func(r *region.Region, _ *region.City, a *region.Army) error {
		return a.DeferMove(r, req.Target, region.ActionArgMove{})
	})
}

func (app *armyApp) Attack(ctx context.Context, req *proto.ArmyAssaultReq) (*proto.None, error) {
	return none, app.app.armyLock('w', req.Id, func(r *region.Region, _ *region.City, a *region.Army) error {
		return a.DeferAttack(r, req.Target, region.ActionArgAssault{})
	})
}

func (app *armyApp) Wait(ctx context.Context, req *proto.ArmyTarget) (*proto.None, error) {
	return none, app.app.armyLock('w', req.Id, func(r *region.Region, _ *region.City, a *region.Army) error {
		return a.DeferWait(r, req.Target)
	})
}

func (app *armyApp) Defend(ctx context.Context, req *proto.ArmyTarget) (*proto.None, error) {
	return none, app.app.armyLock('w', req.Id, func(r *region.Region, _ *region.City, a *region.Army) error {
		return a.DeferDefend(r, req.Target)
	})
}

func (app *armyApp) Disband(ctx context.Context, req *proto.ArmyTarget) (*proto.None, error) {
	return none, app.app.armyLock('w', req.Id, func(r *region.Region, _ *region.City, a *region.Army) error {
		return a.DeferDisband(r, req.Target)
	})
}

func (app *armyApp) Cancel(ctx context.Context, req *proto.ArmyId) (*proto.None, error) {
	return none, app.app.armyLock('w', req, func(r *region.Region, _ *region.City, a *region.Army) error {
		return a.Cancel(r)
	})
}
