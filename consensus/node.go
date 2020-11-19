package consensus

import (
	abcitypes "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/node"
	tmnode "github.com/tendermint/tendermint/node"
	"github.com/tendermint/tendermint/proxy"
	"github.com/tendermint/tendermint/store"
	dbm "github.com/tendermint/tm-db"
)

type Node = tmnode.Node

func TempBlockStore(nodecfg NodeConfig) (blockStore *store.BlockStore, blockStoreDB dbm.DB, err error) {
	config := &nodecfg.CFG
	dbProvider := tmnode.DefaultDBProvider

	blockStoreDB, err = dbProvider(&node.DBContext{"blockstore", config})
	if err != nil {
		return
	}
	blockStore = store.NewBlockStore(blockStoreDB)
	return
}

func NewNode(application abcitypes.Application, nodecfg NodeConfig) (*tmnode.Node, error) {

	clientCreator := proxy.NewLocalClientCreator(application)
	return tmnode.NewNode(
		&nodecfg.CFG,
		nodecfg.privValidator,
		nodecfg.nodeKey,
		clientCreator,
		nodecfg.genesisProvider,
		tmnode.DefaultDBProvider,
		nodecfg.metricsProvider,
		nodecfg.logger)
}
