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
	"github.com/Oneledger/protocol/node/chains/bitcoin"
	"github.com/Oneledger/protocol/node/global"
	"github.com/Oneledger/protocol/node/chains/bitcoin/htlc"
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

// Ensure that all of the base values are at least reasonable.
func (transaction *Swap) Validate() err.Code {
	log.Debug("Validating Swap Transaction")

	if transaction.Party == nil {
		return err.MISSING_DATA
	}
	if transaction.CounterParty == nil {
		return err.MISSING_DATA
	}
	if !transaction.Amount.IsValid() {
		return err.MISSING_DATA
	}
	return err.SUCCESS
}

func (transaction *Swap) ProcessCheck(app interface{}) err.Code {
	log.Debug("Processing Swap Transaction for CheckTx")

	// TODO: Check all of the data to make sure it is valid.

	return err.SUCCESS
}

func (transaction *Swap) ThisNode(app interface{}) bool {
	return true
}

// Start the swap
func (transaction *Swap) ProcessDeliver(app interface{}) err.Code {
	log.Debug("Processing Swap Transaction for DeliverTx")

	commands := transaction.Expand(app)

	Resolve(app, transaction, commands)

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
	chains := GetChains(transaction)

	return GetCommands(SWAP, chains)
}

// Plug in data from the rest of a system into a set of commands
func Resolve(app interface{}, transaction Transaction, commands Commands) {
	identities := GetIdentities(app)
	_ = identities

	utxo := GetUtxo(app)
	_ = utxo

	chains := GetChains(transaction)
	for i := 0; i < len(commands); i++ {
		//TODO: add parameter for actions
		commands[i].Chain = chains[0]
	}
}

// Execute the function
func Execute(app interface{}, command Command) err.Code {
	if command.Execute() {
		return err.SUCCESS
	}
	return err.NOT_IMPLEMENTED
}


func CreateContractBTC(context map[string]string) bool {
	cli := bitcoin.GetBtcClient(global.Current.BTCRpcPort)
	//todo: runCommand(initCmd,cli)

	return true
}

func CreateContractETH(context map[string]string) bool {
	return true
}

func CreateContractOLT(context map[string]string) bool {
	return true
}

 