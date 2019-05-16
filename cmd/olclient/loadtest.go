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

const NODE_ADDRESS = "1234123412341234213412341234"

var loadtestCmd = &cobra.Command{
	Use:   "loadtest",
	Short: "launch load tester",
	Run:   LoadTest,
}

var loadTestArgs = LoadTestArgs{}

type LoadTestArgs struct {
	threads    int
	interval   int
	randomRecv bool
}

func init() {
	RootCmd.AddCommand(loadtestCmd)

	loadtestCmd.Flags().IntVar(&loadTestArgs.interval, "interval",
		1, "interval between successive transactions on a single thread in milliseconds")
	loadtestCmd.Flags().IntVar(&loadTestArgs.threads, "threads",
		1, "number of threads running")
	loadtestCmd.Flags().BoolVar(&loadTestArgs.randomRecv, "random-receiver",
		true, "whether to randomize the receiver everytime")
}

func LoadTest(cmd *cobra.Command, args []string) {

	ctx := NewContext()
	ctx.logger.Debug("Starting Loadtest")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	stopChan := make(chan bool, loadTestArgs.threads)
	waiter := sync.WaitGroup{}
	counterChan := make(chan int, 10000)

	waiter.Add(1)
	go handleSigTerm(ctx, c, counterChan, stopChan, loadTestArgs.threads, &waiter)

	req := data.NewRequestFromData("server.NodeAddress", nil)
	resp := &data.Response{}
	err := ctx.clCtx.Query("server.NodeAddress", *req, resp)
	if err != nil {
		ctx.logger.Fatal("error getting node address", err)
	}

	var nodeAddress = keys.Address(resp.Data)

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

		var accRead = &accounts.Account{}
		serialize.GetSerializer(serialize.CLIENT).Deserialize(resp.Data, &accRead)

		// print details
		ctx.logger.Infof("Created account successfully: %#v", accRead)
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

	recv := ed25519.GenPrivKey().PubKey().Address()

	// populate send arguments
	sendArgsLocal := client.SendArguments{}
	sendArgsLocal.Party = []byte(nodeAddress)
	if randomRev {
		sendArgsLocal.CounterParty = []byte(recv)
	} else {
		sendArgsLocal.CounterParty = []byte(acc.Address()) // receiver is the temp account
	}

	sendArgsLocal.Amount = action.Amount{"OLT", strconv.FormatFloat(amt, 'f', 0, 64)} // make coin from currency object
	sendArgsLocal.Fee = action.Amount{"OLT", strconv.FormatFloat(amt/100, 'f', 0, 64)}

	// Create message
	resp := &data.Response{}
	err := ctx.clCtx.Query("server.SendTx", sendArgsLocal, resp)
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

func getWaitDuration(interval int) time.Duration {
	return time.Millisecond * time.Duration(interval)
}

func handleSigTerm(ctx *Context, c chan os.Signal, counterChan chan int, stopChan chan bool, n int, waiter *sync.WaitGroup) {
	// keeps a running count of messages sent
	msgCounter := 0

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
			fmt.Println("################	Terminating load test	####################")
			fmt.Println("####################################################################")
			fmt.Printf("################	Messages sent: %09d-----	############", msgCounter)

			sig.String()

			waiter.Done()

		case <-counterChan:
			msgCounter++
		}
	}
}
