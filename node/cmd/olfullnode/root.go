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
	RootCmd.PersistentFlags().StringVar(&global.Current.RootDir, "root",
		global.Current.RootDir, "Set root directory")

	RootCmd.PersistentFlags().StringVar(&global.Current.NodeName, "node",
		global.Current.NodeName, "Set a node name")
}
