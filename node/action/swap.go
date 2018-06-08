/*
	Copyright 2017-2018 OneLedger

	An incoming transaction, send, swap, ready, verification, etc.
*/
package action

import (
	"github.com/Oneledger/protocol/node/data"
	"github.com/Oneledger/protocol/node/err"
	"github.com/Oneledger/protocol/node/id"
	"github.com/Oneledger/protocol/node/log"
)

// Synchronize a swap between two users
type Swap struct {
	Base

	Party        id.AccountKey `json:"party"`
	CounterParty id.AccountKey `json:"counter_party"`
	Amount       data.Coin     `json:"amount"`
	Exchange     data.Coin     `json:"exchange"`
	Fee          data.Coin     `json:"fee"`
	Gas          data.Coin     `json:"fee"`
	Nonce        int64         `json:"nonce"`
}

// Issue swaps across other chains, make sure fees are collected
func (transaction *Swap) Validate() err.Code {
	log.Debug("Validating Swap Transaction")
	return err.SUCCESS
}

func (transaction *Swap) ProcessCheck(app interface{}) err.Code {
	log.Debug("Processing Swap Transaction for CheckTx")

	// TODO: Check all of the data to make sure it is valid.
	return err.SUCCESS
}

func (transaction *Swap) ProcessDeliver(app interface{}) err.Code {
	log.Debug("Processing Swap Transaction for DeliverTx")

	commands := transaction.Expand(app)

	Resolve(app, commands)

	for i := 0; i < commands.Count(); i++ {
		status := Execute(app, commands[i])
		if status != err.SUCCESS {
			log.Error("Failed to Execute", "command", commands[i])
			return err.EXPAND_ERROR
		}
	}

	return err.SUCCESS
}

// Given a transaction, expand it into a list of Commands to execute against various chains.
func (transaction *Swap) Expand(app interface{}) Commands {
	// TODO: Table-driven mechanics, probably elsewhere
	chain := GetChain(transaction)
	return GetCommands(SWAP, chain)
}

// Plug in data from the rest of a system into a set of commands
func Resolve(app interface{}, commands Commands) {
	// TODO: Pick the chain
	// TODO: Fill in all of the necessary data
}

// Execute the function
func Execute(app interface{}, command Command) err.Code {
	return err.SUCCESS
}
