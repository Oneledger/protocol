/*
	Copyright 2017-2018 OneLedger

	Cli to init a node (server)
*/
package main

import (
	"os"
	"path/filepath"

	"github.com/Oneledger/protocol/node/config"
	"github.com/Oneledger/protocol/node/global"
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/privval"
	"github.com/tendermint/tendermint/types"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize node (server)",
	RunE:  initNode,
}

type InitCmdArguments struct {
	genesis  string
	nodeName string
}

var initCmdArguments = &InitCmdArguments{}

func init() {
	RootCmd.AddCommand(initCmd)
	initCmd.Flags().StringVar(&initCmdArguments.nodeName, "node_name", "Newton-Node", "Name of the node")
	initCmd.Flags().StringVar(&initCmdArguments.genesis, "genesis", "", "Genesis file to use to generate new node key file")
}

func initNode(cmd *cobra.Command, _ []string) error {
	args := initCmdArguments

	if _, err := os.Stat(global.Current.RootDir); os.IsNotExist(err) {
		err = os.Mkdir(global.Current.RootDir, config.DirPerms)
		if err != nil {
			return err
		}
	}
	// Generate new configuration file
	cfg := config.DefaultServerConfig()
	cfg.Node.NodeName = args.nodeName
	err := cfg.SaveFile(filepath.Join(global.Current.RootDir, config.FileName))
	if err != nil {
		return err
	}

	// If the genesis path given is absolute, just search that path, otherwise join it
	// with the currently set root directory
	var genesisPath string
	if filepath.IsAbs(args.genesis) {
		genesisPath = args.genesis
	} else {
		genesisPath = filepath.Join(global.Current.RootDir, args.genesis)
	}
	genesisdoc, err := types.GenesisDocFromFile(genesisPath)
	if err != nil {
		return err
	}
	configDir := filepath.Join(global.Current.RootDir, "consensus", "config")
	dataDir := filepath.Join(global.Current.RootDir, "consensus", "data")
	nodeDataDir := filepath.Join(global.Current.RootDir, "nodedata")

	err = os.MkdirAll(configDir, config.DirPerms)
	if err != nil {
		return err
	}

	err = os.MkdirAll(dataDir, config.DirPerms)
	if err != nil {
		return err
	}

	err = os.MkdirAll(nodeDataDir, config.DirPerms)
	if err != nil {
		return err
	}

	err = genesisdoc.SaveAs(filepath.Join(configDir, "genesis.json"))
	if err != nil {
		return err
	}
	// Make node key
	_, err = p2p.LoadOrGenNodeKey(filepath.Join(configDir, "node_key.json"))
	if err != nil {
		return err
	}

	// Make private validator file
	pvFile := privval.GenFilePV(filepath.Join(configDir, "priv_validator_key.json"), filepath.Join(dataDir, "priv_validator_state.json"))
	pvFile.Save()

	return nil
}
