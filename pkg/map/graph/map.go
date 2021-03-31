// Copyright (c) 2018-2021 Contributors as noted in the AUTHORS file
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package mapgraph

import (
	"encoding/json"
	"github.com/juju/errors"
	"io"
	"sort"
	"strings"
)

// Edge is ... an edge of the transportation directed graph
type Edge struct {
	// Unique identifier of the source Cell
	S uint64 `json:"src"`

	// Unique identifier of the destination Cell
	D uint64 `json:"dst"`
}

// Vertex is a vertex in the transportation directed graph
type Vertex struct {
	// The unique identifier of the current cell.
	ID uint64 `json:"id"`

	// Location of the Cell on the map. Used for rendering
	X uint64 `json:"x"`
	Y uint64 `json:"y"`

	// Should the current location carry a city when the region starts,
	// and if yes, what should be the name of that city.
	City string `json:"city,omitempty"`
}

// Map is a directed graph destined to be used as a transport network,
// organised as an adjacency list.
type Map struct {
	// The unique name of the map
	ID    string        `json:"id"`
	Cells SetOfVertices `json:"sites"`
	Roads SetOfEdges    `json:"roads"`
	steps map[vector]uint64
}

//go:generate go run github.com/hegemonie-rpg/engine/pkg/gen-set ./map_auto.go mapgraph:SetOfVertices:*Vertex ID:uint64

//go:generate go run github.com/hegemonie-rpg/engine/pkg/gen-set ./map_auto.go mapgraph:SetOfEdges:*Edge S:uint64 D:uint64

//go:generate go run github.com/hegemonie-rpg/engine/pkg/gen-set ./map_auto.go mapgraph:SetOfMaps:*Map ID:string

// EmptyMap builds a graph with no edge and no vertex.
func EmptyMap() Map {
	return Map{
		ID:    "",
		Cells: make(SetOfVertices, 0),
		Roads: make(SetOfEdges, 0),
		steps: make(map[vector]uint64),
	}
}

// NewMap returns a pointer to an EmptyMap
func NewMap() *Map {
	m := EmptyMap()
	return &m
}

// Load consumes a mapclient.MapRaw from os.Stdin and populates the current
// Map accordingly.
func (m *Map) Load(in io.Reader) error {
	err := json.NewDecoder(in).Decode(m)
	if err != nil {
		return errors.NewNotValid(err, "invalid input")
	}

	m.canonize()
	m.rehash()
	return m.check()
}

// LoadJSON is for testing purpose.
func (m *Map) LoadJSON(in string) error {
	return m.Load(strings.NewReader(in))
}

// CellGet is a shortcut to get a node / vertex given its ID
func (m *Map) CellGet(id uint64) *Vertex {
	return m.Cells.Get(id)
}

// CellHas is a shortcut to check the presence of a node / vertex given its ID
func (m *Map) CellHas(id uint64) bool {
	return m.Cells.Has(id)
}

// RoadHas is a shortcut to check the presence of an edge given its source
// and destination.
func (m *Map) RoadHas(src, dst uint64) bool {
	return m.Roads.Has(src, dst)
}

// PathNextStep computes the next step on the path from the source (given by src)
// to the destination (given by dst).
func (m *Map) PathNextStep(src, dst uint64) (uint64, error) {
	if src == dst || src == 0 || dst == 0 {
		return 0, errors.BadRequestf("EINVAL")
	}

	next, ok := m.steps[vector{src, dst}]
	if ok {
		return next, nil
	}
	return 0, errors.New("no route")
}

// CellAdjacency returns the adjacency list of the cell with the given ID,
// in other words, an array of the ID of the target of the roads leaving
// the given cell.
func (m *Map) CellAdjacency(id uint64) []uint64 {
	adj := make([]uint64, 0)

	for i := m.Roads.First(id); i < len(m.Roads); i++ {
		r := m.Roads[i]
		if r.S != id {
			break
		}
		adj = append(adj, r.D)
	}

	return adj
}

func (m *Map) reset() *Map {
	*m = EmptyMap()
	return m
}

func (m *Map) canonize() {
	sort.Sort(&m.Cells)
	sort.Sort(&m.Roads)
}

// Validate the invariants of a Map, on the current map
func (m *Map) check() error {
	if m.ID == "" {
		return errors.BadRequestf("bad name")
	}
	if err := m.Cells.Check(); err != nil {
		return errors.Annotate(err, "error on cells")
	}
	if err := m.Roads.Check(); err != nil {
		return errors.Annotate(err, "error of roads")
	}

	if m.Cells.Len() > 1 {
		for _, s0 := range m.Cells {
			for _, s1 := range m.Cells {
				if _, ok := m.steps[vector{s0.ID, s1.ID}]; !ok {
					return errors.NotValidf("unreachability [%v][%v]", s0.ID, s1.ID)
				}
			}
		}
	}

	for idx, r := range m.Roads {
		if r.S <= 0 {
			return errors.NotValidf("bad source")
		}
		if r.D <= 0 {
			return errors.NotValidf("bad destination")
		}
		if !m.Cells.Has(r.S) {
			return errors.NotValidf("dangling source")
		}
		if !m.Cells.Has(r.D) {
			return errors.NotValidf("dangling destination")
		}
		if r.D == r.S {
			return errors.NotValidf("loop detected")
		}
		if idx > 0 && r.equals(*m.Roads[idx-1]) {
			return errors.NotValidf("duplicated road")
		}
	}
	return nil
}

// Build a new "Next Step" index for the current Map, and replace the previous index.
func (m *Map) rehash() {
	next := make(map[vector]uint64)

	// Fill with the immediate neighbors
	for _, r := range m.Roads {
		next[vector{r.S, r.D}] = r.D
	}

	add := func(src, dst, step uint64) {
		_, found := next[vector{src, dst}]
		if !found {
			next[vector{src, dst}] = step
		}
	}

	// Call one DFS per node and shortcut when possible
	for _, cell := range m.Cells {
		already := make(map[uint64]bool)
		q := newQueue()

		// Bootstrap the DFS with adjacent nodes
		for _, next := range m.CellAdjacency(cell.ID) {
			q.push(next, next)
			already[next] = true
			// No need to add this in the known routes, we already did it
			// with an iteration on the roads (much faster)
		}

		for !q.empty() {
			current, first := q.pop()
			neighbors := m.CellAdjacency(current)
			// TODO(jfs): shuffle the neighbors
			for _, next := range neighbors {
				if !already[next] {
					// Avoid passing again in the neighbor
					already[next] = true
					// Tell to continue at that neighbor
					q.push(next, first)
					// We already learned the shortest path to that neighbor
					add(cell.ID, next, first)
				}
			}
		}
	}

	m.steps = next
}

func (v Vertex) equals(other Vertex) bool { return v.ID == other.ID }

func (e Edge) equals(other Edge) bool { return e.S == other.S && e.D == other.D }

type vector struct {
	src uint64
	dst uint64
}

type dfsTrack struct {
	current uint64
	first   uint64
}

type queue struct {
	tab   []dfsTrack
	start int
}

func newQueue() queue {
	var q queue
	q.tab = make([]dfsTrack, 0)
	return q
}

func (q *queue) push(node, first uint64) {
	q.tab = append(q.tab, dfsTrack{node, first})
}

func (q *queue) pop() (uint64, uint64) {
	v := q.tab[q.start]
	q.start++
	return v.current, v.first
}

func (q *queue) empty() bool {
	return q.start == len(q.tab)
}
