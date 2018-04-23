/*
	Copyright 2017-2018 OneLedger

	Cli to interact with a with the chain.
*/
package main

import (
	"fmt"

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
	app.Log.Debug("Register Account")

	identity, err := app.FindIdentity(arguments.name)
	if err != 0 {
		app.Log.Error("Invalid Identity string")
		return
	}

	if identity == nil {
		CreateIdentity()
	} else {
		UpdateIdentity(identity)
	}
}

func CreateIdentity() {
	chainType, err := app.FindIdentityType(arguments.chain)
	if err != 0 {
		fmt.Println("Invalid identity type")
	}

	identity := app.NewIdentity(arguments.name, chainType)

	// TODO: Need to conver the input...
	identity.AddPublicKey(app.PublicKey{})   //arguments.pubkey)
	identity.AddPrivateKey(app.PrivateKey{}) //arguments.privkey)

	fmt.Println("New Account has been created")
}

func UpdateIdentity(identity app.Identity) {
	// TODO: Need to conver the input...
	identity.AddPublicKey(app.PublicKey{})   //arguments.pubkey)
	identity.AddPrivateKey(app.PrivateKey{}) //arguments.privkey)

	fmt.Println("New Account has been updated")
}
