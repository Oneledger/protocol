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
type ExternalLock struct {
	Base

	Gas     data.Coin    `json:"gas"`
	Fee     data.Coin    `json:"fee"`
	Inputs  []SendInput  `json:"inputs"`
	Outputs []SendOutput `json:"outputs"`
}

func init() {
	serial.Register(ExternalLock{})
}

func (transaction *ExternalLock) Validate() status.Code {
	log.Debug("Validating ExternalLock Transaction")

	if transaction.Fee.LessThan(0) {
		return status.MISSING_DATA
	}
	return status.SUCCESS
}

func (transaction *ExternalLock) ProcessCheck(app interface{}) status.Code {
	log.Debug("Processing ExternalLock Transaction for CheckTx")

	commands := transaction.Resolve(app)

	return commands.Execute(app)
}

func (transaction *ExternalLock) ShouldProcess(app interface{}) bool {
	return true
}

func (transaction *ExternalLock) ProcessDeliver(app interface{}) status.Code {
	log.Debug("Processing ExternalLock Transaction for DeliverTx")

	// TODO: // Update in final copy of Merkle Tree
	return status.SUCCESS
}

func (transaction *ExternalLock) Resolve(app interface{}) Commands {
	return []Command{}
}
