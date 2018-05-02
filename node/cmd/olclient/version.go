/*
	Copyright 2017-2018 OneLedger

	Cli to interact with a with the chain.
*/
package main

import (
	"github.com/Oneledger/prototype/node/version"
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
	Console.Info("Client Version is " + version.String())

	// TODO: Way better to ask the node, than to assume
	Console.Info("Node Version is " + version.String())
}
