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
	"github.com/Oneledger/protocol/data"
	"github.com/Oneledger/protocol/data/balance"
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
	accountKey   string
}

var balArgs *Balance = &Balance{}

func init() {
	RootCmd.AddCommand(balanceCmd)

	// Transaction Parameters
	// TODO either get by identity or read base64 of account key
	balanceCmd.Flags().StringVar(&balArgs.accountKey, "address", "", "account address")
}

// IssueRequest sends out a sendTx to all of the nodes in the chain
func BalanceNode(cmd *cobra.Command, args []string) {

	resp := &data.Response{}
	err := Ctx.Query("NodeName", nil, resp)
	if err != nil {
		logger.Fatal("error in getting nodename", err)
	}

	nodeName := string(resp.Data)

	// assuming we have public key
	bal := balance.NewBalance()
	err = Ctx.Query("Balance", []byte(balArgs.accountKey), bal)
	if err != nil || !resp.Success {
		logger.Fatal("error in getting nodename", err, resp.ErrorMsg)
	}

	printBalance(nodeName, balArgs.accountKey, bal)
}

func printBalance(nodeName string, address string, bal *balance.Balance) {

	logger.Infof("\t Balance for address %x on %s", address, nodeName)
	logger.Info("\t Balance: ", bal.String())
	logger.Info()
}
