/*
	Copyright 2017-2018 OneLedger

	Cli to start a node (server) running.
*/
package main

import (
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"

	"github.com/Oneledger/protocol/node/app" // Import namespace
	"github.com/Oneledger/protocol/node/global"
	"github.com/Oneledger/protocol/node/log"
	"github.com/Oneledger/protocol/node/persist"

	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/abci/server"
)

var nodeCmd = &cobra.Command{
	Use:   "node",
	Short: "Start up node (server)",
	Run:   StartNode,
}

// Setup the command and flags in Cobra
func init() {
	RootCmd.AddCommand(nodeCmd)
}

// Start a node to run continously
func StartNode(cmd *cobra.Command, args []string) {

	// Catch any underlying panics, for now just print out the details properly and stop
	defer func() {
		if r := recover(); r != nil {
			log.Error("Fullnode Fatal Panic, shutting down", "r", r)
			debug.PrintStack()
			if service != nil {
				service.Stop()
			}
			os.Exit(-1)
		}
	}()

	log.Debug("Starting", "appAddress", global.Current.AppAddress, "on", global.Current.NodeName)

	node := app.NewApplication()
	node.Initialize()

	global.Current.SetApplication(persist.Access(node))
	app.SetNodeName(node)
	LogSettings()

	CatchSigterm()

	// TODO: Switch on config
	//service = server.NewGRPCServer("unix://data.sock", types.NewGRPCApplication(*node))
	//service = server.NewSocketServer("tcp://127.0.0.1:46658", *node)
	service = server.NewSocketServer(global.Current.AppAddress, *node)
	service.SetLogger(log.GetLogger())

	// Set it running
	err := service.Start()
	if err != nil {
		os.Exit(-1)
	}

	/*
		// Wait until it is started
		if err := service.OnStart(); err != nil {
			log.Fatal("Startup Failed", "err", err)
			os.Exit(-1)
		}
	*/

	/*
		for {
			if service.IsRunning() {
				break
			}
			log.Debug("Retrying to see if node is up...")
			time.Sleep(1 * time.Second)

		}

		if !service.IsRunning() {
			log.Fatal("Startup is not running")
			os.Exit(-1)
		}
	*/

	log.Debug("Waiting forever...")
	select {}
}

// A polite way of bring down the service on a SIGTERM
func CatchSigterm() {
	// Catch a SIGTERM and stop
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
	go func() {
		for sig := range sigs {
			log.Info("Shutting down from Signal", "signal", sig)
			if service != nil {
				service.Stop()
			}
			os.Exit(-1)
		}
	}()

}

// Log all of the global settings
func LogSettings() {
	log.Info("Diagnostics", "Debug", global.Current.Debug, "DisablePasswords", global.Current.DisablePasswords)
	log.Info("Ownership", "NodeName", global.Current.NodeName, "NodeAccountName", global.Current.NodeAccountName,
		"NodeIdentity", global.Current.NodeIdentity)
	log.Info("Locations", "RootDir", global.Current.RootDir)
	log.Info("Addresses", "RpcAddress", global.Current.RpcAddress, "AppAddress", global.Current.AppAddress)
}
