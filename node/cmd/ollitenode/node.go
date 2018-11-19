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

	// Import namespace
	"github.com/Oneledger/protocol/node/cmd/shared"
	"github.com/Oneledger/protocol/node/global"
	"github.com/Oneledger/protocol/node/log"

	"github.com/spf13/cobra"
)

var nodeCmd = &cobra.Command{
	Use:   "node",
	Short: "Start up a light node",
	Run:   StartNode,
}

// Declare a shared arguments struct
var arguments = &shared.RegisterArguments{}

// Setup the command and flags in Cobra
func init() {
	RootCmd.AddCommand(nodeCmd)
}

// Start a node to run continously
func StartNode(cmd *cobra.Command, args []string) {

	// Catch any underlying panics, for now just print out the details properly and stop
	defer func() {
		if r := recover(); r != nil {
			log.Error("Litenode Fatal Panic, shutting down", "r", r)
			debug.PrintStack()
			if service != nil {
				service.Stop()
			}
			os.Exit(-1)
		}
	}()

	/*
		log.Debug("Starting", "appAddress", global.Current.AppAddress, "on", global.Current.NodeName)

		node := app.NewApplication()
		node.Initialize()

		global.Current.SetApplication(persist.Access(node))
		app.SetNodeName(node)
		LogSettings()
	*/

	CatchSigterm()

	// TODO: Kinda like this, but this example was taken from a later version of Cosmos. Only
	// the proxy.StartProxy call seems to exist in 23???
	/*
		logger := log.GetLogger()
		cacheSize := 120000

		node := rpcclient.NewHttp(nodeAddr, "/websocket")
		verifier, err := proxy.NewVerifier(chainId, home, node, logger, cacheSize)
		if err != nil {
			log.Fatal("Failed to start", "err", err)
		}

		verifier.SetLogger(logger)

		client := proxy.SecureClient(node, verifier)
		err = proxy.StartProxy(client, listenAddr, logger, maxConnections)
		if err != nil {
			log.Fatal("Proxy cannot be created", "err", err)
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
