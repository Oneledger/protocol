/*
	Copyright 2017-2018 OneLedger

	Cli to start a node (server) running.
*/
package main

import (
	"os"
	"runtime/debug"
	"time"

	"github.com/Oneledger/protocol/node/app" // Import namespace
	"github.com/Oneledger/protocol/node/cmd/shared"
	"github.com/Oneledger/protocol/node/config"
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

	/*
		//time.Sleep(5 * time.Second)
		if arguments.Identity != "" {
			log.Debug("Have Register Request", "arguments", arguments)

			// TODO: Maybe Tendermint isn't ready for transactions...
			//time.Sleep(10 * time.Second)

			packet := shared.CreateRegisterRequest(arguments)
			result := comm.Broadcast(packet)

			log.Debug("######## Register Broadcast", "result", result)
		} else {
			log.Debug("Nothing to Register")
		}
	*/
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
	config.LogSettings()

	shared.CatchSigterm(func() {
		if service != nil {
			service.Stop()
		}
	})

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

	// TODO: Sleep until the node is connected and running
	time.Sleep(10 * time.Second)
	log.Debug("################### STARTED UP ######################")

	// If the register flag is set, do that before waiting
	Register()

	log.Debug("Waiting forever...")
	select {}
}
