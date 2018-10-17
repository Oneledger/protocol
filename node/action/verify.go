/*
	Copyright 2017-2018 OneLedger

	An incoming transaction, send, swap, ready, verification, etc.
*/
package action

import (
	"bytes"
	"github.com/Oneledger/protocol/node/id"
	"github.com/Oneledger/protocol/node/log"
	"github.com/Oneledger/protocol/node/serial"
	"github.com/Oneledger/protocol/node/status"
)

// instead of a model for other-chain actions, we make this available by event, when a action (external/internal) is
// done, a event related to other chain is stored with status(false/true, representing finished or not), this verify
// just check the event status.
type Verify struct {
	Base

	Target  id.AccountKey `json:"target"`
	Event   Event         `json:"event"`
	Message Message       `json:"Message"`
}

func init() {
	serial.Register(Verify{})
}

func (transaction Verify) Validate() status.Code {
	log.Debug("Validating Verify Transaction")
	if transaction.Target == nil {
		log.Debug("Missing Target")
		return status.MISSING_DATA
	}

	if &transaction.Event == nil {
		log.Debug("Missing Event")
		return status.MISSING_DATA
	}

	log.Debug("Verify is validated!")
	return status.SUCCESS
}

func (transaction Verify) ProcessCheck(app interface{}) status.Code {
	log.Debug("Processing Verify Transaction for CheckTx")
	//todo : check the data ?
	return status.SUCCESS
}

func (transaction Verify) ShouldProcess(app interface{}) bool {
	account := GetNodeAccount(app)

	if bytes.Equal(transaction.Target, account.AccountKey()) {
		log.Debug("Is verify target", "target", transaction.Base.Owner, "me", account.AccountKey())

		return true
	}
	log.Debug("Not the verify target", "target", transaction.Base.Owner, "me", account.AccountKey())
	return false
}

func (transaction Verify) ProcessDeliver(app interface{}) status.Code {
	log.Debug("Processing Verify Transaction for DeliverTx")

	commands := transaction.Expand(app)

	transaction.Resolve(app, commands)

	//before loop of execute, lastResult is nil
	var lastResult FunctionValues

	for i := 0; i < commands.Count(); i++ {
		status, result := Execute(app, commands[i], lastResult)
		if status != status.SUCCESS {
			log.Error("Failed to Execute", "command", commands[i])
			return status.EXPAND_ERROR
		}
		lastResult = result
	}
	return status.SUCCESS
}

func (transaction Verify) Resolve(app interface{}) Commands {
	eventStatus := FindEvent(app, transaction.Event)
	if !eventStatus {
		status := GetStatus(app)
		swap := FindSwap(status, transaction.Event.Key)
		if swap != nil {
			entry := swap.(*Swap)
			account := GetNodeAccount(app)
			isParty := entry.IsParty(account)
			role := entry.getRole(*isParty)
			for i := 0; i < len(commands); i++ {
				if role == INITIATOR {
					commands[i].Order = 1
				} else if role == PARTICIPANT {
					commands[i].Order = 2
				}
			}
			entry.Resolve(app, commands)
		} else {
			log.Error("Failed to find stored transaction", "key", transaction.Event.Key)
		}

	} else {
		commands = nil
	}
	return commands
}

// Given a transaction, expand it into a list of Commands to execute against various chains.
func (transaction Verify) Expand(app interface{}) Commands {
	// TODO: Table-driven mechanics, probably elsewhere
	return GetCommands(VERIFY, ALL, nil)
}
