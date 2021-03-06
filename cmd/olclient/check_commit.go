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
	hash  string
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

	// Transaction Parameters
	checkCmd.Flags().StringVar(&checkCommit.hash, "hash", "", "tx hash")
	checkCmd.Flags().BoolVar(&checkCommit.prove, "prove", true, "prove")
}

// IssueRequest sends out a sendTx to all of the nodes in the chain
func CheckTransaction(cmd *cobra.Command, args []string) {
	Ctx := NewContext()
	result, b := checkTransactionResult(Ctx, checkCommit.hash, checkCommit.prove)
	if result == nil || b == false {
		fmt.Print("Tx has not been commited yet, please try again later")
		return
	}
	fmt.Println("Hash:", result.Hash)
	fmt.Println("Height:", result.Height)
	fmt.Println("Index:", result.Index)
	fmt.Println("Proof:", result.Proof)
	fmt.Println("TX Result:", result.TxResult.Code)
}
