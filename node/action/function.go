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
			Function: WAIT_FOR_CHAIN, //will create a delay transaction
			Order:    1,
		},
		Command{
			Function: INITIATE,
			Order:    1,
		},
		Command{
			Function: SUBMIT_TRANSACTION,
			Order:    1,
		},
	},
	[]Object{
		SWAP,
		PARTICIPANT,
		Command{
			Function: WAIT_FOR_CHAIN, //will create a delay transaction
			Order:    2,
		},
	},
	[]Object{
		PUBLISH,
		INITIATOR,
		Command{
			Function: AUDITCONTRACT,
			Order:    2,
		},
		Command{
			Function: REDEEM,
			Order:    2,
		},
		Command{
			Function: SUBMIT_TRANSACTION,
			Order:    2,
		},
		Command{
			Function: FINISH,
		},
	},
	[]Object{
		PUBLISH,
		PARTICIPANT,
		Command{
			Function: AUDITCONTRACT,
			Order:    1,
		},
		Command{
			Function: PARTICIPATE,
			Order:    2,
		},
		Command{
			Function: SUBMIT_TRANSACTION,
			Order:    2,
		},
	},
	[]Object{
		PUBLISH,
		ALL,
		Command{
			Function: EXTRACTSECRET,
			Order:    2,
		},
		Command{
			Function: REDEEM,
			Order:    1,
		},
		Command{
			Function: FINISH,
		},
	},
	[]Object{
		SEND,
		ALL,
		Command{
			Function: PREPARE_TRANSACTION,
		},
	},
	[]Object{
		VERIFY,
		ALL,
		Command{
			Function: REFUND,
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
				copy.Order = orig.Order

				result[j] = copy
			}
			return result
		}
	}

	log.Debug("No Commands", "action", action, "chains", chains)

	return []Command{} // Empty Commands
}
