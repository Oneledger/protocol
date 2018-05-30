/*
	Copyright 2017-2018 OneLedger

	Register identities and accounts. This only works if the fullnode is not running, and
	needs a client implementation as well.
*/
package main

import (
	"github.com/Oneledger/protocol/node/app"
	"github.com/Oneledger/protocol/node/id"
	"github.com/spf13/cobra"
)

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Register identities and accounts",
	Run:   RegisterUsers,
}

// Arguments to the command
type RegisterArguments struct {
	name    string
	chain   string
	pubkey  string
	privkey string
}

var arguments = &RegisterArguments{}

func init() {
	RootCmd.AddCommand(registerCmd)

	// Transaction Parameters
	registerCmd.Flags().StringVarP(&arguments.name, "name", "n", "Me", "User's Identity")
	registerCmd.Flags().StringVarP(&arguments.chain, "chain", "c", "OneLedger-Root", "Specify the chain")
	registerCmd.Flags().StringVarP(&arguments.pubkey, "pubkey", "k", "0x00000000", "Specify a public key")
	registerCmd.Flags().StringVarP(&arguments.privkey, "privkey", "p", "0x00000000", "Specify a private key")
}

// IssueRequest sends out a sendTx to all of the nodes in the chain
func RegisterUsers(cmd *cobra.Command, args []string) {

	// TODO: We can't do this, need to be 'light-client' instead...
	node := app.NewApplication()

	app.Register(node, arguments.name, arguments.name+"-OneLedger", id.ParseAccountType(arguments.chain))
}
