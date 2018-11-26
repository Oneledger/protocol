/*
	Copyright 2017-2018 OneLedger

	Cli to interact with a with the chains.
*/
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Oneledger/protocol/node/cmd/shared"
	"github.com/Oneledger/protocol/node/config"
	"github.com/Oneledger/protocol/node/global"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/p2p"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a OneLedger configuration parameter",
	Run:   HandleArguments,
}

// TODO: typing should be way better, see if cobra can help with this...
type GetArguments struct {
	names []string
	peers bool
}

var arguments *GetArguments = &GetArguments{}

func init() {
	RootCmd.AddCommand(getCmd)

	// TODO: I want to have a default account?
	// Transaction Parameters
	getCmd.Flags().StringArrayP("parameter", "p", arguments.names, "account name")
	getCmd.Flags().Bool("peers", false, "handle peers")

}

// IssueRequest sends out a sendTx to all of the nodes in the chain
func HandleArguments(cmd *cobra.Command, args []string) {
	names, err := cmd.Flags().GetStringArray("parameter")
	if err != nil {
		return
	}

	peers, err := cmd.Flags().GetBool("peers")
	if err != nil {
		return
	}

	if peers {
		GetPeers(names)
	} else {
		GetParams(names)
	}
}

// Get the set of static peers
func GetPeers(nodes []string) {
	buffer := ""
	for _, name := range nodes {
		// Reset the config name, and load the relevant config file
		global.Current.ConfigName = name
		config.ServerConfig()

		// Pick out a couple of the parameters
		nodeName := viper.Get("NodeName").(string)
		address := viper.Get("P2PAddress").(string)
		p2pAddress := strings.TrimPrefix(address, "tcp://")

		// Call Tendermint to get it's node id
		data := filepath.Join(os.Getenv("OLDATA"), nodeName, "tendermint", "config", "node_key.json")
		nodeKey, err := p2p.LoadNodeKey(data)
		if err != nil {
			shared.Console.Error(err)
			return
		}

		// Assemble
		entry := string(nodeKey.ID()) + "@" + p2pAddress
		if buffer == "" {
			buffer = entry
		} else {
			buffer += "," + entry
		}
	}
	fmt.Printf("%s", buffer)
	return
}

func GetParams(names []string) {
	config.ServerConfig()

	// Don't check it, if it is empty, return empty string
	first := true
	for _, name := range names {
		text := viper.Get(name)
		if text == nil {
			text = ""
		}
		// Print directly to the screen, so this can be used in shell variables
		if first {
			first = false
		} else {
			fmt.Printf("\n")
		}
		fmt.Printf("%s", text)
	}
}
