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

var exSendCmd = &cobra.Command{
	Use:   "external send",
	Short: "Issue send transaction to external chain",
	Run:   IssueRequest,
}

var exsendargs *shared.ExSendArguments = &shared.ExSendArguments{}

func init() {
	RootCmd.AddCommand(exSendCmd)

	// Transaction Parameters
	exSendCmd.Flags().StringVar(&exsendargs.SenderId, "senderid", "", "external sender identity")
	exSendCmd.Flags().StringVar(&exsendargs.ReceiverId, "receiverid", "", "external recipient identity")
	exSendCmd.Flags().StringVar(&exsendargs.SenderAddress, "senderaddress", "", "external sender address")
	exSendCmd.Flags().StringVar(&exsendargs.ReceiverAddress, "receiveraddress", "", "external recipient address")
	exSendCmd.Flags().StringVar(&exsendargs.Amount, "amount", "0", "specify an amount")
	exSendCmd.Flags().StringVar(&exsendargs.Currency, "currency", "-1", "the currency")

	exSendCmd.Flags().StringVar(&exsendargs.Fee, "fee", "4", "include a fee")
	exSendCmd.Flags().StringVar(&exsendargs.Gas, "gas", "1", "include gas")

	exSendCmd.Flags().StringVar(&exsendargs.Chain, "chain", "0", "destination chain")
	exSendCmd.Flags().StringVar(&exsendargs.ExFee, "exfee", "0", "include a external fee")
	exSendCmd.Flags().StringVar(&exsendargs.ExGas, "exgas", "0", "include external gas")
}

// IssueRequest sends out a sendTx to all of the nodes in the chain
func IssueExSend(cmd *cobra.Command, args []string) {
	log.Debug("Have External Send Request", "exsendargs", exsendargs)

	// Create message
	packet := shared.CreateExSendRequest(exsendargs)

	result := comm.Broadcast(packet)

	log.Debug("Returned Successfully", "result", result)
}
