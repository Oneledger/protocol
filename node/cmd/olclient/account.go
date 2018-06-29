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
	"github.com/spf13/cobra"
)

var accountCmd = &cobra.Command{
	Use:   "account",
	Short: "Check account status",
	Run:   CheckAccount,
}

// TODO: typing should be way better, see if cobr can help with this...
type AccountArguments struct {
	user string
}

var account *AccountArguments = &AccountArguments{}

func init() {
	RootCmd.AddCommand(accountCmd)

	// TODO: I want to have a default account?
	// Transaction Parameters
	accountCmd.Flags().StringVar(&account.user, "identity", "", "identity name")
}

// Format the request into a query structure
func FormatRequest() []byte {
	return action.Message("Account=" + account.user)
}

// IssueRequest sends out a sendTx to all of the nodes in the chain
func CheckAccount(cmd *cobra.Command, args []string) {
	//log.Debug("Checking Account", "account", account)
	request := FormatRequest()
	response := comm.Query("/account", request)
	if response != nil {
		// var accountQuery app.AccountQuery
		var prototype app.AccountQuery
		result, err := comm.Deserialize(response.Response.Value, &prototype)
		if err != nil {
			shared.Console.Error("Failed to deserialize AccountQuery")
			shared.Console.Warning("Query failed")
			return
		}
		printQuery(result.(*app.AccountQuery))
	} else {
		shared.Console.Warning("Query Failed")
	}
}

func printQuery(accountQuery *app.AccountQuery) {
	exports := accountQuery.Accounts

	name := "      Name:"
	balance := "   Balance:"
	accountType := "      Type:"
	accountKey := "AccountKey:"
	nodeName := "  NodeName:"

	shared.Console.Info("\nCheckAccount response: \n")

	for _, export := range exports {
		shared.Console.Info(fmt.Sprintf(nodeName+" %s", export.NodeName))
		shared.Console.Info(fmt.Sprintf(name+" %s", export.Name))
		shared.Console.Info(fmt.Sprintf(accountType+" %s", export.Type))
		shared.Console.Info(fmt.Sprintf(accountKey+" %s", export.AccountKey))
		shared.Console.Info(fmt.Sprintf(balance+" %s", export.Balance))
		shared.Console.Info("\n")
	}
}
