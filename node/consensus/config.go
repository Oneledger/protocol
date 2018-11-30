package consensus

import (
	"strings"

	"github.com/tendermint/tendermint/config"
)

// Config is a OneLedger-specific configuration struct used to create a tendermint configuration
type Config struct {
	Moniker         string
	RootDirectory   string
	RPCAddress      string
	P2PAddress      string
	IndexTags       []string
	PersistentPeers string
}

// NewConfig returns a ready-to-go tendermint configuration
func NewConfig(olcfg Config) *config.Config {
	rootDir := olcfg.RootDirectory

	cfg := config.DefaultConfig()
	cfg.BaseConfig.Moniker = olcfg.Moniker
	cfg.BaseConfig.ProxyApp = "OneLedger"
	cfg.BaseConfig.LogLevel = "main:info,state:info,consensus:info,*:error"

	cfg.RPC.ListenAddress = olcfg.RPCAddress
	cfg.P2P.ListenAddress = olcfg.P2PAddress
	cfg.P2P.PersistentPeers = olcfg.PersistentPeers

	// TODO: Turn this off for production
	cfg.P2P.AddrBookStrict = false

	cfg.TxIndex.IndexTags = strings.Join(olcfg.IndexTags, ",") // TODO: Put this in global

	cfg.BaseConfig.RootDir = rootDir
	cfg.RPC.RootDir = rootDir
	cfg.P2P.RootDir = rootDir
	cfg.Mempool.RootDir = rootDir
	cfg.Consensus.RootDir = rootDir
	cfg.Consensus.CreateEmptyBlocks = false
	return cfg
}
