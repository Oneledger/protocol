/*
	Copyright 2017-2018 OneLedger

	Cli to interact with a with the chain.
*/
package main

import (
	"github.com/Oneledger/protocol/node/id"
	"github.com/Oneledger/protocol/node/log"
	"github.com/spf13/cobra"
)

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Create or reuse an account",
	Run:   Register,
}

// Arguments to the command
type RegistrationArguments struct {
	name    string
	chain   string
	pubkey  string
	privkey string
}

var arguments = &RegistrationArguments{}

func init() {
	RootCmd.AddCommand(registerCmd)

	// Transaction Parameters
	registerCmd.Flags().StringVarP(&arguments.name, "name", "n", "Me", "User's Identity")
	registerCmd.Flags().StringVarP(&arguments.chain, "chain", "c", "OneLedger-Root", "Specify the chain")
	registerCmd.Flags().StringVarP(&arguments.pubkey, "pubkey", "k", "0x00000000", "Specify a public key")
	registerCmd.Flags().StringVarP(&arguments.privkey, "privkey", "p", "0x00000000", "Specify a private key")
}

// IssueRequest sends out a sendTx to all of the nodes in the chain
func Register(cmd *cobra.Command, args []string) {
	log.Debug("Register Account")

	//identity, err := app.FindIdentity(arguments.name)
	identity, err := id.FindIdentity(arguments.name)
	if err != 0 {
		Console.Error("Invalid Identity String")
		return
	}

	if identity == nil {
		CreateIdentity()
	} else {
		UpdateIdentity(identity)
	}
}

func ConvertPublicKey(keystring string) id.PublicKey {
	return id.PublicKey{}
}

// TODO: This should be moved out of cmd and into it's own package
func CreateIdentity() {
	chainType, err := id.FindIdentityType(arguments.chain)
	if err != 0 {
		Console.Error("Invalid Identity Type")
		return
	}

	pubkey := ConvertPublicKey(arguments.pubkey)

	identity := id.NewIdentity(chainType, arguments.name, pubkey)

	// TODO: Need to convert the input...
	identity.AddPrivateKey(id.PrivateKey{}) //arguments.privkey)

	Console.Info("New Account has been created")
}

func UpdateIdentity(identity id.Identity) {
	// TODO: Need to convert the input...
	identity.AddPrivateKey(id.PrivateKey{}) //arguments.privkey)

	Console.Info("New Account has been updated")
}
