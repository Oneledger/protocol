/*
	Copyright 2017-2018 OneLedger

	Setup the root command structure for all of cli commands
*/
package main

import (
	"fmt"
	"os"

	"github.com/Oneledger/protocol/node/config"
	"github.com/Oneledger/protocol/node/global"
	"github.com/Oneledger/protocol/node/log"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "olvm",
	Short: "olvm",
	Long:  "OneLedger Contract Engine",
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

// Initialize Cobra, config global arguments
func init() {
	cobra.OnInitialize(environment)

	RootCmd.PersistentFlags().StringVar(&global.Current.RootDir, "root",
		global.Current.RootDir, "Set root directory")

	RootCmd.PersistentFlags().StringVarP(&global.Current.ConfigName, "config", "c",
		global.Current.ConfigName, "Configuration File Name")

	RootCmd.PersistentFlags().StringVar(&global.Current.NodeName, "node",
		global.Current.NodeName, "Set a node name")

	// Get information to connect to a my tendermint node
	RootCmd.PersistentFlags().StringVarP(&global.Current.Config.Network.RPCAddress, "address", "a",
		global.Current.Config.Network.RPCAddress, "consensus address")

	RootCmd.PersistentFlags().StringVarP(&global.Current.Transport, "transport", "t",
		global.Current.Transport, "transport (socket | grpc)")

	RootCmd.PersistentFlags().BoolVarP(&global.Current.Debug, "debug", "d",
		global.Current.Debug, "Set DEBUG mode")

	RootCmd.PersistentFlags().StringVar(&global.Current.Config.Network.BTCAddress, "btcrpc",
		global.Current.Config.Network.BTCAddress, "bitcoin rpc address")

	RootCmd.PersistentFlags().StringVar(&global.Current.Config.Network.ETHAddress, "ethrpc",
		global.Current.Config.Network.ETHAddress, "ethereum rpc address")

	RootCmd.PersistentFlags().StringVar(&global.Current.TendermintRoot, "tendermintRoot",
		global.Current.TendermintRoot, "tendermint root directory")

	RootCmd.PersistentFlags().StringVar(&global.Current.Config.Network.SDKAddress, "sdkrpc",
		global.Current.Config.Network.SDKAddress, "Address for SDK RPC Server")

}

// Initialize Viper
func environment() {
	log.Debug("Loading Environment")
	// Context-loading of OLVMProtocol and OLVMAddress should be handled below
	config.ConfigureServer()
}
