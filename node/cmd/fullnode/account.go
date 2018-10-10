/*
	Copyright 2017-2018 OneLedger

	Gets the account information, this is a node operation (and won't run if a node already is already running)
*/
package main

import (
	"github.com/Oneledger/protocol/node/app"
	"github.com/Oneledger/protocol/node/id"
	"github.com/Oneledger/protocol/node/log"
	"github.com/spf13/cobra"
)

var accountCmd = &cobra.Command{
	Use:   "account",
	Short: "List out account details",
	Run:   ListIdentities,
}

// Arguments to the command
type ListArguments struct {
	identity string
}

var listargs = &ListArguments{}

// Setup the command in Cobra
func init() {
	RootCmd.AddCommand(accountCmd)

	// Transaction Parameters
	accountCmd.Flags().StringVar(&listargs.identity, "identity", "", "user account name")
}

// IssueRequest sends out a sendTx to all of the nodes in the chain
func ListIdentities(cmd *cobra.Command, args []string) {

	node := app.NewApplication()

	if listargs.identity != "" {
		Console.Print("Listing Account Details for", listargs.identity)
		id, err := node.Identities.FindName(listargs.identity)
		if err != 0 {
			log.Error("Not a valid identity", "status", err)
			return
		}
		if id != nil {
			IdentityInfo(node, id)
		} else {
			Console.Print("Unknown Account")
		}
		return
	}

	Console.Print("Listing Account Details for all users")
	for _, id := range node.Identities.FindAll() {
		IdentityInfo(node, id)
	}
}

func IdentityInfo(node *app.Application, id *id.Identity) {
	Console.Print("Identity " + id.Name)
	// TODO: This out the know active accounts.
}
