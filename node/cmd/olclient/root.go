/*
	Copyright 2017-2018 OneLedger

	Setup the root command structure for all of cli commands
*/
package main

import (
	"fmt"
	"os"

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
	log.Debug("Parse Commands")

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

// Initialize Cobra, config global arguments
func init() {
	cobra.OnInitialize(environment)

	RootCmd.PersistentFlags().BoolVarP(&global.Current.Debug, "debug", "d", global.Current.Debug, "Set DEBUG mode")
	RootCmd.PersistentFlags().StringVarP(&global.Current.RootDir, "root", "r", global.Current.RootDir, "Set root directory")
	RootCmd.PersistentFlags().StringVarP(&global.Current.Name, "name", "n", global.Current.Name, "Set a name")
}

// Initialize Viper
func environment() {
	log.Debug("Setting up Environment")

	viper.AutomaticEnv()
}
