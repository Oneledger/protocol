package consensus

import (
	abcitypes "github.com/tendermint/tendermint/abci/types"
	tmnode "github.com/tendermint/tendermint/node"
	"github.com/tendermint/tendermint/proxy"
)

type Node = tmnode.Node

func NewNode(application abcitypes.Application, nodecfg NodeConfig) (*tmnode.Node, error) {

	clientCreator := proxy.NewLocalClientCreator(application)
	return tmnode.NewNode(
		&nodecfg.CFG,
		nodecfg.privValidator,
		nodecfg.nodeKey,
		clientCreator,
		nodecfg.legacyGenesisProvider,
		tmnode.DefaultDBProvider,
		nodecfg.metricsProvider,
		nodecfg.logger)
}
