/*
	Copyright 2017-2018 OneLedger

	An incoming transaction, send, swap, ready, verification, etc.
*/
package action

import (
	"github.com/Oneledger/protocol/node/data"
	"github.com/Oneledger/protocol/node/global"
	"github.com/Oneledger/protocol/node/log"
	"github.com/Oneledger/protocol/node/serial"
	"github.com/Oneledger/protocol/node/status"
)

// Synchronize a swap between two users
type Payment struct {
	Base

	PayTo []SendTo `json:"payto"`

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
	if !CheckPayTo(app, transaction.PayTo) {
		log.Debug("FAILED to ", "payto", transaction.PayTo)
		return status.INVALID
	}

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

	if !CheckPayTo(app, transaction.PayTo) {
		return status.INVALID
	}

	chain := GetBalances(app)

	// Update the database to the final set of entries
	for _, entry := range transaction.PayTo {
		var balance *data.Balance
		result := chain.Get(entry.AccountKey)
		if result == nil {
			result = data.NewBalance()
		}
		balance = result
		balance.AddAmount(entry.Amount)

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

func CheckPayTo(app interface{}, pay []SendTo) bool {
	balances := GetBalances(app)
	accounts := GetAccounts(app)

	payment, ok := accounts.FindName(global.Current.PaymentAccount)
	if ok != status.SUCCESS {
		log.Error("Failed to get payment account", "status", ok)
		return false
	}

	balance := balances.Get(payment.AccountKey())
	total := data.NewBalance()
	for _, v := range pay {
		total.AddAmount(v.Amount)
	}
	if !balance.IsEnoughBalance(*total) {
		return false
	}
	return true
}
