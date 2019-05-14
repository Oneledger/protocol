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
	"fmt"
	"math/rand"
	"time"

	"github.com/Oneledger/protocol/client"
	"github.com/Oneledger/protocol/data"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/serialize"

	"github.com/Oneledger/protocol/data/accounts"
	"github.com/Oneledger/protocol/data/chain"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/spf13/cobra"
)

var loadtestCmd = &cobra.Command{
	Use:   "send",
	Short: "Issue send transaction",
	Run:   LoadTest,
}

type LoadTestArgs struct {
	threads  int
	interval int
}

var loadTestArgs = LoadTestArgs{}

func init() {
	RootCmd.AddCommand(loadtestCmd)

	loadtestCmd.Flags().IntVar(&loadTestArgs.interval, "interval",
		1, "interval between successive transactions on a single thread in seconds")
	loadtestCmd.Flags().IntVar(&loadTestArgs.threads, "threads",
		1, "number of threads running")
}

func LoadTest(cmd *cobra.Command, args []string) {

	ctx := NewContext()
	ctx.logger.Debug("Starting Loadtest")

	// get a list of registered currencies
	currResp := &data.Response{}
	err := ctx.clCtx.Query("server.Currencies", data.Request{}, currResp)
	if err != nil {
		ctx.logger.Error("failed to get currencies from node", err)
		return
	}

	currencies := map[string]balance.Currency{}
	err = serialize.GetSerializer(serialize.CLIENT).Deserialize(currResp.Data, &currencies)
	if err != nil || len(currencies) == 0 {
		ctx.logger.Error("error getting currencies from server:")
	}

	// get OLT currency object
	OLTCurrency, ok := currencies["OLT"]
	if !ok {
		ctx.logger.Fatal("OLT not registered in currency list")
	}

	for i := 0; i < loadTestArgs.threads; i++ {

		// create a temp account to send OLT
		accName := fmt.Sprintf("acc_%03d", i)
		pubKey, privKey, err := keys.NewKeyPairFromTendermint() // generate a ed25519 key pair
		if err != nil {
			ctx.logger.Error(accName, "error generating key from tendermint", err)
		}

		acc, err := accounts.NewAccount(chain.Type(1), accName, &privKey, &pubKey) // create account object
		if err != nil {
			ctx.logger.Error(accName, "Error initializing account", err)
			return
		}

		ctx.logger.Infof("creating account %#v", acc)
		resp := &data.Response{}
		err = ctx.clCtx.Query("server.AddAccount", acc, resp) // create account on node
		if err != nil {
			ctx.logger.Error(accName, "error creating account", err)
			return
		}

		// start a thread to keep sending transactions after some interval
		go func() {
			for true {
				doSendTransaction(ctx, i, &acc, &OLTCurrency) // send OLT to temp account

				time.Sleep(time.Second * time.Duration(loadTestArgs.interval)) // wait for some time
			}
		}()

	}
}

// doSendTransaction takes in an account and currency object and sends random amounts of coin from the
// node account. It prints any errors to ctx.logger and returns
func doSendTransaction(ctx *Context, threadNo int, acc *accounts.Account, curr *balance.Currency) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("panic in doSendTransaction: thread", threadNo, r)
		}
	}()

	// generate a random amount
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	amt := r.Float64() * 10.0 // amount is a random float between [0, 10)

	// populate send arguments
	sendArgs := &client.SendArguments{}
	sendArgs.Party = NODE_ADDRESS
	sendArgs.CounterParty = []byte(acc.Address())  // receiver is the temp account
	sendArgs.Amount = curr.NewCoinFromFloat64(amt) // make coin from currency object
	sendArgs.Fee = curr.NewCoinFromFloat64(amt / 10.0)

	// Create message
	resp := &data.Response{}
	err := ctx.clCtx.Query("server.SendTx", *sendargs, resp)
	if err != nil {
		ctx.logger.Error(acc.Name, "error executing SendTx", err)
		return
	}

	packet := resp.Data
	if packet == nil {
		ctx.logger.Error(acc.Name, "Error in sending ", resp.ErrorMsg)
		return
	}

	result, err := ctx.clCtx.BroadcastTxAsync(packet)
	if err != nil {
		ctx.logger.Error(acc.Name, "error in BroadcastTxAsync:", err)
		return
	}
	ctx.logger.Info(acc.Name, "Result Data", "log", string(result.Log))
}
