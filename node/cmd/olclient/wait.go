/*
	Copyright 2017-2018 OneLedger

	Cli to interact with a with the chain.
*/
package main

import (
	"time"

	"github.com/Oneledger/protocol/node/action"
	"github.com/Oneledger/protocol/node/chains/bitcoin"
	"github.com/Oneledger/protocol/node/comm"
	"github.com/Oneledger/protocol/node/id"
	"github.com/Oneledger/protocol/node/log"
	"github.com/spf13/cobra"
)

var waitCmd = &cobra.Command{
	Use:   "wait",
	Short: "Wait for something to happen",
	Run:   Wait,
}

type WaitArgs struct {
	completed  string
	parties    []string
	identities []string
}

var waitargs WaitArgs = WaitArgs{}

func init() {
	RootCmd.AddCommand(waitCmd)

	waitCmd.Flags().StringVar(&waitargs.completed, "completed", "", "wait for completion")
	waitCmd.Flags().StringArrayVarP(&waitargs.parties, "party", "p", waitargs.parties, "recipient")
	waitCmd.Flags().StringArrayVarP(&waitargs.identities, "identity", "i", waitargs.identities, "recipient")
}

func Wait(cmd *cobra.Command, args []string) {
	switch waitargs.completed {

	case "swap":
		WaitForSwap()

	case "identity":
		WaitForIdentity(args)
	}
}

func IdentityExists(name string) bool {
	request := action.Message("Identity=" + name)
	response := comm.Query("/accountKey", request)
	if response == nil {
		return false
	}
	switch response.(type) {
	case id.AccountKey:
		return true

	case string:
		return false
	}
	return false
}

// Wait for a set of identities to get created
func WaitForIdentity(args []string) {
	// Left over args get passed in, allows for '--identity x y z'
	waitargs.identities = append(waitargs.identities, args...)

	size := len(waitargs.identities)

	var found []bool
	found = make([]bool, size)

	log.Debug("Waiting for", "identities", waitargs.identities)

	for {
		count := 0
		for i := 0; i < size; i++ {
			if found[i] {
				count++
				continue
			}
			if IdentityExists(waitargs.identities[i]) {
				found[i] = true
			}
		}
		if count == size {
			log.Info("All Identities have been created")
			return
		}
		time.Sleep(time.Second)
	}
}

// Wait for a swap (really just 60 secs of Bitcoin block generation)
func WaitForSwap() {
	log.Debug("Waiting")

	cli := bitcoin.GetBtcClient("127.0.0.1:18833")

	stop := bitcoin.ScheduleBlockGeneration(*cli, 1)
	time.Sleep(60 * time.Second)
	bitcoin.StopBlockGeneration(stop)
}
