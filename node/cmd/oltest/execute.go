/*
	Copyright 2017-2018 OneLedger

	Cli to interact with a with the chains.
*/
package main

import (
	"time"

	"github.com/Oneledger/protocol/node/comm"
	"github.com/Oneledger/protocol/node/log"
	"github.com/spf13/cobra"
)

var exeCmd = &cobra.Command{
	Use:   "execute",
	Short: "Execute Script Test",
	Run:   ExecuteTest,
}

type ExeArgs struct {
	Test string
}

var exeargs ExeArgs = ExeArgs{}

func init() {
	exeCmd.Flags().StringVar(&exeargs.Test, "test", "", "test name")

	RootCmd.AddCommand(exeCmd)
}

// IssueRequest sends out a sendTx to all of the nodes in the chain
func ExecuteTest(cmd *cobra.Command, args []string) {

	start := time.Now()
	result := comm.Query("/testScript", []byte(exeargs.Test))
	elapsed := time.Since(start)

	log.Dump("Test Results", result)
	log.Dump("Query Time", elapsed)
}
