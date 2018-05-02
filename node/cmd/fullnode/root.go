/*
	Copyright 2017-2018 OneLedger

	Setup the root command structure for all of cli commands
*/
package main

import (
	"fmt"
	"os"

	"github.com/Oneledger/prototype/node/app"
	"github.com/Oneledger/prototype/node/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var RootCmd = &cobra.Command{
	Use:   "fullnode",
	Short: "fullnode",
	Long:  "OneLedger chain, fullnode",
}

func Execute() {
	log.Debug("Parse Commamnds")

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

// Initialize Cobra, config global arguments
func init() {
	cobra.OnInitialize(environment)

	RootCmd.PersistentFlags().BoolVarP(&app.Current.Debug, "debug", "d", app.Current.Debug, "Set DEBUG mode")
	RootCmd.PersistentFlags().StringVarP(&app.Current.RootDir, "root", "r", app.Current.RootDir, "Set root directory")
	RootCmd.PersistentFlags().StringVarP(&app.Current.Name, "name", "n", app.Current.Name, "Set a name")
}

// Initialize Viper
func environment() {
	log.Debug("Setting up Environment")
	viper.AutomaticEnv()
}
