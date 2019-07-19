package consensus

import (
	"path/filepath"

	"github.com/Oneledger/protocol/config"
	tmconfig "github.com/tendermint/tendermint/config"
	tmlog "github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/node"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/privval"
	"github.com/tendermint/tendermint/types"
)

type Config = tmconfig.Config

// config is used to provider the right arguments for spinning up a new consensus.Node
type NodeConfig struct {
	CFG             tmconfig.Config
	genesisProvider node.GenesisDocProvider
	privValidator   types.PrivValidator
	nodeKey         *p2p.NodeKey
	dbProvider      node.DBProvider
	metricsProvider node.MetricsProvider
	logger          tmlog.Logger
}

// TMConfig returns a ready to go config for starting a new tendermint node,
// fields like logging and metrics still need to be handled before starting the new node

func ParseConfig(cfg *config.Server) (NodeConfig, error) {
	return parseConfig(cfg)
}

// ParseConfig reads Tendermint level config and return as
func parseConfig(cfg *config.Server) (NodeConfig, error) {
	// Proper consensus dir
	rootDir := Dir(cfg.RootDir())
	tmcfg := cfg.TMConfig()

	genesisProvider := func() (*types.GenesisDoc, error) {
		return types.GenesisDocFromFile(filepath.Join(rootDir, "config", "genesis.json"))
	}

	// If either file path does not exist, the program will exit !
	privValidator := privval.LoadFilePV(tmcfg.PrivValidatorKeyFile(), tmcfg.PrivValidatorStateFile())

	nodeKey, err := p2p.LoadNodeKey(tmcfg.NodeKeyFile())
	if err != nil {
		return NodeConfig{}, err
	}

	// Decide how we want to log the consensus
	var logger tmlog.Logger
	if cfg.Consensus.LogOutput == "stdout" {
		logger, err = newStdOutLogger(tmcfg)
	} else {
		logOutput := filepath.Join(cfg.RootDir(), cfg.Consensus.LogOutput)
		logger, err = newFileLogger(logOutput, tmcfg)
	}
	if err != nil {
		return NodeConfig{}, err
	}

	//for testing propose, we should change it for production
	cfg.P2P.AddrBookStrict = false
	cfg.P2P.AllowDuplicateIP = true

	var dbProvider node.DBProvider = node.DefaultDBProvider

	nilMetricsProvider := node.DefaultMetricsProvider(tmcfg.Instrumentation)
	// TODO: Switch DB provider depending on the value of CFG.Node.DB
	return NodeConfig{
		CFG:             tmcfg,
		genesisProvider: genesisProvider,
		metricsProvider: nilMetricsProvider,
		privValidator:   privValidator,
		nodeKey:         nodeKey,
		dbProvider:      dbProvider,
		logger:          logger,
	}, nil
}
