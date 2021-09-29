package types

import (
	"github.com/Oneledger/protocol/app/node"
	"github.com/Oneledger/protocol/config"
	"github.com/Oneledger/protocol/consensus"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/evm"
	"github.com/Oneledger/protocol/data/fees"
	"github.com/Oneledger/protocol/log"
	"github.com/Oneledger/protocol/vm"
	cs "github.com/tendermint/tendermint/consensus"
	"github.com/tendermint/tendermint/mempool"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/store"
	"github.com/tendermint/tendermint/types"
)

// Web3Service interface define service
type Web3Service interface{}

type EthService interface {
	Web3Service
	GetBlockStore() *store.BlockStore
	GetStateDB() *vm.CommitStateDB
}

// Web3Context interface to define required elements for the API Context
type Web3Context interface {
	// propagation structures
	GetLogger() *log.Logger
	GetNode() *consensus.Node
	GetBlockStore() *store.BlockStore
	GetEventBus() *types.EventBus
	GetMempool() mempool.Mempool
	GetGenesisDoc() *types.GenesisDoc
	GetConsensusReactor() *cs.Reactor
	GetSwitch() *p2p.Switch
	GetContractStore() *evm.ContractStore
	GetAccountKeeper() balance.AccountKeeper
	GetFeePool() *fees.Store
	GetNodeContext() *node.Context
	GetConfig() *config.Server

	// service registry
	RegisterService(name string, srv Web3Service)
	ServiceList() map[string]Web3Service
}
