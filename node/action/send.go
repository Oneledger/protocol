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
type Send struct {
	Base

	Gas     data.Coin    `json:"gas"`
	Fee     data.Coin    `json:"fee"`
	Inputs  []SendInput  `json:"inputs"`
	Outputs []SendOutput `json:"outputs"`
}

func (transaction *Send) Validate() err.Code {
	log.Debug("Validating Send Transaction")

	// TODO: Make sure all of the parameters are there
	// TODO: Check all signatures and keys
	// TODO: Vet that the sender has the values
	return err.SUCCESS
}

func (transaction *Send) ProcessCheck(app interface{}) err.Code {
	log.Debug("Processing Send Transaction for CheckTx")

	// TODO: // Update in memory copy of Merkle Tree
	return err.SUCCESS
}

func (transaction *Send) ProcessDeliver(app interface{}) err.Code {
	log.Debug("Processing Send Transaction for DeliverTx")

	// TODO: // Update in final copy of Merkle Tree
	return err.SUCCESS
}
