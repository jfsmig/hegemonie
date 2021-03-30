// Copyright (c) 2018-2021 Contributors as noted in the AUTHORS file
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"flag"
	"github.com/juju/errors"
	"log"
	"os"
	"strings"
	"text/template"
	"time"
)

var headerTemplate = `// Code generated : DO NOT EDIT.

// Copyright (c) 2018-2021 Contributors as noted in the AUTHORS file
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package {{.Package}}

import (
	"github.com/juju/errors"
	"math/rand"
	"sort"
)
`

var bodyTemplateCommon = `
type {{.SetName}} []{{.ItemType}}

func (s {{.SetName}}) CheckThenFail() {
	if err := s.Check(); err != nil {
		panic(err.Error())
	}
}

func (s {{.SetName}}) Len() int {
	return len(s)
}

func (s {{.SetName}}) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s *{{.SetName}}) Add(a {{.ItemType}}) {
	*s = append(*s, a)
	switch nb := len(*s); nb {
	case 0:
		panic("yet another attack of a solar eruption")
	case 1:
		return
	case 2:
		sort.Sort(s)
	default:
		if !sort.IsSorted((*s)[nb-2:]) {
			sort.Sort(s)
		}
	}
}

func (s {{.SetName}}) Check() error {
	if !sort.IsSorted(s) {
		return errors.NotValidf("sorting (%v) %v", s.Len(), s)
	}
	if !s.areItemsUnique() {
		return errors.NotValidf("unicity")
	}
	return nil
}

func (s *{{.SetName}}) testRandomVacuum() {
	for s.Len() > 0 {
		idx := rand.Intn(s.Len())
		s.Remove((*s)[idx])
		s.CheckThenFail()
	}
}
`

var bodyTemplate1 = `
func (s {{.SetName}}) Less(i, j int) bool {
	return s[i]{{.F0}} < s[j]{{.F0}}
}

func (s {{.SetName}}) areItemsUnique() bool {
	var lastId {{.T0}}
	for _, a := range s {
		if lastId == a{{.F0}} {
			return false
		}
		lastId = a{{.F0}}
	}
	return true
}

func (s {{.SetName}}) Slice(marker {{.T0}}, max uint32) []{{.ItemType}} {
	if max == 0 {
		max = 1000
	} else if max > 100000 {
		max = 100000
	}
	start := sort.Search(len(s), func(i int) bool {
		return s[i]{{.F0}} > marker
	})
	if start < 0 || start >= s.Len() {
		return s[:0]
	}
	remaining := uint32(s.Len() - start)
	if remaining > max {
		remaining = max
	}
	return s[start : uint32(start)+remaining]
}

func (s {{.SetName}}) getIndex(id {{.T0}}) int {
	i := sort.Search(len(s), func(i int) bool {
		return s[i]{{.F0}} >= id
	})
	if i < len(s) && s[i]{{.F0}} == id {
		return i
	}
	return -1
}

func (s {{.SetName}}) Get(id {{.T0}}) {{.ItemType}} {
	var out {{.ItemType}}
	idx := s.getIndex(id)
	if idx >= 0 {
		out = s[idx]
	}
	return out
}

func (s {{.SetName}}) Has(id {{.T0}}) bool {
	return s.getIndex(id) >= 0
}

func (s *{{.SetName}}) Remove(a {{.ItemType}}) {
	s.RemovePK(a{{.F0}})
}

func (s *{{.SetName}}) RemovePK(pk {{.T0}}) {
	idx := s.getIndex(pk)
	if idx >= 0 && idx < len(*s) {
		if len(*s) == 1 {
			*s = (*s)[:0]
		} else {
			s.Swap(idx, s.Len()-1)
			*s = (*s)[:s.Len()-1]
			sort.Sort(*s)
		}
	}
}
`

var bodyTemplate2 = `
func (s {{.SetName}}) Less(i, j int) bool {
	p0, p1 := s[i], s[j]
	return p0{{.F0}} < p1{{.F0}} || (p0{{.F0}} == p1{{.F0}} && p0{{.F1}} < p1{{.F1}})
}

func (s {{.SetName}}) First(at {{.T0}}) int {
	return sort.Search(len(s), func(i int) bool { return s[i]{{.F0}} >= at })
}

func (s {{.SetName}}) areItemsUnique() bool {
	var l0 {{.T0}}
	var l1 {{.T1}}
	for _, a := range s {
		if l0 == a{{.F0}} && l1 == a{{.F1}} {
			return false
		}
		l0 = a{{.F0}}
	}
	return true
}

func (s {{.SetName}}) Slice(m0 {{.T0}}, m1 {{.T1}}, max uint32) []{{.ItemType}} {
	if max == 0 {
		max = 1000
	} else if max > 100000 {
		max = 100000
	}

	iMax := s.Len()
	start := s.First(m0)
	for start < iMax && s[start]{{.F0}} == m0 && s[start]{{.F1}} <= m1 {
		start++
	}

	remaining := uint32(iMax - start)
	if remaining > max {
		remaining = max
	}
	return s[start : uint32(start)+remaining]
}

func (s {{.SetName}}) getIndex(f0 {{.T0}}, f1 {{.T1}}) int {
	i := sort.Search(len(s), func(i int) bool {
		return s[i]{{.F0}} >= f0 || (s[i]{{.F0}} == f0 && s[i]{{.F1}} >= f1)
	})
	if i < len(s) && s[i]{{.F0}} == f0 && s[i]{{.F1}} == f1 {
		return i
	}
	return -1
}

func (s {{.SetName}}) Get(f0 {{.T0}}, f1 {{.T1}}) {{.ItemType}} {
	var out {{.ItemType}}
	idx := s.getIndex(f0, f1)
	if idx >= 0 {
		out = s[idx]
	}
	return out
}

func (s {{.SetName}}) Has(f0 {{.T0}}, f1 {{.T1}}) bool {
	return s.getIndex(f0, f1) >= 0
}

func (s *{{.SetName}}) Remove(a {{.ItemType}}) {
	s.RemovePK(a{{.F0}}, a{{.F1}})
}

func (s *{{.SetName}}) RemovePK(f0 {{.T0}}, f1 {{.T1}}) {
	idx := s.getIndex(f0, f1)
	if idx >= 0 && idx < len(*s) {
		if len(*s) == 1 {
			*s = (*s)[:0]
		} else {
			s.Swap(idx, s.Len()-1)
			*s = (*s)[:s.Len()-1]
			sort.Sort(*s)
		}
	}
}
`

type arrayInstance struct {
	Path string
	Date string

	Package string
	SetName string

	ItemType string

	// First field
	F0 string
	T0 string

	// Second field
	F1 string
	T1 string
}

func run(args []string) error {
	var err error
	var fout *os.File
	var instance arrayInstance

	instance.Date = time.Now().String()
	instance.Path = args[0]
	instance.F0 = ".ID"
	instance.T0 = "uint64"

	p := args[1]
	tokens := strings.Split(p, ":")

	instance.Package = tokens[0]
	instance.SetName = tokens[1]
	instance.ItemType = tokens[2]

	if len(args) > 2 {
		p = args[2]
		tokens = strings.Split(p, ":")
		instance.F0 = tokens[0]
		instance.T0 = tokens[1]
		if instance.F0 != "" {
			instance.F0 = "." + instance.F0
		}
		if instance.T0 == "" {
			instance.T0 = "uint64"
		}
	}

	if len(args) > 3 {
		p = args[3]
		tokens = strings.Split(p, ":")
		instance.F1 = tokens[0]
		instance.T1 = tokens[1]
		if instance.F1 != "" {
			instance.F1 = "." + instance.F1
		}
		if instance.T1 == "" {
			instance.T1 = "uint64"
		}
	}

	if len(args) > 4 {
		return errors.BadRequestf("Too many args")
	}

	header, err := template.New("header").Parse(headerTemplate)
	if err != nil {
		return errors.Annotate(err, "Invalid template")
	}

	common, err := template.New("common").Parse(bodyTemplateCommon)
	if err != nil {
		return errors.Annotate(err, "Invalid template")
	}

	tpl := bodyTemplate1
	if len(args) == 4 {
		tpl = bodyTemplate2
	}
	body, err := template.New("body").Parse(tpl)
	if err != nil {
		return errors.Annotate(err, "Invalid template")
	}

	_, err = os.Stat(instance.Path)
	if err == nil {
		fout, err = os.OpenFile(instance.Path, os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return errors.Annotate(err, "Invalid output file")
		}
	} else {
		fout, err = os.Create(instance.Path)
		if err != nil {
			return errors.Annotate(err, "Invalid output file")
		}
		header.Execute(fout, &instance)
	}
	defer fout.Close()

	common.Execute(fout, &instance)
	body.Execute(fout, &instance)
	return nil
}

func main() {
	flag.Parse()
	if err := run(flag.Args()); err != nil {
		log.Fatalln("Generation failed", err)
	}
}
