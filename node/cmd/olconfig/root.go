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
	Use:   "olconfig",
	Short: "OneLedger Config Tool",
	Long:  "CLI Access to OneLedger's config files",
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
	RootCmd.PersistentFlags().BoolVarP(&global.Current.Debug, "debug", "d",
		false, "Set DEBUG mode")
}

// Initialize Viper
func environment() {
	log.Debug("Loading Environment")
	config.ClientConfig()

}
