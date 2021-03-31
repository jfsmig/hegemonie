// Copyright (c) 2018-2021 Contributors as noted in the AUTHORS file
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package mapclient

import (
	"encoding/json"
	"fmt"
	"github.com/hegemonie-rpg/engine/pkg/utils"
	"github.com/juju/errors"
	"os"
)

func loadAndDo(action func(m MapRaw) error) error {
	var raw MapRaw
	err := json.NewDecoder(os.Stdin).Decode(&raw)
	if err != nil {
		return errors.NewNotValid(err, "invalid json")
	}
	return action(raw)
}

// ToolNoise consumes a MapRaw on os.Stdin, parses it then alter each node position
// (with a configurable statistical noise) and the dumps the new MapRaw on os.Stdout.
func ToolNoise(noise float64) error {
	return loadAndDo(func(raw MapRaw) error {
		m, err := raw.extractMemMap()
		if err != nil {
			return errors.Trace(err)
		}
		xmin, xmax, ymin, ymax := m.computeBox()
		if noise > 0 {
			m.applyNoiseOnPositions(float64(xmax-xmin)*(noise/100), float64(ymax-ymin)*(noise/100))
		}
		return utils.DumpJSON(m.extractRawMap())
	})
}

// ToolSplit consumes a MapRaw on os.Stdin, parses it then adds nodes until each road (edge)
// is shorter than a threshold (given by maxDist)
func ToolSplit(maxDist float64) error {
	return loadAndDo(func(raw MapRaw) error {
		m, err := raw.extractMemMap()
		if err != nil {
			return errors.Trace(err)
		}
		if maxDist > 0 {
			m = m.splitLongRoads(maxDist)
		}
		return utils.DumpJSON(m.extractRawMap())
	})
}

// ToolNormalize consumes a MapRaw on os.Stdin, parses it a remap each node position
// (keeping the same aspect ratio to fit the whole map into a given bounded box.
func ToolNormalize() error {
	return loadAndDo(func(raw MapRaw) error {
		m, err := raw.extractMemMap()
		if err != nil {
			return errors.Trace(err)
		}
		var xbound, ybound, xPad, yPad uint64 = 1920, 1080, 50, 50
		m.resizeAndAdjust(xbound-2*xPad, ybound-2*yPad)
		m.SiftToTheCenter(xbound, ybound)
		return utils.DumpJSON(m.extractRawMap())
	})
}

// ToolDot consumes a MapRaw on os.Stdin, parses it and dumps to os.Stdout a
// dot representation (cf. graphviz.org) of it.
func ToolDot() error {
	return loadAndDo(func(raw MapRaw) error {
		m, err := raw.extractMemMap()
		if err != nil {
			return errors.Trace(err)
		}

		fmt.Println("graph g {")
		for r := range m.uniqueRoads() {
			fmt.Printf("%s -- %s;\n", r.Src.getDotName(), r.Dst.getDotName())
		}
		fmt.Println("}")
		return nil
	})
}

// ToolSvg consumes a MapRaw on os.Stdin, parses it and dumps a SVG representation
// of it to os.Stdout.
func ToolSvg() error {
	return loadAndDo(func(raw MapRaw) error {
		m, err := raw.extractMemMap()
		if err != nil {
			return errors.Trace(err)
		}

		var xbound, ybound, xPad, yPad uint64 = 1920, 1080, 50, 50
		m.resizeAndAdjust(xbound-2*xPad, ybound-2*yPad)
		m.SiftToTheCenter(xbound, ybound)

		fmt.Printf(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE svg PUBLIC "-//W3C//DTD SVG 1.1//EN" "http://www.w3.org/Graphics/SVG/1.1/DTD/svg11.dtd">
<svg xmlns="http://www.w3.org/2000/svg"
	style="background-color: rgb(255, 255, 255);"
	xmlns:xlink="http://www.w3.org/1999/xlink"
	version="1.1"
	width="%dpx" height="%dpx"
	viewBox="-0.5 -0.5 %d %d">
`, int64(xbound), int64(ybound), int64(xbound), int64(ybound))
		fmt.Println(`<g>`)
		for r := range m.uniqueRoads() {
			fmt.Printf(`<line x1="%d" y1="%d" x2="%d" y2="%d" stroke="black" stroke-width="1"/>
`, int64(r.Src.Raw.X), int64(r.Src.Raw.Y), int64(r.Dst.Raw.X), int64(r.Dst.Raw.Y))
		}
		fmt.Println(`</g>`)
		fmt.Println(`<g>`)
		for s := range m.sortedSites() {
			color := `white`
			radius := 5
			stroke := 1
			if s.Raw.City != "" {
				color = `gray`
				radius = 10
				stroke = 1
			}
			fmt.Printf(`<circle id="%v" class="clickable" cx="%d" cy="%d" r="%d" stroke="black" stroke-width="%d" fill="%s"/>
`, s.Raw.ID, int64(s.Raw.X), int64(s.Raw.Y), radius, stroke, color)
		}
		fmt.Println(`</g>`)
		fmt.Println(`</svg>`)
		return nil
	})
}

// ToolInit consumes a MapSeed on os.Stdin, parses it, extrapolates the MapRaw from it
// and then dumps that MapRaw to os.Stdout.
func ToolInit() error {
	var in MapSeed
	var out MapRaw
	err := json.NewDecoder(os.Stdin).Decode(&in)
	if err != nil {
		return errors.NewNotValid(err, "invalid json")
	}
	out, err = in.extractRawMap()
	if err != nil {
		return errors.Trace(err)

	}
	return utils.DumpJSON(out)
}
