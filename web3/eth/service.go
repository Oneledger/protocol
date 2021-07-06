package eth

import (
	"sync"

	"github.com/Oneledger/protocol/log"

	"github.com/Oneledger/protocol/storage"
	web3types "github.com/Oneledger/protocol/web3/types"

	rpcclient "github.com/Oneledger/protocol/client"
)

type Service struct {
	ctx    web3types.Web3Context
	logger *log.Logger

	mu sync.Mutex
}

func NewService(ctx web3types.Web3Context) *Service {
	return &Service{ctx: ctx, logger: ctx.GetLogger()}
}

func (svc *Service) getTMClient() rpcclient.Client {
	return svc.ctx.GetAPI().RPCClient()
}

func (svc *Service) getState() *storage.State {
	return svc.ctx.GetContractStore().State
}

func (svc *Service) getStateHeight(height int64) int64 {
	switch height {
	case -1, -2:
		return svc.getState().Version()
	}
	return height
}
