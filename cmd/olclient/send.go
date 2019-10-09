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
	"errors"
	"strconv"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/client"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/spf13/cobra"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
)

var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Issue send transaction",
	Run:   IssueRequest,
}

type SendArguments struct {
	Party        []byte `json:"party"`
	CounterParty []byte `json:"counterParty"`
	Amount       string `json:"amount"`
	Currency     string `json:"currency"`
	Fee          string `json:"fee"`
	Gas          int64  `json:"gas"`
}

func (args *SendArguments) ClientRequest(currencies *balance.CurrencySet) (client.SendTxRequest, error) {
	c, ok := currencies.GetCurrencyByName(args.Currency)
	if !ok {
		return client.SendTxRequest{}, errors.New("currency not support:" + args.Currency)
	}
	f, err := strconv.ParseFloat(args.Amount, 64)
	if err != nil {
		return client.SendTxRequest{}, err
	}
	amt := c.NewCoinFromFloat64(f).Amount
	olt, _ := currencies.GetCurrencyByName("OLT")
	fee, err := strconv.ParseFloat(args.Fee, 64)
	feeAmt := olt.NewCoinFromFloat64(fee).Amount
	return client.SendTxRequest{
		From:   args.Party,
		To:     args.CounterParty,
		Amount: action.Amount{Currency: args.Currency, Value: *amt},
		Fee:    action.Amount{Currency: "OLT", Value: *feeAmt},
		Gas:    args.Gas,
	}, nil
}

var sendargs = &SendArguments{}

func init() {
	RootCmd.AddCommand(sendCmd)

	// Transaction Parameters
	sendCmd.Flags().BytesHexVar(&sendargs.Party, "party", []byte{}, "send sender")
	sendCmd.Flags().BytesHexVar(&sendargs.CounterParty, "counterparty", []byte{}, "send recipient")
	sendCmd.Flags().StringVar(&sendargs.Amount, "amount", "0", "specify an amount")
	sendCmd.Flags().StringVar(&sendargs.Currency, "currency", "OLT", "the currency")
	sendCmd.Flags().StringVar(&sendargs.Fee, "fee", "0", "include a fee in OLT")
	sendCmd.Flags().Int64Var(&sendargs.Gas, "gas", 20000, "gas limit")
}

// IssueRequest sends out a sendTx to all of the nodes in the chain
func IssueRequest(cmd *cobra.Command, args []string) {

	ctx := NewContext()
	fullnode := ctx.clCtx.FullNodeClient()
	currencies, err := fullnode.ListCurrencies()
	if err != nil {
		ctx.logger.Error("failed to get currencies", err)
		return
	}
	// Create message
	req, err := sendargs.ClientRequest(currencies.Currencies.GetCurrencySet())
	if err != nil {
		ctx.logger.Error("failed to get request", err)
		return
	}
	reply, err := fullnode.SendTx(req)
	if err != nil {
		ctx.logger.Error("failed to create SendTx", err)
		return
	}
	packet := reply.RawTx

	result, err := ctx.clCtx.BroadcastTxCommit(packet)
	if err != nil {
		ctx.logger.Error("error in BroadcastTxCommit", err)
	}

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
		ctx.logger.Infof("Returned Successfully %#v", result)
		ctx.logger.Info("Result Data", "data", string(result.DeliverTx.Data))
	}
}
