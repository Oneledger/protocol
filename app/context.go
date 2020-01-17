package app

import (
	"io"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/db"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/action/btc"
	"github.com/Oneledger/protocol/action/eth"
	action_ons "github.com/Oneledger/protocol/action/ons"
	"github.com/Oneledger/protocol/action/staking"
	"github.com/Oneledger/protocol/action/transfer"
	"github.com/Oneledger/protocol/app/node"
	"github.com/Oneledger/protocol/client"
	"github.com/Oneledger/protocol/config"
	"github.com/Oneledger/protocol/data/accounts"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/bitcoin"
	"github.com/Oneledger/protocol/data/ethereum"
	"github.com/Oneledger/protocol/data/fees"
	"github.com/Oneledger/protocol/data/governance"
	"github.com/Oneledger/protocol/data/jobs"
	"github.com/Oneledger/protocol/data/ons"
	"github.com/Oneledger/protocol/event"
	"github.com/Oneledger/protocol/identity"
	"github.com/Oneledger/protocol/log"
	"github.com/Oneledger/protocol/rpc"
	"github.com/Oneledger/protocol/service"
	"github.com/Oneledger/protocol/storage"
)

// The base context for the application, holds databases and other stateful information contained by the app.
// Used to derive other package-level Contexts
type context struct {
	node node.Context
	cfg  config.Server

	rpc          *rpc.Server
	actionRouter action.Router

	//db for chain state storage
	db         db.DB
	chainstate *storage.ChainState
	check      *storage.State
	deliver    *storage.State

	balances    *balance.Store
	domains     *ons.DomainStore
	validators  *identity.ValidatorStore // Set of validators currently active
	feePool     *fees.Store
	govern      *governance.Store
	btcTrackers *bitcoin.TrackerStore  // tracker for bitcoin balance UTXO
	ethTrackers *ethereum.TrackerStore // tracker for ethereum tracker store

	currencies *balance.CurrencySet

	//storage which is not a chain state
	accounts accounts.Wallet

	jobStore        *jobs.JobStore
	lockScriptStore *bitcoin.LockScriptStore
	internalService *event.Service
	jobBus          *event.JobBus

	logWriter io.Writer
}

func newContext(logWriter io.Writer, cfg config.Server, nodeCtx *node.Context) (context, error) {
	ctx := context{
		cfg:        cfg,
		logWriter:  logWriter,
		currencies: balance.NewCurrencySet(),
		node:       *nodeCtx,
	}

	ctx.rpc = rpc.NewServer(logWriter, &cfg)

	db, err := storage.GetDatabase("chainstate", ctx.dbDir(), ctx.cfg.Node.DB)
	if err != nil {
		return ctx, errors.Wrap(err, "initial db failed")
	}
	ctx.db = db
	ctx.chainstate = storage.NewChainState("chainstate", db)
	ctx.chainstate.SetupRotation(10, 100, 10)
	ctx.deliver = storage.NewState(ctx.chainstate)
	ctx.check = storage.NewState(ctx.chainstate)

	ctx.validators = identity.NewValidatorStore("v", cfg, storage.NewState(ctx.chainstate))
	ctx.balances = balance.NewStore("b", storage.NewState(ctx.chainstate))
	ctx.domains = ons.NewDomainStore("d", storage.NewState(ctx.chainstate))
	ctx.feePool = fees.NewStore("f", storage.NewState(ctx.chainstate))
	ctx.govern = governance.NewStore("g", storage.NewState(ctx.chainstate))

	ctx.btcTrackers = bitcoin.NewTrackerStore("btct", storage.NewState(ctx.chainstate))

	ctx.ethTrackers = ethereum.NewTrackerStore("etht", storage.NewState(ctx.chainstate))
	ctx.accounts = accounts.NewWallet(cfg, ctx.dbDir())

	// TODO check if validator
	// valAddr := ctx.node.ValidatorAddress()

	ctx.jobStore = jobs.NewJobStore(cfg, ctx.dbDir())
	ctx.lockScriptStore = bitcoin.NewLockScriptStore(cfg, ctx.dbDir())

	ctx.actionRouter = action.NewRouter("action")

	ctx.jobBus = event.NewJobBus(event.Option{
		BtcInterval: 30 * time.Second,
		EthInterval: 3 * time.Second,
	}, ctx.jobStore)

	_ = transfer.EnableSend(ctx.actionRouter)
	_ = staking.EnableApplyValidator(ctx.actionRouter)
	_ = action_ons.EnableONS(ctx.actionRouter)
	_ = btc.EnableBTC(ctx.actionRouter)
	_ = eth.EnableETH(ctx.actionRouter)
	return ctx, nil
}

func (ctx context) dbDir() string {
	return filepath.Join(ctx.cfg.RootDir(), ctx.cfg.Node.DBDir)
}

func (ctx *context) Action(header *Header, state *storage.State) *action.Context {

	actionCtx := action.NewContext(
		ctx.actionRouter,
		header,
		state,
		ctx.accounts,
		ctx.balances.WithState(state),
		ctx.currencies,
		ctx.feePool.WithState(state),
		ctx.validators.WithState(state),
		ctx.domains.WithState(state),

		ctx.btcTrackers.WithState(state),
		ctx.ethTrackers.WithState(state),
		ctx.jobStore,
		ctx.lockScriptStore,
		log.NewLoggerWithPrefix(ctx.logWriter, "action").WithLevel(log.Level(ctx.cfg.Node.LogLevel)),
	)

	return actionCtx
}

func (ctx *context) ID() {}
func (ctx *context) Accounts() accounts.Wallet {
	return ctx.accounts
}

func (ctx *context) ValidatorCtx() *identity.ValidatorContext {
	return identity.NewValidatorContext(
		ctx.balances.WithState(ctx.deliver),
		ctx.feePool.WithState(ctx.deliver),
	)
}

// Returns a balance.Context
func (ctx *context) Balances() *balance.Context {
	return balance.NewContext(
		log.NewLoggerWithPrefix(ctx.logWriter, "balances").WithLevel(log.Level(ctx.cfg.Node.LogLevel)),
		ctx.balances,
		ctx.currencies)
}

func (ctx *context) Services() (service.Map, error) {
	extSvcs, err := client.NewExtServiceContext(ctx.cfg.Network.RPCAddress, ctx.cfg.Network.SDKAddress)
	if err != nil {
		return nil, errors.Wrap(err, "failed to start service context")
	}

	btcTrackers := bitcoin.NewTrackerStore("btct", storage.NewState(ctx.chainstate))
	btcTrackers.SetConfig(ctx.btcTrackers.GetConfig())

	svcCtx := &service.Context{
		Balances:     ctx.balances,
		Accounts:     ctx.accounts,
		Currencies:   ctx.currencies,
		FeePool:      ctx.feePool,
		Cfg:          ctx.cfg,
		NodeContext:  ctx.node,
		ValidatorSet: ctx.validators,
		Domains:      ctx.domains,
		Router:       ctx.actionRouter,
		Logger:       log.NewLoggerWithPrefix(ctx.logWriter, "rpc").WithLevel(log.Level(ctx.cfg.Node.LogLevel)),
		Services:     extSvcs,

		Trackers: btcTrackers,
	}

	return service.NewMap(svcCtx)
}

func (ctx *context) Restful() (service.RestfulRouter, error) {
	extSvcs, err := client.NewExtServiceContext(ctx.cfg.Network.RPCAddress, ctx.cfg.Network.SDKAddress)
	if err != nil {
		return nil, errors.Wrap(err, "failed to start service context")
	}
	svcCtx := &service.Context{
		Balances:     ctx.balances,
		Accounts:     ctx.accounts,
		Currencies:   ctx.currencies,
		FeePool:      ctx.feePool,
		NodeContext:  ctx.node,
		ValidatorSet: ctx.validators,
		Domains:      ctx.domains,
		Router:       ctx.actionRouter,
		Logger:       log.NewLoggerWithPrefix(ctx.logWriter, "restful").WithLevel(log.Level(ctx.cfg.Node.LogLevel)),
		Services:     extSvcs,

		Trackers: ctx.btcTrackers,
	}
	return service.NewRestfulService(svcCtx).Router(), nil
}

type StorageCtx struct {
	Balances   *balance.Store
	Domains    *ons.DomainStore
	Validators *identity.ValidatorStore // Set of validators currently active
	FeePool    *fees.Store
	Govern     *governance.Store

	Currencies *balance.CurrencySet
	FeeOption  *fees.FeeOption
	Hash       []byte
	Version    int64
}

func (ctx *context) Storage() StorageCtx {
	return StorageCtx{
		Version:    ctx.chainstate.Version,
		Hash:       ctx.chainstate.Hash,
		Balances:   ctx.balances,
		Domains:    ctx.domains,
		Validators: ctx.validators,
		FeePool:    ctx.feePool,
		Govern:     ctx.govern,
		Currencies: ctx.currencies,
		FeeOption:  ctx.feePool.GetOpt(),
	}
}

// Close all things that need to be closed
func (ctx *context) Close() {
	closers := []closer{ctx.db, ctx.accounts, ctx.rpc, ctx.jobBus}
	for _, closer := range closers {
		closer.Close()
	}
}

func (ctx *context) Node() node.Context {
	return ctx.node
}

func (ctx *context) Validators() *identity.ValidatorStore {
	return ctx.validators
}

func (ctx *context) Replay(version int64) error {
	return ctx.chainstate.ClearFrom(version)
}

func (ctx *context) JobContext() *event.JobsContext {

	return event.NewJobsContext(
		ctx.cfg,
		ctx.internalService,
		ctx.btcTrackers,
		ctx.validators,
		ctx.node.ValidatorECDSAPrivateKey(),
		ctx.node.ValidatorECDSAPrivateKey(),
		ctx.node.ValidatorAddress(),
		ctx.lockScriptStore,
		ctx.ethTrackers.WithState(ctx.deliver),
	)
}
