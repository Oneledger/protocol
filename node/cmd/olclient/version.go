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
	Short: "List out the version number",
	Run:   Version,
}

func init() {
	RootCmd.AddCommand(versionCmd)
}

func Version(cmd *cobra.Command, args []string) {
	version := version.Current
	shared.Console.Info("Olclient version is " + version.String())

	// TODO: Way better to ask the node, than to assume
	shared.Console.Info("Fullnode version is " + version.String())
}
