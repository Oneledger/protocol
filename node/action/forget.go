/*
	Copyright 2017-2018 OneLedger

	An incoming transaction, send, swap, ready, verification, etc.
*/
package action

import (
	"github.com/Oneledger/protocol/node/err"
	"github.com/Oneledger/protocol/node/log"
)

// TODO: This needs to be filled out properly, need a model for other-chain actions...
type Forget struct {
	Base

	Target string `json:"target"`
}

func (transaction *Forget) Validate() err.Code {
	log.Debug("Validating Forget Transaction")
	return err.SUCCESS
}

func (transaction *Forget) ProcessCheck(app interface{}) err.Code {
	log.Debug("Processing Forget Transaction for CheckTx")
	return err.SUCCESS
}

func (transaction *Forget) ProcessDeliver(app interface{}) err.Code {
	log.Debug("Processing Forget Transaction for DeliverTx")
	return err.SUCCESS
}

// Given a transaction, expand it into a list of Commands to execute against various chains.
func (transaction *Forget) Expand(app interface{}) Commands {
	// TODO: Table-driven mechanics, probably elsewhere
	return []Command{}
}
