package eth

import (
	"os"
	"sync"

	"github.com/Oneledger/protocol/log"
	"github.com/Oneledger/protocol/vm"
	"github.com/tendermint/tendermint/mempool"
	"github.com/tendermint/tendermint/store"

	"github.com/Oneledger/protocol/storage"
	rpcfilters "github.com/Oneledger/protocol/web3/eth/filters"
	rpctypes "github.com/Oneledger/protocol/web3/types"
)

var _ rpctypes.EthService = (*Service)(nil)

type Service struct {
	ctx    rpctypes.Web3Context
	logger *log.Logger

	filterAPI *rpcfilters.PublicFilterAPI
	mu        sync.Mutex
}

func NewService(ctx rpctypes.Web3Context) *Service {
	return &Service{
		ctx:       ctx,
		logger:    log.NewLoggerWithPrefix(os.Stdout, "eth"),
		filterAPI: rpcfilters.NewPublicFilterAPI(ctx.GetBlockStore(), ctx.GetEventBus()),
	}
}

func (svc *Service) GetBlockStore() *store.BlockStore {
	return svc.ctx.GetBlockStore()
}

func (svc *Service) GetMempool() mempool.Mempool {
	return svc.ctx.GetMempool()
}

func (svc *Service) getState() *storage.State {
	return svc.ctx.GetContractStore().State
}

func (svc *Service) getStateHeight(height int64) int64 {
	switch height {
	case rpctypes.LatestBlockNumber, rpctypes.PendingBlockNumber:
		return svc.getState().Version()
	case rpctypes.EarliestBlockNumber:
		return rpctypes.InitialBlockNumber
	}
	return height
}

func (svc *Service) GetStateDB() *vm.CommitStateDB {
	stateDB := vm.NewCommitStateDB(svc.ctx.GetContractStore(), svc.ctx.GetAccountKeeper(), svc.logger)
	stateDB.SetBlockStore(svc.ctx.GetBlockStore())
	return stateDB
}
