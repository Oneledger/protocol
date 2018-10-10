/*
	Copyright 2017-2018 OneLedger

	An incoming transaction, send, swap, ready, verification, etc.
*/
package action

import (
	"github.com/Oneledger/protocol/node/err"
	"github.com/Oneledger/protocol/node/global"
	"github.com/Oneledger/protocol/node/id"
	"github.com/Oneledger/protocol/node/log"
)

type Message = []byte // Contents of a transaction
// ENUM for type
type Type byte
type Role byte

const (
	INVALID       Type = iota
	REGISTER           // Register a new identity with the chain
	SEND               // Do a normal send transaction on local chain
	EXTERNAL_SEND      // Do send on external chain
	EXTERNAL_LOCK      // Lock some data on external chain
	SWAP               // Start a swap between chains
	VERIFY             // Verify if a transaction finished
	PUBLISH            // Exchange data on a chain
	READ               // Read a specific transaction on a chain
	PREPARE            // Do everything, except commit
	COMMIT             // Commit to doing the work
	FORGET             // Rollback and forget that this happened
)

const (
	ALL         Role = iota
	INITIATOR        // Register a new identity with the chain
	PARTICIPANT      // Do a normal send transaction on local chain
	NONE
)

type PublicKey = id.PublicKey

// Polymorphism and Serializable
type Transaction interface {
	Validate() err.Code
	ProcessCheck(interface{}) err.Code
	ShouldProcess(interface{}) bool
	ProcessDeliver(interface{}) err.Code
	Expand(interface{}) Commands
	Resolve(interface{}, Commands)
}

// Base Data for each type
type Base struct {
	Type    Type   `json:"type"`
	ChainId string `json:"chain_id"`

	Owner   id.AccountKey `json:"owner"`
	Signers []PublicKey   `json:"signers"`

	Sequence int64 `json:"sequence"`
	Delay    int64 `json:"delay"` // Pause the transaction in the mempool
}

// Execute the function
func Execute(app interface{}, command Command, lastResult map[Parameter]FunctionValue) (err.Code, map[Parameter]FunctionValue) {
	//make sure the first execute use the context, and later uses last result. so if command are executed in a row, every executed function should only add
	//parameters in the context and return instead of create new context every time
	if len(lastResult) > 0 {
		for key, value := range lastResult {
			command.Data[key] = value
		}
	}
	status, result := command.Execute(app)
	if  status {
		return err.SUCCESS, result
	}

	return err.NOT_IMPLEMENTED, lastResult
}

func GetNodeAccount(app interface{}) id.Account {

	accounts := GetAccounts(app)
	account, _ := accounts.FindName(global.Current.NodeAccountName)
	if account == nil {
		log.Error("Node does not have account", "name", global.Current.NodeAccountName)
		accounts.Dump()
		return nil
	}

	return account
}