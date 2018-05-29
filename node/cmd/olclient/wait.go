/*
	Copyright 2017-2018 OneLedger

	Cli to interact with a with the chain.
*/
package main

import (
	"time"

	"github.com/Oneledger/protocol/node/log"
	"github.com/spf13/cobra"
)

var waitCmd = &cobra.Command{
	Use:   "wait",
	Short: "Wait for something to happen",
	Run:   Wait,
}

func init() {
	RootCmd.AddCommand(waitCmd)

	var completed bool
	var strings []string

	waitCmd.Flags().BoolVar(&completed, "completed", false, "send recipient")
	waitCmd.Flags().StringArrayVar(&strings, "party", strings, "send recipient")
}

func Wait(cmd *cobra.Command, args []string) {
	log.Debug("Waiting")
	time.Sleep(10 * 1024 * 1024)

}
