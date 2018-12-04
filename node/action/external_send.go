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

	// TODO: Make sure all of the parameters are there
	// TODO: Check all signatures and keys
	// TODO: Vet that the sender has the values
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
