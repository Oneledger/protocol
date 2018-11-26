/*
	Copyright 2017-2018 OneLedger

	A fullnode for the OneLedger chain. Includes cli arguments to initialize, restart, etc.
*/
package main

import (
	"os"

	"github.com/Oneledger/protocol/node/global"
	"github.com/Oneledger/protocol/node/log"
	"github.com/tendermint/tendermint/libs/common"
)

// Common to all of the sub-commands
var service common.Service

var context *global.Context // Global runtime context

func main() {
	log.Debug("olconfig", "args", os.Args)

	Execute() // Pass control to Cobra
}

// init starts up right away, so the logging and context is initialized as early as possible
func init() {
	context = global.NewContext("olconfig")
}
