// Copyright (c) 2018-2021 Contributors as noted in the AUTHORS file
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package evtbacklocal

import (
	"bytes"
	"fmt"
	"github.com/jfsmig/hegemonie/pkg/utils"
	"github.com/tecbot/gorocksdb"
	"math"
	"strconv"
	"strings"
	"time"
)

// Backend implements a
type Backend struct {
	db *gorocksdb.DB
}

// Item is an Event record.
type Item struct {
	// UserID the unique ID of the character the event belongs to
	UserID string
	// When the timestamp at which the Event 	occured
	When uint64
	// ID the unique ID of the event.
	ID string
	// Payload the actual data of the encoded parameters to render it.
	Payload []byte
}

// Open returns a Backend that's ready to work or an error
func Open(path string) (*Backend, error) {
	options := gorocksdb.NewDefaultOptions()
	options.SetCreateIfMissing(true)

	db, err := gorocksdb.OpenDb(options, path)
	if err != nil {
		return nil, err
	}

	return &Backend{db: db}, nil
}

// Push1 inserts an event record in the current backend.
// The timestamps is determined by the current backend itself.
func (b *Backend) Push1(UserID string, id string, payload []byte) error {
	opts := gorocksdb.NewDefaultWriteOptions()
	opts.SetSync(false)
	defer opts.Destroy()

	when := math.MaxUint64 - uint64(time.Now().UnixNano())
	k := fmt.Sprintf("%s/%16X/%s", UserID, when, id)
	utils.Logger.Warn().Bytes("key", []byte(k)).Msg("PUSH")
	return b.db.Put(opts, []byte(k), payload)
}

// Ack1 removes makes the event record cannot be listed anymore.
// The current Backend implementation simply deletes the Item.
func (b *Backend) Ack1(UserID string, when uint64, id string) error {
	opts := gorocksdb.NewDefaultWriteOptions()
	opts.SetSync(false)
	defer opts.Destroy()

	w := math.MaxUint64 - when
	k := fmt.Sprintf("%s/%16X/%s", UserID, w, id)
	utils.Logger.Warn().Bytes("key", []byte(k)).Msg("DEL")
	return b.db.Delete(opts, []byte(k))
}

// List returns a sorted array of event records started at the given timestamp
// and unique ID, ut all belongng to the Character whose ID is given.
func (b *Backend) List(UserID string, when uint64, max uint32) ([]Item, error) {
	if max <= 0 {
		max = 100
	} else if max > 1000 {
		max = 1000
	}

	var w uint64
	if when == 0 {
		w = 0
	} else {
		w = math.MaxUint64 - when
	}

	prefix := []byte(fmt.Sprintf("%s/", UserID))
	needle := []byte(fmt.Sprintf("%s/%016X/", UserID, w))

	opts := gorocksdb.NewDefaultReadOptions()
	defer opts.Destroy()
	opts.SetFillCache(true)
	opts.SetVerifyChecksums(false)
	iterator := b.db.NewIterator(opts)
	iterator.Seek(needle)

	var err error
	out := make([]Item, 0)
	for ; iterator.Valid(); iterator.Next() {
		k := iterator.Key().Data()
		if !bytes.HasPrefix(k, prefix) {
			break
		}
		tokens := strings.SplitN(string(k), "/", 3)
		when, err = strconv.ParseUint(tokens[1], 16, 64)
		out = append(out, Item{
			UserID:  UserID,
			When:    when,
			ID:      tokens[2],
			Payload: iterator.Value().Data(),
		})
	}

	if err != nil {
		return nil, err
	}
	return out, nil
}
