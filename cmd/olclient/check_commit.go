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
	"github.com/spf13/cobra"
)

type CheckCommit struct {
	hash string
	prove bool
}

var checkCmd = &cobra.Command{
	Use:   "check_commit",
	Short: "Check the result of tx",
	Run:   CheckTransaction,
}

var checkCommit = &CheckCommit{}

func init() {
	RootCmd.AddCommand(checkCmd)

	// TODO: I want to have a default account?
	// Transaction Parameters
	checkCmd.Flags().StringVar(&checkCommit.hash, "hash", "", "tx hash")
	checkCmd.Flags().BoolVar(&checkCommit.prove, "prove", true, "prove")
}

// IssueRequest sends out a sendTx to all of the nodes in the chain
func CheckTransaction(cmd *cobra.Command, args []string) {
	Ctx := NewContext()
	fullnode := Ctx.clCtx.FullNodeClient()
	result, err := fullnode.CheckCommitResult(checkCommit.hash, checkCommit.prove)
	if err != nil {
		fmt.Print(err)
	}
	fmt.Print(result.Result)

}

