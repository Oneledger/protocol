/*
	Copyright 2017-2018 OneLedger

	Cli to interact with a with the chain.
*/
package main

import (
	"os"

	"github.com/Oneledger/protocol/node/cmd/shared"
	"github.com/Oneledger/protocol/node/comm"
	"github.com/Oneledger/protocol/node/log"
	"github.com/spf13/cobra"
	"os"
)

var applyvalidatorCmd = &cobra.Command{
	Use:   "applyvalidator",
	Short: "Apply a dynamic validator",
	Run:   ApplyValidator,
}

var applyValidatorArgs *comm.ApplyValidatorArguments = &comm.ApplyValidatorArguments{}

func init() {
	RootCmd.AddCommand(applyvalidatorCmd)

	// Transaction Parameters
	applyvalidatorCmd.Flags().StringVar(&applyValidatorArgs.Id, "id", "", "specify identity by name")
	applyvalidatorCmd.Flags().Float64Var(&applyValidatorArgs.Amount, "amount", 0.0, "specify an amount")
	applyvalidatorCmd.Flags().BoolVar(&applyValidatorArgs.Purge, "purge", false, "remove the validator")
}

// IssueRequest sends out a sendTx to all of the nodes in the chain
func ApplyValidator(cmd *cobra.Command, args []string) {
	log.Debug("Have ApplyValidator Request", "applyValidatorArgs", applyValidatorArgs)

	// Create message
	packet := shared.CreateApplyValidatorRequest(applyValidatorArgs)
	if packet == nil {
		os.Exit(-1)
	}

	if packet == nil {
		shared.Console.Error("Error in sending request")
		os.Exit(-1)
	}

	result := comm.Broadcast(packet)
	BroadcastStatus(result)
}
