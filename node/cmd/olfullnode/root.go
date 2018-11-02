/*
	Copyright 2017-2018 OneLedger

	Setup the root command structure for all of cli commands
*/
package main

import (
	"fmt"
	"os"

	"github.com/Oneledger/protocol/node/global"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var RootCmd = &cobra.Command{
	Use:   "olfullnode",
	Short: "olfullnode",
	Long:  "OneLedger chain, olfullnode",
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

	RootCmd.PersistentFlags().StringVar(&global.Current.NodeName, "node",
		global.Current.NodeName, "Set a node name")

	// Get information to connect to an ABCI app (myself)
	RootCmd.PersistentFlags().StringVar(&global.Current.AppAddress, "app",
		global.Current.AppAddress, "app address")

	// Get information to connect to a my tendermint node
	RootCmd.PersistentFlags().StringVarP(&global.Current.RpcAddress, "address", "a",
		global.Current.RpcAddress, "consensus address")

	RootCmd.PersistentFlags().StringVarP(&global.Current.Transport, "transport", "t",
		global.Current.Transport, "transport (socket | grpc)")

	RootCmd.PersistentFlags().BoolVarP(&global.Current.Debug, "debug", "d",
		global.Current.Debug, "Set DEBUG mode")

	RootCmd.PersistentFlags().StringVar(&global.Current.BTCAddress, "btcrpc",
		global.Current.BTCAddress, "bitcoin rpc address")

	RootCmd.PersistentFlags().StringVar(&global.Current.ETHAddress, "ethrpc",
		global.Current.ETHAddress, "etherem rpc address")

	RootCmd.PersistentFlags().StringVar(&global.Current.TendermintRoot, "tendermintRoot",
		global.Current.TendermintRoot, "tendermint root directory")
}

// Initialize Viper
func environment() {
	viper.AutomaticEnv()
}
