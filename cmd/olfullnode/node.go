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
	"sync"
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

type nodeCmdContext struct {
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
func (ctx *nodeCmdContext) init(rootDir string) error {

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

	ctx.cfg = cfg

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

	return nil
}

var nodeCtx = &nodeCmdContext{}

// Setup the command and flags in Cobra
func init() {
	RootCmd.AddCommand(nodeCmd)

	// Get information to connect to a my tendermint node
	nodeCmd.Flags().StringVarP(&nodeCtx.rpc, "address", "a",
		"", "port for rpc")

	nodeCmd.Flags().BoolVarP(&nodeCtx.debug, "debug", "d",
		false, "Set DEBUG mode")

	nodeCmd.Flags().StringArrayVar(&nodeCtx.persistentPeers, "persistent-peers", []string{}, "List of persistent peers to connect to")

	// These could be moved to node persistent flags
	nodeCmd.Flags().StringVar(&nodeCtx.p2p, "p2p", "", "Address to use in P2P network")

	nodeCmd.Flags().StringArrayVar(&nodeCtx.seeds, "seeds", []string{}, "List of seeds to connect to")

	nodeCmd.Flags().BoolVar(&nodeCtx.seedMode, "seed-mode", false, "List of seeds to connect to")

	nodeCmd.Flags().BoolVarP(&nodeCtx.shouldWriteConfig, "write-config", "w", false, "Write all specified flags to configuration file")
}

// Start a node to run continously
func StartNode(cmd *cobra.Command, args []string) error {
	waiter := sync.WaitGroup{}

	ctx := nodeCtx
	err := ctx.init(rootArgs.rootDir)
	if err != nil {
		return errors.Wrap(err, "failed to initialize config")
	}

	appNodeContext, err := app.NewNodeContext(ctx.cfg)
	if err != nil {
		return errors.Wrap(err, "failed to create app's node context")
	}

	application, err := app.NewApp(ctx.cfg, appNodeContext)
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

	waiter.Add(1)
	catchSigTerm(ctx.logger, application.Close, &waiter)

	waiter.Wait()
	return nil
}

func catchSigTerm(logger *log.Logger, close func(), waiter *sync.WaitGroup) {
	// Catch a SIGTERM and stop
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

	go func() {
		for sig := range sigs {
			logger.Info("Stopping due to", sig.String())
			close()

			waiter.Done()
		}
	}()

}

func cfgPath(dir string) string {
	return filepath.Join(dir, config.FileName)
}
