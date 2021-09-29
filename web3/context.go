package web3

import (
	"reflect"
	"unsafe"

	"github.com/Oneledger/protocol/app/node"
	"github.com/Oneledger/protocol/config"
	"github.com/Oneledger/protocol/consensus"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/evm"
	"github.com/Oneledger/protocol/data/fees"
	"github.com/Oneledger/protocol/log"
	"github.com/Oneledger/protocol/storage"
	"github.com/Oneledger/protocol/web3/eth"
	"github.com/Oneledger/protocol/web3/net"
	rpctypes "github.com/Oneledger/protocol/web3/types"
	"github.com/Oneledger/protocol/web3/web3"
	cs "github.com/tendermint/tendermint/consensus"
	"github.com/tendermint/tendermint/mempool"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/state/txindex"
	"github.com/tendermint/tendermint/store"
	"github.com/tendermint/tendermint/types"
)

// Context set up for service processing
type Context struct {
	logger      *log.Logger
	node        *consensus.Node
	feePool     *fees.Store
	nodeContext *node.Context
	cfg         *config.Server
	chainstate  *storage.ChainState
	currencies  *balance.CurrencySet

	services map[string]rpctypes.Web3Service
}

// NewContext Initializing context with properties
func NewContext(
	logger *log.Logger, node *consensus.Node,
	feePool *fees.Store, nodeContext *node.Context, cfg *config.Server,
	chainstate *storage.ChainState, currencies *balance.CurrencySet,
) rpctypes.Web3Context {
	ctx := &Context{logger, node, feePool, nodeContext, cfg, chainstate, currencies, make(map[string]rpctypes.Web3Service, 0)}
	ctx.defaultRegisterForAll()
	return ctx
}

// defaultRegisterForAll regs services from the namespaces
func (ctx *Context) defaultRegisterForAll() {
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

func (ctx *Context) GetNode() *consensus.Node {
	return ctx.node
}

func (ctx *Context) GetBlockStore() *store.BlockStore {
	return ctx.node.BlockStore()
}

func (ctx *Context) GetEventBus() *types.EventBus {
	return ctx.node.EventBus()
}

func (ctx *Context) GetMempool() mempool.Mempool {
	return ctx.node.Mempool()
}

func (ctx *Context) GetGenesisDoc() *types.GenesisDoc {
	return ctx.node.GenesisDoc()
}

func (ctx *Context) GetConsensusReactor() *cs.Reactor {
	return ctx.node.ConsensusReactor()
}

func (ctx *Context) GetSwitch() *p2p.Switch {
	return ctx.node.Switch()
}

func (ctx *Context) GetTxIndexer() txindex.TxIndexer {
	// NOTE: Do not touch this, until it will be public available!
	// hack to get txIndexer as it is private field in tendermint core
	tiField := reflect.ValueOf(ctx.node).Elem().FieldByName("txIndexer")
	// unlock for modification
	tiField = reflect.NewAt(tiField.Type(), unsafe.Pointer(tiField.UnsafeAddr())).Elem()
	return tiField.Interface().(txindex.TxIndexer)
}

func (ctx *Context) getImmortalState() *storage.State {
	return storage.NewState(ctx.chainstate)
}

func (ctx *Context) GetContractStore() *evm.ContractStore {
	return evm.NewContractStore(ctx.getImmortalState())
}

func (ctx *Context) GetAccountKeeper() balance.AccountKeeper {
	balanceStore := balance.NewStore("b", ctx.getImmortalState())
	accountKeeper := balance.NewNesterAccountKeeper(
		ctx.getImmortalState(),
		balanceStore,
		ctx.currencies,
	)
	return accountKeeper
}

func (ctx *Context) GetFeePool() *fees.Store {
	fstore := fees.NewStore("f", ctx.getImmortalState())
	opts := *ctx.feePool.GetOpt()
	fstore.SetupOpt(&opts)
	return fstore
}

func (ctx *Context) GetNodeContext() *node.Context {
	return ctx.nodeContext
}

func (ctx *Context) GetConfig() *config.Server {
	return ctx.cfg
}
