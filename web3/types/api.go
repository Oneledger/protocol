package types

import (
	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/app/node"
	"github.com/Oneledger/protocol/client"
	"github.com/Oneledger/protocol/config"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/evm"
	"github.com/Oneledger/protocol/data/fees"
	"github.com/Oneledger/protocol/identity"
	"github.com/Oneledger/protocol/log"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
)

// Web3Service interface define service
type Web3Service interface{}

type EthService interface {
	Web3Service
	GetStateDB() *action.CommitStateDB
	GetBlockByHash(hash common.Hash, fullTx bool) (*Block, error)
	GetBlockByNumber(blockNrOrHash rpc.BlockNumberOrHash, fullTx bool) (*Block, error)
}

// Web3Context interface to define required elements for the API Context
type Web3Context interface {
	// propagation structures
	GetLogger() *log.Logger
	GetAPI() *client.ExtServiceContext
	GetValidatorStore() *identity.ValidatorStore
	GetContractStore() *evm.ContractStore
	GetAccountKeeper() balance.AccountKeeper
	GetFeePool() *fees.Store
	GetNodeContext() *node.Context
	GetConfig() *config.Server

	// service registry
	DefaultRegisterForAll()
	RegisterService(name string, srv Web3Service)
	ServiceList() map[string]Web3Service
}
