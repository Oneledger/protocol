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
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
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
	sendCmd.Flags().Float64Var(&sendargs.Amount, "amount", 0.0, "specify an amount")
	sendCmd.Flags().StringVar(&sendargs.Currency, "currency", "OLT", "the currency")

	sendCmd.Flags().Float64Var(&sendargs.Fee, "fee", 0.0, "include a fee in OLT")
}

// IssueRequest sends out a sendTx to all of the nodes in the chain
func IssueRequest(cmd *cobra.Command, args []string) {
	log.Debug("Have Send Request", "sendargs", sendargs)

	shared.Console.Info(sendargs)
	// Create message

	ctx := comm.NewClientContext()
	packet := shared.CreateSendRequest(ctx, sendargs)

	if packet == nil {
		shared.Console.Error("Error in sending request")
		os.Exit(-1)
	}

	result, _ := ctx.BroadcastTxCommit(packet)
	BroadcastStatus(result)
}

func BroadcastStatus(result *ctypes.ResultBroadcastTxCommit) {
	if result == nil {
		shared.Console.Error("Invalid Transaction")

	} else if result.CheckTx.Code != 0 {
		if result.CheckTx.Code == 200 {
			shared.Console.Info("Returned Successfully(fullnode query)", result)
			shared.Console.Info("Result Data", "data", string(result.CheckTx.Data))
		} else {
			shared.Console.Error("Syntax, CheckTx Failed", result)
		}

	} else if result.DeliverTx.Code != 0 {
		shared.Console.Error("Transaction, DeliverTx Failed", result)

	} else {
		shared.Console.Info("Returned Successfully", result)
		shared.Console.Info("Result Data", "data", string(result.DeliverTx.Data))
	}
}
