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
type ReadyTransaction struct {
	TransactionBase

	Target string `json:"target"`
}

func (transaction *ReadyTransaction) Validate() err.Code {
	log.Debug("Validating Ready Transaction")
	return err.SUCCESS
}

func (transaction *ReadyTransaction) ProcessCheck(app interface{}) err.Code {
	log.Debug("Processing Ready Transaction for CheckTx")
	return err.SUCCESS
}

func (transaction *ReadyTransaction) ProcessDeliver(app interface{}) err.Code {
	log.Debug("Processing Ready Transaction for DeliverTx")
	return err.SUCCESS
}
