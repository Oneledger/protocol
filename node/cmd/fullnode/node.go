/*
	Copyright 2017-2018 OneLedger

	Cli to start a node (server) running.
*/
package main

import (
	"github.com/Oneledger/prototype/node/app" // Import namespace
	"github.com/Oneledger/prototype/node/log"

	"github.com/spf13/cobra"
	"github.com/tendermint/abci/server"
	"github.com/tendermint/tmlibs/common"
)

var nodeCmd = &cobra.Command{
	Use:   "node",
	Short: "Start up node (server)",
	Run:   StartNode,
}

// TODO: Move to Context
var transport string
var address string

func init() {
	RootCmd.AddCommand(nodeCmd)

	nodeCmd.Flags().StringVarP(&app.Current.Transport, "transport", "t", "socket", "transport (socket | grpc)")
	nodeCmd.Flags().StringVarP(&app.Current.Address, "address", "a", "tcp://127.0.0.1:46658", "full address")
}

func HandleArguments() {
}

func StartNode(cmd *cobra.Command, args []string) {
	log.Info("Starting up a Node")

	node := app.NewApplication()

	// TODO: Switch on config
	//service = server.NewGRPCServer("unix://data.sock", types.NewGRPCApplication(*node))
	//service = server.NewSocketServer("tcp://127.0.0.1:46658", *node)
	service = server.NewSocketServer(app.Current.Address, *node)
	service.SetLogger(log.GetLogger())

	// TODO: catch any panics

	// Set it running
	err := service.Start()
	if err != nil {
		common.Exit(err.Error())
	}

	// Catch any signals, stop nicely
	common.TrapSignal(func() {
		log.Info("Shutting down")
		service.Stop()
	})
}
