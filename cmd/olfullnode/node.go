// /*
// 	Copyright 2017-2018 OneLedger
//
// 	Cli to start a node (server) running.
// */
package main
//
// import (
// 	"github.com/Oneledger/protocol/node/app" // Import namespace
// 	"github.com/Oneledger/protocol/node/cmd/shared"
// 	"github.com/Oneledger/protocol/node/global"
// 	"github.com/Oneledger/protocol/node/log"
// 	"github.com/spf13/cobra"
// )
//
// var nodeCmd = &cobra.Command{
// 	Use:   "node",
// 	Short: "Start up node (server)",
// 	RunE:  StartNode,
// }
//
// var shouldWriteConfig bool
//
// // Setup the command and flags in Cobra
// func init() {
// 	RootCmd.AddCommand(nodeCmd)
// 	// Get information to connect to a my tendermint node
// 	nodeCmd.Flags().StringVarP(&global.Current.Config.Network.RPCAddress, "address", "a",
// 		global.Current.Config.Network.RPCAddress, "consensus address")
//
// 	nodeCmd.Flags().BoolVarP(&global.Current.Debug, "debug", "d",
// 		global.Current.Debug, "Set DEBUG mode")
//
// 	nodeCmd.Flags().StringVar(&global.Current.Config.Network.BTCAddress, "btcrpc",
// 		global.Current.Config.Network.BTCAddress, "bitcoin rpc address")
//
// 	nodeCmd.Flags().StringVar(&global.Current.Config.Network.ETHAddress, "ethrpc",
// 		global.Current.Config.Network.ETHAddress, "ethereum rpc address")
//
// 	nodeCmd.Flags().StringVar(&global.Current.Config.Network.SDKAddress, "sdkrpc",
// 		global.Current.Config.Network.SDKAddress, "Address for SDK RPC Server")
//
// 	nodeCmd.Flags().StringArrayVar(&global.Current.PersistentPeers, "persistent_peers", []string{}, "List of persistent peers to connect to")
//
// 	// These could be moved to node persistent flags
// 	nodeCmd.Flags().StringVar(&global.Current.Config.Network.P2PAddress, "p2p", "", "Address to use in P2P network")
//
// 	nodeCmd.Flags().StringVar(&global.Current.Seeds, "seeds", "", "List of seeds to connect to")
//
// 	nodeCmd.Flags().BoolVar(&global.Current.SeedMode, "seed_mode", false, "List of seeds to connect to")
//
// 	nodeCmd.Flags().BoolVarP(&shouldWriteConfig, "write-config", "w", shouldWriteConfig, "Write all specified flags to configuration file")
// }
//
// // Start a node to run continously
// func StartNode(cmd *cobra.Command, args []string) error {
// 	err := global.Current.ReadConfig()
// 	if err != nil {
// 		log.Error("Failed to read config")
// 		return err
// 	}
//
// 	log.Debug("Starting", "p2pAddress", global.Current.Config.Network.P2PAddress, "on", global.Current.NodeName)
//
// 	node := app.NewApplication()
//
// 	node.Initialize()
//
// 	app.SetNodeName(node)
// 	log.Settings()
//
// 	shared.CatchSigterm(func() {
// 		if service != nil {
// 			service.Stop()
// 		}
// 	})
//
// 	if shouldWriteConfig {
// 		err := global.Current.SaveConfig()
// 		if err != nil {
// 			log.Error("Failed to write command-line flags to configuration file", "err", err)
// 			log.Error("Continuing...")
// 		}
// 	}
//
// 	log.Debug("Waiting forever...")
// 	select {}
// }
//
//
// type FullNode struct {
// 	cfg *cfg.Server
//
// }
