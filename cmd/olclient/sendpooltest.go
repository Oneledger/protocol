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
	"os"
	"os/signal"
	"strconv"
	"sync"
	"time"

	"github.com/spf13/cobra"

	"github.com/Oneledger/protocol/data/accounts"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/chain"
	"github.com/Oneledger/protocol/data/keys"
)

// cobra command to loadtest
var sendpooltestCmd = &cobra.Command{
	Use:   "sendpooltest",
	Short: "launch send pool Test",
	Run:   sendPoolTest,
}

// global to load command line args
var sendpooltestArgs = SendPoolTestArgs{}

// struct to hold the command-line args
type SendPoolTestArgs struct {
	threads    int  // no. of threads in the load test; for concurrency
	interval   int  // interval (in milliseconds) between two successive send transactions on a thread
	randomRecv bool // whether to send tokens to a random address every time or no, the default is false
	maxTx      int  // max transactions after which the load test should stop, default is 10000(10k)
	address    []byte
}

// init function initializes the loadtest command and attaches a bunch of flag parsers
func init() {
	RootCmd.AddCommand(sendpooltestCmd)

	sendpooltestCmd.Flags().IntVar(&sendpooltestArgs.interval, "interval", 1,
		"interval between successive transactions on a single thread in milliseconds")

	sendpooltestCmd.Flags().IntVar(&sendpooltestArgs.threads, "threads", 1,
		"number of threads running")

	sendpooltestCmd.Flags().BoolVar(&sendpooltestArgs.randomRecv, "random-receiver", false,
		"whether to randomize the receiver every time, default false")

	sendpooltestCmd.Flags().IntVar(&sendpooltestArgs.maxTx, "max-tx", 10000,
		"number of max tx in before the load test stop")

	sendpooltestCmd.Flags().BytesHexVar(&sendpooltestArgs.address, "address", []byte(nil),
		"fund address used to fund the pools")
}

// loadTest function spawns a few thread which create an account and execute send transactions on the
// rpc server. This command is used to create a simulated load of send transactions
func sendPoolTest(_ *cobra.Command, _ []string) {

	// create new connection to RPC server
	ctx := NewContext()
	ctx.logger.Info("Starting to send Funds to Pools...")

	// create a channel to catch os.Interrupt from a SIGTERM or similar kill signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// create vars for concurrency management
	stopChan := make(chan bool, sendpooltestArgs.threads)
	waiter := sync.WaitGroup{}
	counterChan := make(chan int, sendpooltestArgs.threads)

	// spawn a goroutine to handle sigterm and max transactions
	waiter.Add(1)
	counter := 0
	go handleSigTerm(c, counterChan, stopChan, sendpooltestArgs.threads, sendpooltestArgs.maxTx, &waiter, &counter)

	fullnode := ctx.clCtx.FullNodeClient()
	defer fullnode.Close()
	currencies, err := fullnode.ListCurrencies()
	if err != nil {
		ctx.logger.Error("failed to get currencies", err)
		return
	}
	fmt.Printf("currencies: %#v", currencies)

	var nodeAddress keys.Address
	// get address of the node
	if sendpooltestArgs.address != nil {
		nodeAddress = sendpooltestArgs.address
	} else {
		reply, err := fullnode.NodeAddress()
		if err != nil {
			ctx.logger.Fatal("error getting node address", err)
		}
		nodeAddress = reply.Address
	}

	accs := make([]accounts.Account, 0, 10)
	accReply, err := fullnode.ListAccounts()
	if err != nil {
		ctx.logger.Fatal("failed to get any address")
	}
	accs = accReply.Accounts
	// start threads
	for i := 0; i < sendpooltestArgs.threads; i++ {

		thLogger := ctx.logger.WithPrefix(fmt.Sprintf("thread: %d", i))
		acc := accounts.Account{}
		if len(accs) <= i+1 {
			// create a temp account to send OLT
			accName := fmt.Sprintf("acc_%03d", i)
			pubKey, privKey, err := keys.NewKeyPairFromTendermint() // generate a ed25519 key pair
			if err != nil {
				thLogger.Error(accName, "error generating key from tendermint", err)
			}

			acc, err = accounts.NewAccount(chain.Type(1), accName, &privKey, &pubKey) // create account object
			if err != nil {
				thLogger.Error(accName, "Error initializing account", err)
				return
			}

			thLogger.Infof("creating account %#v", acc)
			reply, err := fullnode.AddAccount(acc)
			if err != nil {
				thLogger.Error(accName, "error creating account", err)
				return
			}

			// add account to wallet
			wallet, err := accounts.NewWalletKeyStore(keyStorePath)
			if err != nil {
				return
			}
			if !wallet.Open(acc.Address(), "pass") {
				return
			}
			err = wallet.Add(acc)
			if err != nil {
				return
			}
			wallet.Close()

			// praccint details
			thLogger.Infof("Created account successfully: %#v", reply.Account)
			ctx.logger.Infof("Address for the account is: %s", acc.Address().Humanize())
		} else {
			acc = accs[i+1]
		}
		waiter.Add(1)
		poolList := []string{"RewardsPool", "ERRORPool"}
		var pool string
		if sendpooltestArgs.randomRecv {
			rand.Seed(time.Now().Unix()) // initialize global pseudo random generator
			pool = poolList[rand.Intn(len(poolList))]
		} else {
			pool = "RewardsPool" // receiver is the temp account
		}
		// start a thread to keep sending transactions after some interval
		go func(stop chan bool) {
			waitDuration := getWaitDuration(sendpooltestArgs.interval)

			for true {

				doSendPoolTransaction(ctx, i, &acc, nodeAddress, pool, currencies.Currencies) // send OLT to temp account
				counterChan <- 1

				select {
				case <-stop:
					waiter.Done()
					return
				default:
					time.Sleep(waitDuration)
				}
			}

		}(stopChan)
	}

	// wait for all threads to close through sigterm; indefinitely
	waiter.Wait()

	// print stats
	fmt.Println("####################################################################")
	fmt.Println("################     Finished Sending to Pool TESt   ###############")
	fmt.Println("####################################################################")
	fmt.Printf("################       Messages sent: % 9d      ###############\n", counter)
	fmt.Println("####################################################################")
}

// doSendTransaction takes in an account and currency object and sends random amounts of coin from the
// node account. It prints any errors to ctx.logger and returns
func doSendPoolTransaction(ctx *Context, threadNo int, acc *accounts.Account, nodeAddress keys.Address, poolName string, currencies balance.Currencies) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("panic in doSendTransaction: thread", threadNo, r)
		}
	}()

	// generate a random amount
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	num := r.Float64() * 10 // amount is a random float between [0, 10)
	// populate send arguments
	sendPoolArgsLocal := SendPoolArguments{}
	sendPoolArgsLocal.Party = []byte(nodeAddress)
	sendPoolArgsLocal.PoolName = poolName
	// set amount and fee
	sendPoolArgsLocal.Amount = strconv.FormatFloat(num, 'f', 6, 64)
	sendPoolArgsLocal.Fee = strconv.FormatFloat(0.0000003, 'f', 9, 64)
	sendPoolArgsLocal.Currency = "OLT"
	sendPoolArgsLocal.Gas = 200030

	// Create message
	fullnode := ctx.clCtx.FullNodeClient()

	req, err := sendPoolArgsLocal.ClientRequest(currencies.GetCurrencySet())
	if err != nil {
		ctx.logger.Error("failed to get request", err)
	}

	reply, err := fullnode.SendToPoolTx(req)
	if err != nil {
		ctx.logger.Error(acc.Name, "error executing SendTx", err)
		return
	}

	// check packet
	packet := reply.RawTx
	if packet == nil {
		ctx.logger.Error(acc.Name, "error in creating new SendTx but server responded with no error")
		return
	}

	// broadcast packet over tendermint
	result, err := ctx.clCtx.BroadcastTxAsync(packet)
	if err != nil {
		ctx.logger.Error(acc.Name, "error in BroadcastTxAsync:", err)
		return
	}

	ctx.logger.Info(acc.Name, "Result Data", "log", string(result.Log))
}
