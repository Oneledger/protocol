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

var mintCmd = &cobra.Command{
	Use:   "testmint",
	Short: "Issue testmint transaction",
	Run:   IssueMintRequest,
}

var mintargs *shared.SendArguments = &shared.SendArguments{}

func init() {
	RootCmd.AddCommand(mintCmd)

	// Transaction Parameters
	mintCmd.Flags().StringVar(&mintargs.Party, "party", "", "send recipient")
	mintCmd.Flags().StringVar(&mintargs.Amount, "amount", "0", "specify an amount")
	mintCmd.Flags().StringVar(&mintargs.Currency, "currency", "OLT", "the currency")

	mintCmd.Flags().StringVar(&mintargs.Fee, "fee", "1", "include a fee")
	mintCmd.Flags().StringVar(&mintargs.Gas, "gas", "1", "include gas")
}

// IssueRequest sends out a sendTx to all of the nodes in the chain
func IssueMintRequest(cmd *cobra.Command, args []string) {
	log.Debug("Have Testmint Request", "mintargs", mintargs)

	// Create message
	packet := shared.CreateMintRequest(mintargs)

	result := comm.Broadcast(packet)

	log.Debug("Returned Successfully", "result", result)
}
