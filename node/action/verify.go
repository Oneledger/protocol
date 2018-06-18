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

// TODO: This needs to be filled out properly, need a model for other-chain actions...
type Verify struct {
	Base

	Target id.AccountKey `json:"target"`
}

func (transaction Verify) Validate() err.Code {
	log.Debug("Validating Verify Transaction")
	return err.SUCCESS
}

func (transaction Verify) ProcessCheck(app interface{}) err.Code {
	log.Debug("Processing Verify Transaction for CheckTx")
	return err.SUCCESS
}

func (transaction Verify) ShouldProcess(app interface{}) bool {
	return true
}

func (transaction Verify) ProcessDeliver(app interface{}) err.Code {
	log.Debug("Processing Verify Transaction for DeliverTx")

	commands := transaction.Expand(app)
	transaction.Resolve(app, commands)

	for i := 0; i < commands.Count(); i++ {
		status := Execute(app, commands[i])
		if status != err.SUCCESS {
			return err.EXPAND_ERROR
		}
	}
	return err.SUCCESS
}

func (transaction Verify) Resolve(app interface{}, commands Commands) {
}

// Given a transaction, expand it into a list of Commands to execute against various chains.
func (transaction Verify) Expand(app interface{}) Commands {
	// TODO: Table-driven mechanics, probably elsewhere
	return []Command{}
}
