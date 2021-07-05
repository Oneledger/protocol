package eth

import (
	"math/big"

	"github.com/Oneledger/protocol/client"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/evm"
	"github.com/Oneledger/protocol/log"
	"github.com/Oneledger/protocol/utils"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

type Service struct {
	log           *log.Logger
	ext           client.ExtServiceContext
	contracts     *evm.ContractStore
	accountKeeper balance.AccountKeeper
}

func NewService(logger *log.Logger, ext client.ExtServiceContext, contracts *evm.ContractStore, accountKeeper balance.AccountKeeper) *Service {
	return &Service{
		log:           logger,
		ext:           ext,
		contracts:     contracts,
		accountKeeper: accountKeeper,
	}
}

func (svc *Service) ChainId() hexutil.Big {
	block := svc.ext.Block(0).Block
	chainID := utils.HashToBigInt(block.Header.ChainID)
	return hexutil.Big(*chainID)
}

func (svc *Service) BlockNumber() hexutil.Big {
	height := svc.contracts.State.Version()
	blockNumber := big.NewInt(height)
	return hexutil.Big(*blockNumber)
}
