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
type Payment struct {
	Base

	Inputs  []SendInput  `json:"inputs"`
	Outputs []SendOutput `json:"outputs"`

	Gas data.Coin `json:"gas"`
	Fee data.Coin `json:"fee"`
}

func init() {
	serial.Register(Payment{})
}

func (transaction *Payment) Validate() status.Code {
	log.Debug("Validating Payment Transaction")

	if transaction.Fee.LessThan(0) {
		log.Debug("Missing Fee", "payment", transaction)
		return status.MISSING_DATA
	}

	if transaction.Gas.LessThan(0) {
		log.Debug("Missing Gas", "payment", transaction)
		return status.MISSING_DATA
	}

	return status.SUCCESS
}

func (transaction *Payment) ProcessCheck(app interface{}) status.Code {
	log.Debug("Processing Payment Transaction for CheckTx")

	if !CheckAmounts(app, transaction.Inputs, transaction.Outputs) {
		log.Debug("FAILED", "inputs", transaction.Inputs, "outputs", transaction.Outputs)
		return status.INVALID
	}

	// TODO: Validate the transaction against the UTXO database, check tree
	chain := GetBalances(app)
	_ = chain

	return status.SUCCESS
}

func (transaction *Payment) ShouldProcess(app interface{}) bool {
	return true
}

type PaymentRecord struct {
	Amount      data.Coin
	BlockHeight int64
}

func init() {
	serial.Register(PaymentRecord{})
}

func (transaction *Payment) ProcessDeliver(app interface{}) status.Code {
	log.Debug("Processing Payment Transaction for DeliverTx")

	if !CheckAmounts(app, transaction.Inputs, transaction.Outputs) {
		return status.INVALID
	}

	chain := GetBalances(app)

	// Update the database to the final set of entries
	for _, entry := range transaction.Outputs {
		var balance *data.Balance
		result := chain.Get(entry.AccountKey)
		if result == nil {
			tmp := data.NewBalance()
			result = &tmp
		}
		balance = result
		balance.SetAmmount(entry.Amount)

		chain.Set(entry.AccountKey, *balance)
	}

	admin := GetAdmin(app)

	//store payment record in database (O OLT, -1) because delete doesn't work
	var paymentRecordKey data.DatabaseKey = data.DatabaseKey("PaymentRecord")
	var paymentRecord PaymentRecord
	paymentRecord.Amount = data.NewCoin(0, "OLT")
	paymentRecord.BlockHeight = -1

	session := admin.Begin()
	session.Set(paymentRecordKey, paymentRecord)
	session.Commit()

	return status.SUCCESS
}

func (transaction *Payment) Resolve(app interface{}) Commands {
	return []Command{}
}
