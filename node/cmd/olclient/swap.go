/*
	Copyright 2017-2018 OneLedger

	Cli to interact with a with the chain.
*/
package main

import (
	"os"

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

var swapargs = &comm.SwapArguments{}

func init() {
	RootCmd.AddCommand(swapCmd)

	// Transaction Parameters
	swapCmd.Flags().StringVar(&swapargs.Party, "party", "", "base address")
	swapCmd.Flags().StringVar(&swapargs.CounterParty, "counterparty", "", "target address")

	swapCmd.Flags().Float64Var(&swapargs.Amount, "amount", 0.0, "the coins to exchange")
	swapCmd.Flags().StringVar(&swapargs.Currency, "currency", "OLT", "currency of amount")

	swapCmd.Flags().Float64Var(&swapargs.Exchange, "exchange", 0.0, "the value to trade for")
	swapCmd.Flags().StringVar(&swapargs.Excurrency, "excurrency", "ETH", "the currency")

	swapCmd.Flags().Int64Var(&swapargs.Nonce, "nonce", 0, "number used once")

	swapCmd.Flags().Float64Var(&swapargs.Fee, "fee", 0.0, "include a fee in OLT")
	swapCmd.Flags().Int64Var(&swapargs.Gas, "gas", 0, "gas in units")
}

// IssueRequest sends out a sendTx to all of the nodes in the chain
func SwapCurrency(cmd *cobra.Command, args []string) {
	log.Debug("Swap Request", "swapargs", swapargs)

	// Create message
	packet := shared.CreateSwapRequest(swapargs)
	if packet == nil {
		shared.Console.Error("Error in sending request")
		os.Exit(-1)
	}

	result := comm.Broadcast(packet)
	BroadcastStatus(result)
}

// TODO: Check to see that this is a valid currency
func GetCurrency(value string) string {
	return value
}
