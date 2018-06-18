/*
	Copyright 2017-2018 OneLedger

	Cli to start a node (server) running.
*/
package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Oneledger/protocol/node/app" // Import namespace
	"github.com/Oneledger/protocol/node/cmd/shared"
	"github.com/Oneledger/protocol/node/comm"
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

// Declare a shared arguments struct
var arguments = &shared.RegisterArguments{}

// Setup the command and flags in Cobra
func init() {
	RootCmd.AddCommand(nodeCmd)

	nodeCmd.Flags().StringVar(&arguments.Identity, "register", "", "Register this identity")
}

// Use the client side to broadcast an identity to all nodes.
func Register() {

	// Don't let the death of a client stop the node from running
	defer func() {
		if r := recover(); r != nil {
			log.Error("Ignoring Client Panic", "r", r)
			return
		}
	}()

	if arguments.Identity != "" {
		log.Debug("Have Register Request", "arguments", arguments)

		// TODO: Maybe Tendermint isn't ready for transactions...
		time.Sleep(5 * time.Second)

		packet := shared.CreateRegisterRequest(arguments)
		result := comm.Broadcast(packet)

		log.Debug("Registered Successfully", "result", result)
	}
}

// Start a node to run continously
func StartNode(cmd *cobra.Command, args []string) {

	// Catch any underlying panics, for now just print out the details properly and stop
	defer func() {
		if r := recover(); r != nil {
			log.Error("Fatal Panic, shutting down", "r", r)
			if service != nil {
				service.Stop()
			}
			os.Exit(-1)
		}
	}()

	log.Debug("Starting", "appAddress", global.Current.AppAddress)

	node := app.NewApplication()
	node.Initialize()

	if global.Current.NodeAccountName == "" {
		log.Warn("Missing NodeAccount")
	} else {
		log.Info("Node", "account", global.Current.NodeAccountName)
	}

	global.Current.SetApplication(persist.Access(node))

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

	// If the register flag is set, do that before waiting
	Register()

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
