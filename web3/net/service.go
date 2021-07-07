package net

import (
	"github.com/Oneledger/protocol/utils"
	"math/big"
	"sync"

	"github.com/Oneledger/protocol/log"
	web3types "github.com/Oneledger/protocol/web3/types"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

type Service struct {
	ctx    web3types.Web3Context
	logger *log.Logger

	mu sync.Mutex
}

func NewService(ctx web3types.Web3Context) *Service {
	return &Service{ctx: ctx, logger: ctx.GetLogger()}
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

	return utils.HashToBigInt(blockResult.Block.ChainID).String(), nil
}
