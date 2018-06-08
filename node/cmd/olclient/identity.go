/*
	Copyright 2017-2018 OneLedger

	Cli to interact with a with the chain.
*/
package main

import (
	"github.com/Oneledger/protocol/node/action"
	"github.com/Oneledger/protocol/node/comm"
	"github.com/Oneledger/protocol/node/log"
	"github.com/spf13/cobra"
)

var identityCmd = &cobra.Command{
	Use:   "identity",
	Short: "Check identity status",
	Run:   CheckIdentity,
}

// TODO: typing should be way better, see if cobr can help with this...
type IdentityArguments struct {
	identity string
}

var ident *IdentityArguments = &IdentityArguments{}

func init() {
	RootCmd.AddCommand(identityCmd)

	// TODO: I want to have a default account?
	// Transaction Parameters
	identityCmd.Flags().StringVar(&ident.identity, "identity", "", "identity name")
}

// Format the request into a query structure
func FormatIdentityRequest() []byte {
	return action.Message("Identity=" + ident.identity)
}

// IssueRequest sends out a sendTx to all of the nodes in the chain
func CheckIdentity(cmd *cobra.Command, args []string) {
	log.Debug("Checking Identity", "identity", ident)

	request := FormatIdentityRequest()
	response := comm.Query("/identity", request).Response

	log.Debug("Returned Successfully with", "response", string(response.Value))
}
