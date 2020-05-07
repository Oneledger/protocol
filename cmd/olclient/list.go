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
	"fmt"

	accounts2 "github.com/Oneledger/protocol/data/accounts"
	"github.com/Oneledger/protocol/data/keys"

	"github.com/spf13/cobra"
)

type ListArguments struct {
	identityName string
	accountName  string
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List out Node data",
	Run:   ListNode,
}

var list = &ListArguments{}

func init() {
	RootCmd.AddCommand(listCmd)

	// TODO: I want to have a default account?
	// Transaction Parameters
	listCmd.Flags().StringVar(&list.identityName, "identity", "", "identity name")
	listCmd.Flags().StringVar(&list.accountName, "account", "", "account name")
}

// IssueRequest sends out a sendTx to all of the nodes in the chain
func ListNode(cmd *cobra.Command, args []string) {
	Ctx := NewContext()
	fullnode := Ctx.clCtx.FullNodeClient()

	wallet, err := accounts2.NewWalletKeyStore(keyStorePath)
	if err != nil {
		logger.Error("listnode: error creating secure wallet.")
		return
	}

	addresses, err := wallet.ListAddresses()
	if err != nil {
		logger.Error("error in getting all accounts", err)
		return
	}

	logger.Infof("Accounts on node: %s ", Ctx.cfg.Node.NodeName)
	fmt.Println("Accounts from Secure Wallet:")
	for _, a := range addresses {
		fmt.Print("Address: ", a, "    ")
		rep, err := fullnode.Balance(a)
		if err != nil {
			continue
		}
		fmt.Printf("Balance: %s \n\n", rep.Balance)
	}

	addressListReply, err := fullnode.ListAccountAddresses()
	if err != nil {
		logger.Error("listnode: error getting address list from node wallet.", err)
		return
	}

	nodeWalletAddresses := addressListReply.Addresses

	fmt.Println("Accounts from Node Wallet:")
	for _, a := range nodeWalletAddresses {
		fmt.Print("Address: ", a, "    ")
		addr := keys.Address{}
		err = addr.UnmarshalText([]byte(a))
		if err != nil {
			return
		}
		rep, err := fullnode.Balance(addr)
		if err != nil {
			continue
		}
		fmt.Printf("Balance: %s \n\n", rep.Balance)
	}
}
