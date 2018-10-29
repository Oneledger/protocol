/*
	Copyright 2017-2018 OneLedger

	Cli to interact with a with the chain.
*/
package main

import (
	"github.com/Oneledger/protocol/node/action"
	"github.com/Oneledger/protocol/node/comm"
	"github.com/Oneledger/protocol/node/log"
	"github.com/spf13/cobra"
)

var utxoCmd = &cobra.Command{
	Use:   "utxo",
	Short: "Check utxo database",
	Run:   CheckUtxo,
}

// TODO: typing should be way better, see if cobr can help with this...
type UtxoArguments struct {
	account string
}

var utxo *UtxoArguments = &UtxoArguments{}

func init() {
	RootCmd.AddCommand(utxoCmd)
	utxoCmd.Flags().StringVar(&utxo.account, "account", "", "account key")
}

// Format the request into a query structure
func FormatUtxoRequest() []byte {
	return action.Message("Utxo=" + account.user)
}

// IssueRequest sends out a sendTx to all of the nodes in the chain
func CheckUtxo(cmd *cobra.Command, args []string) {

	request := FormatUtxoRequest()
	response := comm.Query("/utxo", request)
	if response != nil {
		log.Debug("Returned Successfully with", "response", response)
	} else {
		log.Debug("No Response from Node", "request", request)
	}
}
