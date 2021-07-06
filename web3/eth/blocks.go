package eth

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

func (svc *Service) BlockNumber() hexutil.Big {
	height := svc.getState().Version()
	blockNumber := big.NewInt(height)
	return hexutil.Big(*blockNumber)
}
