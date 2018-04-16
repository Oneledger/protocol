/*
	Copyright 2017-2018 OneLedger

	Setup the root command structure for all of cli commands
*/
package main

import (
	"fmt"
	"os"

	"github.com/Oneledger/prototype/node/app"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var RootCmd = &cobra.Command{
	Use:   "fullnode",
	Short: "fullnode",
	Long:  "OneLedger chain, fullnode",
}

func Execute() {
	app.Log.Debug("Parse Commamnds")

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

// Initialize Cobra, config global arguments
func init() {
	cobra.OnInitialize(environment)
	RootCmd.PersistentFlags().StringVarP(&app.Current.RootDir, "root", "r", "Invalid", "Set root directory")
	RootCmd.PersistentFlags().BoolVarP(&app.Current.Debug, "debug", "d", false, "Set DEBUG mode")
}

// Initialize Viper
func environment() {
	app.Log.Debug("Setting up Environment")
	viper.AutomaticEnv()
}
