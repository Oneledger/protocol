/*
	Copyright 2017-2018 OneLedger

	An incoming transaction, send, swap, ready, verification, etc.
*/
package action

import (
	"github.com/Oneledger/protocol/node/status"
	"github.com/Oneledger/protocol/node/log"
)

// TODO: This needs to be filled out properly, need a model for other-chain actions...
type Commit struct {
	Base

	Target string `json:"target"`
}

func (transaction *Commit) Validate() status.Code {
	log.Debug("Validating Commit Transaction")
	if transaction.Target == "" {
		return status.MISSING_DATA
	}
	return status.SUCCESS
}

func (transaction *Commit) ProcessCheck(app interface{}) status.Code {
	log.Debug("Processing Commit Transaction for CheckTx")
	return status.SUCCESS
}

func (transaction *Commit) ShouldProcess(app interface{}) bool {
	return true
}

func (transaction *Commit) ProcessDeliver(app interface{}) status.Code {
	log.Debug("Processing Commit Transaction for DeliverTx")

	commands := transaction.Resolve(app)

	return commands.Execute(app)
}

func (transaction *Commit) Resolve(app interface{}) Commands {
	return []Command{}
}

