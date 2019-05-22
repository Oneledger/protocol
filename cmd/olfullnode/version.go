// /*
// 	Copyright 2017-2018 OneLedger
//
// 	Cli to interact with a with the chain.
// */
package main

import (
	"fmt"
	"github.com/Oneledger/protocol/version"
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
	fmt.Println("Protocol version: " + version.Protocol.String())
	fmt.Println("Olfullnode version: " + version.Fullnode.String())
}
