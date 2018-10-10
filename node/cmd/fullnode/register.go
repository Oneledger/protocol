/*
	Copyright 2017-2018 OneLedger

	Register identities and accounts. This only works if the fullnode is not running, and
	needs a client implementation as well.
*/
package main

import (
	"github.com/Oneledger/protocol/node/app"
	"github.com/Oneledger/protocol/node/id"
	"github.com/Oneledger/protocol/node/log"
	"github.com/spf13/cobra"
)

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Register identities and accounts",
	Run:   RegisterUsers,
}

// Arguments to the command
type RegisterArguments struct {
	identity string
	chain    string
	pubkey   string
	privkey  string
}

var regArguments = &RegisterArguments{}

// Initialize the command and flags
func init() {
	RootCmd.AddCommand(registerCmd)

	// Transaction Parameters
	registerCmd.Flags().StringVar(&regArguments.identity, "identity", "", "User's Identity")
	registerCmd.Flags().StringVar(&regArguments.chain, "chain", "OneLedger", "Specify the chain")
	registerCmd.Flags().StringVar(&regArguments.pubkey, "pubkey", "", "Specify a public key")
	registerCmd.Flags().StringVar(&regArguments.privkey, "privkey", "", "Specify a private key")
}

// IssueRequest sends out a sendTx to all of the nodes in the chain
func RegisterUsers(cmd *cobra.Command, args []string) {

	// TODO: We can't do this, need to be 'light-client' instead...
	node := app.NewApplication()

	name := regArguments.identity
	chain := regArguments.chain

	accountType := id.ParseAccountType(chain)

	privateKey, publicKey := id.GenerateKeys([]byte(name + chain)) // TODO: Switch with passphrase
	log.Info("RegisterUsers", "name", name, "publicKey", publicKey, "privateKey", privateKey, "chain", chain)
	app.RegisterLocally(node, name, chain, accountType, publicKey, privateKey)

}
