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

var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Issue send transaction",
	Run:   IssueRequest,
}

// TODO: typing should be way better, see if cobr can help with this...
type SendArguments struct {
	party        string // the recipient
	counterparty string // the source
	amount       string
	fee          string
	gas          string // Optional
	currency     string
}

var sendargs *SendArguments = &SendArguments{}

func init() {
	RootCmd.AddCommand(sendCmd)

	// Transaction Parameters
	sendCmd.Flags().StringVar(&sendargs.party, "party", "undefined", "send recipient")
	sendCmd.Flags().StringVar(&sendargs.counterparty, "counterparty", "undefined", "send recipient")
	sendCmd.Flags().StringVar(&sendargs.amount, "amount", "0", "specify an amount")
	sendCmd.Flags().StringVar(&sendargs.currency, "currency", "OLT", "the currency")

	sendCmd.Flags().StringVar(&sendargs.fee, "fee", "1", "include a fee")
	sendCmd.Flags().StringVar(&sendargs.gas, "gas", "1", "include gas")
}

// CreateRequest builds and signs the transaction based on the arguments
func CreateRequest() []byte {
	signers := GetSigners()

	conv := convert.NewConvert()

	party := id.Address(conv.GetHash(sendargs.party))
	counterparty := id.Address(conv.GetHash(sendargs.counterparty))
	_ = party
	_ = counterparty

	gas := data.Coin{
		Currency: conv.GetCurrency(sendargs.currency),
		Amount:   conv.GetInt64(sendargs.gas),
	}

	if conv.HasErrors() {
		Console.Error(conv.GetErrors())
		os.Exit(-1)
	}

	// Create base transaction
	send := &action.Send{
		Base: action.Base{
			Type:     action.SEND,
			ChainId:  app.ChainId,
			Signers:  signers,
			Sequence: global.Current.Sequence,
		},
		Fee: gas,
		Gas: gas,
	}

	signed := SignTransaction(action.Transaction(send))
	packet := PackRequest(signed)

	return packet
}

// IssueRequest sends out a sendTx to all of the nodes in the chain
func IssueRequest(cmd *cobra.Command, args []string) {
	log.Debug("Have Send Request", "sendargs", sendargs)

	// Create message
	packet := CreateRequest()

	result := Broadcast(packet)

	log.Debug("Returned Successfully", "result", result)
}
