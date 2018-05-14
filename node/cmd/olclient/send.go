/*
	Copyright 2017-2018 OneLedger

	Cli to interact with a with the chain.
*/
package main

import (
	"os"

	"github.com/Oneledger/protocol/node/app"
	"github.com/Oneledger/protocol/node/convert"
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
	user     string
	to       string // the recipient
	from     string // the source
	amount   string
	fee      string
	gas      string // Optional
	currency string
	sequence int // Replay protection
}

var sendargs *SendArguments = &SendArguments{}

func init() {
	RootCmd.AddCommand(sendCmd)

	// Operational Parameters
	// TODO: Should be global flags?
	sendCmd.Flags().StringVarP(&app.Current.Transport, "transport", "t", "socket", "transport (socket | grpc)")
	sendCmd.Flags().StringVarP(&app.Current.Address, "address", "a", "tcp://127.0.0.1:46658", "full address")

	// Transaction Parameters
	sendCmd.Flags().StringVar(&sendargs.user, "user", "undefined", "user name")
	sendCmd.Flags().StringVar(&sendargs.to, "to", "undefined", "send recipient")
	sendCmd.Flags().StringVar(&sendargs.amount, "amount", "0", "specify an amount")
	sendCmd.Flags().StringVar(&sendargs.currency, "currency", "OLT", "the currency")
	sendCmd.Flags().StringVar(&sendargs.fee, "fee", "1", "include a fee")
	sendCmd.Flags().StringVar(&sendargs.gas, "gas", "1", "include gas")
	sendCmd.Flags().IntVarP(&sendargs.sequence, "sequence", "s", 1, "unique sequence number (replay protection)")
}

// CreateRequest builds and signs the transaction based on the arguments
func CreateRequest() []byte {
	signers := GetSigners()

	conv := convert.NewConvert()

	to := app.Address(conv.GetHash(sendargs.to))
	_ = to

	gas := app.Coin{
		Currency: conv.GetCurrency(sendargs.currency),
		Amount:   conv.GetInt64(sendargs.gas),
	}

	if conv.HasErrors() {
		Console.Error(conv.GetErrors())
		os.Exit(-1)
	}

	// Create base transaction
	send := &app.SendTransaction{
		TransactionBase: app.TransactionBase{
			Type:     app.SEND_TRANSACTION,
			ChainId:  app.ChainId,
			Signers:  signers,
			Sequence: sendargs.sequence,
		},
		Gas: gas,
	}

	signed := SignTransaction(app.Transaction(send))
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
