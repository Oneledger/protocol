/*
	Copyright 2017-2018 OneLedger

	A fullnode for the OneLedger chain. Includes cli arguments to initialize, restart, etc.
*/
package main

import (
	"github.com/tendermint/tmlibs/common"
	"github.com/tendermint/tmlibs/log"
	"os"
)

var service common.Service
var logger log.Logger

func main() {
	// Pass control to Cobra
	Execute()
}

// init starts up right away, so the logging is initialized as early as possible
func init() {
	// Setup initial logging
	logger = log.NewTMLogger(log.NewSyncWriter(os.Stdout))
	logger.Debug("Starting")
}
