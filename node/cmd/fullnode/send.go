/*
	Copyright 2017-2018 OneLedger

	Cli to interact with a with the chain.
*/
package main

import (
	"os"

	"github.com/Oneledger/prototype/node/app"
	"github.com/spf13/cobra"
	wire "github.com/tendermint/go-wire"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
)

var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Issue send transaction",
	Run:   IssueRequest,
}

// TODO: typing should be way better, see if cobr can help with this...
type TransactionArguments struct {
	to       string
	amount   string
	fee      string
	gas      uint64
	sequence int
}

var transaction *TransactionArguments

func init() {
	transaction = &TransactionArguments{}

	RootCmd.AddCommand(sendCmd)

	// Operational Parameters
	sendCmd.Flags().StringVarP(&app.Current.Transport, "transport", "t", "socket", "transport (socket | grpc)")
	sendCmd.Flags().StringVarP(&app.Current.Address, "address", "a", "tcp://127.0.0.1:46658", "full address")

	// Transaction Parameters
	sendCmd.Flags().StringVar(&transaction.to, "to", "undefined", "send recipient")
	sendCmd.Flags().StringVar(&transaction.amount, "amount", "0olt", "specify an amount")
	sendCmd.Flags().StringVar(&transaction.fee, "fee", "1olt", "include a fee")
}

// Convert the arguments?
func HandleSendArguments() {
}

// SignTransaction with the local keys
func SignTransaction(full *app.FullSendTransaction) {
}

// Pack a request into a transferable format (wire)
func PackRequest(request *app.FullSendTransaction) []byte {
	packet := wire.BinaryBytes(request)
	return packet
}

// CreateRequest builds and signs the transaction based on the arguments
func CreateRequest() []byte {
	// Create base transaction
	send := &app.SendTransaction{Type: app.SEND_TRANSACTION}
	full := &app.FullSendTransaction{Transaction: send}

	// Sign it
	SignTransaction(full)

	// Encode the message
	packet := PackRequest(full)

	return packet
}

// IssueRequest sends out a sendTx to all of the nodes in the chain
func IssueRequest(cmd *cobra.Command, args []string) {
	app.Log.Info("Issuing a client request")

	app.Log.Debug("Have Request", "tx", transaction)

	// Create message
	packet := CreateRequest()

	app.Log.Debug("Creating Client")
	client := rpcclient.NewHTTP("127.0.0.1:46657", "/websocket")

	result, err := client.BroadcastTxCommit(packet)
	if err != nil {
		app.Log.Error("Error", "err", err)
		os.Exit(-1)
	}
	app.Log.Debug("Returned Successfully", "result", result)

}
