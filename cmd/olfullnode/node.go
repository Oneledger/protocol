// /*
// 	Copyright 2017-2018 OneLedger
//
// 	Cli to start a node (server) running.
// */
package main

import (
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/Oneledger/protocol/app"
	"github.com/Oneledger/protocol/config"
	"github.com/Oneledger/protocol/log"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var nodeCmd = &cobra.Command{
	Use:   "node",
	Short: "Start up node (server)",
	RunE:  StartNode,
}

type nodeContext struct {
	cfg               *config.Server
	logger            *log.Logger
	debug             bool
	rpc               string
	p2p               string
	persistentPeers   []string
	seeds             []string
	seedMode          bool
	shouldWriteConfig bool
	rootDir           string
}

// init reads the configuration file
func (ctx *nodeContext) init(rootDir string) error {

	ctx.logger = log.NewLoggerWithPrefix(os.Stdout, "olfullnode node")

	cfg := &config.Server{}
	rootPath, err := filepath.Abs(rootDir)
	if err != nil {
		return err
	}

	ctx.rootDir = rootPath

	err = cfg.ReadFile(cfgPath(rootPath))
	if err != nil {
		return errors.Wrapf(err, "failed to read configuration file at at %s", cfgPath(rootPath))
	}

	if ctx.rpc != "" {
		ctx.cfg.Network.RPCAddress = ctx.rpc
	}

	if ctx.p2p != "" {
		ctx.cfg.Network.P2PAddress = ctx.p2p
	}

	if len(ctx.persistentPeers) != 0 {
		ctx.cfg.P2P.PersistentPeers = ctx.persistentPeers
	}

	if len(ctx.seeds) != 0 {
		ctx.cfg.P2P.Seeds = ctx.seeds
	}

	if ctx.seedMode {
		ctx.cfg.P2P.SeedMode = ctx.seedMode
	}

	ctx.cfg = cfg

	return nil
}

var nodeCtx = &nodeContext{}

// Setup the command and flags in Cobra
func init() {
	defaults := config.DefaultServerConfig()
	RootCmd.AddCommand(nodeCmd)

	// Get information to connect to a my tendermint node
	nodeCmd.Flags().StringVarP(&nodeCtx.rpc, "address", "a",
		defaults.Network.RPCAddress, "port for rpc")

	nodeCmd.Flags().BoolVarP(&nodeCtx.debug, "debug", "d",
		false, "Set DEBUG mode")

	nodeCmd.Flags().StringArrayVar(&nodeCtx.persistentPeers, "persistent-peers", defaults.P2P.PersistentPeers, "List of persistent peers to connect to")

	// These could be moved to node persistent flags
	nodeCmd.Flags().StringVar(&nodeCtx.p2p, "p2p", defaults.Network.P2PAddress, "Address to use in P2P network")

	nodeCmd.Flags().StringArrayVar(&nodeCtx.seeds, "seeds", defaults.P2P.Seeds, "List of seeds to connect to")

	nodeCmd.Flags().BoolVar(&nodeCtx.seedMode, "seed-mode", defaults.P2P.SeedMode, "List of seeds to connect to")

	nodeCmd.Flags().BoolVarP(&nodeCtx.shouldWriteConfig, "write-config", "w", false, "Write all specified flags to configuration file")
}

// Start a node to run continously
func StartNode(cmd *cobra.Command, args []string) error {
	ctx := nodeCtx
	err := ctx.init(rootArgs.rootDir)
	if err != nil {
		return errors.Wrap(err, "failed to initialize config")
	}
	application, err := app.NewApp(ctx.cfg, ctx.rootDir)
	if err != nil {
		return errors.Wrap(err, "failed to create new app")
	}

	err = application.Start()
	if err != nil {
		return errors.Wrap(err, "failed to start app")
	}

	if ctx.shouldWriteConfig {
		err = ctx.cfg.SaveFile(cfgPath(ctx.rootDir))
		if err != nil {
			ctx.logger.Error("Failed to write command-line flags to configuration file", "err", err)
		}
	}

	catchSigTerm(ctx.logger, application.Close)

	select {}
}

func catchSigTerm(logger *log.Logger, close func()) {
	// Catch a SIGTERM and stop
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

	go func() {
		for sig := range sigs {
			logger.Info("Stopping due to", sig.String())
			close()
			os.Exit(-1)
		}
	}()

}

func cfgPath(dir string) string {
	return filepath.Join(dir, config.FileName)
}
