/*
	Copyright 2017-2018 OneLedger

	Cli to interact with a with the chain.
*/
package main

import (
	"github.com/Oneledger/protocol/node/cmd/shared"
	"github.com/Oneledger/protocol/node/comm"
	"github.com/spf13/cobra"
)

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Register and Identity with the Chain",
	Run:   RegisterIdentity,
}

// Arguments to the command
type RegistrationArguments struct {
	identity string
	account  string
	nodeName string
	pubkey   string
}

var arguments = &RegistrationArguments{}

func init() {
	RootCmd.AddCommand(registerCmd)

	// Transaction Parameters
	registerCmd.Flags().StringVar(&arguments.identity, "identity", "", "User's Identity")
	registerCmd.Flags().StringVar(&arguments.account, "account", "", "User's Default Account")
	registerCmd.Flags().StringVar(&arguments.nodeName, "node", "", "User's Default Node")

	registerCmd.Flags().StringVar(&arguments.pubkey, "pubkey", "", "Specify a public key")
}

func RegisterIdentity(cmd *cobra.Command, args []string) {
	arguments := &shared.RegisterArguments{
		Identity: arguments.identity,
		Account:  arguments.account,
		NodeName: arguments.nodeName,
	}

	register := shared.RegisterIdentityRequest(arguments)

	comm.SDKRequest(register)
}
