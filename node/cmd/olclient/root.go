/*
	Copyright 2017-2018 OneLedger

	Setup the root command structure for all of cli commands
*/
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/Oneledger/protocol/node/cmd/shared"
	"github.com/Oneledger/protocol/node/config"
	"github.com/Oneledger/protocol/node/global"
	"github.com/Oneledger/protocol/node/log"
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

	RootCmd.PersistentFlags().StringVarP(&global.Current.ConfigName, "config", "c",
		global.Current.ConfigName, "Configuration File Name")

	RootCmd.PersistentFlags().StringVar(&global.Current.NodeName, "node",
		global.Current.NodeName, "Set a node name")

	RootCmd.PersistentFlags().BoolVarP(&global.Current.Debug, "debug", "d",
		global.Current.Debug, "Set DEBUG mode")

	RootCmd.PersistentFlags().StringVarP(&global.Current.Transport, "transport", "t",
		global.Current.Transport, "transport (socket | grpc)")

	RootCmd.PersistentFlags().StringVarP(&global.Current.RpcAddress, "address", "a",
		global.Current.RpcAddress, "full address")

	RootCmd.PersistentFlags().StringVar(&global.Current.SDKAddress, "sdkrpc",
		global.Current.SDKAddress, "SDK address")

}

func handleError(err error) {
	shared.Console.Error(err)
	os.Exit(1)
}

func indentJSON(in []byte) bytes.Buffer {
	var out bytes.Buffer
	err := json.Indent(&out, in, "-", "  ")
	if err != nil {
		handleError(err)
	}

	return out
}

// Initialize Viper
func environment() {
	log.Debug("Loading Environment")
	config.ClientConfig()

	// TODO: Static variables vs Dynamic variables :-(
	global.Current.SDKAddress = viper.Get("SDKAddress").(string)
	global.Current.RpcAddress = viper.Get("RpcAddress").(string)

	//viper.AutomaticEnv()
}
