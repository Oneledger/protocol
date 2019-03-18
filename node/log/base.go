/*
	Copyright 2017-2018 OneLedger

	Setup a global logger
*/
package log

import (
	"errors"
	"fmt"
	"os"

	"github.com/Oneledger/protocol/node/global"
	"github.com/davecgh/go-spew/spew"
	"github.com/tendermint/tendermint/libs/log"
)

var current log.Logger

// init a logger right away
func init() {
	current = NewLogger()

	spew.Config = spew.ConfigState{
		Indent:                  "    ",
		DisableMethods:          false,
		DisablePointerMethods:   false,
		DisablePointerAddresses: false,
		DisableCapacities:       false,
		SortKeys:                false,
	}
}

// Dump put the arguments
func Dump(msg string, args ...interface{}) {
	if global.Current.Debug {
		Raw(msg + ":\n")
		spew.Dump(args...)
		Raw("\n")
	}
}

// NewLogger sets in the defaults
func NewLogger() log.Logger {
	// TODO: Replace with non-tendermint logger
	return log.NewTMLogger(log.NewSyncWriter(os.Stdout))
}

// TODO: should be push/pop?
func SetLogger(logger log.Logger) {
	current = logger
}

// GetLogger lets the gobal logger get passed to libraries
func GetLogger() log.Logger {
	return current
}

func Raw(text string) {
	fmt.Print(text)
}

func Info(msg string, args ...interface{}) {
	current.Info(msg, args...)
}

func Debug(msg string, args ...interface{}) {
	if global.Current.Debug {
		current.Debug(msg, args...)
	}
}

func Warn(msg string, args ...interface{}) {
	current.Error("WARNING: "+msg, args...)
}

func Error(msg string, args ...interface{}) {
	current.Error(msg, args...)
}

func Fatal(msg string, args ...interface{}) {
	current.Error("FATAL: "+msg, args...)
	panic(errors.New("Execution stopped due to " + msg))
}

func Settings() {
	Dump("Configuration settings", global.Current.Config)
}
