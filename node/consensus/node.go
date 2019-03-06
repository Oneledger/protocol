package consensus

import (
	abcitypes "github.com/tendermint/tendermint/abci/types"
	tmconfig "github.com/tendermint/tendermint/config"
	tmnode "github.com/tendermint/tendermint/node"
	"github.com/tendermint/tendermint/proxy"
	"github.com/tendermint/tendermint/types"
)

type Node = tmnode.Node

func NewNode(
	application abcitypes.Application,
	configuration Config,
	privValidator types.PrivValidator,
	genesisDoc *types.GenesisDoc,
) (*tmnode.Node, error) {
	tmConfig := NewConfig(configuration)

	clientCreator := proxy.NewLocalClientCreator(application)
	nilMetricsConfig := tmconfig.InstrumentationConfig{false, "", 0}
	noopMetrics := tmnode.DefaultMetricsProvider(&nilMetricsConfig)

	genesisProvider := func() (*types.GenesisDoc, error) {
		return genesisDoc, nil
	}

	logger := NewLogger(tmConfig.RootDir+".log", *tmConfig)

	return tmnode.NewNode(
		tmConfig,
		privValidator,
		clientCreator,
		genesisProvider,
		tmnode.DefaultDBProvider,
		noopMetrics,
		logger)
}
