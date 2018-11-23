/*
	Copyright 2017-2018 OneLedger

	Cli to interact with a with the chain.
*/
package main

import (
	"github.com/Oneledger/protocol/node/cmd/shared"
	"github.com/Oneledger/protocol/node/version"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "List out the version numbers",
	Run:   Version,
}

func init() {
	RootCmd.AddCommand(versionCmd)
}

func Version(cmd *cobra.Command, args []string) {
	version := version.Current
	shared.Console.Info("Olclient version is " + version.String())

	// TODO: Query the fullnode, get it's version number
	shared.Console.Info("Olfullnode version is " + version.String())
}
