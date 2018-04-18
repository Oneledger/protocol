/*
	Copyright 2017-2018 OneLedger

	Cli to interact with a with the chain.
*/
package main

import (
	"github.com/Oneledger/prototype/node/app"

	crypto "github.com/tendermint/go-crypto"
	wire "github.com/tendermint/go-wire"
	"github.com/tendermint/go-wire/data"
	rpcclient "github.com/tendermint/tendermint/rpc/client"

	//"golang.org/x/net/trace"

	"github.com/spf13/cobra"
)

var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Issue send transaction",
	Run:   IssueRequest,
}

type Coin struct {
	Denom  string `json:"denom"`
	Amount int64  `json:"amount"`
}

type Coins []Coin

type TxInput struct {
	Address   data.Bytes       `json:"address"`   // Hash of the PubKey
	Coins     Coins            `json:"coins"`     //
	Sequence  int              `json:"sequence"`  // Must be 1 greater than the last committed TxInput
	Signature crypto.Signature `json:"signature"` // Depends on the PubKey type and the whole Tx
	PubKey    crypto.PubKey    `json:"pub_key"`   // Is present iff Sequence == 0
}

type TxOutput struct {
	Address data.Bytes `json:"address"` // Hash of the PubKey
	Coins   Coins      `json:"coins"`   //
}

type SendTx struct {
	Gas     int64      `json:"gas"`
	Fee     Coin       `json:"fee"`
	Inputs  []TxInput  `json:"inputs"`
	Outputs []TxOutput `json:"outputs"`
}

// TODO: Basecoin structure, should revise
// TODO: typing should be way better
type Transaction struct {
	to       string
	amount   string
	fee      string
	gas      uint64
	sequence int
}

var transaction *Transaction

type FullSendTx struct {
	ChainId string
	Signers []crypto.PubKey
	Tx      *SendTx
}

func init() {
	transaction = &Transaction{}

	RootCmd.AddCommand(sendCmd)

	// Operational Parameters
	sendCmd.Flags().StringVarP(&app.Current.Transport, "transport", "t", "socket", "transport (socket | grpc)")
	sendCmd.Flags().StringVarP(&app.Current.Address, "address", "a", "tcp://127.0.0.1:46658", "full address")

	// Transaction Parameters
	sendCmd.Flags().StringVar(&transaction.to, "to", "undefined", "send recipient")
	sendCmd.Flags().StringVar(&transaction.amount, "amount", "0olt", "specify an amount")
	sendCmd.Flags().StringVar(&transaction.fee, "fee", "1olt", "include a fee")
}

func HandleSendArguments() {
}

func SignTransaction(full *FullSendTx) {
}

// IssueRequest sends out a sendTx to all of the nodes in the chain
func IssueRequest(cmd *cobra.Command, args []string) {
	app.Log.Info("Issuing a client request")

	app.Log.Debug("Have Request", "tx", transaction)

	// Create base transaction
	send := &SendTx{}
	full := &FullSendTx{Tx: send}

	// Sign it
	SignTransaction(full)

	// Create message
	packet := wire.BinaryBytes(full)
	_ = packet

	app.Log.Debug("Creating Client")
	/*
	 */
	client := rpcclient.NewHTTP("127.0.0.1:46657", "/websocket")
	result, err := client.BroadcastTxCommit(packet)
	if err != nil {
		app.Log.Error("Error", "err", err)
	}
	app.Log.Debug("Returned", "result", result)

}
