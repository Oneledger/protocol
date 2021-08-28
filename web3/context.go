package web3

import (
	"github.com/Oneledger/protocol/app/node"
	"github.com/Oneledger/protocol/client"
	"github.com/Oneledger/protocol/config"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/evm"
	"github.com/Oneledger/protocol/data/fees"
	"github.com/Oneledger/protocol/identity"
	"github.com/Oneledger/protocol/log"
	"github.com/Oneledger/protocol/storage"
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
	feePool        *fees.Store
	nodeContext    *node.Context
	cfg            *config.Server
	chainstate     *storage.ChainState
	currencies     *balance.CurrencySet

	services map[string]rpctypes.Web3Service
}

// NewContext Initializing context with properties
func NewContext(
	logger *log.Logger, extAPI *client.ExtServiceContext,
	validatorStore *identity.ValidatorStore, feePool *fees.Store, nodeContext *node.Context, cfg *config.Server,
	chainstate *storage.ChainState, currencies *balance.CurrencySet,
) rpctypes.Web3Context {
	return &Context{logger, extAPI, validatorStore, feePool, nodeContext, cfg, chainstate, currencies, make(map[string]rpctypes.Web3Service, 0)}
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

func (ctx *Context) getImmortalState() *storage.State {
	ctx.logger.Debug("chainstate version", ctx.chainstate.Version, "hash", ctx.chainstate.Hash)
	return storage.NewState(ctx.chainstate)
}

func (ctx *Context) GetContractStore() *evm.ContractStore {
	return evm.NewContractStore(ctx.getImmortalState())
}

func (ctx *Context) GetAccountKeeper() balance.AccountKeeper {
	state := ctx.getImmortalState()
	balanceStore := balance.NewStore("b", state)
	accountKeeper := balance.NewNesterAccountKeeper(
		state,
		balanceStore,
		ctx.currencies,
	)
	return accountKeeper
}

func (ctx *Context) GetFeePool() *fees.Store {
	return ctx.feePool
}

func (ctx *Context) GetNodeContext() *node.Context {
	return ctx.nodeContext
}

func (ctx *Context) GetConfig() *config.Server {
	return ctx.cfg
}
