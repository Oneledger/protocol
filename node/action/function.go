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
		INITIATOR,
		Command{
			Function: INITIATE,
			Side:     0,
		},
		Command{
			Function: AUDITCONTRACT,
			Side:     1,
		},
		Command{
			Function: REDEEM,
			Side:     1,
		},
		Command{
			Function: WAIT_FOR_CHAIN,
			Side:     0,
		},
	},
	[]Object{
		SWAP,
		PARTICIPANT,
		Command{
			Function: AUDITCONTRACT,
			Side:     0,
		},
		Command{
			Function: PARTICIPATE,
			Side:     1,
		},
		Command{
			Function: EXTRACTSECRET,
			Side:     1,
		},
		Command{
			Function: REDEEM,
			Side: 	  0,
		},
		Command{
			Function: WAIT_FOR_CHAIN,
			Side:     1,
		},
	},
	[]Object{
		SEND,
		ALL,
		Command{
			Function: SUBMIT_TRANSACTION,
		},
	},
}

// Given an action and a chain, return a list of commands
func GetCommands(action Type, role Role, chains []data.ChainType) Commands {

	for i := 0; i < len(FunctionMapping); i++ {
		transactionType := FunctionMapping[i][0].(Type)
		transactionRole := FunctionMapping[i][1].(Role)

		// The asymmetric start of the list of commands
		offset := 2

		if action == transactionType && role == transactionRole {
			size := len(FunctionMapping[i]) - offset
			result := make(Commands, size, size)
			for j := 0; j < size; j++ {
				var copy Command
				orig := FunctionMapping[i][j+offset].(Command)

				// Make sure we take a copy of this, not the original
				copy.Function = orig.Function
				copy.Chain = orig.Chain
				copy.Data = make(map[Parameter]FunctionValue)

				result[j] = copy
			}
			return result
		}
	}

	log.Debug("No Commands", "action", action, "chains", chains)

	return []Command{} // Empty Commands
}
