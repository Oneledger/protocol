/*
	Copyright 2017-2018 OneLedger

	An incoming transaction, send, swap, ready, verification, etc.
*/
package action

import (
	"github.com/Oneledger/protocol/node/data"
	"github.com/Oneledger/protocol/node/err"
	"github.com/Oneledger/protocol/node/id"
	"github.com/Oneledger/protocol/node/log"
)

// Synchronize a swap between two users
type Swap struct {
	TransactionBase

	Party1   id.Address `json:"party1"`
	Party2   id.Address `json:"party2"`
	Fee      data.Coin  `json:"fee"`
	Gas      data.Coin  `json:"fee"`
	Amount   data.Coin  `json:"amount"`
	Exchange data.Coin  `json:"exchange"`
}

// Issue swaps across other chains, make sure fees are collected
func (transaction *Swap) Validate() err.Code {
	log.Debug("Validating Swap Transaction")
	return err.SUCCESS
}

func (transaction *Swap) ProcessCheck(app interface{}) err.Code {
	log.Debug("Processing Swap Transaction for CheckTx")
	return err.SUCCESS
}

func (transaction *Swap) ProcessDeliver(app interface{}) err.Code {
	log.Debug("Processing Swap Transaction for DeliverTx")
	return err.SUCCESS
}
