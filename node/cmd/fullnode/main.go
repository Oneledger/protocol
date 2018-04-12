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
	Initialize()
	ParseArgs()
}

func ParseArgs() {
	// TODO: Initialize Cobra (cli) and Viper (config)

	// TODO: Cobra will call this...
	//StartNode(nil, nil)
	Execute()
}

func Initialize() {
	// Setup initial logging
	logger = log.NewTMLogger(log.NewSyncWriter(os.Stdout))
	logger.Info("Starting")
}

/*
func StartNode(cmd *cobra.Command, args []string) {
	fmt.Println("Starting up a Node")
	node := app.NewApplicationContext()
	service = server.NewGRPCServer("unix://data.sock", types.NewGRPCApplication(*node))
	service.SetLogger(logger)
}
*/
