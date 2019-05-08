/*
   ____             _              _                      _____           _                  _
  / __ \           | |            | |                    |  __ \         | |                | |
 | |  | |_ __   ___| |     ___  __| | __ _  ___ _ __     | |__) | __ ___ | |_ ___   ___ ___ | |
 | |  | | '_ \ / _ \ |    / _ \/ _` |/ _` |/ _ \ '__|    |  ___/ '__/ _ \| __/ _ \ / __/ _ \| |
 | |__| | | | |  __/ |___|  __/ (_| | (_| |  __/ |       | |   | | | (_) | || (_) | (_| (_) | |
  \____/|_| |_|\___|______\___|\__,_|\__, |\___|_|       |_|   |_|  \___/ \__\___/ \___\___/|_|
                                      __/ |
                                     |___/

	Copyright 2017 - 2019 OneLedger

*/

package main

import (
	"strings"

	"github.com/Oneledger/protocol/data"
	"github.com/Oneledger/protocol/data/accounts"
	"github.com/Oneledger/protocol/data/chain"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update an account",
	Run:   UpdateAccount,
}

// Arguments to the command
type UpdateArguments struct {
	account     string
	chain       string
	pubkey      string
	privkey     string
	chainkey    string
	nodeaccount bool
}

var updateArgs = &UpdateArguments{}

func init() {
	RootCmd.AddCommand(updateCmd)

	// Transaction Parameters
	updateCmd.Flags().StringVar(&updateArgs.account, "account", "", "Account Name")
	updateCmd.Flags().StringVar(&updateArgs.chain, "chain", "OneLedger", "Specify the chain")

	updateCmd.Flags().StringVar(&updateArgs.pubkey, "pubkey", "0x00000000", "Specify a public key")
	updateCmd.Flags().StringVar(&updateArgs.privkey, "privkey", "0x00000000", "Specify a private key")
	updateCmd.Flags().StringVar(&updateArgs.chainkey, "chainkey", "<empty>", "Specify the chain key")
	updateCmd.Flags().BoolVar(&updateArgs.nodeaccount, "nodeaccount", false, "Specify whether it's a node account or not")
}

func UpdateAccount(cmd *cobra.Command, args []string) {
	logger.Debug("UPDATING ACCOUNT")

	typ, err := chain.TypeFromName(updateArgs.chain)
	if err != nil {
		logger.Error("chain not registered")
		return
	}

	var privKey keys.PrivateKey
	var pubKey  keys.PublicKey

	if updateArgs.privkey == "" || updateArgs.pubkey == ""{
		// if a public key or a private key is not passed; generate a pair of keys
		tmPrivKey := ed25519.GenPrivKey()
		tmPublicKey := tmPrivKey.PubKey()

		privKey = keys.PrivateKey{keys.ED25519, tmPrivKey.Bytes()}
		pubKey = keys.PublicKey{keys.ED25519, tmPublicKey.Bytes()}
	} else {
		// parse keys passed through commandline

		pubKeyStr := strings.TrimPrefix(updateArgs.pubkey, "0x")
		pubKey, err = keys.GetPublicKeyFromBytes([]byte(pubKeyStr), keys.ED25519)
		if err != nil {
			logger.Error("incorrect public key", err)
			return
		}

		privKeyStr := strings.TrimPrefix(updateArgs.privkey, "0x")
		privKey, err = keys.GetPrivateKeyFromBytes([]byte(privKeyStr), keys.ED25519)
		if err != nil {
			logger.Error("incorrect private key", err)
			return
		}
	}

	acc, err := accounts.NewAccount(typ, updateArgs.account, privKey, pubKey)
	if err != nil {
		logger.Error("Error initializing account", err)
		return
	}

	resp := &data.Response{}
	err = Ctx.Query("AddAcount", acc, resp)
	if err != nil {
		logger.Error("error creating account", err)
		return
	}

	logger.Info("created account", "account", acc)
}
