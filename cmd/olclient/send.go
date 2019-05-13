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
	"github.com/Oneledger/protocol/client"
	"github.com/Oneledger/protocol/data"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/serialize"

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
	sendCmd.Flags().Float64Var(&sendargs.AmountFloat, "amount", 0.0, "specify an amount")
	sendCmd.Flags().StringVar(&sendargs.CurrencyStr, "currency", "OLT", "the currency")
	sendCmd.Flags().Float64Var(&sendargs.FeeFloat, "fee", 0.0, "include a fee in OLT")
}

// IssueRequest sends out a sendTx to all of the nodes in the chain
func IssueRequest(cmd *cobra.Command, args []string) {

	ctx := NewContext()

	ctx.logger.Debug("Have Send Request", "sendargs", sendargs)

	currResp := &data.Response{}
	err := ctx.clCtx.Query("RPCServerCtx.Currencies", data.Request{}, currResp)
	if err != nil {
		ctx.logger.Error("failed to get currencies from node", err)
		return
	}

	currencies := map[string]balance.Currency{}
	err = serialize.GetSerializer(serialize.CLIENT).Deserialize(currResp.Data, &currencies)
	if err != nil {

	}

	ctx.logger.Infof("arguments for send transaction: %#v", sendargs)

	currency, ok := currencies[sendargs.CurrencyStr]
	if !ok {
		ctx.logger.Errorf("currency %s not registered", sendargs.CurrencyStr)
		return
	}

	sendargs.Amount = currency.NewCoinFromFloat64(sendargs.AmountFloat)
	sendargs.Fee = currency.NewCoinFromFloat64(sendargs.FeeFloat)

	// Create message
	resp := &data.Response{}
	err = ctx.clCtx.Query("RPCServerCtx.SendTx", *sendargs, resp)
	if err != nil {
		ctx.logger.Error("error executing SendTx", err)
		return
	}

	packet := resp.Data

	if packet == nil {
		ctx.logger.Error("Error in sending ", resp.ErrorMsg)
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
