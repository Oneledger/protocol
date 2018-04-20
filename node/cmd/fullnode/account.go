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
	user      string
	account   string
	chainType string
	pubkey    string
	privkey   string
}

var listargs = &ListArguments{}

func init() {
	RootCmd.AddCommand(accountCmd)

	// Operational Parameters
	//sendCmd.Flags().StringVarP(&app.Current.Transport, "transport", "t", "socket", "transport (socket | grpc)")
	//sendCmd.Flags().StringVarP(&app.Current.Address, "address", "a", "tcp://127.0.0.1:46658", "full address")

	// Transaction Parameters
	accountCmd.Flags().StringVarP(&listargs.account, "account", "a", "undefined", "account")
	accountCmd.Flags().StringVarP(&listargs.user, "user", "u", "undefined", "user")
}

// IssueRequest sends out a sendTx to all of the nodes in the chain
func ListAccount(cmd *cobra.Command, args []string) {
	app.Log.Debug("Listing Account Details")
}

// Verify that the account actually has access to the chain in question
func VerifyAccess() {
}
