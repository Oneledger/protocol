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
type ExternalSend struct {
	Base

	Gas data.Coin `json:"gas"`
	Fee data.Coin `json:"fee"`

	ExGas    data.Coin      `json:"exgas"`
	ExFee    data.Coin      `json:"exfee"`
	Chain    data.ChainType `json:"chain"`
	Sender   string         `json:"sender"`
	Receiver string         `json:"receiver"`
	Amount   data.Coin      `json:"amount"`
}

func init() {
	serial.Register(ExternalSend{})
}

func (transaction *ExternalSend) Validate() status.Code {
	log.Debug("Validating ExternalSend Transaction")

	baseValidate := transaction.Base.Validate()

	if baseValidate != status.SUCCESS {
		return baseValidate
	}

	if transaction.Fee.LessThan(0) {
		log.Debug("Invalid Fee", "transaction", transaction)
		return status.BAD_VALUE
	}

	if transaction.Gas.LessThan(0) {
		log.Debug("Invalid Gas", "transaction", transaction)
		return status.BAD_VALUE
	}

	if transaction.ExGas.LessThan(0) {
		log.Debug("Invalid ExGas", "transaction", transaction)
		return status.BAD_VALUE
	}

	if transaction.ExFee.LessThan(0) {
		log.Debug("Invalid ExFee", "transaction", transaction)
		return status.BAD_VALUE
	}

	if transaction.Sender == "" {
		log.Debug("Missing Sender", "transaction", transaction)
		return status.MISSING_DATA
	}

	if transaction.Receiver == "" {
		log.Debug("Missing Receiver", "transaction", transaction)
		return status.MISSING_DATA
	}

	if transaction.Amount.LessThan(0) {
		log.Debug("Invalid Amount", "transaction", transaction)
		return status.BAD_VALUE
	}

	return status.SUCCESS
}

func (transaction *ExternalSend) ProcessCheck(app interface{}) status.Code {
	log.Debug("Processing ExternalSend Transaction for CheckTx")

	// TODO: // Update in memory copy of Merkle Tree
	return status.SUCCESS
}

func (transaction *ExternalSend) ShouldProcess(app interface{}) bool {
	return true
}

func (transaction *ExternalSend) ProcessDeliver(app interface{}) status.Code {
	log.Debug("Processing ExternalSend Transaction for DeliverTx")

	commands := transaction.Resolve(app)

	return commands.Execute(app)
}

func (transaction *ExternalSend) Resolve(app interface{}) Commands {
	return []Command{}
}
