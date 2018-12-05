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

var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Issue send transaction",
	Run:   IssueRequest,
}

var sendargs *comm.SendArguments = &comm.SendArguments{}

func init() {
	RootCmd.AddCommand(sendCmd)

	// Transaction Parameters
	sendCmd.Flags().StringVar(&sendargs.Party, "party", "", "send sender")
	sendCmd.Flags().StringVar(&sendargs.CounterParty, "counterparty", "", "send recipient")
	sendCmd.Flags().StringVar(&sendargs.Amount, "amount", "0", "specify an amount")
	sendCmd.Flags().StringVar(&sendargs.Currency, "currency", "OLT", "the currency")

	sendCmd.Flags().StringVar(&sendargs.Fee, "fee", "4", "include a fee")
	sendCmd.Flags().StringVar(&sendargs.Gas, "gas", "1", "include gas")
}

// IssueRequest sends out a sendTx to all of the nodes in the chain
func IssueRequest(cmd *cobra.Command, args []string) {
	log.Debug("Have Send Request", "sendargs", sendargs)

	// Create message
	packet := shared.CreateSendRequest(sendargs)

	result := comm.Broadcast(packet)

	log.Debug("Returned Successfully", "result", result)
}
