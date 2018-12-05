/*
	Copyright 2017-2018 OneLedger

	An incoming transaction, send, swap, ready, verification, etc.
*/
package action

import (
	"bytes"
	"github.com/Oneledger/protocol/node/data"
	"github.com/Oneledger/protocol/node/id"
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
	serial.Register(SendInput{})
	serial.Register(SendOutput{})
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

//todo: We probably don't need SendInput and SendOutput anymore

// inputs into a send transaction (similar to Bitcoin)
type SendInput struct {
	AccountKey id.AccountKey `json:"account_key"`
	PubKey     PublicKey     `json:"pub_key"`
	Signature  id.Signature  `json:"signature"`

	Amount data.Coin `json:"coin"`

	// TODO: Is sequence needed per input?
	Sequence int `json:"sequence"`
}

func NewSendInput(accountKey id.AccountKey, amount data.Coin) SendInput {

	if bytes.Equal(accountKey, []byte("")) {
		log.Fatal("Missing AccountKey", "key", accountKey, "amount", amount)
		// TODO: Error handling should be better
		return SendInput{}
	}

	if !amount.IsValid() {
		log.Fatal("Missing Amount", "key", accountKey, "amount", amount)
	}

	return SendInput{
		AccountKey: accountKey,
		Amount:     amount,
	}
}

// outputs for a send transaction (similar to Bitcoin)
type SendOutput struct {
	AccountKey id.AccountKey `json:"account_key"`
	Amount     data.Coin     `json:"coin"`
}

func NewSendOutput(accountKey id.AccountKey, amount data.Coin) SendOutput {

	if bytes.Equal(accountKey, []byte("")) {
		log.Fatal("Missing AccountKey", "key", accountKey, "amount", amount)
		// TODO: Error handling should be better
		return SendOutput{}
	}

	if !amount.IsValid() {
		log.Fatal("Missing Amount", "key", accountKey, "amount", amount)
	}

	return SendOutput{
		AccountKey: accountKey,
		Amount:     amount,
	}
}

func CheckAmountsAbsolute(app interface{}, inputs []SendInput, outputs []SendOutput) bool {
	total := data.NewCoin(0, "OLT")
	for _, input := range inputs {
		if input.Amount.LessThan(0) {
			log.Debug("FAILED: Less Than 0", "input", input)
			return false
		}

		if !input.Amount.IsCurrency("OLT") {
			log.Debug("FAILED: Send_Abusolute on Currency isn't implement yet")
			return false
		}

		if bytes.Compare(input.AccountKey, []byte("")) == 0 {
			log.Debug("FAILED: Key is Empty", "input", input)
			return false
		}
		if !CheckBalance(app, input.AccountKey, input.Amount) {
			log.Warn("Balance is incorrect", "input", input)
			//return false
		}
		total.Plus(input.Amount)
	}
	for _, output := range outputs {

		if output.Amount.LessThan(0) {
			log.Debug("FAILED: Less Than 0", "output", output)
			return false
		}

		if !output.Amount.IsCurrency("OLT") {
			log.Debug("FAILED: Send_Abusolute on Currency isn't implement yet")
			return false
		}

		if bytes.Compare(output.AccountKey, []byte("")) == 0 {
			log.Debug("FAILED: Key is Empty", "output", output)
			return false
		}
		total.Minus(output.Amount)
	}
	if !total.Equals(data.NewCoin(0, "OLT")) {
		log.Debug("FAILED: Doesn't add up", "inputs", inputs, "outputs", outputs)
		return false
	}
	return true
}

func CheckBalance(app interface{}, accountKey id.AccountKey, amount data.Coin) bool {
	balances := GetBalances(app)

	balance := balances.Get(accountKey)
	if balance == nil {
		// New accounts don't have a balance until the first transaction
		log.Debug("New Balance", "key", accountKey, "amount", amount, "balance", balance)
		interim := data.NewBalanceFromString(0, amount.Currency.Name)
		balance = &interim
		if !balance.GetAmountByName(amount.Currency.Name).Equals(amount) {
			return false
		}
		return true
	}

	if !balance.GetAmountByName(amount.Currency.Name).Equals(amount) {
		log.Warn("Balance Mismatch", "key", accountKey, "amount", amount, "balance", balance)
		return false
	}
	return true
}
