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
	Use:   "olclient",
	Short: "OneLedger client",
	Long:  "Client access to the OneLedger chain",
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

	RootCmd.PersistentFlags().BoolVarP(&global.Current.Debug, "debug", "d",
		global.Current.Debug, "Set DEBUG mode")

	RootCmd.PersistentFlags().StringVarP(&global.Current.Transport, "transport", "t",
		global.Current.Transport, "transport (socket | grpc)")

	RootCmd.PersistentFlags().StringVarP(&global.Current.RpcAddress, "address", "a",
		global.Current.RpcAddress, "full address")

	RootCmd.PersistentFlags().Int64VarP(&global.Current.Sequence, "sequence", "s",
		global.Current.Sequence, "unique sequence id")
}

// Initialize Viper
func environment() {
	viper.AutomaticEnv()
}
