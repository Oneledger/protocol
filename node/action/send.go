/*
	Copyright 2017-2018 OneLedger

	An incoming transaction, send, swap, ready, verification, etc.
*/
package action

import (
	"github.com/Oneledger/protocol/node/data"
	"github.com/Oneledger/protocol/node/log"
	"github.com/Oneledger/protocol/node/serial"
	"github.com/Oneledger/protocol/node/status"
)

// Synchronize a swap between two users
type Send struct {
	Base

	Inputs  []SendInput  `json:"inputs"`
	Outputs []SendOutput `json:"outputs"`

	Gas data.Coin `json:"gas"`
	Fee data.Coin `json:"fee"`
}

func init() {
	serial.Register(Send{})
}

func (transaction *Send) Validate() status.Code {
	log.Debug("Validating Send Transaction")

	if transaction.Fee.LessThan(0) {
		log.Debug("Missing Fee", "send", transaction)
		return status.MISSING_DATA
	}

	if transaction.Gas.LessThan(0) {
		log.Debug("Missing Gas", "send", transaction)
		return status.MISSING_DATA
	}

	return status.SUCCESS
}

func (transaction *Send) ProcessCheck(app interface{}) status.Code {
	log.Debug("Processing Send Transaction for CheckTx")

	if !CheckAmounts(app, transaction.Inputs, transaction.Outputs) {
		log.Debug("FAILED", "inputs", transaction.Inputs, "outputs", transaction.Outputs)
		return status.INVALID
		//return status.SUCCESS
	}

	// TODO: Validate the transaction against the UTXO database, check tree
	balances := GetBalances(app)
	_ = balances

	return status.SUCCESS
}

func (transaction *Send) ShouldProcess(app interface{}) bool {
	return true
}

func (transaction *Send) ProcessDeliver(app interface{}) status.Code {
	log.Debug("Processing Send Transaction for DeliverTx")

	if !CheckAmounts(app, transaction.Inputs, transaction.Outputs) {
		return status.INVALID
	}

	balances := GetBalances(app)

	// Update the database to the final set of entries
	for _, entry := range transaction.Outputs {
		var balance *data.Balance
		result := balances.Get(entry.AccountKey)
		if result == nil {
			tmp := data.NewBalance()
			result = &tmp
		}
		balance = result
		balance.SetAmmount(entry.Amount)

		balances.Set(entry.AccountKey, *balance)
	}

	return status.SUCCESS
}

func (transaction *Send) Resolve(app interface{}) Commands {
	return []Command{}
}
