/*
	Copyright 2017-2018 OneLedger

	A fullnode for the OneLedger chain. Includes cli arguments to initialize, restart, etc.
*/
package main

import (
	"github.com/Oneledger/protocol/node/comm"
	"os"

	"github.com/Oneledger/protocol/node/global"
	"github.com/Oneledger/protocol/node/log"
	"github.com/tendermint/tendermint/libs/common"
)

// Common to all of the sub-commands
var service common.Service

var context *global.Context // Global runtime context

var rpcclient comm.ClientInterface

func main() {
	log.Debug("olclient", "args", os.Args)

	rpcclient = comm.GetClient()

	Execute() // Pass control to Cobra
}

// init starts up right away, so the logging and context is initialized as early as possible
func init() {
	context = global.NewContext("olclient")
}
