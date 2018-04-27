/*
	Copyright 2017-2018 OneLedger

	Cli to interact with a with the chain.
*/
package main

import (
	"github.com/Oneledger/prototype/node/app"
	"github.com/Oneledger/prototype/node/log"
	"github.com/spf13/cobra"
)

var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Issue send transaction",
	Run:   IssueRequest,
}

// TODO: typing should be way better, see if cobr can help with this...
type TransactionArguments struct {
	to       string // the recipient
	from     string
	amount   string
	fee      string
	gas      string // TODO: Not sure this is necessary, unless the chain is like Ethereum
	currency string
	sequence int // Replay protection
}

var transaction *TransactionArguments = &TransactionArguments{}

func init() {
	RootCmd.AddCommand(sendCmd)

	// Operational Parameters
	sendCmd.Flags().StringVarP(&app.Current.Transport, "transport", "t", "socket", "transport (socket | grpc)")
	sendCmd.Flags().StringVarP(&app.Current.Address, "address", "a", "tcp://127.0.0.1:46658", "full address")

	// Transaction Parameters
	sendCmd.Flags().StringVar(&transaction.to, "to", "undefined", "send recipient")
	sendCmd.Flags().StringVar(&transaction.amount, "amount", "0", "specify an amount")
	sendCmd.Flags().StringVar(&transaction.currency, "currency", "OLT", "the currency")
	sendCmd.Flags().StringVar(&transaction.fee, "fee", "1", "include a fee")
}

// CreateRequest builds and signs the transaction based on the arguments
func CreateRequest() []byte {
	signers := GetSigners()

	// Create base transaction
	transaction := &app.SwapTransaction{
		TransactionBase: app.TransactionBase{
			Type:    app.SWAP_TRANSACTION,
			ChainId: app.ChainId,
			Signers: signers,
		},
	}

	signed := SignTransaction(app.Transaction(transaction))
	packet := PackRequest(signed)

	return packet
}

// IssueRequest sends out a sendTx to all of the nodes in the chain
func IssueRequest(cmd *cobra.Command, args []string) {
	log.Debug("Have Request", "tx", transaction)

	// Create message
	packet := CreateRequest()

	result := Broadcast(packet)

	log.Debug("Returned Successfully", "result", result)
}
