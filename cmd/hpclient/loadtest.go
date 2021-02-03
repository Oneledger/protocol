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
	"crypto/sha256"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/Oneledger/protocol/log"

	"github.com/Oneledger/protocol/data/accounts"
	"github.com/spf13/cobra"
)

// cobra command to loadtest
var loadtestCmd = &cobra.Command{
	Use:   "loadtest",
	Short: "launch load tester",
	Run:   loadTest,
}
var logger = log.NewLoggerWithPrefix(os.Stdout, "hpclient")

// global to load command line args
var loadTestArgs = LoadTestArgs{}

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func String(length int) string {
	return StringWithCharset(length, charset)
}

// struct to hold the command-line args
type LoadTestArgs struct {
	threads            int // no. of threads in the load test; for concurrency
	interval           int // interval (in milliseconds) between two successive send transactions on a thread
	maxTx              int // max transactions after which the load test should stop, default is 10000(10k)
	noOfNodes          int
	noOfValidatorNodes int
	superAdminKeyPath  string
}

// init function initializes the loadtest command and attaches a bunch of flag parsers
func init() {
	RootCmd.AddCommand(loadtestCmd)

	loadtestCmd.Flags().IntVar(&loadTestArgs.interval, "interval", 1,
		"interval between successive transactions on a single thread in milliseconds")

	loadtestCmd.Flags().IntVar(&loadTestArgs.threads, "threads", 1,
		"number of threads running")

	loadtestCmd.Flags().IntVar(&loadTestArgs.maxTx, "max-tx", 10000,
		"number of max tx in before the load test stop")
	loadtestCmd.Flags().IntVar(&loadTestArgs.noOfNodes, "nodes", 9,
		"Number of nodes to run the loadtest on")

	loadtestCmd.Flags().IntVar(&loadTestArgs.noOfValidatorNodes, "noodvalidatornodes", 4,
		"Number of nodes to run the loadtest on")

	loadtestCmd.Flags().StringVar(&loadTestArgs.superAdminKeyPath, "superadmins", "",
		"Folder which contains the super-admin Keys")
}

// loadTest function spawns a few thread which create an account and execute send transactions on the
// rpc server. This command is used to create a simulated load of send transactions
func loadTest(_ *cobra.Command, _ []string) {
	// Get Super Admins
	walletAdmin, err := accounts.NewWalletKeyStore(filepath.Clean(loadTestArgs.superAdminKeyPath))
	if err != nil {
		return
	}
	addresses, err := walletAdmin.ListAddresses()
	if err != nil {
		logger.Fatal(err)
	}
	superAdminNames := make([]string, len(addresses))
	walletAdmin.Open(addresses[0], "1234")
	for i, addr := range addresses {
		_, err := walletAdmin.GetAccount(addr)
		if err != nil {
			logger.Fatal(err)
		}
		sum := sha256.Sum256([]byte(addr.Humanize()))
		userID := fmt.Sprintf("%x", sum)
		superAdminNames[i] = userID
	}
	walletAdmin.Close()
	// reduce by 1 to get index

	noofSuperAdmins := len(addresses) - 1
	loadTestArgs.noOfNodes = loadTestArgs.noOfNodes - 1

	var nodelist []LoadTestNode
	defer Cleanup(&nodelist)

	loadTestController := NewController(loadTestArgs.threads, getWaitDuration(loadTestArgs.interval))
	loadTestController.AddExecutorFunction(CreateHospitalAdminRequest, nil)
	loadTestController.AddExecutorFunction(CreateScreenerAdminRequest, nil)
	loadTestController.AddExecutorFunction(AddTestInfoRequest, AddTestInfoInit)
	loadTestController.AddExecutorFunction(ReadTestInfoRequest, ReadTestInfoInit)

	// create a channel to catch os.Interrupt from a SIGTERM or similar kill signal
	signal.Notify(loadTestController.osChan, os.Interrupt)

	// spawn a goroutine to handle sigterm and max transactions
	counter := 0
	go handleSigTerm(loadTestController.osChan, loadTestController.counterChan, loadTestController.stopChan, loadTestArgs.threads, loadTestArgs.maxTx, &loadTestController.waiter, &counter)
	loadTestController.waiter.Add(1)

	// init threads
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < loadTestArgs.threads; i++ {
		//Prepare a loadtest Node for each thread
		//Always pick a non validator node
		nodenumber := rand.Intn(loadTestArgs.noOfNodes-loadTestArgs.noOfValidatorNodes+1) + loadTestArgs.noOfValidatorNodes
		//Pick a superadmin key
		superadminnumber := rand.Intn(noofSuperAdmins-0+1) + 0
		nodeRoot := filepath.Join(rootArgs.rootDir, strconv.Itoa(nodenumber)+"-Node")

		node := NewLoadTestNode(addresses[superadminnumber], nodeRoot, superAdminNames[superadminnumber])
		//Keep Track of generated nodes for cleanup later
		nodelist = append(nodelist, node)
		// start a thread to keep sending transactions after some interval
		loadTestController.InitThread(node, i)
	}

	// run threads
	loadTestController.functionindex = 0
	for i := 0; i < loadTestArgs.threads; i++ {
		node := nodelist[i]
		loadTestController.AddThread(node, i)
	}

	// wait for all threads to close through sigterm; indefinitely
	loadTestController.waiter.Wait()

	fmt.Printf("################       Messages sent: % 9d      ###############\n", counter)
}

func Cleanup(nodelist *[]LoadTestNode) {
	for _, node := range *nodelist {
		err := os.RemoveAll(node.keypath)
		if err != nil {
			logger.Info("Unable to delete temp directory : ", node.keypath)
		}
	}
}

func getWaitDuration(interval int) time.Duration {
	return time.Millisecond * time.Duration(interval)
}

func handleSigTerm(c chan os.Signal, counterChan chan int, stopChan chan bool,
	n int, maxTx int, waiter *sync.WaitGroup, cnt *int) {

	// indefinite loop listens over the counter and os.Signal for interrupt signal
	for true {
		select {
		//c Channel populated when we reach maxTX
		case <-c:
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
			if *cnt == maxTx {
				c <- os.Interrupt
			}
		}
	}
}
