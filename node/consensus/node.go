package consensus

import (
	"github.com/Oneledger/protocol/node/app"
	tmlog "github.com/tendermint/tendermint/libs/log"
	tmnode "github.com/tendermint/tendermint/node"
	"github.com/tendermint/tendermint/proxy"
	"github.com/tendermint/tendermint/types"
)

type genesisProvider tmnode.GenesisDocProvider
type node tmnode.DBProvider

func NewNode(
	application app.Application,
	configuration *Config,
	privValidator types.PrivValidator,
	genesisDoc *types.GenesisDoc,
	logger tmlog.Logger,
) (*tmnode.Node, error) {
	tmConfig := NewConfig(configuration)

	clientCreator := proxy.NewLocalClientCreator(application)
	noopMetrics := tmnode.DefaultMetricsProvider(nil)

	genesisProvider := func() (*types.GenesisDoc, error) {
		return genesisDoc, nil
	}

	return tmnode.NewNode(
		tmConfig,
		privValidator,
		clientCreator,
		genesisProvider,
		tmnode.DefaultDBProvider,
		noopMetrics,
		// TODO: Separate logging
		logger)
}

// func NewNode(config *cfg.Config,
// 	privValidator types.PrivValidator,
// 	clientCreator proxy.ClientCreator,
// 	genesisDocProvider GenesisDocProvider,
// 	dbProvider DBProvider,
// 	metricsProvider MetricsProvider,
// 	logger log.Logger) (*Node, error) {
