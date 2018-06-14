/*
	Copyright 2017-2018 OneLedger

	An incoming transaction, send, swap, ready, verification, etc.
*/
package action

import (
	"github.com/Oneledger/protocol/node/data"
	"github.com/Oneledger/protocol/node/err"
	"github.com/Oneledger/protocol/node/log"
)

// Synchronize a swap between two users
type ExternalLock struct {
	Base

	Gas     data.Coin    `json:"gas"`
	Fee     data.Coin    `json:"fee"`
	Inputs  []SendInput  `json:"inputs"`
	Outputs []SendOutput `json:"outputs"`
}

func (transaction *ExternalLock) Validate() err.Code {
	log.Debug("Validating ExternalLock Transaction")

	// TODO: Make sure all of the parameters are there
	// TODO: Check all signatures and keys
	// TODO: Vet that the sender has the values
	return err.SUCCESS
}

func (transaction *ExternalLock) ProcessCheck(app interface{}) err.Code {
	log.Debug("Processing ExternalLock Transaction for CheckTx")

	// TODO: // Update in memory copy of Merkle Tree
	return err.SUCCESS
}

func (transaction *ExternalLock) ShouldProcess(app interface{}) bool {
	return true
}

func (transaction *ExternalLock) ProcessDeliver(app interface{}) err.Code {
	log.Debug("Processing ExternalLock Transaction for DeliverTx")

	// TODO: // Update in final copy of Merkle Tree
	return err.SUCCESS
}

// Given a transaction, expand it into a list of Commands to execute against various chains.
func (transaction *ExternalLock) Expand(app interface{}) Commands {
	// TODO: Table-driven mechanics, probably elsewhere
	return []Command{}
}
