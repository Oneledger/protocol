/*
	Copyright 2017-2018 OneLedger

	Cli to interact with a with the chain.
*/
package main

import (
	"fmt"

	"github.com/Oneledger/protocol/node/action"
	"github.com/Oneledger/protocol/node/app"
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

	request := FormatIdentityRequest()
	response := comm.Query("/identity", request)
	if response != nil {
		var prototype app.IdentityQuery
		result, err := comm.Deserialize(response.Response.Value, &prototype)
		if err != nil {
			shared.Console.Error("Failed to deserialize IdentityQuery:")
			return
		}
		printResponse(result.(*app.IdentityQuery))
	}
}

func printResponse(idQuery *app.IdentityQuery) {
	shared.Console.Info("\nCheckIdentity Response:\n")

	for _, identity := range idQuery.Identities {
		printIdentity(&identity)
	}
}

func printIdentity(export *id.IdentityExport) {
	// Right-align fieldnames in console
	name := "      Name:"
	scope := "     Scope:"
	accountKey := "AccountKey:"

	var scopeOutput string
	if export.External {
		scopeOutput = "External"
	} else {
		scopeOutput = "Local"
	}

	shared.Console.Info(fmt.Sprintf(name+" %s", export.Name))
	shared.Console.Info(fmt.Sprintf(scope+" %s", scopeOutput))
	shared.Console.Info(fmt.Sprintf(accountKey+" %s", export.AccountKey))
	shared.Console.Info()
}
