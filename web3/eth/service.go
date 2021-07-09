package eth

import (
	"os"
	"sync"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/log"

	"github.com/Oneledger/protocol/storage"
	rpctypes "github.com/Oneledger/protocol/web3/types"

	rpcclient "github.com/Oneledger/protocol/client"
)

var _ rpctypes.EthService = (*Service)(nil)

type Service struct {
	ctx    rpctypes.Web3Context
	logger *log.Logger

	mu sync.Mutex
}

func NewService(ctx rpctypes.Web3Context) *Service {
	return &Service{ctx: ctx, logger: log.NewLoggerWithPrefix(os.Stdout, "eth")}
}

func (svc *Service) getTMClient() rpcclient.Client {
	return svc.ctx.GetAPI().RPCClient()
}

func (svc *Service) getState() *storage.State {
	return svc.ctx.GetContractStore().State
}

func (svc *Service) getStateHeight(height int64) int64 {
	switch height {
	case rpctypes.LatestBlockNumber, rpctypes.PendingBlockNumber:
		return svc.getState().Version()
	case rpctypes.EarliestBlockNumber:
		return 1
	}
	return height
}

func (svc *Service) GetStateDB() *action.CommitStateDB {
	return action.NewCommitStateDB(svc.ctx.GetContractStore(), svc.ctx.GetAccountKeeper(), svc.logger)
}
