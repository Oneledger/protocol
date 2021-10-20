package consensus

import (
	abcitypes "github.com/tendermint/tendermint/abci/types"
	tmnode "github.com/tendermint/tendermint/node"
	"github.com/tendermint/tendermint/proxy"
	"github.com/tendermint/tendermint/store"
	dbm "github.com/tendermint/tm-db"
)

type Node = tmnode.Node

func NewNode(application abcitypes.Application, nodecfg NodeConfig, loadedBlockStore chan<- *store.BlockStore) (*tmnode.Node, error) {
	clientCreator := proxy.NewLocalClientCreator(application)
	node, err := tmnode.NewNode(
		&nodecfg.CFG,
		nodecfg.privValidator,
		nodecfg.nodeKey,
		clientCreator,
		nodecfg.legacyGenesisProvider,
		// tmnode.DefaultDBProvider,
		HackWrapDBProvider(loadedBlockStore),
		nodecfg.metricsProvider,
		nodecfg.logger)
	return node, err
}

// HackWrapDBProvider replace default db provider to get block store befor replay
func HackWrapDBProvider(loadedBlockStore chan<- *store.BlockStore) func(ctx *tmnode.DBContext) (dbm.DB, error) {
	return func(ctx *tmnode.DBContext) (dbm.DB, error) {
		dbType := dbm.BackendType(ctx.Config.DBBackend)
		db := dbm.NewDB(ctx.ID, dbType, ctx.Config.DBDir())
		// TODO: Remove hardcore
		if ctx.ID == "blockstore" {
			loadedBlockStore <- store.NewBlockStore(db)
		}
		return db, nil
	}
}
