// Copyright (c) 2018-2021 Contributors as noted in the AUTHORS file
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package regagent

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/google/uuid"
	hegemonie_rpevent_proto "github.com/jfsmig/hegemonie/pkg/event/proto"
	"github.com/jfsmig/hegemonie/pkg/region/model"
	"github.com/jfsmig/hegemonie/pkg/utils"
	"github.com/juju/errors"
	"google.golang.org/grpc"
)

type EventStore struct {
	cnx *grpc.ClientConn
}

// NewEventStoreClient instantiates an event notifier that targets the events service
// returned by utils.DefaultDiscovery
func NewEventStoreClient(ctx context.Context) (region.Notifier, error) {
	eventEndpoint, err := utils.DefaultDiscovery.Event()
	if err != nil {
		return nil, errors.Annotate(err, "discovery error")
	}

	var cnxEvent *grpc.ClientConn
	cnxEvent, err = utils.Dial(ctx, eventEndpoint)
	if err != nil {
		return nil, errors.Annotate(err, "dial error")
	}

	return &EventStore{cnx: cnxEvent}, nil
}

type EventArmy struct {
	store  *EventStore
	charID string

	SourceCityID uint64 `json:"SourceCityId"`
	SourceCity   string `json:"SourceCity"`

	ArmyID   string `json:"ArmyId"`
	ArmyName string `json:"Army"`

	ArmyCityID   uint64 `json:"ArmyCityId"`
	ArmyCityName string `json:"ArmyCity"`

	Src uint64 `json:"Src"`
	Dst uint64 `json:"Dst"`

	Action string `json:"action"`
}

type EventKnowledge struct {
	store *EventStore
}

type EventUnits struct {
	store *EventStore
}

func (es *EventStore) Army(log *region.City) region.EventArmy {
	return &EventArmy{
		store:        es,
		charID:       log.Owner,
		SourceCity:   log.Name,
		SourceCityID: log.ID,
	}
}

func (es *EventStore) Knowledge(log *region.City) region.EventKnowledge {
	return &EventKnowledge{store: es}
}

func (es *EventStore) Units(log *region.City) region.EventUnits {
	return &EventUnits{store: es}
}

func (evt *EventArmy) Item(a *region.Army) region.EventArmy {
	evt.ArmyID = a.ID
	evt.ArmyName = a.Name
	evt.ArmyCityID = a.City.ID
	evt.ArmyCityName = a.City.Name
	return evt
}

func (evt *EventArmy) Move(src, dst uint64) region.EventArmy {
	evt.Src, evt.Dst = src, dst
	evt.Action = "Move"
	return evt
}

func (evt *EventArmy) NoRoute(src, dst uint64) region.EventArmy {
	evt.Src, evt.Dst = src, dst
	evt.Action = "NoRoute"
	return evt
}

func (evt *EventArmy) Send() {
	var buffer bytes.Buffer
	enc := json.NewEncoder(&buffer)
	enc.SetIndent("", "")
	enc.Encode(evt)

	client := hegemonie_rpevent_proto.NewProducerClient(evt.store.cnx)
	client.Push1(context.Background(), &hegemonie_rpevent_proto.Push1Req{
		CharId:  evt.charID,
		EvtId:   uuid.New().String(),
		Payload: buffer.Bytes(),
	})
}

func (evt *EventKnowledge) Item(c *region.City, kt *region.KnowledgeType) region.EventKnowledge {
	// TODO FIXME
	return evt
}

func (evt *EventKnowledge) Step(current, max uint64) region.EventKnowledge {
	// TODO FIXME
	return evt
}

func (evt *EventKnowledge) Send() {
	// TODO FIXME
}

func (evt *EventUnits) Item(c *region.City, ut *region.UnitType) region.EventUnits {
	// TODO FIXME
	return evt
}

func (evt *EventUnits) Step(current, max uint64) region.EventUnits {
	// TODO FIXME
	return evt
}

func (evt *EventUnits) Send() {
	// TODO FIXME
}
