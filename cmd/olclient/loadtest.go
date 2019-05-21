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

	"github.com/tendermint/tendermint/crypto/ed25519"

	"github.com/Oneledger/protocol/serialize"

	"github.com/Oneledger/protocol/action"

	"github.com/Oneledger/protocol/client"
	"github.com/Oneledger/protocol/data"
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
	counterChan := make(chan int)

	// spawn a goroutine to handle sigterm and max transactions
	waiter.Add(1)
	go handleSigTerm(c, counterChan, stopChan, loadTestArgs.threads, loadTestArgs.maxTx, &waiter)

	// get address of the node
	req := data.NewRequestFromData("server.NodeAddress", nil)
	resp := &data.Response{}
	err := ctx.clCtx.Query("server.NodeAddress", *req, resp)
	if err != nil {
		ctx.logger.Fatal("error getting node address", err)
	}

	var nodeAddress = keys.Address(resp.Data)

	// start threads
	for i := 0; i < loadTestArgs.threads; i++ {

		thLogger := ctx.logger.WithPrefix(fmt.Sprintf("thread: %d", i))

		// create a temp account to send OLT
		accName := fmt.Sprintf("acc_%03d", i)
		pubKey, privKey, err := keys.NewKeyPairFromTendermint() // generate a ed25519 key pair
		if err != nil {
			thLogger.Error(accName, "error generating key from tendermint", err)
		}

		acc, err := accounts.NewAccount(chain.Type(1), accName, &privKey, &pubKey) // create account object
		if err != nil {
			thLogger.Error(accName, "Error initializing account", err)
			return
		}

		thLogger.Infof("creating account %#v", acc)
		resp := &data.Response{}
		err = ctx.clCtx.Query("server.AddAccount", acc, resp) // create account on node
		if err != nil {
			thLogger.Error(accName, "error creating account", err)
			return
		}

		var accRead = &accounts.Account{}
		err = serialize.GetSerializer(serialize.CLIENT).Deserialize(resp.Data, &accRead)
		if err != nil {
			thLogger.Error("error de-serializing account data", err)
			return
		}

		// print details
		thLogger.Infof("Created account successfully: %#v", accRead)
		ctx.logger.Infof("Address for the account is: %s", acc.Address().Humanize())

		waiter.Add(1)
		// start a thread to keep sending transactions after some interval
		go func(stop chan bool) {
			waitDuration := getWaitDuration(loadTestArgs.interval)

			for true {
				doSendTransaction(ctx, i, &acc, nodeAddress, loadTestArgs.randomRecv) // send OLT to temp account
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
}

// doSendTransaction takes in an account and currency object and sends random amounts of coin from the
// node account. It prints any errors to ctx.logger and returns
func doSendTransaction(ctx *Context, threadNo int, acc *accounts.Account, nodeAddress keys.Address, randomRev bool) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("panic in doSendTransaction: thread", threadNo, r)
		}
	}()

	// generate a random amount
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	amt := r.Float64() * 10.0 // amount is a random float between [0, 10)

	// populate send arguments
	sendArgsLocal := client.SendArguments{}
	sendArgsLocal.Party = []byte(nodeAddress)
	if randomRev {
		recv := ed25519.GenPrivKey().PubKey().Address()
		sendArgsLocal.CounterParty = []byte(recv)
	} else {
		sendArgsLocal.CounterParty = []byte(acc.Address()) // receiver is the temp account
	}

	// set amount and fee
	sendArgsLocal.Amount = *action.NewAmount("OLT", strconv.FormatFloat(amt, 'f', 0, 64))
	sendArgsLocal.Fee = *action.NewAmount("OLT", strconv.FormatFloat(amt/100, 'f', 0, 64))

	// Create message
	resp := &data.Response{}
	err := ctx.clCtx.Query("server.SendTx", sendArgsLocal, resp)
	if err != nil {
		ctx.logger.Error(acc.Name, "error executing SendTx", err)
		return
	}

	// check packet
	packet := resp.Data
	if packet == nil {
		ctx.logger.Error(acc.Name, "Error in sending ", resp.ErrorMsg)
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
	n int, maxTx int, waiter *sync.WaitGroup) {

	// keeps a running count of messages sent
	msgCounter := 0

	// indefinite loop listens over the counter and os.Signal for interrupt signal
	for true {
		select {
		case sig := <-c:
			// signal the goroutines to stop
			for i := 0; i < n; i++ {
				stopChan <- true
			}
			// wait for the goroutines to stop
			time.Sleep(time.Second)

			// print stats
			fmt.Println("####################################################################")
			fmt.Println("################        Terminating load test        ###############")
			fmt.Println("####################################################################")
			fmt.Printf("################       Messages sent: % 9d      ###############\n", msgCounter)
			fmt.Println("####################################################################")

			sig.String()

			waiter.Done()

		case <-counterChan:
			// increment counter
			msgCounter++

			// send shutdown signal if max no. of transactions is reached
			if msgCounter >= maxTx {
				c <- os.Interrupt
			}
		}
	}
}
