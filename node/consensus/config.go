package consensus

import (
	"strings"

	"github.com/tendermint/tendermint/config"
)

// Parameters to configure a tendermint node
type Config struct {
	Moniker       string
	RootDirectory string
	RPCAddress    string
	P2PAddress    string
	ProxyAddress  string
	IndexTags     []string
}

func NewConfig(olConfig *Config) *config.Config {
	cfg := config.DefaultConfig()
	cfg.BaseConfig.Moniker = olConfig.Moniker
	cfg.BaseConfig.ProxyApp = "OneLedger"

	cfg.RPC.ListenAddress = olConfig.RPCAddress
	cfg.P2P.ListenAddress = olConfig.P2PAddress

	cfg.TxIndex.IndexTags = strings.Join(olConfig.IndexTags, ",") // TODO: Put this in global

	cfg.BaseConfig.RootDir = olConfig.RootDirectory
	cfg.RPC.RootDir = olConfig.RootDirectory
	cfg.P2P.RootDir = olConfig.RootDirectory
	cfg.Mempool.RootDir = olConfig.RootDirectory
	cfg.Consensus.RootDir = olConfig.RootDirectory

	return cfg
}
