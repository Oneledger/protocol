/*
	Copyright 2017-2018 OneLedger

	Cli to interact with a with the chain.
*/
package main

import (
	"github.com/Oneledger/protocol/node/log"
	"github.com/Oneledger/protocol/node/transaction"
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
	accountCmd.Flags().StringVar(&account.user, "user", "undefined", "send recipient")
}

// Format the request into a query structure
func FormatRequest() []byte {
	return transaction.Message("User=" + account.user)
}

// IssueRequest sends out a sendTx to all of the nodes in the chain
func CheckAccount(cmd *cobra.Command, args []string) {
	log.Debug("Checking Acccount", "account", account)

	request := FormatRequest()

	// TODO: path was a partial URL path? Need to check to see if that is still required.
	response := Query("/account", request).Response

	log.Debug("Returned Successfully with", "response", string(response.Value))
}
