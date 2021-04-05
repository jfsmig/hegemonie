// Copyright (c) 2018-2021 Contributors as noted in the AUTHORS file
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package regagent

import (
	"context"
	"github.com/hegemonie-rpg/engine/pkg/region/model"
	"github.com/hegemonie-rpg/engine/pkg/region/proto"
	"github.com/juju/errors"
	"io"
)

type gmApp struct {
	proto.UnimplementedGameMasterServer

	app *regionBackend
}

func (app *gmApp) Produce(ctx context.Context, req *proto.RegionId) (*proto.None, error) {
	return none, app.app._regLock('w', req.Region, func(r *region.Region) error {
		r.Produce(ctx)
		return nil
	})
}

func (app *gmApp) Move(ctx context.Context, req *proto.RegionId) (*proto.None, error) {
	return none, app.app._regLock('w', req.Region, func(r *region.Region) error {
		r.Move(ctx)
		return nil
	})
}

func (app *gmApp) GetStats(req *proto.RegionId, stream proto.GameMaster_GetStatsServer) error {
	return app.app._regLock('r', req.Region, func(r *region.Region) error {
		for _, c := range r.Cities {
			// FIXME(jfs): Calling Send() from a critical section is a bad idea
			err := stream.Send(showCityStatsRecord(r, c))
			if err == io.EOF {
				return nil
			}
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (app *gmApp) GetDetail(ctx context.Context, req *proto.CityId) (reply *proto.CityView, err error) {
	err = app.app.cityLock('r', req, func(r *region.Region, c *region.City) error {
		reply = showCity(r, c)
		return nil
	})
	return reply, err
}

func (app *gmApp) LifecycleConfigure(ctx context.Context, req *proto.LifecycleConfigureReq) (*proto.None, error) {
	cityID := proto.CityId{Region: req.Region, City: req.City}
	err := app.app.cityLock('w', &cityID, func(r *region.Region, c *region.City) error {
		doAndRestate := func(newState uint32) error {
			e := c.ApplyNamedModel(r, req.Model)
			if e != nil {
				return errors.Trace(e)
			}
			c.State = newState
			return nil
		}
		switch c.State {
		case region.CityStatePrivIdle:
			return doAndRestate(region.CityStatePrivConfigured)
		case region.CityStatePrivAssigned:
			return doAndRestate(region.CityStatePrivReady)
		case region.CityStatePrivConfigured:
			return doAndRestate(region.CityStatePrivConfigured)
		case region.CityStatePubAssigned:
			return doAndRestate(region.CityStatePubActive)
		default:
			return errors.NotValidf("city")
		}
	})
	return none, err
}

func (app *gmApp) LifecycleAssign(ctx context.Context, req *proto.LifecycleAssignReq) (*proto.None, error) {
	cityID := proto.CityId{Region: req.City.Region, City: req.City.City}
	err := app.app.cityLock('w', &cityID, func(r *region.Region, c *region.City) error {
		switch c.State {
		case region.CityStatePrivIdle:
			c.State = region.CityStatePrivAssigned
		case region.CityStatePrivConfigured:
			c.State = region.CityStatePrivReady
		case region.CityStatePrivHeadless:
			c.State = region.CityStatePrivAuto

		case region.CityStatePubIdle:
			c.State = region.CityStatePubAssigned
		case region.CityStatePubHeadless:
			c.State = region.CityStatePubAuto
		default:
			return errors.NotValidf("city")
		}
		c.Owner = req.User
		return nil
	})
	return none, err
}

func (app *gmApp) LifecycleResume(ctx context.Context, req *proto.LifecycleAbstractReq) (*proto.None, error) {
	err := app.app.cityLock('w', req.City, func(r *region.Region, c *region.City) error {
		switch c.State {
		case region.CityStatePrivReady,
			region.CityStatePrivActive:
			c.State = region.CityStatePrivActive
		case region.CityStatePrivSuspended:
			c.State = region.CityStatePrivAuto

		case region.CityStatePubSuspended,
			region.CityStatePubAuto:
			c.State = region.CityStatePubAuto
		case region.CityStatePubActive:
			c.State = region.CityStatePubActive
		default:
			return errors.NotValidf("city")
		}
		return nil
	})
	return none, err
}

func (app *gmApp) LifecycleDismiss(ctx context.Context, req *proto.LifecycleAbstractReq) (*proto.None, error) {
	err := app.app.cityLock('w', req.City, func(r *region.Region, c *region.City) error {
		switch c.State {
		case region.CityStatePrivAssigned:
			c.State = region.CityStatePrivIdle
		case region.CityStatePrivReady:
			c.State = region.CityStatePrivConfigured
		case region.CityStatePrivActive,
			region.CityStatePrivHeadless,
			region.CityStatePrivSuspended,
			region.CityStatePrivAuto:
			c.State = region.CityStatePrivHeadless

		case region.CityStatePubAssigned:
			c.State = region.CityStatePubIdle
		case region.CityStatePubActive,
			region.CityStatePubHeadless,
			region.CityStatePubSuspended,
			region.CityStatePubAuto:
			c.State = region.CityStatePubHeadless
		default:
			return errors.NotValidf("city")
		}
		c.Owner = ""
		return nil
	})
	return none, err
}

func (app *gmApp) LifecycleSuspend(ctx context.Context, req *proto.LifecycleAbstractReq) (*proto.None, error) {
	err := app.app.cityLock('w', req.City, func(r *region.Region, c *region.City) error {
		switch c.State {
		case region.CityStatePrivSuspended,
			region.CityStatePrivAuto,
			region.CityStatePrivActive:
			c.State = region.CityStatePrivSuspended

		case region.CityStatePubSuspended,
			region.CityStatePubAuto,
			region.CityStatePubActive:
			c.State = region.CityStatePubSuspended
		default:
			return errors.NotValidf("city")
		}
		return nil
	})
	return none, err
}

func (app *gmApp) LifecycleReset(ctx context.Context, req *proto.LifecycleAbstractReq) (*proto.None, error) {
	err := app.app.cityLock('w', req.City, func(r *region.Region, c *region.City) error {
		switch c.State {
		case region.CityStatePrivIdle,
			region.CityStatePrivConfigured,
			region.CityStatePrivSuspended:
			c.State = region.CityStatePrivIdle

		case region.CityStatePubIdle,
			region.CityStatePubSuspended,
			region.CityStatePubHeadless:
			c.State = region.CityStatePubIdle
		default:
			return errors.NotValidf("city")
		}
		return nil
	})
	return none, err
}
