/*
	Copyright 2017-2018 OneLedger

	Cli to interact with a with the chain.
*/
package main

import (
	"os"

	"github.com/Oneledger/protocol/node/action"
	"github.com/Oneledger/protocol/node/app"
	"github.com/Oneledger/protocol/node/convert"
	"github.com/Oneledger/protocol/node/data"
	"github.com/Oneledger/protocol/node/global"
	"github.com/Oneledger/protocol/node/id"
	"github.com/Oneledger/protocol/node/log"
	"github.com/spf13/cobra"
)

var swapCmd = &cobra.Command{
	Use:   "swap",
	Short: "Setup or confirm a currency swap",
	Run:   SwapCurrency,
}

// Arguments to the command
type SwapArguments struct {
	party        string
	counterparty string
	amount       string
	fee          string
	gas          string // TODO: Not sure this is necessary, unless the chain is like Ethereum
	currency     string
	exchange     string
	excurrency   string
	nonce        int64
}

var swapargs = &SwapArguments{}

func init() {
	RootCmd.AddCommand(swapCmd)

	// Transaction Parameters
	swapCmd.Flags().StringVar(&swapargs.party, "party", "unknown", "base address")
	swapCmd.Flags().StringVar(&swapargs.counterparty, "counterparty", "unknown", "target address")
	swapCmd.Flags().StringVar(&swapargs.amount, "amount", "0", "the coins to exchange")
	swapCmd.Flags().StringVar(&swapargs.currency, "currency", "OLT", "currency of amount")
	swapCmd.Flags().StringVar(&swapargs.exchange, "exchange", "0", "the value to trade for")
	swapCmd.Flags().StringVar(&swapargs.excurrency, "excurrency", "ETH", "the currency")
	swapCmd.Flags().Int64Var(&swapargs.nonce, "nonce", 1001, "number used once")

	swapCmd.Flags().StringVar(&swapargs.fee, "fee", "1", "fees in coins")
	swapCmd.Flags().StringVar(&swapargs.gas, "gas", "1", "gas, if necessary")
}

func CreateSwapRequest() []byte {
	log.Debug("swap args", "swapargs", swapargs)

	// TODO: Need better validation and error handling...

	conv := convert.NewConvert()

	party := id.Address(conv.GetHash(swapargs.party))
	counterparty := id.Address(conv.GetHash(swapargs.counterparty))

	// TOOD: a clash with the basic data model
	signers := GetSigners()

	fee := data.Coin{
		Currency: conv.GetCurrency(swapargs.currency),
		Amount:   conv.GetInt64(swapargs.fee),
	}

	gas := data.Coin{
		Currency: conv.GetCurrency(swapargs.currency),
		Amount:   conv.GetInt64(swapargs.gas),
	}

	amount := data.Coin{
		Currency: conv.GetCurrency(swapargs.currency),
		Amount:   conv.GetInt64(swapargs.amount),
	}

	exchange := data.Coin{
		Currency: conv.GetCurrency(swapargs.excurrency),
		Amount:   conv.GetInt64(swapargs.exchange),
	}

	if conv.HasErrors() {
		Console.Error(conv.GetErrors())
		os.Exit(-1)
	}

	swap := &action.Swap{
		Base: action.Base{
			Type:     action.SWAP,
			ChainId:  app.ChainId,
			Signers:  signers,
			Sequence: global.Current.Sequence,
		},
		Party:        party,
		CounterParty: counterparty,
		Fee:          fee,
		Gas:          gas,
		Amount:       amount,
		Exchange:     exchange,
		Nonce:        swapargs.nonce,
	}

	signed := SignTransaction(action.Transaction(swap))
	packet := PackRequest(signed)

	return packet
}

// IssueRequest sends out a sendTx to all of the nodes in the chain
func SwapCurrency(cmd *cobra.Command, args []string) {
	log.Debug("Swap Request", "tx", swapargs)

	// Create message
	packet := CreateSwapRequest()

	result := Broadcast(packet)

	log.Debug("Returned Successfully", "result", result)
}

func GetAddress(value string) id.Address {
	return id.Address{}
}

func GetCurrency(value string) string {
	// TODO: Check to see that this is a valid currency
	return value
}

func GetInteger(value string) int64 {
	return -1
}
