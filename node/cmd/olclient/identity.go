/*
	Copyright 2017-2018 OneLedger

	Cli to interact with a with the chain.
*/
package main

import (
	"github.com/Oneledger/protocol/node/action"
	"github.com/Oneledger/protocol/node/cmd/shared"
	"github.com/Oneledger/protocol/node/comm"
	"github.com/Oneledger/protocol/node/id"
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

	nodeName := shared.GetNodeName()

	request := FormatIdentityRequest()
	response := comm.Query("/identity", request)
	if response == nil {
		shared.Console.Error("Node", nodeName, "unavailable")
	}
	if str := comm.IsError(response); str != nil {
		shared.Console.Error("Node", nodeName, str)
	}

	printResponse(nodeName, response)
}

func printResponse(nodeName string, idQuery interface{}) {
	identities := idQuery.([]id.Identity)

	shared.Console.Info("\nOneLedger Identities on", nodeName, ":\n")

	for _, identity := range identities {
		printIdentity(identity)
	}
}

func printIdentity(identity id.Identity) {
	// Right-align fieldnames in console
	name := "             Name:"
	scope := "            Scope:"
	accountKey := "              Key:"
	tendermintAddress := "TendermintAddress:"
	tendermintPubKey := " TendermintPubKey:"

	var scopeOutput string
	if identity.External {
		scopeOutput = "External"
	} else {
		scopeOutput = "Local"
	}

	shared.Console.Info(name, identity.Name)
	shared.Console.Info(scope, scopeOutput)
	shared.Console.Info(accountKey, identity.AccountKey.AsString())
	shared.Console.Info(tendermintAddress, identity.TendermintAddress)
	shared.Console.Info(tendermintPubKey, identity.TendermintPubKey)
	shared.Console.Info()
}
