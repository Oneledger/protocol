/*
	Copyright 2017-2018 OneLedger

	Setup the root command structure for all of cli commands
*/
package main

import (
	"fmt"
	"os"

	"github.com/Oneledger/protocol/log"
	"github.com/spf13/cobra"
)

var logger *log.Logger

var rootArgs = &rootArguments{}

type rootArguments struct {
	rootDir string
}

var RootCmd = &cobra.Command{
	Use:   "olfullnode",
	Short: "olfullnode",
	Long:  "OneLedger chain, olfullnode",
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}

// Initialize Cobra, config global arguments
func init() {
	// Initialize logger
	logger = log.NewLoggerWithPrefix(os.Stdout, "olfullnode:")
	RootCmd.PersistentFlags().StringVar(&rootArgs.rootDir, "root",
		"./", "Set root directory")
}
