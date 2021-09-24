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
	netInfo, err := svc.ctx.GetAPI().RPCClient().NetInfo()
	if err != nil {
		svc.logger.Error(err)
		return false
	}
	return netInfo.Listening
}

func (svc *Service) PeerCount() hexutil.Big {
	netInfo, err := svc.ctx.GetAPI().RPCClient().NetInfo()
	if err != nil {
		svc.logger.Error(err)
		return hexutil.Big(*big.NewInt(int64(0)))
	}

	return hexutil.Big(*big.NewInt(int64(netInfo.NPeers)))
}

func (svc *Service) Version() (string, error) {
	svc.logger.Debug("net_version")
	blockResult, err := svc.ctx.GetAPI().RPCClient().Block(nil)
	if err != nil {
		return "", err
	}
	if blockResult.Block == nil {
		return "", errors.New("not loaded")
	}
	return utils.HashToBigInt(blockResult.Block.ChainID).String(), nil
}
