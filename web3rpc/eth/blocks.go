package eth

import (
	"math/big"

	"github.com/Oneledger/protocol/web3rpc"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

func (svc *web3rpc.Service) BlockNumber() hexutil.Big {
	height := svc.getState().Version()
	blockNumber := big.NewInt(height)
	return hexutil.Big(*blockNumber)
}
