package web3

import (
	"github.com/Oneledger/protocol/app/node"
	"github.com/Oneledger/protocol/client"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/evm"
	"github.com/Oneledger/protocol/identity"
	"github.com/Oneledger/protocol/log"

	"github.com/Oneledger/protocol/web3/eth"
	web3types "github.com/Oneledger/protocol/web3/types"
)

// Context set up for service processing
type Context struct {
	logger         *log.Logger
	externalAPI    *client.ExtServiceContext
	validatorStore *identity.ValidatorStore
	contractStore  *evm.ContractStore
	accountKeeper  balance.AccountKeeper
	nodeContext    node.Context

	services map[string]interface{}
}

// NewContext Initializing context with properties
func NewContext(
	logger *log.Logger, extAPI *client.ExtServiceContext,
	validatorStore *identity.ValidatorStore, contractStore *evm.ContractStore,
	accountKeeper balance.AccountKeeper, nodeContext node.Context,
) web3types.Web3Context {
	return &Context{logger, extAPI, validatorStore, contractStore, accountKeeper, nodeContext, make(map[string]interface{}, 0)}
}

// DefaultRegisterForAll regs services from the namespaces
func (ctx *Context) DefaultRegisterForAll() {
	ctx.RegisterService("eth", eth.NewService(ctx))
}

// RegisterService used to register service. NOTE: Must be called by service
func (ctx *Context) RegisterService(name string, service interface{}) {
	ctx.services[name] = service
}

// ServiceList represents a cuurent list of registered servises
func (ctx *Context) ServiceList() map[string]interface{} {
	return ctx.services
}

func (ctx *Context) GetLogger() *log.Logger {
	return ctx.logger
}

func (ctx *Context) GetAPI() *client.ExtServiceContext {
	return ctx.externalAPI
}

func (ctx *Context) GetValidatorStore() *identity.ValidatorStore {
	return ctx.validatorStore
}

func (ctx *Context) GetContractStore() *evm.ContractStore {
	return ctx.contractStore
}

func (ctx *Context) GetAccountKeeper() balance.AccountKeeper {
	return ctx.accountKeeper
}

func (ctx *Context) GetNodeContext() node.Context {
	return ctx.nodeContext
}
