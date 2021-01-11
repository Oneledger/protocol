/*
   ____             _              _                      _____           _                  _
  / __ \           | |            | |                    |  __ \         | |                | |
 | |  | |_ __   ___| |     ___  __| | __ _  ___ _ __     | |__) | __ ___ | |_ ___   ___ ___ | |
 | |  | | '_ \ / _ \ |    / _ \/ _` |/ _` |/ _ \ '__|    |  ___/ '__/ _ \| __/ _ \ / __/ _ \| |
 | |__| | | | |  __/ |___|  __/ (_| | (_| |  __/ |       | |   | | | (_) | || (_) | (_| (_) | |
  \____/|_| |_|\___|______\___|\__,_|\__, |\___|_|       |_|   |_|  \___/ \__\___/ \___\___/|_|
                                      __/ |
                                     |___/


Copyright 2017 - 2020 OneLedger
*/

package main

import (
	"fmt"
	"github.com/Oneledger/protocol/config"
	"github.com/Oneledger/protocol/data/accounts"

	"github.com/spf13/cobra"
)

type MergeArguments struct {
	oldDBDir     string
	currentDBDir string
	dbType       string
}

var mergeCmd = &cobra.Command{
	Use:   "merge",
	Short: "Merge node wallet from previous node to this node",
	Run:   MergeNodeWallet,
}

var merge = &MergeArguments{}

func init() {
	RootCmd.AddCommand(mergeCmd)

	// Parameters
	mergeCmd.Flags().StringVar(&merge.oldDBDir, "old", "", "db location for old wallet")
	mergeCmd.Flags().StringVar(&merge.currentDBDir, "current", "./nodedata", "db location for current wallet")
	mergeCmd.Flags().StringVar(&merge.dbType, "type", "cleveldb", "db type, goleveldb|cleveldb")
}

func MergeNodeWallet(cmd *cobra.Command, args []string) {
	// create new wallet session
	config := config.Server{Node: &config.NodeConfig{}}
	config.Node.DB = merge.dbType
	oldWallet := accounts.NewWallet(config, merge.oldDBDir)
	oldAccountList := oldWallet.AccountsWithKey()

	fmt.Println("Accounts in old wallet:")
	for _, a := range oldAccountList {
		fmt.Print("Privkey: ", a.PrivateKey, "    ")
		fmt.Print("Pubkey: ", a.PublicKey, "    ")
		fmt.Println("Address: ", a.Address())
	}

	newWallet := accounts.NewWallet(config, merge.currentDBDir)
	newAccountList := newWallet.AccountsWithKey()

	fmt.Println("Accounts in new wallet before merging:")
	for _, a := range newAccountList {
		fmt.Print("Privkey: ", a.PrivateKey, "    ")
		fmt.Print("Pubkey: ", a.PublicKey, "    ")
		fmt.Println("Address: ", a.Address())
	}

	//merge
	for _, a := range oldAccountList {
		err := newWallet.Add(a)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	newAccountListAfterMerge := newWallet.AccountsWithKey()
	fmt.Println("Accounts in new wallet after merging:")
	for _, a := range newAccountListAfterMerge {
		fmt.Print("Privkey: ", a.PrivateKey, "    ")
		fmt.Print("Pubkey: ", a.PublicKey, "    ")
		fmt.Println("Address: ", a.Address())
	}


}
