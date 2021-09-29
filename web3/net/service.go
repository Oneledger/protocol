package net

import (
	"errors"
	"math/big"
	"os"
	"sync"

	"github.com/Oneledger/protocol/utils"

	"github.com/Oneledger/protocol/log"
	rpctypes "github.com/Oneledger/protocol/web3/types"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

var _ rpctypes.Web3Service = (*Service)(nil)

type Service struct {
	ctx    rpctypes.Web3Context
	logger *log.Logger

	mu sync.Mutex
}

func NewService(ctx rpctypes.Web3Context) *Service {
	return &Service{ctx: ctx, logger: log.NewLoggerWithPrefix(os.Stdout, "net")}
}

func (svc *Service) Listening() bool {
	return svc.ctx.GetNode().IsListening()
}

func (svc *Service) PeerCount() hexutil.Big {
	return hexutil.Big(*big.NewInt(int64(len(svc.ctx.GetSwitch().Peers().List()))))
}

func (svc *Service) Version() (string, error) {
	svc.logger.Debug("net_version")
	genesis := svc.ctx.GetGenesisDoc()
	if genesis == nil {
		return "", errors.New("version not found")
	}
	return utils.HashToBigInt(genesis.ChainID).String(), nil
}
