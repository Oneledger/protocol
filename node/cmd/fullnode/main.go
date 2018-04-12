/*
	Copyright 2017-2018 OneLedger
*/

package main

import (
	"fmt"
	"github.com/Oneledger/prototype/node/app"
	"github.com/spf13/cobra"
	"github.com/tendermint/abci/server"
	"github.com/tendermint/abci/types"
	"github.com/tendermint/tmlibs/common"
	"github.com/tendermint/tmlibs/log"
	"os"
)

// Setup a basic command for testing
var start = &cobra.Command{
	Run:   Start,
	Use:   "start",
	Short: "Start",
	Long:  "Start",
}

var service common.Service

func main() {
	//logger := log.TestingLogger()
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))

	logger.Info("Starting")
	//fmt.Println("Staring OneLedger Node")

	// TODO: Initialize Cobra (cli) and Viper (config)

	node := app.NewApplicationContext()
	_ = node

	// TODO: Just use a generic base, not hooked up yet.
	//app := types.NewBaseApplication()
	app := *node

	service = server.NewGRPCServer("unix://data.sock", types.NewGRPCApplication(app))
	service.SetLogger(logger)
}

func Start(cmd *cobra.Command, args []string) {
	fmt.Println("Start")
}
