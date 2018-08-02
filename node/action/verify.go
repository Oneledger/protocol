/*
	Copyright 2017-2018 OneLedger

	An incoming transaction, send, swap, ready, verification, etc.
*/
package action

import (
	"github.com/Oneledger/protocol/node/err"
	"github.com/Oneledger/protocol/node/id"
	"github.com/Oneledger/protocol/node/log"
    "bytes"
    )


// instead of a model for other-chain actions, we make this available by event, when a action (external/internal) is
// done, a event related to other chain is stored with status(false/true, representing finished or not), this verify
// just check the event status.
type Verify struct {
	Base

	Target  id.AccountKey   `json:"target"`
	Event   Event           `json:"event"`
	Message Message         `json:"Message"`
}

func (transaction Verify) Validate() err.Code {
	log.Debug("Validating Verify Transaction")
    if transaction.Target == nil {
        log.Debug("Missing Target")
        return err.MISSING_DATA
    }

    if &transaction.Event == nil {
        log.Debug("Missing Event")
        return err.MISSING_DATA
    }

    log.Debug("Publish is validated!")
	return err.SUCCESS
}

func (transaction Verify) ProcessCheck(app interface{}) err.Code {
	log.Debug("Processing Verify Transaction for CheckTx")
	//todo : check the data ?
	return err.SUCCESS
}

func (transaction Verify) ShouldProcess(app interface{}) bool {
    account := GetNodeAccount(app)
    log.Debug("Not the publish target", "target", transaction.Base.Owner, "me", account.AccountKey() )

    if bytes.Equal(transaction.Target, account.AccountKey()) {
        return true
    }

    return false
}

func (transaction Verify) ProcessDeliver(app interface{}) err.Code {
	log.Debug("Processing Verify Transaction for DeliverTx")

    commands := transaction.Expand(app)

    transaction.Resolve(app, commands)

    //before loop of execute, lastResult is nil
    var lastResult map[Parameter]FunctionValue

    for i := 0; i < commands.Count(); i++ {
        status, result := Execute(app, commands[i], lastResult)
        if status != err.SUCCESS {
            log.Error("Failed to Execute", "command", commands[i])
            return err.EXPAND_ERROR
        }
        lastResult = result
    }
    return err.SUCCESS
}

func (transaction Verify) Resolve(app interface{}, commands Commands) {
    eventStatus := FindEvent(app, transaction.Event)
    if !eventStatus {
        status := GetStatus(app)
        swap := FindSwap(status, transaction.Event.Key)
        swap.Resolve(app, commands)
    } else {
        commands = nil
    }
    return
}

// Given a transaction, expand it into a list of Commands to execute against various chains.
func (transaction Verify) Expand(app interface{}) Commands {
	// TODO: Table-driven mechanics, probably elsewhere
    return GetCommands(VERIFY, ALL, nil)
}
