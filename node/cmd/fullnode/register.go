/*
	Copyright 2017-2018 OneLedger

	Cli to interact with a with the chain.
*/
package main

import (
	"github.com/Oneledger/prototype/node/app"
	"github.com/spf13/cobra"
)

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Create or reuse an account",
	Run:   Register,
}

// Arguments to the command
type RegistrationArguments struct {
	chain     string
	account   string
	chainType string
	pubkey    string
	privkey   string
}

var arguments = &RegistrationArguments{}

func init() {
	RootCmd.AddCommand(registerCmd)

	// Transaction Parameters
	registerCmd.Flags().StringVarP(&arguments.chain, "chain", "c", "OneLedger-Root", "Specify the chain")
	registerCmd.Flags().StringVarP(&arguments.pubkey, "pubkey", "k", "0x00000000", "Specify a public key")
	registerCmd.Flags().StringVarP(&arguments.privkey, "privkey", "p", "0x00000000", "Specify a private key")
}

// IssueRequest sends out a sendTx to all of the nodes in the chain
func Register(cmd *cobra.Command, args []string) {
	app.Log.Debug("Register Account")
}
