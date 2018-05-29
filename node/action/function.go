/*
	Copyright 2017 - 2018 OneLedger
*/
package action

import "github.com/Oneledger/protocol/node/id"

type Object interface{}

// Given an action and a chain, return a list of commands
func GetCommands(action Type, chain id.AccountType) Commands {
	for i := 0; i < len(FunctionMapping); i++ {
		transactionType := FunctionMapping[i][0].(Type)
		target := FunctionMapping[i][1].(id.AccountType)
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
		id.BITCOIN,
		Command{
			Function: CREATE_LOCKBOX,
		},
		Command{
			Function: WAIT_FOR_CHAIN,
		},
	},
	[]Object{
		SWAP,
		id.ETHEREUM,
		Command{
			Function: CREATE_LOCKBOX,
		},
		Command{
			Function: WAIT_FOR_CHAIN,
		},
	},
	[]Object{
		SEND,
		id.BITCOIN,
		Command{
			Function: SUBMIT_TRANSACTION,
		},
	},
}
