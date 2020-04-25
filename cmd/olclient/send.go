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
	"fmt"
	"os"
	"strconv"
	"strings"

	accounts2 "github.com/Oneledger/protocol/data/accounts"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/serialize"
	"github.com/spf13/cobra"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/client"
	"github.com/Oneledger/protocol/data/balance"
)

type SendArguments struct {
	Party        []byte `json:"party"`
	CounterParty []byte `json:"counterParty"`
	Amount       string `json:"amount"`
	Currency     string `json:"currency"`
	Fee          string `json:"fee"`
	Gas          int64  `json:"gas"`
	Password     string `json:"password"`
}

var (
	sendCmd = &cobra.Command{
		Use:   "send",
		Short: "Issue send transaction",
		Run:   IssueRequest,
	}

	sendFundsCmd = &cobra.Command{
		Use:   "sendfunds",
		Short: "Send funds to a given address",
		RunE:  sendFunds,
	}

	sendargs      = &SendArguments{}
	sendfundsargs = &SendArguments{}

	testenv = "OLTEST"
)

func (args *SendArguments) ClientRequest(currencies *balance.CurrencySet) (client.SendTxRequest, error) {
	c, ok := currencies.GetCurrencyByName(args.Currency)
	if !ok {
		return client.SendTxRequest{}, errors.New("currency not support:" + args.Currency)
	}
	padZero := func(s string) string {
		ss := strings.Split(s, ".")
		if len(ss) == 2 {
			ss = []string{strings.TrimLeft(ss[0], "0"), strings.TrimLeft(ss[1], "0"), strings.Repeat("0", 18-len(ss[1]))}
		} else {
			ss = []string{strings.TrimLeft(ss[0], "0"), strings.Repeat("0", 18)}
		}
		s = strings.Join(ss, "")
		return s
	}
	_, err := strconv.ParseFloat(args.Amount, 64)
	if err != nil {
		return client.SendTxRequest{}, err
	}
	amt := c.NewCoinFromString(padZero(args.Amount)).Amount

	olt, _ := currencies.GetCurrencyByName("OLT")

	_, err = strconv.ParseFloat(args.Fee, 64)
	if err != nil {
		return client.SendTxRequest{}, err
	}
	feeAmt := olt.NewCoinFromString(padZero(args.Fee)).Amount

	return client.SendTxRequest{
		From:     args.Party,
		To:       args.CounterParty,
		Amount:   action.Amount{Currency: args.Currency, Value: *amt},
		GasPrice: action.Amount{Currency: "OLT", Value: *feeAmt},
		Gas:      args.Gas,
	}, nil
}

func init() {

	RootCmd.AddCommand(sendCmd)
	setArgs(sendCmd, sendargs)

	testEnv := os.Getenv(testenv)
	if testEnv == "1" {
		RootCmd.AddCommand(sendFundsCmd)
		setArgs(sendFundsCmd, sendfundsargs)
	}
}

func setArgs(command *cobra.Command, sendArgs *SendArguments) {
	// Transaction Parameters
	command.Flags().BytesHexVar(&sendArgs.Party, "party", []byte{}, "send sender")
	command.Flags().BytesHexVar(&sendArgs.CounterParty, "counterparty", []byte{}, "send recipient")
	command.Flags().StringVar(&sendArgs.Amount, "amount", "0", "specify an amount")
	command.Flags().StringVar(&sendArgs.Currency, "currency", "OLT", "the currency")
	command.Flags().StringVar(&sendArgs.Fee, "fee", "0", "include a fee in OLT")
	command.Flags().StringVar(&sendArgs.Password, "password", "", "password to access secure wallet.")
	command.Flags().Int64Var(&sendArgs.Gas, "gas", 20000, "gas limit")
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
	fmt.Println(req)

	//Prompt for password
	if len(sendargs.Password) == 0 {
		sendargs.Password = PromptForPassword()
	}

	//Create new Wallet and User Address
	wallet, err := accounts2.NewWalletKeyStore(keyStorePath)
	if err != nil {
		ctx.logger.Error("failed to create secure wallet", err)
		return
	}

	//Verify User Password
	usrAddress := keys.Address(sendargs.Party)
	authenticated, err := wallet.VerifyPassphrase(usrAddress, sendargs.Password)
	if !authenticated {
		ctx.logger.Error("authentication error", err)
		return
	}

	//Get Raw "Send" Transaction
	reply, err := fullnode.CreateRawSend(req)
	if err != nil {
		ctx.logger.Error("failed to create SendTx", err)
		return
	}
	rawTx := &action.RawTx{}
	err = serialize.GetSerializer(serialize.NETWORK).Deserialize(reply.RawTx, rawTx)
	if err != nil {
		ctx.logger.Error("failed to deserialize RawTx", err)
		return
	}

	if !wallet.Open(usrAddress, sendargs.Password) {
		ctx.logger.Error("failed to open secure wallet")
		return
	}

	//Sign Raw "Send" Transaction Using Secure Wallet.
	pub, signature, err := wallet.SignWithAddress(reply.RawTx, usrAddress)
	if err != nil {
		ctx.logger.Error("error signing transaction", err)
	}

	signatures := []action.Signature{{pub, signature}}
	signedTx := &action.SignedTx{
		RawTx:      *rawTx,
		Signatures: signatures,
	}

	packet, err := serialize.GetSerializer(serialize.NETWORK).Serialize(signedTx)
	if err != nil {
		ctx.logger.Error("failed to serialize signedTx", err)
		return
	}

	//Broadcast Transaction
	result, err := ctx.clCtx.BroadcastTxCommit(packet)
	if err != nil {
		ctx.logger.Error("error in BroadcastTxCommit", err)
	}

	BroadcastStatus(ctx, result)
}

// IssueRequest sends out a sendTx to all of the nodes in the chain
func sendFunds(cmd *cobra.Command, args []string) error {

	ctx := NewContext()
	fullnode := ctx.clCtx.FullNodeClient()
	currencies, err := fullnode.ListCurrencies()
	if err != nil {
		ctx.logger.Error("failed to get currencies", err)
		return err
	}
	// Create message
	req, err := sendfundsargs.ClientRequest(currencies.Currencies.GetCurrencySet())
	if err != nil {
		ctx.logger.Error("failed to get request", err)
		return err
	}
	fmt.Println(req)
	reply, err := fullnode.SendTx(req)
	if err != nil {
		ctx.logger.Error("failed to create SendTx", err)
		return err
	}
	packet := reply.RawTx

	result, err := ctx.clCtx.BroadcastTxCommit(packet)
	if err != nil {
		ctx.logger.Error("error in BroadcastTxCommit", err)
	}

	BroadcastStatus(ctx, result)

	return nil
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
