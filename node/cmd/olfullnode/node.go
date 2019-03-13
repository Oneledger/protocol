/*
	Copyright 2017-2018 OneLedger

	Cli to start a node (server) running.
*/
package main

import (
	"github.com/Oneledger/protocol/node/app" // Import namespace
	"github.com/Oneledger/protocol/node/cmd/shared"
	"github.com/Oneledger/protocol/node/config"
	"github.com/Oneledger/protocol/node/consensus"

	"github.com/Oneledger/protocol/node/global"
	"github.com/Oneledger/protocol/node/log"
	"github.com/Oneledger/protocol/node/persist"
	"github.com/spf13/cobra"
)

var nodeCmd = &cobra.Command{
	Use:   "node",
	Short: "Start up node (server)",
	RunE:  StartNode,
}

// Setup the command and flags in Cobra
func init() {
	RootCmd.AddCommand(nodeCmd)
	// Get information to connect to a my tendermint node
	nodeCmd.Flags().StringVarP(&global.Current.RpcAddress, "address", "a",
		global.Current.RpcAddress, "consensus address")

	// DELETEME:
	nodeCmd.Flags().StringVarP(&global.Current.Transport, "transport", "t",
		global.Current.Transport, "transport (socket | grpc)")

	nodeCmd.Flags().BoolVarP(&global.Current.Debug, "debug", "d",
		global.Current.Debug, "Set DEBUG mode")

	nodeCmd.Flags().StringVar(&global.Current.BTCAddress, "btcrpc",
		global.Current.BTCAddress, "bitcoin rpc address")

	nodeCmd.Flags().StringVar(&global.Current.ETHAddress, "ethrpc",
		global.Current.ETHAddress, "ethereum rpc address")

	// DELETEME: Should be consistent, always derived from specified RootDir
	nodeCmd.Flags().StringVar(&global.Current.TendermintRoot, "tendermintRoot",
		global.Current.TendermintRoot, "tendermint root directory")

	nodeCmd.Flags().StringVar(&global.Current.SDKAddress, "sdkrpc",
		global.Current.SDKAddress, "Address for SDK RPC Server")

	// TODO: Put this in configuration file
	nodeCmd.Flags().StringArrayVar(&global.Current.PersistentPeers, "persistent_peers", []string{}, "List of persistent peers to connect to")

	// These could be moved to node persistent flags
	nodeCmd.Flags().StringVar(&global.Current.P2PAddress, "p2p", "", "Address to use in P2P network")

	// TODO: Add external listening address address

	nodeCmd.Flags().StringVar(&global.Current.Seeds, "seeds", "", "List of seeds to connect to")

	nodeCmd.Flags().BoolVar(&global.Current.SeedMode, "seed_mode", false, "List of seeds to connect to")
}

// Start a node to run continously
func StartNode(cmd *cobra.Command, args []string) error {

	log.Debug("Starting", "p2pAddress", global.Current.P2PAddress, "on", global.Current.NodeName)

	node := app.NewApplication()
	//if node.CheckIfInitialized() == false {
	//	log.Fatal("Node was not properly initialized")
	//}

	node.Initialize()

	global.Current.SetApplication(persist.Access(node))
	app.SetNodeName(node)
	config.LogSettings()

	shared.CatchSigterm(func() {
		if service != nil {
			service.Stop()
		}
	})

	tmDir := global.ConsensusDir()
	tmConfig := consensus.Config{
		Moniker:            global.Current.NodeName,
		RootDirectory:      tmDir,
		RPCAddress:         global.Current.RpcAddress,
		P2PAddress:         global.Current.P2PAddress,
		ExternalP2PAddress: global.Current.ExternalP2PAddress,
		IndexTags:          []string{"tx.owner", "tx.type", "tx.swapkey", "tx.participants"},
		PersistentPeers:    global.Current.PersistentPeers,
		Seeds:              global.Current.Seeds,
		SeedMode:           global.Current.SeedMode,
	}

	// TODO: change the the priv_validator locaiton
	service, err := consensus.NewNode(*node, tmConfig)
	if err != nil {
		log.Error("Failed to create NewNode", "err", err)
		return err
	}

	// Set it running
	err = service.Start()
	if err != nil {
		log.Error("Can't start up node", "err", err)
		return err
	}

	global.Current.SetConsensusNode(service)
	log.Debug("Waiting forever...")
	select {}
}
