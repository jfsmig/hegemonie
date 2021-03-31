// Copyright (c) 2018-2021 Contributors as noted in the AUTHORS file
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package region

import (
	"github.com/hegemonie-rpg/engine/pkg/utils"
	"github.com/rs/zerolog"
)

type Notifier interface {
	// Prepare a notification context to inform :to: of the movement of the army :a:.
	Army(log *City) EventArmy
	// Prepare a notification context to inform :to: of the evolution of the knowledge of someone.
	Knowledge(log *City) EventKnowledge
	// Prepare a notification context to inform :to: of someone hiring troops
	Units(log *City) EventUnits
}

type EventArmy interface {
	Item(a *Army) EventArmy
	// Notify the movement
	Move(src, dst uint64) EventArmy
	// Notify the movement is not possible
	NoRoute(src, dst uint64) EventArmy
	Send()
}

// EventKnowledge defines the builder of an event that informs about the progression of a knowledge
type EventKnowledge interface {
	// Item collects the City 'who' which progresses in the learning of the Knowledge 'k'
	Item(who *City, k *KnowledgeType) EventKnowledge

	// Step collects about the progression in the earning of the Knowledge referenced by Item
	Step(current, max uint64) EventKnowledge

	// Send emits the event to the collector.
	Send()
}

type EventUnits interface {
	Item(who *City, u *UnitType) EventUnits
	Step(current, max uint64) EventUnits
	Send()
}

type noEvt struct{}
type noEvtArmy struct{}
type noEvtKnowledge struct{}
type noEvtUnits struct{}

func (n *noEvt) Army(to *City) EventArmy           { return &noEvtArmy{} }
func (n *noEvt) Knowledge(to *City) EventKnowledge { return &noEvtKnowledge{} }
func (n *noEvt) Units(to *City) EventUnits         { return &noEvtUnits{} }

func (ctx *noEvtArmy) Item(a *Army) EventArmy            { return ctx }
func (ctx *noEvtArmy) Move(src, dst uint64) EventArmy    { return ctx }
func (ctx *noEvtArmy) NoRoute(src, dst uint64) EventArmy { return ctx }
func (ctx *noEvtArmy) Send()                             {}

func (ctx *noEvtKnowledge) Item(c *City, k *KnowledgeType) EventKnowledge { return ctx }
func (ctx *noEvtKnowledge) Step(current, max uint64) EventKnowledge       { return ctx }
func (ctx *noEvtKnowledge) Send()                                         {}

func (ctx *noEvtUnits) Item(c *City, k *UnitType) EventUnits { return ctx }
func (ctx *noEvtUnits) Step(current, max uint64) EventUnits  { return ctx }
func (ctx *noEvtUnits) Send()                                {}

func LogEvent(n Notifier) Notifier {
	return &eventLogger{sub: n}
}

type eventLogger struct {
	sub Notifier
}

type logEvtArmy struct {
	log *zerolog.Event
	sub EventArmy
}

type logEvtKnowledge struct {
	log *zerolog.Event
	sub EventKnowledge
}

type logEvtUnits struct {
	log *zerolog.Event
	sub EventUnits
}

func logger(to *City) *zerolog.Event {
	return utils.Logger.Info().
		Str("logChar", to.Owner).
		Uint64("logCity", to.ID)
}

func (n *eventLogger) Army(to *City) EventArmy {
	return &logEvtArmy{log: logger(to), sub: n.sub.Army(to)}
}

func (n *eventLogger) Knowledge(to *City) EventKnowledge {
	return &logEvtKnowledge{log: logger(to), sub: n.sub.Knowledge(to)}
}

func (n *eventLogger) Units(to *City) EventUnits {
	return &logEvtUnits{log: logger(to), sub: n.sub.Units(to)}
}

func (evt *logEvtArmy) Item(a *Army) EventArmy {
	evt.sub.Item(a)
	evt.log.Str("army", a.ID)
	return evt
}

func (evt *logEvtArmy) Move(src, dst uint64) EventArmy {
	evt.sub.Move(src, dst)
	evt.log.Uint64("src", src).Uint64("dst", dst)
	return evt
}

func (evt *logEvtArmy) NoRoute(src, dst uint64) EventArmy {
	evt.sub.NoRoute(src, dst)
	evt.log.Uint64("src", src).Uint64("dst", dst)
	return evt
}

func (evt *logEvtArmy) Send() {
	evt.sub.Send()
	evt.log.Send()
}

func (evt *logEvtKnowledge) Item(c *City, k *KnowledgeType) EventKnowledge {
	evt.sub.Item(c, k)
	evt.log.Uint64("city", c.ID).Uint64("id", k.ID)
	return evt
}

func (evt *logEvtKnowledge) Step(current, max uint64) EventKnowledge {
	evt.sub.Step(current, max)
	evt.log.Uint64("cur", current).Uint64("max", max)
	return evt
}

func (evt *logEvtKnowledge) Send() {
	evt.sub.Send()
	evt.log.Send()
}

func (evt *logEvtUnits) Item(c *City, u *UnitType) EventUnits {
	evt.sub.Item(c, u)
	evt.log.Uint64("city", c.ID).Uint64("id", u.ID)
	return evt
}

func (evt *logEvtUnits) Step(current, max uint64) EventUnits {
	evt.sub.Step(current, max)
	evt.log.Uint64("cur", current).Uint64("max", max)
	return evt
}

func (evt *logEvtUnits) Send() {
	evt.sub.Send()
	evt.log.Send()
}
