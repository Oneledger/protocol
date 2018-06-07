/*
	Copyright 2017 - 2018 OneLedger

	Table-driven list of all of the possible functions associated with their transactions
*/
package action

import (
	"github.com/Oneledger/protocol/node/data"
)

type Object interface{}

// Given an action and a chain, return a list of commands
func GetCommands(action Type, chain data.ChainType) Commands {
	for i := 0; i < len(FunctionMapping); i++ {
		transactionType := FunctionMapping[i][0].(Type)
		target := FunctionMapping[i][1].(data.ChainType)
		if action == transactionType && target == chain {
			size := len(FunctionMapping[i]) - 2
			result := make(Commands, size, size)
			for j := 0; j < size; j++ {
				result[j] = FunctionMapping[i][j+2].(Command)
			}
			return result
		}
	}
	return []Command{} // Empty Commands
}

var FunctionMapping = [][]Object{
	[]Object{
		SWAP,
		data.BITCOIN,
		Command{
			Function: CREATE_LOCKBOX,
		},
		Command{
			Function: WAIT_FOR_CHAIN,
		},
	},
	[]Object{
		SWAP,
		data.ETHEREUM,
		Command{
			Function: CREATE_LOCKBOX,
		},
		Command{
			Function: WAIT_FOR_CHAIN,
		},
	},
	[]Object{
		SEND,
		data.BITCOIN,
		Command{
			Function: SUBMIT_TRANSACTION,
		},
	},
}
