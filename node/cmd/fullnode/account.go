/*
	Copyright 2017-2018 OneLedger

	Cli to interact with a with the chain.
*/
package main

import (
	"fmt"

	"github.com/Oneledger/prototype/node/app"
	"github.com/spf13/cobra"
)

var accountCmd = &cobra.Command{
	Use:   "account",
	Short: "List out account details",
	Run:   ListAccount,
}

// Arguments to the command
type ListArguments struct {
	identity string
}

var listargs = &ListArguments{}

func init() {
	RootCmd.AddCommand(accountCmd)

	// Transaction Parameters
	accountCmd.Flags().StringVarP(&listargs.identity, "identity", "u", "undefined", "user account name")
}

// IssueRequest sends out a sendTx to all of the nodes in the chain
func ListAccount(cmd *cobra.Command, args []string) {
	app.Log.Debug("Listing Account Details")

	identity, err := app.FindIdentity(listargs.identity)
	if err != 0 {
		app.Log.Error("Not a valid identity", "err", err)
		return
	}

	account, err := app.GetAccount(identity)
	if err != 0 {
		app.Log.Error("Invalid Account", "err", err)
		return
	}

	PrintAccount(identity, account)
}

func PrintAccount(identity app.Identity, account *app.Account) {
	fmt.Println("Identity")
}
