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
	fee      float64
}

var arguments = &RegistrationArguments{}

func init() {
	RootCmd.AddCommand(registerCmd)

	// Transaction Parameters
	registerCmd.Flags().StringVar(&arguments.identity, "identity", "", "User's Identity")
	registerCmd.Flags().StringVar(&arguments.account, "account", "", "User's Default Account")
	// This should be set to the global version of this
	registerCmd.Flags().StringVar(&arguments.nodeName, "node", "", "User's Default Node")

	registerCmd.Flags().StringVar(&arguments.pubkey, "pubkey", "", "Specify a public key")
	registerCmd.Flags().Float64Var(&arguments.fee, "fee", 0.0, "Transaction Fee in OLT")
}

func RegisterIdentity(cmd *cobra.Command, args []string) {
	if arguments.identity == "" && arguments.account == "" {
		shared.Console.Fatal("Registration missing an identity or an account argument")
	}

	if arguments.fee == 0.0 {
		shared.Console.Fatal("Registration must include a fee")
	}

	arguments := &shared.RegisterArguments{
		Identity: arguments.identity,
		Account:  arguments.account,
		NodeName: arguments.nodeName,
		Fee:      arguments.fee,
	}

	register := shared.RegisterIdentityRequest(arguments)

	if register == nil {
		shared.Console.Fatal("Invalid Registration arguments")
	}

	result := comm.SDKRequest(register)
	_ = result

	shared.Console.Info("Identity has been broadcast to chain")
}
