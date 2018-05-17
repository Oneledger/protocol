/*
	Copyright 2017-2018 OneLedger

	An incoming transaction, send, swap, ready, verification, etc.
*/
package action

import (
	"github.com/Oneledger/protocol/node/err"
	"github.com/Oneledger/protocol/node/id"
	"github.com/Oneledger/protocol/node/log"
)

// Synchronize a swap between two users
type SwapTransaction struct {
	TransactionBase

	Party1   id.Address `json:"party1"`
	Party2   id.Address `json:"party2"`
	Fee      Coin       `json:"fee"`
	Gas      Coin       `json:"fee"`
	Amount   Coin       `json:"amount"`
	Exchange Coin       `json:"exchange"`
}

// Issue swaps across other chains, make sure fees are collected
func (transaction *SwapTransaction) Validate() err.Code {
	log.Debug("Validating Swap Transaction")
	return err.SUCCESS
}

func (transaction *SwapTransaction) ProcessCheck(app interface{}) err.Code {
	log.Debug("Processing Swap Transaction for CheckTx")
	return err.SUCCESS
}

func (transaction *SwapTransaction) ProcessDeliver(app interface{}) err.Code {
	log.Debug("Processing Swap Transaction for DeliverTx")
	return err.SUCCESS
}
