package consensus

import (
	"strings"

	"github.com/Oneledger/protocol/node/global"
	"github.com/Oneledger/protocol/node/log"
	"github.com/tendermint/tendermint/config"
)

// Config is a OneLedger-specific configuration struct used to create a tendermint configuration
type Config struct {
	Moniker         string
	RootDirectory   string
	RPCAddress      string
	P2PAddress      string
	IndexTags       []string
	PersistentPeers []string
	Seeds           string
	SeedMode        bool
}

// NewConfig returns a ready-to-go tendermint configuration
func NewConfig(olcfg Config) *config.Config {
	rootDir := olcfg.RootDirectory

	log.Dump("Initial Config", olcfg)

	cfg := config.DefaultConfig()
	cfg.BaseConfig.Moniker = olcfg.Moniker
	cfg.BaseConfig.ProxyApp = "OneLedger"

	// TODO: Tendermint should have it's own debug flag
	if global.Current.Debug {
		cfg.BaseConfig.LogLevel = "debug"
	} else {
		cfg.BaseConfig.LogLevel = "main:info,state:info,consensus:info,*:error"
	}

	cfg.RPC.ListenAddress = olcfg.RPCAddress
	cfg.P2P.ListenAddress = olcfg.P2PAddress

	cfg.P2P.PersistentPeers = strings.Join(olcfg.PersistentPeers, ",")
	cfg.P2P.Seeds = olcfg.Seeds
	cfg.P2P.SeedMode = olcfg.SeedMode
	cfg.P2P.PexReactor = true

	// When debug is set, loosen the constraints
	if global.Current.Debug {

	}
	//for testing propose, we should change it for production
	cfg.P2P.AddrBookStrict = false
	cfg.P2P.AllowDuplicateIP = true

	cfg.TxIndex.IndexTags = strings.Join(olcfg.IndexTags, ",") // TODO: Put this in global

	cfg.BaseConfig.RootDir = rootDir
	cfg.RPC.RootDir = rootDir
	cfg.P2P.RootDir = rootDir
	cfg.Mempool.RootDir = rootDir
	cfg.Consensus.RootDir = rootDir
	cfg.Consensus.CreateEmptyBlocks = false
	return cfg
}
