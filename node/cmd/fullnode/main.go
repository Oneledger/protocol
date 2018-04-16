/*
	Copyright 2017-2018 OneLedger

	A fullnode for the OneLedger chain. Includes cli arguments to initialize, restart, etc.
*/
package main

import (
	"github.com/Oneledger/prototype/node/app"
	"github.com/tendermint/tmlibs/common"
)

var service common.Service

var context *app.Context // Global runtime context

func main() {
	Execute() // Pass control to Cobra
}

// init starts up right away, so the logging and context is initialized as early as possible
func init() {
	context = app.NewContext("Fullnode")
}
