/*
	Copyright 2017-2018 OneLedger

	Cli to init a node (server)
*/
package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Oneledger/protocol/data/keys"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	"github.com/Oneledger/protocol/config"
	"github.com/Oneledger/protocol/consensus"
	"github.com/Oneledger/protocol/log"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/privval"
	"github.com/tendermint/tendermint/types"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize node (server)",
	RunE:  runInitNode,
}

type InitCmdArguments struct {
	genesis  string
	nodeName string
}

var initCmdArgs = &InitCmdArguments{}

func init() {
	RootCmd.AddCommand(initCmd)
	initCmd.Flags().StringVar(&initCmdArgs.nodeName, "node_name", "Newton-Node", "Name of the node")
	initCmd.Flags().StringVar(&initCmdArgs.genesis, "genesis", "", "Genesis file to use to generate new node key file")
}

type initContext struct {
	genesis  *consensus.GenesisDoc
	logger   *log.Logger
	rootDir  string
	nodeName string
}

// Given the path of a genesis file and a specified root directory, initNode creates all the configuration files
// needed to run a fullnode inside that specified directory
func runInitNode(cmd *cobra.Command, _ []string) error {
	rootDir, err := filepath.Abs(rootArgs.rootDir)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Invalid root directory specified %s", rootArgs))
	}

	if _, err := os.Stat(rootDir); os.IsNotExist(err) {
		err = os.Mkdir(rootDir, config.DirPerms)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("Failed to create the specified rootdir at %s", rootDir))
		}
	}

	// Generate new configuration file
	cfg := config.DefaultServerConfig()
	cfg.Node.NodeName = initCmdArgs.nodeName
	cfgPath := filepath.Join(rootDir, config.FileName)
	err = cfg.SaveFile(cfgPath)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Failed to create config file at %s", cfgPath))
	}

	csRoot := consensus.RootDirName
	csConfig := consensus.ConfigDirName
	csData := consensus.DataDirName

	configDir := filepath.Join(rootDir, csRoot, csConfig)
	dataDir := filepath.Join(rootDir, csRoot, csData)
	nodeDataDir := filepath.Join(rootDir, "nodedata")

	dirs := []string{configDir, dataDir, nodeDataDir}
	for _, dir := range dirs {
		err = os.MkdirAll(dir, config.DirPerms)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("Dir creation failed at %s", dir))
		}
	}

	// Put the genesis file in the right place
	if initCmdArgs.genesis != "" {

		genesisPath, err := filepath.Abs(initCmdArgs.genesis)
		if err != nil {
			return errors.Wrap(err, "invalid genesis file path")
		}
		fmt.Println("verifying genesis file provided")
		genesis, err := types.GenesisDocFromFile(genesisPath)
		if err != nil {
			return err
		}
		err = genesis.SaveAs(filepath.Join(configDir, consensus.GenesisFilename))
		if err != nil {
			return errors.Wrap(err, "Failed to save genesis file")
		}
		// Make node key

	} else {
		fmt.Println("no genesis file provided, node is not runnable until genesis file is provided at: ", filepath.Join(configDir, consensus.GenesisFilename))
	}

	nodekey, err := p2p.LoadOrGenNodeKey(filepath.Join(configDir, consensus.NodeKeyFilename))
	if err != nil {
		return errors.Wrap(err, "Failed to generate node key")
	}
	fmt.Println("node key address: ", nodekey.PubKey().Address().String())

	// Make private validator file
	pvFile := privval.GenFilePV(filepath.Join(configDir, consensus.PrivValidatorKeyFilename),
		filepath.Join(dataDir, consensus.PrivValidatorStateFilename))
	pvFile.Save()
	fmt.Println("validator key address: ", pvFile.GetAddress().String())
	fmt.Println("validator public key: ", pvFile.GetPubKey())

	ecdsaPrivKey := secp256k1.GenPrivKey()
	ecdsaPrivKeyBytes := base64.StdEncoding.EncodeToString([]byte(ecdsaPrivKey[:]))
	_, err = keys.GetPrivateKeyFromBytes([]byte(ecdsaPrivKey[:]), keys.SECP256K1)
	if err != nil {
		return errors.Wrap(err, "error generating secp256k1 private key")
	}

	ecdsaFile := strings.Replace(consensus.PrivValidatorKeyFilename, ".json", "_ecdsa.json", 1)

	f, err := os.Create(filepath.Join(configDir, ecdsaFile))
	if err != nil {
		return errors.Wrap(err, "failed to open file to write validator ecdsa private key")
	}
	n, err := f.Write([]byte(ecdsaPrivKeyBytes))
	if err != nil && n != len(ecdsaPrivKeyBytes) {
		return errors.Wrap(err, "failed to write validator ecdsa private key")
	}
	err = f.Close()
	if err != nil && n != len(ecdsaPrivKeyBytes) {
		return errors.Wrap(err, "failed to save validator ecdsa private key")
	}
	fmt.Println("witness key address: ", ecdsaPrivKey.PubKey().Address().String())

	return nil
}
