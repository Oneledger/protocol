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
type Prepare struct {
	Base

	Target string `json:"target"`
}

func (transaction *Prepare) Validate() err.Code {
	log.Debug("Validating Prepare Transaction")
	return err.SUCCESS
}

func (transaction *Prepare) ProcessCheck(app interface{}) err.Code {
	log.Debug("Processing Prepare Transaction for CheckTx")
	return err.SUCCESS
}

func (transaction *Prepare) ShouldProcess(app interface{}) bool {
	return true
}

func (transaction *Prepare) ProcessDeliver(app interface{}) err.Code {

	commands := transaction.Resolve(app)

	//before loop of execute, lastResult is nil
	return commands.Execute(app)
}

func (transaction *Prepare) Resolve(app interface{}) Commands {
	return []Command{}
}
