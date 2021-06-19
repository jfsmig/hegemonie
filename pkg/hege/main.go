// Copyright (c) 2018-2021 Contributors as noted in the AUTHORS file
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"github.com/hegemonie-rpg/engine/pkg/utils"
	"github.com/juju/errors"
	"github.com/spf13/cobra"
	"log"
)

func main() {
	cmd := &cobra.Command{
		Use:   "hege",
		Short: "Hegemonie CLI",
		Long:  "Hegemonie client with subcommands for the several agents of an Hegemonie system.",
		Args:  cobra.MinimumNArgs(1),
		RunE:  nonLeaf,
	}
	utils.AddLogFlagsToCommand(cmd)
	cmd.AddCommand(clients(), servers(), tools())
	if err := cmd.Execute(); err != nil {
		log.Fatalln(errors.ErrorStack(err))
	}
}

func loadDiscovery(config *utils.MainConfig, path string, must bool) error {
	if err := config.LoadFile(path, must); err != nil {
		return err
	}

	utils.Logger.Debug().Str("path", path).RawJSON("cfg", utils.JSON2Buf(config)).Msg("Loaded")
	return nil
}

func nonLeaf(_ *cobra.Command, _ []string) error { return errors.New("missing subcommand") }

func first(args []string) string {
	if len(args) <= 0 {
		return ""
	}
	return args[0]
}

func second(args []string) string {
	if len(args) <= 1 {
		return ""
	}
	return args[1]
}
