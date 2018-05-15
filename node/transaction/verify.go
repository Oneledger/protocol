/*
	Copyright 2017-2018 OneLedger

	An incoming transaction, send, swap, ready, verification, etc.
*/
package transaction

import (
	"github.com/Oneledger/protocol/node/err"
	"github.com/Oneledger/protocol/node/log"
)

// TODO: This needs to be filled out properly, need a model for other-chain actions...
type VerifyTransaction struct {
	TransactionBase

	Target string `json:"target"`
}

func (transaction *VerifyTransaction) Validate() err.Code {
	log.Debug("Validating Verify Transaction")
	return err.SUCCESS
}

func (transaction *VerifyTransaction) ProcessCheck(app interface{}) err.Code {
	log.Debug("Processing Verify Transaction for CheckTx")
	return err.SUCCESS
}

func (transaction *VerifyTransaction) ProcessDeliver(app interface{}) err.Code {
	log.Debug("Processing Verify Transaction for DeliverTx")
	return err.SUCCESS
}
