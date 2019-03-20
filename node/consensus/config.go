package consensus

import (
	"path/filepath"
	"strings"

	"github.com/Oneledger/protocol/node/config"
	"github.com/Oneledger/protocol/node/global"
	tmconfig "github.com/tendermint/tendermint/config"
	tmlog "github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/node"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/privval"
	"github.com/tendermint/tendermint/types"
)

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
	var leveldb string
	if cfg.Node.DB == "goleveldb" {
		leveldb = "leveldb"
	} else {
		leveldb = cfg.Node.DB
	}
	nilMetricsConfig := tmconfig.InstrumentationConfig{Namespace: "metrics"}
	nilMetricsProvider := node.DefaultMetricsProvider(&nilMetricsConfig)

	baseConfig := tmconfig.DefaultBaseConfig()
	baseConfig.RootDir = global.Current.ConsensusDir()
	baseConfig.ProxyApp = "OneLedgerProtocol"
	baseConfig.Moniker = cfg.Node.NodeName
	baseConfig.FastSync = cfg.Node.FastSync
	baseConfig.DBBackend = leveldb
	baseConfig.DBPath = "data"
	baseConfig.LogLevel = cfg.Consensus.LogLevel

	p2pConfig := cfg.P2P.TMConfig()
	p2pConfig.ListenAddress = cfg.Network.P2PAddress
	if cfg.Network.ExternalP2PAddress == "" {
		p2pConfig.ExternalAddress = cfg.Network.P2PAddress
	} else {
		p2pConfig.ExternalAddress = cfg.Network.ExternalP2PAddress
	}

	genesisProvider := func() (*types.GenesisDoc, error) {
		return types.GenesisDocFromFile(filepath.Join(global.Current.ConsensusDir(), "config", "genesis.json"))
	}
	privValidator := privval.LoadFilePV(baseConfig.PrivValidatorKeyFile(), baseConfig.PrivValidatorStateFile())

	nodeKey, err := p2p.LoadNodeKey(baseConfig.NodeKeyFile())
	if err != nil {
		return NodeConfig{}, err
	}

	rpcConfig := tmconfig.DefaultRPCConfig()
	rpcConfig.ListenAddress = global.Current.Config.Network.RPCAddress

	tmcfg := tmconfig.Config{
		BaseConfig: baseConfig,
		RPC:        rpcConfig,
		P2P:        p2pConfig,
		Mempool:    cfg.Mempool.TMConfig(),
		Consensus:  cfg.Consensus.TMConfig(),
		TxIndex: &tmconfig.TxIndexConfig{
			Indexer:      "kv",
			IndexTags:    strings.Join(cfg.Node.IndexTags, ","),
			IndexAllTags: cfg.Node.IndexAllTags,
		},
		Instrumentation: &nilMetricsConfig,
	}

	// Decide how we want to log the consensus
	var logger tmlog.Logger
	if cfg.Consensus.LogOutput == "stdout" {
		logger, err = newStdOutLogger(tmcfg)
	} else {
		logOutput := filepath.Join(global.Current.RootDir, cfg.Consensus.LogOutput)
		logger, err = newFileLogger(logOutput, tmcfg)
	}
	if err != nil {
		return NodeConfig{}, err
	}
	//for testing propose, we should change it for production
	cfg.P2P.AddrBookStrict = false
	cfg.P2P.AllowDuplicateIP = true

	var dbProvider node.DBProvider
	dbProvider = node.DefaultDBProvider

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
