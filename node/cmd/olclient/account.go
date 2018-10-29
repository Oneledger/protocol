/*
	Copyright 2017-2018 OneLedger

	Cli to interact with a with the chain.
*/
package main

import (
	"github.com/Oneledger/protocol/node/action"
	"github.com/Oneledger/protocol/node/cmd/shared"
	"github.com/Oneledger/protocol/node/comm"
	"github.com/Oneledger/protocol/node/data"
	"github.com/Oneledger/protocol/node/id"
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

	// Blank user implies all users
	return action.Message("Account=" + account.user)
}

// IssueRequest sends out a sendTx to all of the nodes in the chain
func CheckAccount(cmd *cobra.Command, args []string) {
	//log.Debug("Checking Account", "account", account)
	request := FormatRequest()
	response := comm.Query("/account", request)
	if response != nil {
		printQuery(response)
	} else {
		shared.Console.Warning("No Response from Node for:", string(request))
	}
}

func printQuery(accountQuery interface{}) {
	nodeName := shared.GetNodeName()

	accounts := accountQuery.([]id.Account)

	if len(accounts) < 1 {
		shared.Console.Info("No Accounts on", nodeName)
		return
	}

	name := "      Name:"
	balance := "   Balance:"
	accountType := "      Type:"
	accountKey := "       Key:"

	first := true

	for _, account := range accounts {
		if first {
			shared.Console.Info("\nAccount(s) on", nodeName+":")
			first = false
		}

		shared.Console.Info(name, account.Name())
		shared.Console.Info(accountType, account.Chain().String())
		shared.Console.Info(accountKey, account.AccountKey().AsString())

		if account.Chain() == data.ONELEDGER {
			value := shared.GetBalance(account.AccountKey())
			shared.Console.Info(balance, value.AsString())
		}

		shared.Console.Info()
	}
}
