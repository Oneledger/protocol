package web3

import (
	"github.com/Oneledger/protocol/app/node"
	"github.com/Oneledger/protocol/client"
	"github.com/Oneledger/protocol/config"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/evm"
	"github.com/Oneledger/protocol/identity"
	"github.com/Oneledger/protocol/log"
	"github.com/Oneledger/protocol/web3/eth"
	"github.com/Oneledger/protocol/web3/net"
	rpctypes "github.com/Oneledger/protocol/web3/types"
	"github.com/Oneledger/protocol/web3/web3"
)

// Context set up for service processing
type Context struct {
	logger         *log.Logger
	externalAPI    *client.ExtServiceContext
	validatorStore *identity.ValidatorStore
	contractStore  *evm.ContractStore
	accountKeeper  balance.AccountKeeper
	nodeContext    *node.Context
	cfg            *config.Server

	services map[string]rpctypes.Web3Service
}

// NewContext Initializing context with properties
func NewContext(
	logger *log.Logger, extAPI *client.ExtServiceContext,
	validatorStore *identity.ValidatorStore, contractStore *evm.ContractStore,
	accountKeeper balance.AccountKeeper, nodeContext *node.Context, cfg *config.Server,
) rpctypes.Web3Context {
	return &Context{logger, extAPI, validatorStore, contractStore, accountKeeper, nodeContext, cfg, make(map[string]rpctypes.Web3Service, 0)}
}

// DefaultRegisterForAll regs services from the namespaces
func (ctx *Context) DefaultRegisterForAll() {
	ctx.RegisterService("eth", eth.NewService(ctx))
	ctx.RegisterService("net", net.NewService(ctx))
	ctx.RegisterService("web3", web3.NewService(ctx))
}

// RegisterService used to register service. NOTE: Must be called by service
func (ctx *Context) RegisterService(name string, service rpctypes.Web3Service) {
	ctx.services[name] = service
}

// ServiceList represents a cuurent list of registered servises
func (ctx *Context) ServiceList() map[string]rpctypes.Web3Service {
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

func (ctx *Context) GetNodeContext() *node.Context {
	return ctx.nodeContext
}

func (ctx *Context) GetConfig() *config.Server {
	return ctx.cfg
}
