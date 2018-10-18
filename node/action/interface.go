/*
	Copyright 2017 - 2018 OneLedger

	Interface to specific chain functions
*/

package action

import (
	"github.com/Oneledger/protocol/node/data"
	"github.com/Oneledger/protocol/node/log"
)

func Noop(app interface{}, chain data.ChainType, context FunctionValues) (bool, FunctionValues) {
	log.Info("Executing Noop Command", "chain", chain, "context", context)
	return true, nil
}
