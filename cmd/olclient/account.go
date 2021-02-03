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
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Oneledger/protocol/data/accounts"
	"github.com/Oneledger/protocol/data/chain"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/spf13/cobra"
)

// Arguments to the update command
type AddArguments struct {
	account  string
	chain    string
	password string
	pubkey   []byte
	privkey  []byte
}

//Arguments to the delete command
type DeleteArguments struct {
	Address  string `json:"address"`
	Password string `json:"password"`
}

//Arguments to the get command
type GetArguments struct {
	Address  string `json:"address"`
	Password string `json:"password"`
}

var (
	accountCmd = &cobra.Command{
		Use:   "account",
		Short: "handling an account",
		Long:  "local account handling through secure wallet",
	}

	deleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "delete an account",
		RunE:  Delete,
	}

	addCmd = &cobra.Command{
		Use:   "add",
		Short: "update or create an account",
		RunE:  Add,
	}

	getCmd = &cobra.Command{
		Use:   "get",
		Short: "retrieve account data",
		RunE:  Get,
	}

	addArgs    = &AddArguments{}
	deleteArgs = &DeleteArguments{}
	getArgs    = &GetArguments{}
)

func parseUpdateArgs() {
	addCmd.Flags().StringVar(&addArgs.account, "name", "", "Account Name")
	addCmd.Flags().StringVar(&addArgs.chain, "chain", "OneLedger", "Specify the chain")
	addCmd.Flags().BytesBase64Var(&addArgs.pubkey, "pubkey", []byte{}, "Specify a base64 public key")
	addCmd.Flags().BytesBase64Var(&addArgs.privkey, "privkey", []byte{}, "Specify a base64 private key")
	addCmd.Flags().StringVar(&addArgs.password, "password", "", "password to access secure wallet")
}

func parseDeleteArgs() {
	deleteCmd.Flags().StringVar(&deleteArgs.Address, "address", "", "address to delete")
	deleteCmd.Flags().StringVar(&deleteArgs.Password, "password", "", "password to access secure wallet")
}

func parseGetArgs() {
	getCmd.Flags().StringVar(&getArgs.Address, "address", "", "address of account to retrieve")
	getCmd.Flags().StringVar(&getArgs.Password, "password", "", "password to access secure wallet")
}

func init() {
	RootCmd.AddCommand(accountCmd)
	accountCmd.AddCommand(deleteCmd)
	accountCmd.AddCommand(addCmd)
	accountCmd.AddCommand(getCmd)

	//Add Transaction Parameters
	parseUpdateArgs()

	//Delete Transaction Parameters
	parseDeleteArgs()

	//Get Transaction Parameters
	parseGetArgs()
}

func Add(cmd *cobra.Command, args []string) error {
	wallet, err := accounts.NewWalletKeyStore(rootArgs.rootDir + keyStorePath)
	if err != nil {
		return err
	}

	typ, err := chain.TypeFromName(addArgs.chain)
	if err != nil {
		return errors.New("chain not registered: " + addArgs.chain)
	}

	// get the kys for the new account
	var privKey keys.PrivateKey
	var pubKey keys.PublicKey

	if len(addArgs.privkey) == 0 || len(addArgs.pubkey) == 0 {
		// if a public key or a private key is not passed; generate a pair of keys
		pubKey, privKey, err = keys.NewKeyPairFromTendermint()
		if err != nil {
			return errors.New("error generating key from tendermint" + err.Error())
		}
	} else {
		// parse keys passed through commandline

		pubKey, err = keys.GetPublicKeyFromBytes(addArgs.pubkey, keys.ED25519)
		if err != nil {
			fmt.Println("incorrect public key" + err.Error())
			return err
		}

		privKey, err = keys.GetPrivateKeyFromBytes(addArgs.privkey, keys.ED25519)
		if err != nil {
			fmt.Println("incorrect private key" + err.Error())
			return err
		}
	}

	// create the account
	acc, err := accounts.NewAccount(typ, addArgs.account, &privKey, &pubKey)
	if err != nil {
		return errors.New("Error initializing account" + err.Error())
	}

	//Prompt for password update account
	if len(addArgs.password) == 0 {
		addArgs.password = PromptForPassword()
	}

	if !wallet.Open(acc.Address(), addArgs.password) {
		return errors.New("error opening wallet")
	}

	err = wallet.Add(acc)
	if err != nil {
		return err
	}

	fmt.Println("Successfully added account to secure wallet.")
	fmt.Println("Address for the account is: ", acc.Address())

	wallet.Close()

	return nil
}

func Delete(cmd *cobra.Command, args []string) error {
	if len(deleteArgs.Address) <= 0 {
		return errors.New("error: invalid address")
	}

	wallet, err := accounts.NewWalletKeyStore(keyStorePath)
	if err != nil {
		return err
	}

	//Get Address
	usrAddress := keys.Address{}
	err = usrAddress.UnmarshalText([]byte(deleteArgs.Address))
	if err != nil {
		return err
	}

	if !wallet.KeyExists(usrAddress) {
		return errors.New("address does not exist")
	}

	//Prompt for password
	if len(deleteArgs.Password) == 0 {
		deleteArgs.Password = PromptForPassword()
	}

	//Verify User Password
	authenticated, _ := wallet.VerifyPassphrase(usrAddress, deleteArgs.Password)
	if !authenticated {
		return errors.New("authentication error")
	}

	if !wallet.Open(usrAddress, deleteArgs.Password) {
		return errors.New("error opening wallet")
	}

	err = wallet.Delete(usrAddress)
	if err != nil {
		return err
	}

	fmt.Printf("Successfully deleted account from secure wallet: %v\n", deleteArgs.Address)

	wallet.Close()

	return nil
}

func Get(cmd *cobra.Command, args []string) error {
	if len(getArgs.Address) <= 0 {
		return errors.New("error: invalid address")
	}

	wallet, err := accounts.NewWalletKeyStore(keyStorePath)
	if err != nil {
		return err
	}

	//Get Address
	usrAddress := keys.Address{}
	err = usrAddress.UnmarshalText([]byte(getArgs.Address))
	if err != nil {
		return err
	}

	//If Account already exists, Verify Password
	if !wallet.KeyExists(usrAddress) {
		return errors.New("address not found")
	}

	if len(getArgs.Password) == 0 {
		getArgs.Password = PromptForPassword()
	}
	auth, _ := wallet.VerifyPassphrase(usrAddress, getArgs.Password)
	if !auth {
		return errors.New("authentication failed")
	}

	if !wallet.Open(usrAddress, getArgs.Password) {
		return errors.New("error opening wallet")
	}

	account, err := wallet.GetAccount(usrAddress)
	if err != nil {
		return err
	}

	out, err := json.MarshalIndent(account, "", " ")
	fmt.Println("\n" + string(out))

	wallet.Close()

	return nil
}
