/*
	Copyright 2017-2018 OneLedger

	Cli to interact with a with the chains.
*/
package main

import (
	"github.com/Oneledger/protocol/node/comm"
	"github.com/Oneledger/protocol/node/log"
	"github.com/spf13/cobra"
)

var exeCmd = &cobra.Command{
	Use:   "execute",
	Short: "Execute Script Test",
	Run:   ExecuteTest,
}

func init() {
	RootCmd.AddCommand(exeCmd)
}

// IssueRequest sends out a sendTx to all of the nodes in the chain
func ExecuteTest(cmd *cobra.Command, args []string) {
	result := comm.Query("/testScript", []byte(""))
	log.Dump("Test Results", result)
}
