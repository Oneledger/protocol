/*
	Copyright 2017-2018 OneLedger

	Cli to interact with a with the chain.
*/
package main

import (
	"github.com/Oneledger/protocol/node/cmd/shared"
	"github.com/Oneledger/protocol/node/comm"
	"github.com/Oneledger/protocol/node/id"
	"github.com/Oneledger/protocol/node/log"
	"github.com/spf13/cobra"
	"reflect"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update an account",
	Run:   UpdateAccount,
}

// Arguments to the command
type UpdateArguments struct {
	account string
	chain   string
	pubkey  string
	privkey string
}

var updateArgs = &UpdateArguments{}

func init() {
	RootCmd.AddCommand(updateCmd)

	// Transaction Parameters
	updateCmd.Flags().StringVar(&updateArgs.account, "account", "", "Account Name")
	updateCmd.Flags().StringVar(&updateArgs.chain, "chain", "OneLedger", "Specify the chain")
	updateCmd.Flags().StringVar(&updateArgs.pubkey, "pubkey", "0x00000000", "Specify a public key")
	updateCmd.Flags().StringVar(&updateArgs.privkey, "privkey", "0x00000000", "Specify a private key")
}

func UpdateAccount(cmd *cobra.Command, args []string) {
	log.Debug("UPDATING ACCOUNT")

	// TODO: Don't need two levels of structures here
	request := &shared.AccountArguments{
		Account:    updateArgs.account,
		Chain:      updateArgs.chain,
		PublicKey:  updateArgs.pubkey,
		PrivateKey: updateArgs.privkey,
	}

	update := shared.UpdateAccountRequest(request)

	result := comm.SDKRequest(update)

	switch value := result.(type) {
	case string:
		shared.Console.Info(value)
	case id.Account:
		shared.Console.Info("Created account: ", value.Name())
	default:
		shared.Console.Info("Invalid type: ", reflect.TypeOf(value).String())
	}
}
