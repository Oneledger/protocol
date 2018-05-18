/*
	Copyright 2017-2018 OneLedger

	Gets the account information, this is a node operation (and won't run if a node already exists)
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
	user string
}

var listargs = &ListArguments{}

func init() {
	RootCmd.AddCommand(accountCmd)

	// Transaction Parameters
	accountCmd.Flags().StringVarP(&listargs.user, "user", "u", "", "user account name")
}

// IssueRequest sends out a sendTx to all of the nodes in the chain
func ListIdentities(cmd *cobra.Command, args []string) {

	// TODO: We can't do this, need to be 'light-client' instead...
	node := app.NewApplication()

	if listargs.user != "" {
		Console.Print("Listing Account Details for", listargs.user)
		identity, err := node.Identities.FindIdentity(listargs.user)
		if err != 0 {
			log.Error("Not a valid identity", "err", err)
			return
		}
		IdentityInfo(node, identity)
		return
	}

	Console.Print("Listing Account Details for all users")
	for _, identity := range node.Identities.AllIdentities() {
		IdentityInfo(node, &identity)
	}
}

func IdentityInfo(node *app.Application, identity *id.Identity) {
	Console.Print("Identity")
}
