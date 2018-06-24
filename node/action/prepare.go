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
	log.Debug("Processing Prepare Transaction for DeliverTx")

	commands := transaction.Expand(app)
	transaction.Resolve(app, commands)

	//before loop of execute, lastResult is nil
	var lastResult map[Parameter]FunctionValue
	var status err.Code

	for i := 0; i < commands.Count(); i++ {
		status, lastResult = Execute(app, commands[i], lastResult)
		if status != err.SUCCESS {
			return err.EXPAND_ERROR
		}
	}
	return err.SUCCESS
}

func (transaction *Prepare) Resolve(app interface{}, commands Commands) {
}

// Given a transaction, expand it into a list of Commands to execute against various chains.
func (transaction *Prepare) Expand(app interface{}) Commands {
	// TODO: Table-driven mechanics, probably elsewhere
	return []Command{}
}
