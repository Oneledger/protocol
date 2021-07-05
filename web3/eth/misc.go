package eth

import (
	"github.com/Oneledger/protocol/utils"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

func (svc *Service) ChainId() hexutil.Big {
	block := svc.ext.Block(0).Block
	chainID := utils.HashToBigInt(block.Header.ChainID)
	return hexutil.Big(*chainID)
}
