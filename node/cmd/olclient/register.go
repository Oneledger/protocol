/*
	Copyright 2017-2018 OneLedger

	Cli to interact with a with the chain.
*/
package main

import (
	"github.com/Oneledger/protocol/node/app"
	"github.com/Oneledger/protocol/node/comm"
	"github.com/Oneledger/protocol/node/log"
	"github.com/Oneledger/protocol/node/serial"
	"github.com/spf13/cobra"
)

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Create or reuse an account",
	Run:   Register,
}

// Arguments to the command
type RegistrationArguments struct {
	identity string
	chain    string
	pubkey   string
	privkey  string
}

var arguments = &RegistrationArguments{}

func init() {
	RootCmd.AddCommand(registerCmd)

	// Transaction Parameters
	registerCmd.Flags().StringVar(&arguments.identity, "identity", "Unknown", "User's Identity")
	registerCmd.Flags().StringVar(&arguments.chain, "chain", "OneLedger", "Specify the chain")
	registerCmd.Flags().StringVar(&arguments.pubkey, "pubkey", "0x00000000", "Specify a public key")
	registerCmd.Flags().StringVar(&arguments.privkey, "privkey", "0x00000000", "Specify a private key")
}

// IssueRequest sends out a sendTx to all of the nodes in the chain
func Register(cmd *cobra.Command, args []string) {
	log.Debug("Client Register Account via SetOption...")

	cli := &app.RegisterArguments{
		Identity:   arguments.identity,
		Chain:      arguments.chain,
		PublicKey:  arguments.pubkey,
		PrivateKey: arguments.privkey,
	}

	buffer, err := serial.Serialize(cli, serial.CLIENT)
	if err != nil {
		log.Error("Register Failed", "err", err)
		return
	}
	comm.SetOption("Register", string(buffer))
}

/*
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
	//identity.AddPrivateKey(id.PrivateKey{}) //arguments.privkey)

	Console.Info("New Account has been created")
}

func UpdateIdentity(identity id.Identity) {
	// TODO: Need to convert the input...
	//identity.AddPrivateKey(id.PrivateKey{}) //arguments.privkey)

	Console.Info("New Account has been updated")
}
*/
