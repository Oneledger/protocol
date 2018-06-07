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
type Commit struct {
	Base

	Target string `json:"target"`
}

func (transaction *Commit) Validate() err.Code {
	log.Debug("Validating Commit Transaction")
	return err.SUCCESS
}

func (transaction *Commit) ProcessCheck(app interface{}) err.Code {
	log.Debug("Processing Commit Transaction for CheckTx")
	return err.SUCCESS
}

func (transaction *Commit) ProcessDeliver(app interface{}) err.Code {
	log.Debug("Processing Commit Transaction for DeliverTx")
	return err.SUCCESS
}

// Given a transaction, expand it into a list of Commands to execute against various chains.
func (transaction *Commit) Expand(app interface{}) Commands {
	// TODO: Table-driven mechanics, probably elsewhere
	return []Command{}
}
