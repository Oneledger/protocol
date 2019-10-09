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

	"github.com/Oneledger/protocol/data/balance"

	"github.com/tendermint/tendermint/crypto/ed25519"

	"github.com/Oneledger/protocol/data/accounts"
	"github.com/Oneledger/protocol/data/chain"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/spf13/cobra"
)

// cobra command to loadtest
var loadtestCmd = &cobra.Command{
	Use:   "loadtest",
	Short: "launch load tester",
	Run:   loadTest,
}

// global to load command line args
var loadTestArgs = LoadTestArgs{}

// struct to hold the command-line args
type LoadTestArgs struct {
	threads    int  // no. of threads in the load test; for concurrency
	interval   int  // interval (in milliseconds) between two successive send transactions on a thread
	randomRecv bool // whether to send tokens to a random address every time or no, the default is false
	maxTx      int  // max transactions after which the load test should stop, default is 10000(10k)
}

// init function initializes the loadtest command and attaches a bunch of flag parsers
func init() {
	RootCmd.AddCommand(loadtestCmd)

	loadtestCmd.Flags().IntVar(&loadTestArgs.interval, "interval", 1,
		"interval between successive transactions on a single thread in milliseconds")

	loadtestCmd.Flags().IntVar(&loadTestArgs.threads, "threads", 1,
		"number of threads running")

	loadtestCmd.Flags().BoolVar(&loadTestArgs.randomRecv, "random-receiver", false,
		"whether to randomize the receiver every time, default false")

	loadtestCmd.Flags().IntVar(&loadTestArgs.maxTx, "max-tx", 10000,
		"number of max tx in before the load test stop")
}

// loadTest function spawns a few thread which create an account and execute send transactions on the
// rpc server. This command is used to create a simulated load of send transactions
func loadTest(_ *cobra.Command, _ []string) {

	// create new connection to RPC server
	ctx := NewContext()
	ctx.logger.Info("Starting Loadtest...")

	// create a channel to catch os.Interrupt from a SIGTERM or similar kill signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// create vars for concurrency management
	stopChan := make(chan bool, loadTestArgs.threads)
	waiter := sync.WaitGroup{}
	counterChan := make(chan int, loadTestArgs.threads)

	// spawn a goroutine to handle sigterm and max transactions
	waiter.Add(1)
	counter := 0
	go handleSigTerm(c, counterChan, stopChan, loadTestArgs.threads, loadTestArgs.maxTx, &waiter, &counter)

	fullnode := ctx.clCtx.FullNodeClient()
	defer fullnode.Close()
	currencies, err := fullnode.ListCurrencies()
	if err != nil {
		ctx.logger.Error("failed to get currencies", err)
		return
	}
	fmt.Printf("currencies: %#v", currencies)

	// get address of the node
	nodeAddress, err := fullnode.NodeAddress()
	if err != nil {
		ctx.logger.Fatal("error getting node address", err)
	}

	accs := make([]accounts.Account, 0, 10)
	accReply, err := fullnode.ListAccounts()
	if err != nil {
		ctx.logger.Error("failed to get any address")
	}
	accs = accReply.Accounts
	// start threads
	for i := 0; i < loadTestArgs.threads; i++ {

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

			accRead := reply.Account

			// praccint details
			thLogger.Infof("Created account successfully: %#v", accRead)
			ctx.logger.Infof("Address for the account is: %s", acc.Address().Humanize())
		} else {
			acc = accs[i+1]
		}
		waiter.Add(1)
		// start a thread to keep sending transactions after some interval
		go func(stop chan bool) {
			waitDuration := getWaitDuration(loadTestArgs.interval)

			for true {
				doSendTransaction(ctx, i, &acc, nodeAddress, loadTestArgs.randomRecv, currencies.Currencies) // send OLT to temp account
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
	fmt.Println("################        Terminating load test        ###############")
	fmt.Println("####################################################################")
	fmt.Printf("################       Messages sent: % 9d      ###############\n", counter)
	fmt.Println("####################################################################")
}

// doSendTransaction takes in an account and currency object and sends random amounts of coin from the
// node account. It prints any errors to ctx.logger and returns
func doSendTransaction(ctx *Context, threadNo int, acc *accounts.Account, nodeAddress keys.Address, randomRev bool, currencies balance.Currencies) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("panic in doSendTransaction: thread", threadNo, r)
		}
	}()

	// generate a random amount
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	num := r.Float64() * 10 // amount is a random float between [0, 10)
	// populate send arguments
	sendArgsLocal := SendArguments{}
	sendArgsLocal.Party = []byte(nodeAddress)
	if randomRev {
		recv := ed25519.GenPrivKey().PubKey().Address()
		sendArgsLocal.CounterParty = []byte(recv)
	} else {
		sendArgsLocal.CounterParty = []byte(acc.Address()) // receiver is the temp account
	}

	// set amount and fee
	sendArgsLocal.Amount = strconv.FormatFloat(num, 'f', 6, 64)
	sendArgsLocal.Fee = strconv.FormatFloat(0.0000003, 'f', 9, 64)
	sendArgsLocal.Currency = "OLT"
	sendArgsLocal.Gas = 200030

	// Create message
	fullnode := ctx.clCtx.FullNodeClient()

	req, err := sendArgsLocal.ClientRequest(currencies.GetCurrencySet())
	if err != nil {
		ctx.logger.Error("failed to get request", err)
	}

	reply, err := fullnode.SendTx(req)
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

func getWaitDuration(interval int) time.Duration {
	return time.Millisecond * time.Duration(interval)
}

// handleSigTerm keeps a count of messages sent and if the maximum number of transactions is reached it stops
// all threads and proceeds to shut down the main thread. If it catches a SIGTERM or a CTRL C it similarly shuts down
// gracefully. This function is blocking and is called as a go routine.
func handleSigTerm(c chan os.Signal, counterChan chan int, stopChan chan bool,
	n int, maxTx int, waiter *sync.WaitGroup, cnt *int) {

	// indefinite loop listens over the counter and os.Signal for interrupt signal
	for true {
		select {
		case <-c:
			// signal the goroutines to stop
			for i := 0; i < n; i++ {
				stopChan <- true
			}
			// wait for the goroutines to stop
			time.Sleep(time.Second)

			waiter.Done()

		case <-counterChan:
			// increment counter
			*cnt++

			// send shutdown signal if max no. of transactions is reached
			if *cnt >= maxTx {
				c <- os.Interrupt
			}
		}
	}
}
