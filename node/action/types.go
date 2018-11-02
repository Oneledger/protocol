/*
	Copyright 2017-2018 OneLedger

	Declare basic types used by the Application

	If a type requires functions or a few types are intertwinded, then should be in their own file.
*/
package action

import (
	"bytes"

	"strconv"

	"github.com/Oneledger/protocol/node/data"
	"github.com/Oneledger/protocol/node/id"
	"github.com/Oneledger/protocol/node/log"
	"github.com/Oneledger/protocol/node/serial"
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

func init() {
	serial.Register(SendInput{})
	serial.Register(SendOutput{})
	serial.Register(Event{})
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

func CheckBalance(app interface{}, accountKey id.AccountKey, amount data.Coin) bool {
	utxo := GetUtxo(app)

	balance := utxo.Get(accountKey)
	if balance == nil {
		// New accounts don't have a balance until the first transaction
		log.Debug("New Balance", "key", accountKey, "amount", amount, "balance", balance)
		interim := data.NewBalance(0, "OLT")
		balance = &interim
		if !balance.Amount.Equals(amount) {
			return false
		}
		return true
	}

	if !balance.Amount.Equals(amount) {
		log.Warn("Balance Mismatch", "key", accountKey, "amount", amount, "balance", balance)
		return false
	}
	return true
}

func GetHeight(app interface{}) int64 {
	utxo := GetUtxo(app)

	height := int64(utxo.Version)
	return height
}

func CheckAmounts(app interface{}, inputs []SendInput, outputs []SendOutput) bool {
	total := data.NewCoin(0, "OLT")
	for _, input := range inputs {
		if input.Amount.LessThan(0) {
			log.Debug("Less Than 0", "input", input)
			return false
		}

		if !input.Amount.IsCurrency("OLT") {
			log.Debug("Send on Currency isn't implement yet")
			return false
		}

		if bytes.Compare(input.AccountKey, []byte("")) == 0 {
			log.Debug("Key is Empty", "input", input)
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
			log.Debug("Less Than 0", "output", output)
			return false
		}

		if !output.Amount.IsCurrency("OLT") {
			log.Debug("Send on Currency isn't implement yet")
			return false
		}

		if bytes.Compare(output.AccountKey, []byte("")) == 0 {
			log.Debug("Key is Empty", "output", output)
			return false
		}
		total.Minus(output.Amount)
	}
	if !total.Equals(data.NewCoin(0, "OLT")) {
		log.Debug("Doesn't add up", "inputs", inputs, "outputs", outputs)
		return false
	}
	return true
}

type Event struct {
	Type        Type   `json:"type"`
	SwapKeyHash []byte `json:"swapkeyhash"`
	Stage       int    `json:"stage"`
}

func (e Event) ToKey() []byte {
	buffer, err := serial.Serialize(e, serial.CLIENT)
	if err != nil {
		log.Error("Failed to Serialize event key")
	}
	return buffer
}

func SaveEvent(app interface{}, eventKey Event, status bool) {
	events := GetEvent(app)

	log.Debug("Save Event", "key", eventKey)

	session := events.Begin()
	session.Set(eventKey.ToKey(), []byte(strconv.FormatBool(status)))
	session.Commit()
}

func FindEvent(app interface{}, eventKey Event) bool {
	log.Debug("Load Event", "key", eventKey)
	events := GetEvent(app)
	result := events.Get(eventKey.ToKey())
	if result == nil {
		return false
	}

	r, err := strconv.ParseBool(result.(string))
	if err != nil {
		return false
	}

	return r
}

func SaveContract(app interface{}, contractKey []byte, nonce int64, contract []byte) {
	contracts := GetContracts(app)
	n := strconv.AppendInt(contractKey, nonce, 10)
	log.Debug("Save contract", "key", contractKey, "afterNonce", n)
	session := contracts.Begin()
	session.Set(contractKey, contract)
	session.Commit()
}

func FindContract(app interface{}, contractKey []byte, nonce int64) []byte {
	log.Debug("Load Contract", "key", contractKey)
	contracts := GetContracts(app)
	n := strconv.AppendInt(contractKey, nonce, 10)
	result := contracts.Get(n)
	if result == nil {
		return nil
	}
	return result.([]byte)
}

func DeleteContract(app interface{}, contractKey []byte, nonce int64) {
	contracts := GetContracts(app)
	n := strconv.AppendInt(contractKey, nonce, 10)
	log.Debug("Delete contract", "key", contractKey, "afterNonce", n)
	session := contracts.Begin()
	session.Delete(n)
	session.Commit()
}
