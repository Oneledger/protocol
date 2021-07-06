package net

import (
	"github.com/Oneledger/protocol/web3rpc"
	"math/big"

	"github.com/Oneledger/protocol/utils"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

func (svc *web3rpc.Service) Listening() bool {
	netInfo, err := svc.ctx.GetAPI().RPCClient().NetInfo()
	if err != nil {
		svc.logger.Error(err)
		return false
	}
	return netInfo.Listening
}

func (svc *web3rpc.Service) PeerCount() hexutil.Big {
	netInfo, err := svc.ctx.GetAPI().RPCClient().NetInfo()
	if err != nil {
		svc.logger.Error(err)
		return hexutil.Big(*big.NewInt(int64(0)))
	}

	return hexutil.Big(*big.NewInt(int64(netInfo.NPeers)))
}

func (svc *web3rpc.Service) Version() string {
	svc.logger.Debug("net_version??")

	block := svc.ctx.GetAPI().Block(1).Block
	svc.logger.Debug("???????  ", block.Header.ChainID)
	return utils.HashToBigInt(block.Header.ChainID).String()
}
