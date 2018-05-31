/*
	Copyright 2017-2018 OneLedger

	Cli to start a node (server) running.
*/
package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/Oneledger/protocol/node/app" // Import namespace
	"github.com/Oneledger/protocol/node/global"
	"github.com/Oneledger/protocol/node/log"
	"github.com/Oneledger/protocol/node/persist"

	"github.com/spf13/cobra"
	"github.com/tendermint/abci/server"
)

var nodeCmd = &cobra.Command{
	Use:   "node",
	Short: "Start up node (server)",
	Run:   StartNode,
}

func init() {
	RootCmd.AddCommand(nodeCmd)
}

func HandleArguments() {
}

func StartNode(cmd *cobra.Command, args []string) {
	log.Info("Starting up a Node")

	// Catch any underlying panics, for now just print out the details properly and stop
	defer func() {
		if r := recover(); r != nil {
			log.Error("Fatal Panic, shutting down", "r", r)
			os.Exit(-1)
		}
	}()

	node := app.NewApplication()
	global.Current.SetApplication(persist.Access(node))

	// TODO: Switch on config
	//service = server.NewGRPCServer("unix://data.sock", types.NewGRPCApplication(*node))
	//service = server.NewSocketServer("tcp://127.0.0.1:46658", *node)

	log.Debug("Starting", "address", global.Current.Address)

	CatchSigterm()

	service = server.NewSocketServer(global.Current.Address, *node)
	service.SetLogger(log.GetLogger())

	// Set it running
	err := service.Start()
	if err != nil {
		os.Exit(-1)
	}

	select {} // Wait forever
}

func CatchSigterm() {
	// Catch a SIGTERM and stop
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
	go func() {
		for sig := range sigs {
			log.Info("Shutting down from Signal", "signal", sig)
			if service != nil {
				service.Stop()
				os.Exit(-1)
			}
		}
	}()

}
