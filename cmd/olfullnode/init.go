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

func newInitContext(args *InitCmdArguments, rootArgs *rootArguments) (*initContext, error) {
	logger := log.NewLoggerWithPrefix(os.Stdout, "olfullnode init")

	rootDir, err := filepath.Abs(rootArgs.rootDir)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Invalid root directory specified %s", rootArgs))
	}

	genesisPath, err := filepath.Abs(args.genesis)
	if err != nil {
		return nil, err
	}
	genesis, err := types.GenesisDocFromFile(genesisPath)
	if err != nil {
		return nil, err
	}

	return &initContext{
		rootDir: rootDir,
		logger:  logger,
		genesis: genesis,
	}, nil
}

func runInitNode(cmd *cobra.Command, _ []string) error {
	ctx, err := newInitContext(initCmdArgs, rootArgs)
	if err != nil {
		return err
	}
	return initNode(ctx)
}

// Given the path of a genesis file and a specified root directory, initNode creates all the configuration files
// needed to run a fullnode inside that specified directory
func initNode(ctx *initContext) error {
	if _, err := os.Stat(ctx.rootDir); os.IsNotExist(err) {
		err = os.Mkdir(ctx.rootDir, config.DirPerms)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("Failed to create the specified rootdir at %s", ctx.rootDir))
		}
	}

	// Generate new configuration file
	cfg := config.DefaultServerConfig()
	cfg.Node.NodeName = ctx.nodeName
	cfgPath := filepath.Join(ctx.rootDir, config.FileName)
	err := cfg.SaveFile(cfgPath)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Failed to create config file at %s", cfgPath))
	}

	csRoot := consensus.RootDirName
	csConfig := consensus.ConfigDirName
	csData := consensus.DataDirName

	configDir := filepath.Join(ctx.rootDir, csRoot, csConfig)
	dataDir := filepath.Join(ctx.rootDir, csRoot, csData)
	nodeDataDir := filepath.Join(ctx.rootDir, "nodedata")

	dirs := []string{configDir, dataDir, nodeDataDir}
	for _, dir := range dirs {
		err = os.MkdirAll(dir, config.DirPerms)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("Dir creation failed at %s", dir))
		}
	}

	// Put the genesis file in the right place
	err = ctx.genesis.SaveAs(filepath.Join(configDir, consensus.GenesisFilename))
	if err != nil {
		return errors.Wrap(err, "Failed to save genesis file")
	}
	// Make node key
	_, err = p2p.LoadOrGenNodeKey(filepath.Join(configDir, consensus.NodeKeyFilename))
	if err != nil {
		return errors.Wrap(err, "Failed to generate node key")
	}

	// Make private validator file
	pvFile := privval.GenFilePV(filepath.Join(configDir, consensus.PrivValidatorKeyFilename), filepath.Join(dataDir, consensus.PrivValidatorStateFilename))
	pvFile.Save()

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

	return nil
}
