/*
   ____             _              _                      _____           _                  _
  / __ \           | |            | |                    |  __ \         | |                | |
 | |  | |_ __   ___| |     ___  __| | __ _  ___ _ __     | |__) | __ ___ | |_ ___   ___ ___ | |
 | |  | | '_ \ / _ \ |    / _ \/ _` |/ _` |/ _ \ '__|    |  ___/ '__/ _ \| __/ _ \ / __/ _ \| |
 | |__| | | | |  __/ |___|  __/ (_| | (_| |  __/ |       | |   | | | (_) | || (_) | (_| (_) | |
  \____/|_| |_|\___|______\___|\__,_|\__, |\___|_|       |_|   |_|  \___/ \__\___/ \___\___/|_|
                                      __/ |
                                     |___/


Copyright 2017 - 2019 OneLedger
*/

package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "olclient",
	Short: "OneLedger client",
	Long:  "Client access to the OneLedger chain",
}

var rootArgs = rootArguments{}

type rootArguments struct {
	rootDir string
}

// Initialize Cobra, config global arguments
func init() {

	RootCmd.PersistentFlags().StringVar(&rootArgs.rootDir, "root",
		"./", "Set root directory")

	/*
		RootCmd.PersistentFlags().StringVar(&global.Current.NodeName, "node",
			global.Current.NodeName, "Set a node name")

		RootCmd.PersistentFlags().BoolVarP(&global.Current.Debug, "debug", "d",
			global.Current.Debug, "Set DEBUG mode")

		RootCmd.PersistentFlags().StringVarP(&global.Current.Config.Network.RPCAddress, "address", "a",
			global.Current.Config.Network.RPCAddress, "full address")

		RootCmd.PersistentFlags().StringVar(&global.Current.Config.Network.SDKAddress, "sdkrpc",
			global.Current.Config.Network.SDKAddress, "SDK address")
	*/
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
