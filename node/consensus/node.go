package consensus

import (
	"fmt"
	abcitypes "github.com/tendermint/tendermint/abci/types"
	tmconfig "github.com/tendermint/tendermint/config"
	tmnode "github.com/tendermint/tendermint/node"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/privval"
	"github.com/tendermint/tendermint/proxy"
	"github.com/tendermint/tendermint/types"
	"path/filepath"
)

type Node = tmnode.Node

func NewNode(application abcitypes.Application, configuration Config) (*tmnode.Node, error) {
	tmConfig := NewConfig(configuration)

	genesisDoc, err := types.GenesisDocFromFile(filepath.Join(tmConfig.RootDir, "config", "genesis.json"))
	if err != nil {
		return nil, fmt.Errorf("couldn't read genesis file: %v", err)
	}
	privValidator := privval.LoadFilePV(tmConfig.PrivValidatorKeyFile(), tmConfig.PrivValidatorStateFile())
	clientCreator := proxy.NewLocalClientCreator(application)
	nilMetricsConfig := tmconfig.InstrumentationConfig{false, "", 0, "metrics"}
	noopMetrics := tmnode.DefaultMetricsProvider(&nilMetricsConfig)

	genesisProvider := func() (*types.GenesisDoc, error) {
		return genesisDoc, nil
	}

	nodekey, err := p2p.LoadNodeKey(tmConfig.NodeKeyFile())
	if err != nil {
		return nil, fmt.Errorf("failed to locate the node key file: %v", err)
	}

	logger := NewLogger(tmConfig.RootDir+".log", *tmConfig)

	return tmnode.NewNode(
		tmConfig,
		privValidator,
		nodekey,
		clientCreator,
		genesisProvider,
		tmnode.DefaultDBProvider,
		noopMetrics,
		logger)
}
