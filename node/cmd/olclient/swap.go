/*
	Copyright 2017-2018 OneLedger

	Cli to interact with a with the chain.
*/
package main

import (
	"github.com/Oneledger/protocol/node/cmd/shared"
	"github.com/Oneledger/protocol/node/comm"
	"github.com/Oneledger/protocol/node/log"
	"github.com/spf13/cobra"
)

var swapCmd = &cobra.Command{
	Use:   "swap",
	Short: "Setup or confirm a currency swap",
	Run:   SwapCurrency,
}

var swapargs = &shared.SwapArguments{}

func init() {
	RootCmd.AddCommand(swapCmd)

	// Transaction Parameters
	swapCmd.Flags().StringVar(&swapargs.Party, "party", "unknown", "base address")
	swapCmd.Flags().StringVar(&swapargs.CounterParty, "counterparty", "unknown", "target address")
	swapCmd.Flags().StringVar(&swapargs.Amount, "amount", "0", "the coins to exchange")
	swapCmd.Flags().StringVar(&swapargs.Currency, "currency", "OLT", "currency of amount")
	swapCmd.Flags().StringVar(&swapargs.Exchange, "exchange", "0", "the value to trade for")
	swapCmd.Flags().StringVar(&swapargs.Excurrency, "excurrency", "ETH", "the currency")
	swapCmd.Flags().Int64Var(&swapargs.Nonce, "nonce", 1001, "number used once")

	swapCmd.Flags().StringVar(&swapargs.Fee, "fee", "1", "fees in coins")
	swapCmd.Flags().StringVar(&swapargs.Gas, "gas", "1", "gas, if necessary")
}

// IssueRequest sends out a sendTx to all of the nodes in the chain
func SwapCurrency(cmd *cobra.Command, args []string) {
	log.Debug("Swap Request", "tx", swapargs)

	// Create message
	packet := shared.CreateSwapRequest(swapargs)

	result := comm.Broadcast(packet)

	log.Debug("Returned Successfully", "result", result)
}

// TODO: Check to see that this is a valid currency
func GetCurrency(value string) string {
	return value
}
