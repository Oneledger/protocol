/*
	Copyright 2017-2018 OneLedger

	An incoming transaction, send, swap, ready, verification, etc.
*/
package action

import (
	"github.com/Oneledger/protocol/node/data"
	"github.com/Oneledger/protocol/node/id"
	"github.com/Oneledger/protocol/node/log"
	"github.com/Oneledger/protocol/node/serial"
	"github.com/Oneledger/protocol/node/status"
)

// Apply a dynamic validator
type ApplyValidator struct {
	Base

	AccountKey id.AccountKey

	TendermintAddress string
	TendermintPubKey  string

	Stake data.Coin
}

func init() {
	serial.Register(ApplyValidator{})
}

func (transaction *ApplyValidator) Validate() status.Code {
	log.Debug("Validating ApplyValidator Transaction")

	if transaction.Stake.LessThan(0) {
		log.Debug("Missing Stake", "ApplyValidator", transaction)
		return status.MISSING_DATA
	}

	return status.SUCCESS
}

func (transaction *ApplyValidator) ProcessCheck(app interface{}) status.Code {
	log.Debug("Processing ApplyValidator Transaction for CheckTx")

	return status.SUCCESS
}

func (transaction *ApplyValidator) ShouldProcess(app interface{}) bool {
	return true
}

func (transaction *ApplyValidator) ProcessDeliver(app interface{}) status.Code {
	log.Debug("Processing ApplyValidator Transaction for DeliverTx")

	return status.SUCCESS
}

func (transaction *ApplyValidator) Resolve(app interface{}) Commands {
	return []Command{}
}
