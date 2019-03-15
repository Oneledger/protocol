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

	"github.com/Oneledger/protocol/node/global"
	"github.com/Oneledger/protocol/node/log"
	"github.com/Oneledger/protocol/node/olvm/interpreter/vm"

	"github.com/spf13/cobra"
)

var exeCmd = &cobra.Command{
	Use:   "execute",
	Short: "Startup a service for contract execution",
	Run:   StartEngine,
}

// Setup the command and flags in Cobra
func init() {
	RootCmd.AddCommand(exeCmd)
}

// Start a node to run continously
func StartEngine(cmd *cobra.Command, args []string) {

	// Catch any underlying panics, for now just print out the details properly and stop
	defer func() {
		if r := recover(); r != nil {
			log.Error("OLVM Fatal Panic, shutting down", "r", r)
			debug.PrintStack()
			if service != nil {
				service.Stop()
			}
			os.Exit(-1)
		}
	}()

	log.Debug("Starting", "p2p address", global.Current.Config.Network.P2PAddress, "on", global.Current.NodeName)

	// TODO: Switch with config and shared versions
	log.Settings()
	CatchSigterm()

	service := vm.NewOLVMService()
	service.StartService()

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
