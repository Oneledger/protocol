/*
	Copyright 2017-2018 OneLedger

	An incoming transaction, send, swap, ready, verification, etc.
*/
package action

import (
	"bytes"

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

	// TODO: Make sure all of the parameters are there
	// TODO: Check all signatures and keys
	// TODO: Vet that the sender has the values
	return err.SUCCESS
}

func (transaction *Send) ProcessCheck(app interface{}) err.Code {
	log.Debug("Processing Send Transaction for CheckTx")

	if !CheckAmounts(transaction.Inputs, transaction.Outputs) {
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

	if !CheckAmounts(transaction.Inputs, transaction.Outputs) {
		return err.INVALID
	}

	// TODO: Revalidate the transaction
	// TODO: Need to rollback if any errors occur

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

// Make sure the inputs and outputs all add up correctly.
func CheckAmounts(inputs []SendInput, outputs []SendOutput) bool {
	for _, input := range inputs {
		if input.Amount.LessThanEqual(0) {
			return false
		}
		if bytes.Compare(input.AccountKey, []byte("")) == 0 {
			return false
		}
	}
	for _, output := range outputs {
		if output.Amount.LessThanEqual(0) {
			return false
		}
		if bytes.Compare(output.AccountKey, []byte("")) == 0 {
			return false
		}
	}
	return true
}

// Given a transaction, expand it into a list of Commands to execute against various chains.
func (transaction *Send) Expand(app interface{}) Commands {
	// TODO: Table-driven mechanics, probably elsewhere
	return []Command{}
}
