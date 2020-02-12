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
	"github.com/spf13/cobra"

	"github.com/Oneledger/protocol/client"
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
	currencyName string
}

var balArgs *Balance = &Balance{}

func init() {
	RootCmd.AddCommand(balanceCmd)

	// Transaction Parameters
	// TODO either get by identity or read base64 of account key
	balanceCmd.Flags().BytesHexVar(&balArgs.accountKey, "address", []byte{}, "account address")

	balanceCmd.Flags().StringVar(&balArgs.currencyName, "currency", "", "currency name")

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
	if balArgs.currencyName == "" {
		bal, err := fullnode.Balance(balArgs.accountKey)
		if err != nil {
			logger.Fatal("error in getting balance", err)
		}

		printBalance(nodeName.Name, balArgs.accountKey, bal)

	} else {
		bal, err := fullnode.CurrBalance(balArgs.accountKey, balArgs.currencyName)
		if err != nil {
			logger.Fatal("error in getting balance", err)
		}

		printCurrBalance(nodeName.Name, balArgs.accountKey, bal)
	}
}

func printBalance(nodeName string, address []byte, bal client.BalanceReply) {
	logger.Infof("\t Balance for address %x on %s", address, nodeName)
	logger.Info("\t Balance:", bal.Balance)
	logger.Info("\t Height:", bal.Height)
}

func printCurrBalance(nodeName string, address []byte, bal client.CurrencyBalanceReply) {
	logger.Infof("\t Balance for address %x on %s", address, nodeName)
	logger.Info("\t Currency:", bal.Currency)
	logger.Info("\t Balance:", bal.Balance)
	logger.Info("\t Height:", bal.Height)
}
