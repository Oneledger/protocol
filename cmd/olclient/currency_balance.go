/*

 */

package main

import (
	"github.com/Oneledger/protocol/client"
	"github.com/spf13/cobra"
)

var currBalanceCmd = &cobra.Command{
	Use:   "curr_balance",
	Short: "Print out balance for account",
	Run:   BalanceNode,
}

type CurrBalance struct {
	accountKey   []byte
	currencyName string
}

var currBalArgs *CurrBalance = &CurrBalance{}

func init() {
	RootCmd.AddCommand(currBalanceCmd)

	// Transaction Parameters
	// TODO either get by identity or read base64 of account key
	currBalanceCmd.Flags().BytesHexVar(&currBalArgs.accountKey, "address", []byte{}, "account address")

	currBalanceCmd.Flags().StringVar(&currBalArgs.currencyName, "currency", "OLT", "currency name")
}

// IssueRequest sends out a sendTx to all of the nodes in the chain
func CurrBalanceNode(cmd *cobra.Command, args []string) {
	Ctx := NewContext()

	if len(currBalArgs.accountKey) == 0 {
		logger.Error("missing address")
		return
	}

	fullnode := Ctx.clCtx.FullNodeClient()
	nodeName, err := fullnode.NodeName()
	if err != nil {
		logger.Fatal(err)
	}

	// assuming we have public key
	bal, err := fullnode.CurrBalance(currBalArgs.accountKey, currBalArgs.currencyName)
	if err != nil {
		logger.Fatal("error in getting balance", err)
	}
	printCurrBalance(nodeName, currBalArgs.accountKey, bal)
}

func printCurrBalance(nodeName string, address []byte, bal client.CurrencyBalanceReply) {
	logger.Infof("\t Balance for address %x on %s", address, nodeName)
	logger.Info("\t Currency:", bal.Currency)
	logger.Info("\t Balance:", bal.Balance)
	logger.Info("\t Height:", bal.Height)
}
