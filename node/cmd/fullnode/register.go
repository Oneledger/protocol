/*
	Copyright 2017-2018 OneLedger

	Register identities and accounts. This only works if the fullnode is not running, and
	needs a client implementation as well.
*/
package main

import (
	"github.com/Oneledger/protocol/node/action"
	"github.com/Oneledger/protocol/node/app"
	"github.com/Oneledger/protocol/node/global"
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
	identity string
	chain    string
	pubkey   string
	privkey  string
}

var regArguments = &RegisterArguments{}

func init() {
	RootCmd.AddCommand(registerCmd)

	// Transaction Parameters
	registerCmd.Flags().StringVar(&regArguments.identity, "identity", "unknown", "User's Identity")
	registerCmd.Flags().StringVar(&regArguments.chain, "chain", "OneLedger-Root", "Specify the chain")
	registerCmd.Flags().StringVar(&regArguments.pubkey, "pubkey", "0x00000000", "Specify a public key")
	registerCmd.Flags().StringVar(&regArguments.privkey, "privkey", "0x00000000", "Specify a private key")
}

// IssueRequest sends out a sendTx to all of the nodes in the chain
func RegisterUsers(cmd *cobra.Command, args []string) {

	// TODO: We can't do this, need to be 'light-client' instead...
	node := app.NewApplication()

	var signers []id.PublicKey

	app.Register(node, regArguments.identity, regArguments.identity+"-OneLedger", id.ParseAccountType(regArguments.chain))
	transaction := action.Register{
		Base: action.Base{
			Type:     action.REGISTER,
			ChainId:  app.ChainId,
			Signers:  signers,
			Sequence: global.Current.Sequence,
		},
	}
	action.SubmitTransaction(action.Transaction(transaction))
}
