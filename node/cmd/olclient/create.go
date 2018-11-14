/*
	Copyright 2017-2018 OneLedger

	Cli to interact with a with the chain.
*/
package main

import (
	"github.com/Oneledger/protocol/node/cmd/shared"
	"github.com/Oneledger/protocol/node/comm"
	"github.com/Oneledger/protocol/node/log"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create or reuse an account",
	Run:   CreateAccount,
}

// Arguments to the command
type CreateArguments struct {
	name    string
	chain   string
	pubkey  string
	privkey string
}

var createArgs = &CreateArguments{}

func init() {
	RootCmd.AddCommand(createCmd)

	// Transaction Parameters
	createCmd.Flags().StringVar(&createArgs.name, "name", "", "Account Name")
	createCmd.Flags().StringVar(&createArgs.chain, "chain", "OneLedger", "Specify the chain")
	createCmd.Flags().StringVar(&createArgs.pubkey, "pubkey", "0x00000000", "Specify a public key")
	createCmd.Flags().StringVar(&createArgs.privkey, "privkey", "0x00000000", "Specify a private key")
}

func CreateAccount(cmd *cobra.Command, args []string) {
	request := &shared.AccountArguments{}

	register := shared.CreateAccountRequest(request)

	result := comm.SDKRequest(register)
	log.Dump("The results are", result)
}

/*
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
*/
