/*
	Copyright 2017-2018 OneLedger

	An incoming transaction, send, swap, ready, verification, etc.
*/
package action

import (
	"github.com/Oneledger/protocol/node/comm"
	"github.com/Oneledger/protocol/node/data"
	"github.com/Oneledger/protocol/node/err"
	"github.com/Oneledger/protocol/node/log"
)

// Synchronize a swap between two users
type Send struct {
	Base

	Inputs  []SendInput  `json:"inputs"`
	Outputs []SendOutput `json:"outputs"`

	Gas data.Coin `json:"gas"`
	Fee data.Coin `json:"fee"`
}

func (transaction *Send) Validate() err.Code {
	log.Debug("Validating Send Transaction")

	if transaction.Fee.LessThan(0) {
		return err.MISSING_DATA
	}
	if transaction.Gas.LessThan(0) {
		return err.MISSING_DATA
	}

	return err.SUCCESS
}

func (transaction *Send) ProcessCheck(app interface{}) err.Code {
	log.Debug("Processing Send Transaction for CheckTx")

	if !CheckAmounts(app, transaction.Inputs, transaction.Outputs) {
		log.Debug("FAILED", "inputs", transaction.Inputs, "outputs", transaction.Outputs)
		return err.INVALID
	}

	// TODO: Validate the transaction against the UTXO database, check tree
	chain := GetUtxo(app)
	_ = chain

	return err.SUCCESS
}

func (transaction *Send) ShouldProcess(app interface{}) bool {
	return true
}

func (transaction *Send) ProcessDeliver(app interface{}) err.Code {
	log.Debug("Processing Send Transaction for DeliverTx")

	if !CheckAmounts(app, transaction.Inputs, transaction.Outputs) {
		return err.INVALID
	}

	chain := GetUtxo(app)

	// Update the database to the final set of entries
	for _, entry := range transaction.Outputs {
		balance := data.Balance{
			Amount: entry.Amount,
		}
		buffer, _ := comm.Serialize(balance)
		chain.Delivered.Set(entry.AccountKey, buffer)
	}

	return err.SUCCESS
}

func (transaction *Send) Resolve(app interface{}, commands Commands) {
}

// Given a transaction, expand it into a list of Commands to execute against various chains.
func (transaction *Send) Expand(app interface{}) Commands {
	// TODO: Table-driven mechanics, probably elsewhere
	return []Command{}
}
