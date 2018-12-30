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

	result := comm.Broadcast(packet)

	if result == nil {
		shared.Console.Error("Invalid Transaction")
	} else if result.CheckTx.Code != 0 {
		shared.Console.Error("Syntax, CheckTx Failed", result)
	} else if result.DeliverTx.Code != 0 {
		shared.Console.Error("Transaction, DeliverTx Failed", result)
	} else {
		shared.Console.Info("Returned Successfully", result)
	}
}
