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

	SendTo []SendTo `json:"payto"`
}

func init() {
	serial.Register(Payment{})
}

func (transaction *Payment) Validate() status.Code {
	log.Debug("Validating Payment Transaction")

	baseValidate := transaction.Base.Validate()

	if baseValidate != status.SUCCESS {
		return baseValidate
	}

	if len(transaction.SendTo) < 3 {
		return status.INVALID
	}
	return status.SUCCESS
}

func (transaction *Payment) ProcessCheck(app interface{}) status.Code {
	log.Debug("Processing Payment Transaction for CheckTx")

	if !CheckValidatorList(app, transaction.SendTo) {
		return status.INVALID
	}

	if !CheckPayTo(app, transaction.SendTo) {
		log.Debug("FAILED to ", "payto", transaction.SendTo)
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

	if !CheckValidatorList(app, transaction.SendTo) {
		return status.INVALID
	}

	if !CheckPayTo(app, transaction.SendTo) {
		return status.INVALID
	}

	chain := GetBalances(app)
	accounts := GetAccounts(app)
	payment, err := accounts.FindName(global.Current.PaymentAccount)
	if err != status.SUCCESS {
		log.Error("Failed to get payment account", "status", err)
		return err
	}
	paymentBalance := chain.Get(payment.AccountKey(), false)

	// Update the database to the final set of entries
	for _, entry := range transaction.SendTo {
		var balance *data.Balance
		result := chain.Get(entry.AccountKey, false)
		if result == nil {
			result = data.NewBalance()
		}
		balance = result
		balance.AddAmount(entry.Amount)

		chain.Set(entry.AccountKey, balance)
		paymentBalance.MinusAmount(entry.Amount)
	}

	chain.Set(payment.AccountKey(), paymentBalance)

	admin := GetAdmin(app)

	//store payment record in database (O OLT, -1) because delete doesn't work
	var paymentRecordKey data.DatabaseKey = data.DatabaseKey("PaymentRecord")
	var paymentRecord PaymentRecord
	paymentRecord.Amount = data.NewCoinFromInt(0, "OLT")
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

	balance := balances.Get(payment.AccountKey(), false)
	total := data.NewBalance()
	for _, v := range pay {
		if v.Amount.LessThan(0) {
			return false
		}
		total.AddAmount(v.Amount)
	}
	if !balance.IsEnoughBalance(*total) {
		return false
	}
	return true
}

func CheckValidatorList(app interface{}, SendTo []SendTo) bool {
	validators := GetValidators(app)

	for i, v := range SendTo {
		if !validators.IsValidAccountKey(v.AccountKey, i) {
			return false
		}
	}

	return true
}
