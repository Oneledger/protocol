/*
	Copyright 2017-2018 OneLedger

	An incoming transaction, send, swap, ready, verification, etc.
*/
package action

import (
	"github.com/Oneledger/protocol/node/log"
	"github.com/Oneledger/protocol/node/serial"
	"github.com/Oneledger/protocol/node/status"
)

// TODO: This needs to be filled out properly, need a model for other-chain actions...
type Prepare struct {
	Base

	Target string `json:"target"`
}

func init() {
	serial.Register(Prepare{})
}

func (transaction *Prepare) Validate() status.Code {
	log.Debug("Validating Prepare Transaction")
	return status.SUCCESS
}

func (transaction *Prepare) ProcessCheck(app interface{}) status.Code {
	log.Debug("Processing Prepare Transaction for CheckTx")
	return status.SUCCESS
}

func (transaction *Prepare) ShouldProcess(app interface{}) bool {
	return true
}

func (transaction *Prepare) ProcessDeliver(app interface{}) status.Code {

	commands := transaction.Resolve(app)

	//before loop of execute, lastResult is nil
	return commands.Execute(app)
}

func (transaction *Prepare) Resolve(app interface{}) Commands {
	return []Command{}
}
