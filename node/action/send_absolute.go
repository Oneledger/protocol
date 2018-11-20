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
type Send_Abusolute struct {
	Base

	Inputs  []SendInput  `json:"inputs"`
	Outputs []SendOutput `json:"outputs"`

	Gas data.Coin `json:"gas"`
	Fee data.Coin `json:"fee"`
}

func init() {
	serial.Register(Send_Abusolute{})
}

func (transaction *Send_Abusolute) TransactionType() Type {
	return transaction.Base.Type
}

func (transaction *Send_Abusolute) Validate() status.Code {
	log.Debug("Validating Send_Abusolute Transaction")

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

func (transaction *Send_Abusolute) ProcessCheck(app interface{}) status.Code {
	log.Debug("Processing Send_Abusolute Transaction for CheckTx")

	if !CheckAmountsAbsolute(app, transaction.Inputs, transaction.Outputs) {
		log.Debug("FAILED", "inputs", transaction.Inputs, "outputs", transaction.Outputs)
		return status.INVALID
		//return status.SUCCESS
	}

	// TODO: Validate the transaction against the UTXO database, check tree
	balances := GetBalances(app)
	_ = balances

	return status.SUCCESS
}

func (transaction *Send_Abusolute) ShouldProcess(app interface{}) bool {
	return true
}

func (transaction *Send_Abusolute) ProcessDeliver(app interface{}) status.Code {
	log.Debug("Processing Send_Abusolute Transaction for DeliverTx")

	if !CheckAmountsAbsolute(app, transaction.Inputs, transaction.Outputs) {
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

func (transaction *Send_Abusolute) Resolve(app interface{}) Commands {
	return []Command{}
}
