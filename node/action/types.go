/*
	Copyright 2017-2018 OneLedger

	Declare basic types used by the Application

	If a type requires functions or a few types are intertwinded, then should be in their own file.
*/
package action

import (
	"bytes"

	"github.com/Oneledger/protocol/node/comm"
	"github.com/Oneledger/protocol/node/data"
	"github.com/Oneledger/protocol/node/id"
	"github.com/Oneledger/protocol/node/log"
)

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
		log.Fatal("Missing Account")
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
		log.Fatal("Missing Account")
	}

	return SendOutput{
		AccountKey: accountKey,
		Amount:     amount,
	}
}

func CheckBalance(app interface{}, accountKey id.AccountKey, amount data.Coin) bool {
	utxo := GetUtxo(app)

	version := utxo.Delivered.Version64()
	_, value := utxo.Delivered.GetVersioned(accountKey, version)
	if value != nil {
		return false
	}

	var bal data.Balance
	buffer, _ := comm.Deserialize(value, bal)
	if buffer == nil {
		return false
	}

	balance := buffer.(data.Balance)
	if !balance.Amount.Equals(amount) {
		return false
	}
	return true
}

func CheckAmounts(app interface{}, inputs []SendInput, outputs []SendOutput) bool {
	total := data.NewCoin(0, "OLT")
	for _, input := range inputs {
		if input.Amount.LessThan(0) {
			return false
		}
		if bytes.Compare(input.AccountKey, []byte("")) == 0 {
			return false
		}
		if !CheckBalance(app, input.AccountKey, input.Amount) {
			return false
		}
		total.Plus(input.Amount)
	}
	for _, output := range outputs {
		if output.Amount.LessThan(0) {
			return false
		}
		if bytes.Compare(output.AccountKey, []byte("")) == 0 {
			return false
		}
		total.Minus(output.Amount)
	}
	if !total.Equals(data.NewCoin(0, "OLT")) {
		return false
	}
	return true
}
