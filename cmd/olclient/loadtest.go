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
	"os/exec"
	"os/signal"
	"strconv"
	"sync"
	"time"

	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/crypto/ed25519"

	"github.com/Oneledger/protocol/data/accounts"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/chain"
	"github.com/Oneledger/protocol/data/keys"
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
	address    []byte
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

	loadtestCmd.Flags().BytesHexVar(&loadTestArgs.address, "address", []byte(nil),
		"fund address that loadtest uses")
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
	//fmt.Printf("currencies SFS: %#v", currencies)

	var nodeAddress keys.Address
	// get address of the node
	if loadTestArgs.address != nil {
		nodeAddress = loadTestArgs.address
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

			//thLogger.Infof("creating account %#v", acc)
			reply, err := fullnode.AddAccount(acc)
			if err != nil {
				thLogger.Error(accName, "error creating account", err)
				return
			}

			accRead := reply.Account

			// praccint details
			thLogger.Infof("Created account successfully: %#s", accRead.Address().Humanize())
			//ctx.logger.Infof("Address for the New account is: %s", acc.Address().Humanize())
		} else {
			acc = accs[i+1]
		}
		waiter.Add(1)
		// start a thread to keep sending transactions after some interval
		go func(stop chan bool) {
			waitDuration := getWaitDuration(loadTestArgs.interval)

			for true {
				// send OLT to temp account
				//doOnsTrasactions()
				doSendTransaction(ctx, i, &acc, nodeAddress, loadTestArgs.randomRecv, currencies.Currencies)

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

		go func(stop chan bool) {
			waitDuration := getWaitDuration(loadTestArgs.interval)
			for true {
				// Do ONS Create Domain , Create Sub Domain ,Send To Domain
				doOnsTrasactions()
				counterChan <- 3
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

func doOnsTrasactions() {
	//fmt.Println("Running python script")
	cmd := exec.Command("python", "ons/create_sub_domain.py")
	out, err := cmd.Output()

	if err != nil {
		fmt.Println(string(out))
		fmt.Println("Error in python script : ", err.Error())
		return
	}

	//fmt.Println(string(out))

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
	_, err = ctx.clCtx.BroadcastTxAsync(packet)
	if err != nil {
		ctx.logger.Error(acc.Name, "error in BroadcastTxAsync:", err)
		return
	}

	//ctx.logger.Info(acc.Name, "Broadcasted TX")
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

		case txsentiniteration := <-counterChan:
			// increment counter
			*cnt += txsentiniteration

			// send shutdown signal if max no. of transactions is reached
			if *cnt >= maxTx {
				c <- os.Interrupt
			}
		}
	}
}

//I[2020-03-26T15:24:10-04:00] app: Loadtest metric height=10, total_tx =6463, tx/b=646, blktime=1.036688 , tps=889.922322
//I[2020-03-26T15:24:27-04:00] app: Loadtest metric height=30, total_tx =24799, tx/b=826, blktime=1.009512 , tps=909.644127
//I[2020-03-26T15:24:50-04:00] app: Loadtest metric height=50, total_tx =46567, tx/b=931, blktime=1.005464 , tps=985.297080
//I[2020-03-26T15:25:11-04:00] app: Loadtest metric height=70, total_tx =65776, tx/b=939, blktime=1.011063 , tps=970.915257
//I[2020-03-26T15:25:32-04:00] app: Loadtest metric height=90, total_tx =86403, tx/b=960, blktime=1.019297 , tps=974.279540
//I[2020-03-26T15:25:53-04:00] app: Loadtest metric height=110, total_tx =107656, tx/b=978, blktime=1.033474 , tps=973.497410
//I[2020-03-26T15:26:15-04:00] app: Loadtest metric height=130, total_tx =128185, tx/b=986, blktime=1.030771 , tps=979.161571
//I[2020-03-26T15:26:40-04:00] app: Loadtest metric height=150, total_tx =151940, tx/b=1012, blktime=1.060506 , tps=974.602512
//I[2020-03-26T15:27:02-04:00] app: Loadtest metric height=170, total_tx =172922, tx/b=1017, blktime=1.066872 , tps=970.530044
//I[2020-03-26T15:27:26-04:00] app: Loadtest metric height=190, total_tx =196138, tx/b=1032, blktime=1.080966 , tps=970.279690
