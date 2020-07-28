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

	"github.com/spf13/cobra"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/client"
	accounts2 "github.com/Oneledger/protocol/data/accounts"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/serialize"
)

func (args *SendPoolArguments) ClientRequest(currencies *balance.CurrencySet) (client.SendPoolTxRequest, error) {
	c, ok := currencies.GetCurrencyByName(args.Currency)
	if !ok {
		return client.SendPoolTxRequest{}, errors.New("currency not support:" + args.Currency)
	}
	_, err := strconv.ParseFloat(args.Amount, 64)
	if err != nil {
		return client.SendPoolTxRequest{}, err
	}
	amt := c.NewCoinFromString(padZero(args.Amount)).Amount

	olt, _ := currencies.GetCurrencyByName("OLT")

	_, err = strconv.ParseFloat(args.Fee, 64)
	if err != nil {
		return client.SendPoolTxRequest{}, err
	}
	feeAmt := olt.NewCoinFromString(padZero(args.Fee)).Amount

	return client.SendPoolTxRequest{
		From:     args.Party,
		PoolName: args.PoolName,
		Amount:   action.Amount{Currency: args.Currency, Value: *amt},
		GasPrice: action.Amount{Currency: "OLT", Value: *feeAmt},
		Gas:      args.Gas,
	}, nil
}

//Send funds to a Pool
func sendFundsPool(cmd *cobra.Command, args []string) error {

	ctx := NewContext()
	fullnode := ctx.clCtx.FullNodeClient()
	currencies, err := fullnode.ListCurrencies()
	if err != nil {
		ctx.logger.Error("failed to get currencies", err)
		return err
	}
	// Create message
	req, err := sendpoolargs.ClientRequest(currencies.Currencies.GetCurrencySet())
	if err != nil {
		ctx.logger.Error("failed to get request", err)
		return err
	}

	if len(sendpoolargs.Password) == 0 {
		sendpoolargs.Password = PromptForPassword()
	}
	wallet, err := accounts2.NewWalletKeyStore(keyStorePath)
	if err != nil {
		ctx.logger.Error("failed to create secure wallet", err)
		return err
	}
	usrAddress := keys.Address(sendpoolargs.Party)
	authenticated, err := wallet.VerifyPassphrase(usrAddress, sendpoolargs.Password)
	if !authenticated {
		ctx.logger.Error("authentication error", err)
		return err
	}

	reply, err := fullnode.CreateRawSendPool(req)
	if err != nil {
		ctx.logger.Error("failed to create SendPoolTx", err)
		return err
	}
	rawTx := &action.RawTx{}
	err = serialize.GetSerializer(serialize.NETWORK).Deserialize(reply.RawTx, rawTx)
	if err != nil {
		ctx.logger.Error("failed to deserialize RawTx", err)
		return err
	}

	if !wallet.Open(usrAddress, sendpoolargs.Password) {
		ctx.logger.Error("failed to open secure wallet")
		return err
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
		return err
	}

	//Broadcast Transaction
	result, err := ctx.clCtx.BroadcastTxCommit(packet)
	if err != nil {
		ctx.logger.Error("error in BroadcastTxCommit", err)
	}

	BroadcastStatus(ctx, result)

	return nil
}
