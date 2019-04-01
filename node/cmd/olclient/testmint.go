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

var mintargs *comm.SendArguments = &comm.SendArguments{}

func init() {
	RootCmd.AddCommand(mintCmd)

	// Transaction Parameters
	mintCmd.Flags().StringVar(&mintargs.Party, "party", "", "send recipient")
	mintCmd.Flags().Float64Var(&mintargs.Amount, "amount", 0.0, "specify an amount")
	mintCmd.Flags().StringVar(&mintargs.Currency, "currency", "OLT", "the currency")

	mintCmd.Flags().Float64Var(&mintargs.Fee, "fee", 0.0, "include a fee in OLT")
	//mintCmd.Flags().Int64Var(&mintargs.Gas, "gas", 0, "include gas in units")
}

// IssueRequest sends out a sendTx to all of the nodes in the chain
func IssueMintRequest(cmd *cobra.Command, args []string) {
	log.Debug("Have Testmint Request", "mintargs", mintargs)

	// Create message
	packet := shared.CreateMintRequest(mintargs)

	if packet == nil {
		log.Info("Bad Request", "mintargs", mintargs)
		return
	}

	result, _ := rpcclient.BroadcastTxCommit(packet)
	BroadcastStatus(result)
}
