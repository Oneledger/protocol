/*
	Copyright 2017 - 2018 OneLedger

	Table-driven list of all of the possible functions associated with their transactions

	Need to fill in the target chain later, since for any set of instructions it changes...
*/
package action

import (
	"github.com/Oneledger/protocol/node/data"
	"github.com/Oneledger/protocol/node/log"
)

type Object interface{}

// Table-Driven Mapping between transactions and the specific actions to be performed on a set of chains
var FunctionMapping = [][]Object{
	[]Object{
		SWAP,
		Command{
			Function: CREATE_LOCKBOX,
		},
		Command{
			Function: SIGN_LOCKBOX,
		},
		Command{
			Function: WAIT_FOR_CHAIN,
		},
	},
	[]Object{
		SEND,
		Command{
			Function: SUBMIT_TRANSACTION,
		},
	},
}

// Given an action and a chain, return a list of commands
func GetCommands(action Type, chains []data.ChainType) Commands {

	for i := 0; i < len(FunctionMapping); i++ {
		transactionType := FunctionMapping[i][0].(Type)

		// The asymmetric start of the list of commands
		offset := 1

		if action == transactionType {
			size := len(FunctionMapping[i]) - offset
			result := make(Commands, size, size)
			for j := 0; j < size; j++ {
				result[j] = FunctionMapping[i][j+offset].(Command)
			}
			return result
		}
	}

	log.Debug("No Commands", "action", action, "chains", chains)

	return []Command{} // Empty Commands
}
