/*
	Copyright 2017-2018 OneLedger

	An incoming transaction, send, publish, ready, verification, etc.
*/
package action

import (
	"bytes"
	"github.com/Oneledger/protocol/node/id"
	"github.com/Oneledger/protocol/node/log"
	"github.com/Oneledger/protocol/node/serial"
	"github.com/Oneledger/protocol/node/status"
)

// Synchronize a publish between two users
type Publish struct {
	Base

	Target     id.AccountKey `json:"target"`
	Contract   Message       `json:"message"` //message converted from HTLContract
	SecretHash [32]byte      `json:"secrethash"`
	Count      int           `json:"count"`
}

func init() {
	serial.Register(Publish{})
}

// Ensure that all of the base values are at least reasonable.
func (publish *Publish) Validate() status.Code {
	log.Debug("Validating Publish Transaction")

	if publish.Target == nil {
		log.Debug("Missing Target")
		return status.MISSING_DATA
	}

	if publish.Contract == nil {
		log.Debug("Missing Contract")
		return status.MISSING_DATA
	}

	log.Debug("Publish is validated!")
	return status.SUCCESS
}

func (publish *Publish) ProcessCheck(app interface{}) status.Code {
	log.Debug("Processing Publish Transaction for CheckTx")

	// TODO: Check all of the data to make sure it is valid.

	return status.SUCCESS
}

// Start the publish
func (publish *Publish) ProcessDeliver(app interface{}) status.Code {
	log.Debug("Processing Publish Transaction for DeliverTx")

	commands := publish.Resolve(app)
	commands.Execute(app)
	return status.SUCCESS
}

// Is this node one of the partipants in the publish
func (publish *Publish) ShouldProcess(app interface{}) bool {
	account := GetNodeAccount(app)

	log.Debug("Not the publish target", "target", publish.Target, "me", account.AccountKey())

	if bytes.Equal(publish.Target, account.AccountKey()) {
		return true
	}

	return false
}

// Plug in data from the rest of a system into a set of commands
func (publish *Publish) Resolve(app interface{}) Commands {
	return Commands{}
}
