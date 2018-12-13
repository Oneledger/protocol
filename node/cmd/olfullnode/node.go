/*
	Copyright 2017-2018 OneLedger

	Cli to start a node (server) running.
*/
package main

import (
	"os"
	"path/filepath"
	"runtime/debug"

	"github.com/Oneledger/protocol/node/app" // Import namespace
	"github.com/Oneledger/protocol/node/cmd/shared"
	"github.com/Oneledger/protocol/node/config"
	"github.com/Oneledger/protocol/node/consensus"

	"github.com/Oneledger/protocol/node/global"
	"github.com/Oneledger/protocol/node/log"
	"github.com/Oneledger/protocol/node/persist"
	"github.com/tendermint/tendermint/privval"
	"github.com/tendermint/tendermint/types"

	"github.com/spf13/cobra"
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

	if node.CheckIfInitialized() == false {
		log.Fatal("Node was not properly initialized")
	}

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

	tmDir := global.ConsensusDir()
	privValidator := privval.LoadFilePV(filepath.Join(tmDir, "config", "priv_validator.json"))
	genesisDoc, err := types.GenesisDocFromFile(filepath.Join(tmDir, "config", "genesis.json"))
	if err != nil {
		log.Fatal("Couldn't read genesis file", "location", filepath.Join(tmDir, "genesis.json"))
	}

	// TODO: Source this from static file
	tmConfig := consensus.Config{
		Moniker:         global.Current.NodeName,
		RootDirectory:   tmDir,
		RPCAddress:      global.Current.RpcAddress,
		P2PAddress:      global.Current.P2PAddress,
		IndexTags:       []string{"tx.owner", "tx.type"},
		PersistentPeers: global.Current.PersistentPeers,
		Seeds:           global.Current.Seeds,
		SeedMode:        global.Current.SeedMode,
	}

	// TODO: change the the priv_validator locaiton
	service, err := consensus.NewNode(*node, tmConfig, privValidator, genesisDoc)
	if err != nil {
		log.Error("Failed to create NewNode", "err", err)
		os.Exit(1)
	}

	// Set it running
	err = service.Start()
	if err != nil {
		log.Error("Can't start up node", "err", err)
		os.Exit(-1)
	}

	global.Current.SetConsensusNode(service)
	log.Debug("Waiting forever...")
	select {}
}
