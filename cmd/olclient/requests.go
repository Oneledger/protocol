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
	"github.com/spf13/cobra"

	"github.com/Oneledger/protocol/client"
	"github.com/Oneledger/protocol/data/evidence"
)

var requestsCmd = &cobra.Command{
	Use:   "requests",
	Short: "Print out vote requests",
	Run:   GetVoteRequests,
}

type RequestsArguments struct {
	address []byte
}

var rArgs *RequestsArguments = &RequestsArguments{}

func init() {
	EvidencesCmd.AddCommand(requestsCmd)
}

func GetVoteRequests(cmd *cobra.Command, args []string) {
	ctx := NewContext()

	fullnode := ctx.clCtx.FullNodeClient()

	request := client.VoteRequestRequest{}
	vs, err := fullnode.VoteRequests(request)
	if err != nil {
		logger.Fatal("error in getting validator status", err)
	}
	logger.Info("-------------------")
	for i := len(vs.Requests) - 1; i >= 0; i-- {
		printVoteRequest(vs.Requests[i])
	}

}

func printVoteRequest(ar evidence.AllegationRequest) {
	logger.Info("\t ID:", ar.ID)
	logger.Info("\t ReporterAddress:", ar.ReporterAddress)
	logger.Info("\t MaliciousAddress:", ar.MaliciousAddress)
	logger.Info("\t Block height:", ar.BlockHeight)
	logger.Info("\t Status:", evidence.VoteToString(ar.Status))
	logger.Info("\t Proof message:", ar.ProofMsg)

	yesCount := 0
	noCount := 0

	for i := range ar.Votes {
		vote := ar.Votes[i]
		switch vote.Choice {
		case evidence.YES:
			yesCount++
		case evidence.NO:
			noCount++
		}
	}
	logger.Info("\t Yes count:", yesCount)
	logger.Info("\t No count:", noCount)
	logger.Info("-------------------")
}
