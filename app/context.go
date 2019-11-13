package app

import (
	"io"
	"path/filepath"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/action/btc"
	action_ons "github.com/Oneledger/protocol/action/ons"
	"github.com/Oneledger/protocol/action/staking"
	"github.com/Oneledger/protocol/action/transfer"
	"github.com/Oneledger/protocol/app/node"
	"github.com/Oneledger/protocol/client"
	"github.com/Oneledger/protocol/config"
	"github.com/Oneledger/protocol/data/accounts"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/bitcoin"
	"github.com/Oneledger/protocol/data/chain"
	"github.com/Oneledger/protocol/data/fees"
	"github.com/Oneledger/protocol/data/governance"
	"github.com/Oneledger/protocol/data/jobs"
	"github.com/Oneledger/protocol/data/ons"
	"github.com/Oneledger/protocol/identity"
	"github.com/Oneledger/protocol/log"
	"github.com/Oneledger/protocol/rpc"
	"github.com/Oneledger/protocol/service"
	"github.com/Oneledger/protocol/storage"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/db"
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

	balances   *balance.Store
	domains    *ons.DomainStore
	validators *identity.ValidatorStore // Set of validators currently active
	feePool    *fees.Store
	govern     *governance.Store
	trackers   *bitcoin.TrackerStore // tracker for bitcoin balance UTXO

	currencies *balance.CurrencySet
	feeOption  *fees.FeeOption

	//storage which is not a chain state
	accounts accounts.Wallet

	jobStore        *jobs.JobStore
	lockScriptStore *bitcoin.LockScriptStore
	internalService *action.Service

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
	ctx.trackers = bitcoin.NewTrackerStore("btct", storage.NewState(ctx.chainstate))

	ctx.accounts = accounts.NewWallet(cfg, ctx.dbDir())

	// TODO check if validator
	if true {
		ctx.jobStore = jobs.NewJobStore(cfg, ctx.dbDir())
		ctx.lockScriptStore = bitcoin.NewLockScriptStore(cfg, ctx.dbDir())

	}

	ctx.actionRouter = action.NewRouter("action")
	ctx.feeOption = &fees.FeeOption{
		FeeCurrency:   balance.Currency{Name: "OLT", Chain: chain.Type(0), Decimal: 18},
		MinFeeDecimal: 0,
	}

	_ = transfer.EnableSend(ctx.actionRouter)
	_ = staking.EnableApplyValidator(ctx.actionRouter)
	_ = action_ons.EnableONS(ctx.actionRouter)
	_ = btc.EnableBTC(ctx.actionRouter)
	return ctx, nil
}

func (ctx context) dbDir() string {
	return filepath.Join(ctx.cfg.RootDir(), ctx.cfg.Node.DBDir)
}

func (ctx *context) Action(header *Header, state *storage.State) *action.Context {

	var params *chaincfg.Params
	switch ctx.cfg.ChainDriver.BitcoinChainType {
	case "mainnet":
		params = &chaincfg.MainNetParams
	case "testnet3":
		params = &chaincfg.TestNet3Params
	case "regtest":
		params = &chaincfg.RegressionNetParams
	case "simnet":
		params = &chaincfg.SimNetParams
	default:
		params = &chaincfg.TestNet3Params
	}

	actionCtx := action.NewContext(
		ctx.actionRouter,
		header,
		state,
		ctx.accounts,
		ctx.balances.WithState(state),
		ctx.currencies,
		ctx.feeOption,
		ctx.feePool.WithState(state),
		ctx.validators.WithState(state),
		ctx.domains.WithState(state),

		ctx.trackers.WithState(state),
		ctx.jobStore,
		params,

		log.NewLoggerWithPrefix(ctx.logWriter, "action").WithLevel(log.Level(ctx.cfg.Node.LogLevel)))

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
	svcCtx := &service.Context{
		Balances:     ctx.balances,
		Accounts:     ctx.accounts,
		Currencies:   ctx.currencies,
		FeeOpt:       ctx.feeOption,
		Cfg:          ctx.cfg,
		NodeContext:  ctx.node,
		ValidatorSet: ctx.validators,
		Domains:      ctx.domains,
		Router:       ctx.actionRouter,
		Logger:       log.NewLoggerWithPrefix(ctx.logWriter, "rpc").WithLevel(log.Level(ctx.cfg.Node.LogLevel)),
		Services:     extSvcs,

		Trackers: ctx.trackers,
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
		FeeOpt:       ctx.feeOption,
		Cfg:          ctx.cfg,
		NodeContext:  ctx.node,
		ValidatorSet: ctx.validators,
		Domains:      ctx.domains,
		Router:       ctx.actionRouter,
		Logger:       log.NewLoggerWithPrefix(ctx.logWriter, "restful").WithLevel(log.Level(ctx.cfg.Node.LogLevel)),
		Services:     extSvcs,

		Trackers: ctx.trackers,
	}
	return service.NewRestfulService(svcCtx).Router(), nil
}

// Close all things that need to be closed
func (ctx *context) Close() {
	closers := []closer{ctx.db, ctx.accounts, ctx.rpc}
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
