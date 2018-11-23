/*
	Copyright 2017-2018 OneLedger

	Cli to interact with a with the chain.
*/
package main

import (
	"github.com/Oneledger/protocol/node/cmd/shared"
	"github.com/Oneledger/protocol/node/comm"
	"github.com/Oneledger/protocol/node/log"
	"github.com/spf13/cobra"
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
	applyvalidatorCmd.Flags().StringVar(&applyValidatorArgs.Amount, "amount", "0", "specify an amount")
}

// IssueRequest sends out a sendTx to all of the nodes in the chain
func ApplyValidator(cmd *cobra.Command, args []string) {
	log.Debug("Have ApplyValidator Request", "applyValidatorArgs", applyValidatorArgs)

	// Create message
	packet := shared.CreateApplyValidatorRequest(applyValidatorArgs)

	result := comm.Broadcast(packet)

	log.Debug("Returned Successfully", "result", result)
}
