/*
	Copyright 2017-2018 OneLedger

	Cli to interact with a with the chain.
*/
package main

import (
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
	user string
}

var listargs = &ListArguments{}

func init() {
	RootCmd.AddCommand(accountCmd)

	// Transaction Parameters
	accountCmd.Flags().StringVarP(&listargs.user, "user", "u", "", "user account name")
}

// IssueRequest sends out a sendTx to all of the nodes in the chain
func ListAccount(cmd *cobra.Command, args []string) {

	// TODO: We can't do this, need to be 'light-client' instead...
	node := app.NewApplication()

	if listargs.user != "" {
		Console.Print("Listing Account Details for", listargs.user)
		identity, err := app.FindIdentity(listargs.user)
		if err != 0 {
			app.Log.Error("Not a valid identity", "err", err)
			return
		}
		AccountInfo(node, identity)
		return
	}

	Console.Print("Listing Account Details for all users")
	for _, identity := range node.Accounts.AllAccounts() {
		AccountInfo(node, identity)
	}
}

func AccountInfo(node *app.Application, identity app.Identity) {

	name, err := app.GetAccount(identity)
	if err != 0 {
		app.Log.Error("Invalid Account", "err", err)
		return
	}

	PrintAccount(identity, name)
}

func PrintAccount(identity app.Identity, name string) {
	Console.Print("Identity")
}
