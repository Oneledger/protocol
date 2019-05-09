/*
   ____             _              _                      _____           _                  _
  / __ \           | |            | |                    |  __ \         | |                | |
 | |  | |_ __   ___| |     ___  __| | __ _  ___ _ __     | |__) | __ ___ | |_ ___   ___ ___ | |
 | |  | | '_ \ / _ \ |    / _ \/ _` |/ _` |/ _ \ '__|    |  ___/ '__/ _ \| __/ _ \ / __/ _ \| |
 | |__| | | | |  __/ |___|  __/ (_| | (_| |  __/ |       | |   | | | (_) | || (_) | (_| (_) | |
  \____/|_| |_|\___|______\___|\__,_|\__, |\___|_|       |_|   |_|  \___/ \__\___/ \___\___/|_|
                                      __/ |
                                     |___/


Copyright 2017 - 2019 OneLedger
*/

package main

import (
	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/client"
	"github.com/Oneledger/protocol/data"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/spf13/cobra"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
)

var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Issue send transaction",
	Run:   IssueRequest,
}

var sendargs *client.SendArguments = &client.SendArguments{}

func init() {
	RootCmd.AddCommand(sendCmd)

	// Transaction Parameters
	sendCmd.Flags().BytesHexVar(&sendargs.Party, "party", []byte{}, "send sender")
	sendCmd.Flags().BytesHexVar(&sendargs.CounterParty, "counterparty", []byte{}, "send recipient")
	sendCmd.Flags().Float64Var(&sendargs.Amount, "amount", 0.0, "specify an amount")
	sendCmd.Flags().StringVar(&sendargs.Currency, "currency", "OLT", "the currency")

	sendCmd.Flags().Float64Var(&sendargs.Fee, "fee", 0.0, "include a fee in OLT")
}

// IssueRequest sends out a sendTx to all of the nodes in the chain
func IssueRequest(cmd *cobra.Command, args []string) {

	ctx := NewContext()

	ctx.logger.Debug("Have Send Request", "sendargs", sendargs)

	ctx.logger.Infof("%=v", sendargs)

	// Create message
	resp := &data.Response{}
	err := ctx.clCtx.Query("sendTx", *sendargs, resp)
	if err != nil {

	}

	packet := resp.Data

	if packet == nil {
		ctx.logger.Error("Error in sending request")
		return
	}

	result, _ := ctx.clCtx.BroadcastTxCommit(packet)
	BroadcastStatus(ctx, result)
}

func BroadcastStatus(ctx *Context, result *ctypes.ResultBroadcastTxCommit) {
	if result == nil {
		ctx.logger.Error("Invalid Transaction")

	} else if result.CheckTx.Code != 0 {
		if result.CheckTx.Code == 200 {
			ctx.logger.Info("Returned Successfully(fullnode query)", result)
			ctx.logger.Info("Result Data", "data", string(result.CheckTx.Data))
		} else {
			ctx.logger.Error("Syntax, CheckTx Failed", result)
		}

	} else if result.DeliverTx.Code != 0 {
		ctx.logger.Error("Transaction, DeliverTx Failed", result)

	} else {
		ctx.logger.Info("Returned Successfully", result)
		ctx.logger.Info("Result Data", "data", string(result.DeliverTx.Data))
	}
}
