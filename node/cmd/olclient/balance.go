/*
	Copyright 2017-2018 OneLedger

	Cli to interact with a with the chains.
*/
package main

import (
	"encoding/hex"
	"os"

	"github.com/Oneledger/protocol/node/action"
	"github.com/Oneledger/protocol/node/cmd/shared"
	"github.com/Oneledger/protocol/node/comm"
	"github.com/Oneledger/protocol/node/data"
	"github.com/Oneledger/protocol/node/id"
	"github.com/spf13/cobra"
)

var balanceCmd = &cobra.Command{
	Use:   "balance",
	Short: "Print out balance for account",
	Run:   BalanceNode,
}

// TODO: typing should be way better, see if cobra can help with this...
type Balance struct {
	identityName string
	accountName  string
}

var balance *Balance = &Balance{}

func init() {
	RootCmd.AddCommand(balanceCmd)

	// TODO: I want to have a default account?
	// Transaction Parameters
	balanceCmd.Flags().StringVar(&balance.identityName, "identity", "", "identity name")
	balanceCmd.Flags().StringVar(&balance.accountName, "account", "", "account name")
}

func GetName() string {
	if balance.identityName != "" {
		return balance.identityName
	}
	if balance.accountName != "" {
		return balance.accountName
	}
	shared.Console.Error("Invalid Query, missing Account or Identity")
	os.Exit(-1)

	return ""
}

// Format the request into a query structure
func AccountKeyRequest(name string) []byte {
	return action.Message("Identity=" + name)
}

func BalanceRequest(accountKey string) []byte {
	return action.Message("Balance=" + accountKey)
}

// IssueRequest sends out a sendTx to all of the nodes in the chain
func BalanceNode(cmd *cobra.Command, args []string) {
	nodeName := shared.GetNodeName()
	name := GetName()

	accountKeyRequest := AccountKeyRequest(name)
	accountKey := comm.Query("/accountKey", accountKeyRequest)

	balanceRequest := BalanceRequest(hex.EncodeToString(accountKey.(id.AccountKey)))

	balance := comm.Query("/balance", balanceRequest)

	if balance == nil {
		shared.Console.Warning("No Response from Node for:", string(balanceRequest))
		return
	}

	printBalance(nodeName, name, balance.(*data.Balance))
}

func printBalance(nodeName string, name string, balance *data.Balance) {

	balanceLabel := "          Balance:"

	shared.Console.Info("Balance for", name, "on", nodeName+":\n")

	shared.Console.Info(balanceLabel, balance.String())
	shared.Console.Info()
}
