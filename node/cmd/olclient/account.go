/*
	Copyright 2017-2018 OneLedger

	Cli to interact with a with the chain.
*/
package main

import (
	"github.com/Oneledger/protocol/node/action"
	"github.com/Oneledger/protocol/node/app"
	"github.com/Oneledger/protocol/node/cmd/shared"
	"github.com/Oneledger/protocol/node/comm"
	"github.com/Oneledger/protocol/node/serial"
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
		result, err := serial.Deserialize(response.Response.Value, prototype, serial.CLIENT)
		if err != nil {
			shared.Console.Error("Failed to deserialize AccountQuery")
			shared.Console.Warning("Query failed")
			return
		}
		final := result.(app.AccountQuery)
		printQuery(&final)
	} else {
		shared.Console.Warning("Query Failed")
	}
}

func printQuery(accountQuery *app.AccountQuery) {
	exports := accountQuery.Accounts

	if len(exports) < 1 {
		return
	}

	name := "      Name:"
	balance := "   Balance:"
	accountType := "      Type:"
	accountKey := "       Key:"

	first := true

	for _, export := range exports {
		if first {
			shared.Console.Info("\nAccount(s) on", export.NodeName+":")
			first = false
		}

		shared.Console.Info(name, export.Name)
		shared.Console.Info(accountType, export.Type)
		shared.Console.Info(accountKey, export.AccountKey)
		if export.Type == "OneLedger" {
			shared.Console.Info(balance, export.Balance)
		}
		shared.Console.Info()
	}
}
