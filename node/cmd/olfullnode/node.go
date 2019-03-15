/*
	Copyright 2017-2018 OneLedger

	Cli to start a node (server) running.
*/
package main

import (
	"github.com/Oneledger/protocol/node/app" // Import namespace
	"github.com/Oneledger/protocol/node/cmd/shared"
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
	nodeCmd.Flags().StringVarP(&global.Current.Config.Network.RPCAddress, "address", "a",
		global.Current.Config.Network.RPCAddress, "consensus address")

	// DELETEME:
	nodeCmd.Flags().StringVarP(&global.Current.Transport, "transport", "t",
		global.Current.Transport, "transport (socket | grpc)")

	nodeCmd.Flags().BoolVarP(&global.Current.Debug, "debug", "d",
		global.Current.Debug, "Set DEBUG mode")

	nodeCmd.Flags().StringVar(&global.Current.Config.Network.BTCAddress, "btcrpc",
		global.Current.Config.Network.BTCAddress, "bitcoin rpc address")

	nodeCmd.Flags().StringVar(&global.Current.Config.Network.ETHAddress, "ethrpc",
		global.Current.Config.Network.ETHAddress, "ethereum rpc address")

	// DELETEME: Should be consistent, always derived from specified RootDir
	nodeCmd.Flags().StringVar(&global.Current.ConsensusDir(), "tendermintRoot",
		global.Current.ConsensusDir(), "tendermint root directory")

	nodeCmd.Flags().StringVar(&global.Current.Config.Network.SDKAddress, "sdkrpc",
		global.Current.Config.Network.SDKAddress, "Address for SDK RPC Server")

	nodeCmd.Flags().StringArrayVar(&global.Current.PersistentPeers, "persistent_peers", []string{}, "List of persistent peers to connect to")

	// These could be moved to node persistent flags
	nodeCmd.Flags().StringVar(&global.Current.Config.Network.P2PAddress, "p2p", "", "Address to use in P2P network")

	nodeCmd.Flags().StringVar(&global.Current.Seeds, "seeds", "", "List of seeds to connect to")

	nodeCmd.Flags().BoolVar(&global.Current.SeedMode, "seed_mode", false, "List of seeds to connect to")
}

// Start a node to run continously
func StartNode(cmd *cobra.Command, args []string) error {

	log.Debug("Starting", "p2pAddress", global.Current.Config.Network.P2PAddress, "on", global.Current.NodeName)

	node := app.NewApplication()
	//if node.CheckIfInitialized() == false {
	//	log.Fatal("Node was not properly initialized")
	//}

	node.Initialize()

	global.Current.SetApplication(persist.Access(node))
	app.SetNodeName(node)
	log.Settings()

	shared.CatchSigterm(func() {
		if service != nil {
			service.Stop()
		}
	})

	tmDir := global.Current.ConsensusDir()
	tmConfig := consensus.Config{
		Moniker:         global.Current.NodeName,
		RootDirectory:   tmDir,
		RPCAddress:      global.Current.Config.Network.RPCAddress,
		P2PAddress:      global.Current.Config.Network.P2PAddress,
		IndexTags:       []string{"tx.owner", "tx.type", "tx.swapkey", "tx.participants"},
		PersistentPeers: global.Current.PersistentPeers,
		Seeds:           global.Current.Seeds,
		SeedMode:        global.Current.SeedMode,
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
