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

type Forget struct {
	Base

	Target string `json:"target"`
}

func init() {
	serial.Register(Forget{})
}

func (transaction *Forget) Validate() status.Code {
	log.Debug("Validating Forget Transaction")
	return status.SUCCESS
}

func (transaction *Forget) ProcessCheck(app interface{}) status.Code {
	log.Debug("Processing Forget Transaction for CheckTx")
	return status.SUCCESS
}

func (transaction *Forget) ShouldProcess(app interface{}) bool {
	return true
}

func (transaction *Forget) ProcessDeliver(app interface{}) status.Code {

	commands := transaction.Resolve(app)

	//before loop of execute, lastResult is nil
	return commands.Execute(app)
}

func (transaction *Forget) Resolve(app interface{}) Commands {
	return []Command{}
}
