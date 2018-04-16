/*
	Copyright 2017-2018 OneLedger

	Cli to start a node (server) running.
*/
package main

import (
	"github.com/Oneledger/prototype/node/app" // Import namespace

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
var name string
var transport string
var address string

func init() {
	RootCmd.AddCommand(nodeCmd)

	nodeCmd.Flags().StringVarP(&name, "name", "n", "Fullnode", "node name")
	nodeCmd.Flags().StringVarP(&transport, "transport", "t", "socket", "transport (socket | grpc)")
	nodeCmd.Flags().StringVarP(&address, "address", "a", "tcp://127.0.0.1:46658", "full address")
}

func HandleArguments() {
}

func StartNode(cmd *cobra.Command, args []string) {
	app.Log.Info("Starting up a Node")

	node := app.NewApplication()

	// TODO: Switch on config
	//service = server.NewGRPCServer("unix://data.sock", types.NewGRPCApplication(*node))
	service = server.NewSocketServer("tcp://127.0.0.1:46658", *node)
	service.SetLogger(app.GetLogger())

	// TODO: catch any panics

	// Set it running
	err := service.Start()
	if err != nil {
		common.Exit(err.Error())
	}

	// Catch any signals, stop nicely
	common.TrapSignal(func() {
		app.Log.Info("Shutting down")
		service.Stop()
	})
}
