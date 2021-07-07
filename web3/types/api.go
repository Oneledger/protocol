package types

import (
	"github.com/Oneledger/protocol/app/node"
	"github.com/Oneledger/protocol/client"
	"github.com/Oneledger/protocol/config"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/evm"
	"github.com/Oneledger/protocol/identity"
	"github.com/Oneledger/protocol/log"
)

// Interface to define required elements for the API Context
type Web3Context interface {
	// propagation structures
	GetLogger() *log.Logger
	GetAPI() *client.ExtServiceContext
	GetValidatorStore() *identity.ValidatorStore
	GetContractStore() *evm.ContractStore
	GetAccountKeeper() balance.AccountKeeper
	GetNodeContext() *node.Context
	GetConfig() *config.Server

	// service registry
	DefaultRegisterForAll()
	RegisterService(name string, srv interface{})
	ServiceList() map[string]interface{}
}
