/*
	Copyright 2017-2018 OneLedger

	Setup a global logger
*/
package app

import (
	"os"

	"github.com/tendermint/tmlibs/log"
)

var Log log.Logger

// init a logger right away
func init() {
	Log = NewLogger()
}

// NewLogger sets in the defaults
func NewLogger() log.Logger {
	return log.NewTMLogger(log.NewSyncWriter(os.Stdout))
}

// TODO: should be push/pop?
func SetLogger(logger log.Logger) {
	Log = logger
}

// GetLogger lets the gobal logger get passed to libraries
func GetLogger() log.Logger {
	return Log
}
