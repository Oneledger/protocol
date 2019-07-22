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
	"github.com/Oneledger/protocol/client"
	"github.com/spf13/cobra"
)

var balanceCmd = &cobra.Command{
	Use:   "balance",
	Short: "Print out balance for account",
	Run:   BalanceNode,
}

type Balance struct {
	identityName string
	accountName  string
	accountKey   []byte
}

var balArgs *Balance = &Balance{}

func init() {
	RootCmd.AddCommand(balanceCmd)

	// Transaction Parameters
	// TODO either get by identity or read base64 of account key
	balanceCmd.Flags().BytesHexVar(&balArgs.accountKey, "address", []byte{}, "account address")
}

// IssueRequest sends out a sendTx to all of the nodes in the chain
func BalanceNode(cmd *cobra.Command, args []string) {
	Ctx := NewContext()

	if len(balArgs.accountKey) == 0 {
		logger.Error("missing address")
		return
	}

	fullnode := Ctx.clCtx.FullNodeClient()
	nodeName, err := fullnode.NodeName()
	if err != nil {
		logger.Fatal(err)
	}

	// assuming we have public key
	bal, err := fullnode.Balance(balArgs.accountKey)
	if err != nil {
		logger.Fatal("error in getting balance", err)
	}
	printBalance(nodeName, balArgs.accountKey, bal)
}

func printBalance(nodeName string, address []byte, bal client.BalanceReply) {
	logger.Infof("\t Balance for address %x on %s", address, nodeName)
	logger.Info("\t Balance:", bal.Balance)
	logger.Info("\t Height:", bal.Height)
}
