package app

import (
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/Oneledger/protocol/data"

	"github.com/pkg/errors"
	db "github.com/tendermint/tm-db"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/action/eth"
	action_gov "github.com/Oneledger/protocol/action/governance"
	action_ons "github.com/Oneledger/protocol/action/ons"
	"github.com/Oneledger/protocol/action/staking"
	"github.com/Oneledger/protocol/action/transfer"
	"github.com/Oneledger/protocol/app/node"
	"github.com/Oneledger/protocol/client"
	"github.com/Oneledger/protocol/config"
	"github.com/Oneledger/protocol/data/accounts"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/bitcoin"
	"github.com/Oneledger/protocol/data/delegation"
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

	extStores   data.Router //External Stores
	balances    *balance.Store
	domains     *ons.DomainStore
	validators  *identity.ValidatorStore // Set of validators currently active
	witnesses   *identity.WitnessStore   // Set of witnesses currently active
	feePool     *fees.Store
	govern      *governance.Store
	btcTrackers *bitcoin.TrackerStore  // tracker for bitcoin balance UTXO
	ethTrackers *ethereum.TrackerStore // Tracker store for ongoing ethereum trackers
	currencies  *balance.CurrencySet
	//storage which is not a chain state
	accounts accounts.Wallet

	jobStore        *jobs.JobStore
	lockScriptStore *bitcoin.LockScriptStore
	internalRouter  action.Router
	internalService *event.Service
	jobBus          *event.JobBus
	proposalMaster  *governance.ProposalMasterStore
	delegators      *delegation.DelegationStore
	logWriter       io.Writer
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
	errRotation := ctx.chainstate.SetupRotation(ctx.cfg.Node.ChainStateRotation)
	if errRotation != nil {
		return ctx, errors.Wrap(errRotation, "error in loading chain state rotation config")
	}
	ctx.deliver = storage.NewState(ctx.chainstate)
	ctx.check = storage.NewState(ctx.chainstate)

	ctx.validators = identity.NewValidatorStore("v", storage.NewState(ctx.chainstate))
	ctx.witnesses = identity.NewWitnessStore("w", storage.NewState(ctx.chainstate))
	ctx.balances = balance.NewStore("b", storage.NewState(ctx.chainstate))
	ctx.domains = ons.NewDomainStore("d", storage.NewState(ctx.chainstate))
	ctx.feePool = fees.NewStore("f", storage.NewState(ctx.chainstate))
	ctx.govern = governance.NewStore("g", storage.NewState(ctx.chainstate))
	ctx.proposalMaster = NewProposalMasterStore(ctx.chainstate)
	ctx.delegators = delegation.NewDelegationStore("st", storage.NewState(ctx.chainstate))
	ctx.btcTrackers = bitcoin.NewTrackerStore("btct", storage.NewState(ctx.chainstate))

	ctx.ethTrackers = ethereum.NewTrackerStore("etht", "ethfailed", "ethsuccess", storage.NewState(ctx.chainstate))
	ctx.accounts = accounts.NewWallet(cfg, ctx.dbDir())

	// TODO check if validator
	// valAddr := ctx.node.ValidatorAddress()

	ctx.jobStore = jobs.NewJobStore(cfg, ctx.dbDir())
	ctx.lockScriptStore = bitcoin.NewLockScriptStore(cfg, ctx.dbDir())
	ctx.actionRouter = action.NewRouter("action")
	ctx.internalRouter = action.NewRouter("internal")
	ctx.extStores = data.NewStorageRouter()

	testEnv := os.Getenv("OLTEST")

	btime := 600 * time.Second
	ttime := 30 * time.Second
	oltime := 3 * time.Second
	if testEnv == "1" {
		btime = 30 * time.Second
		ttime = 3 * time.Second
	}
	ctx.jobBus = event.NewJobBus(event.Option{
		BtcInterval: btime,
		EthInterval: ttime,
		OltInterval: oltime,
	}, ctx.jobStore)

	_ = transfer.EnableSend(ctx.actionRouter)
	_ = action_ons.EnableONS(ctx.actionRouter)

	//"btc" service temporarily disabled
	//_ = btc.EnableBTC(ctx.actionRouter)

	_ = eth.EnableETH(ctx.actionRouter)
	_ = eth.EnableInternalETH(ctx.internalRouter)

	_ = action_gov.EnableGovernance(ctx.actionRouter)
	_ = action_gov.EnableInternalGovernance(ctx.internalRouter)
	_ = staking.EnableStaking(ctx.actionRouter)

	return ctx, nil
}

func NewProposalMasterStore(chainstate *storage.ChainState) *governance.ProposalMasterStore {
	proposals := governance.NewProposalStore("propActive", "propPassed", "propFailed", storage.NewState(chainstate))
	proposalFunds := governance.NewProposalFundStore("propFunds", storage.NewState(chainstate))
	proposalVotes := governance.NewProposalVoteStore("propVotes", storage.NewState(chainstate))
	return governance.NewProposalMasterStore(proposals, proposalFunds, proposalVotes)
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
		ctx.witnesses.WithState(state),
		ctx.domains.WithState(state),
		ctx.govern.WithState(state),
		ctx.delegators.WithState(state),

		ctx.btcTrackers.WithState(state),
		ctx.ethTrackers.WithState(state),
		ctx.jobStore,
		ctx.lockScriptStore,
		log.NewLoggerWithPrefix(ctx.logWriter, "action").WithLevel(log.Level(ctx.cfg.Node.LogLevel)),
		ctx.proposalMaster.WithState(state),
		ctx.extStores,
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
		ctx.delegators.WithState(ctx.deliver),
		ctx.govern.WithState(ctx.deliver),
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

	feePool := fees.NewStore("f", storage.NewState(ctx.chainstate))
	feePool.SetupOpt(ctx.feePool.GetOpt())

	ethTracker := ethereum.NewTrackerStore("etht", "ethfailed", "ethsuccess", storage.NewState(ctx.chainstate))
	ethTracker.SetupOption(ctx.ethTrackers.GetOption())

	ons := ons.NewDomainStore("d", storage.NewState(ctx.chainstate))
	ons.SetOptions(ctx.domains.GetOptions())

	proposalMaster := NewProposalMasterStore(ctx.chainstate)
	proposalMaster.Proposal.SetOptions(ctx.proposalMaster.Proposal.GetOptions())

	svcCtx := &service.Context{
		Balances:       balance.NewStore("b", storage.NewState(ctx.chainstate)),
		Accounts:       ctx.accounts,
		Currencies:     ctx.currencies,
		FeePool:        feePool,
		Cfg:            ctx.cfg,
		NodeContext:    ctx.node,
		ValidatorSet:   identity.NewValidatorStore("v", storage.NewState(ctx.chainstate)),
		WitnessSet:     identity.NewWitnessStore("w", storage.NewState(ctx.chainstate)),
		Domains:        ons,
		Govern:         ctx.govern,
		Delegators:     ctx.delegators,
		ProposalMaster: proposalMaster,
		ExtStores:      ctx.extStores,
		Router:         ctx.actionRouter,
		Logger:         log.NewLoggerWithPrefix(ctx.logWriter, "rpc").WithLevel(log.Level(ctx.cfg.Node.LogLevel)),
		Services:       extSvcs,
		EthTrackers:    ethTracker,
		Trackers:       btcTrackers,
	}

	return service.NewMap(svcCtx)
}

func (ctx *context) Restful() (service.RestfulRouter, error) {
	extSvcs, err := client.NewExtServiceContext(ctx.cfg.Network.RPCAddress, ctx.cfg.Network.SDKAddress)
	if err != nil {
		return nil, errors.Wrap(err, "failed to start service context")
	}
	svcCtx := &service.Context{
		Cfg:            ctx.cfg,
		Balances:       ctx.balances,
		Accounts:       ctx.accounts,
		Currencies:     ctx.currencies,
		FeePool:        ctx.feePool,
		NodeContext:    ctx.node,
		ValidatorSet:   ctx.validators,
		Domains:        ctx.domains,
		ProposalMaster: ctx.proposalMaster,
		Router:         ctx.actionRouter,
		Logger:         log.NewLoggerWithPrefix(ctx.logWriter, "restful").WithLevel(log.Level(ctx.cfg.Node.LogLevel)),
		Services:       extSvcs,

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
	Trackers   *ethereum.TrackerStore //TODO: Create struct to contain all tracker types including Bitcoin.

	Currencies *balance.CurrencySet
	FeeOption  *fees.FeeOption
	Hash       []byte
	Version    int64
	Chainstate *storage.ChainState
}

func (ctx *context) Storage() StorageCtx {
	return StorageCtx{
		Version:    ctx.chainstate.Version,
		Hash:       ctx.chainstate.Hash,
		Chainstate: ctx.chainstate,
		Balances:   ctx.balances,
		Domains:    ctx.domains,
		Validators: ctx.validators,
		FeePool:    ctx.feePool,
		Govern:     ctx.govern,
		Currencies: ctx.currencies,
		FeeOption:  ctx.feePool.GetOpt(),
		Trackers:   ctx.ethTrackers,
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
		ctx.node.ValidatorECDSAPrivateKey(), // BTC private key
		ctx.node.ValidatorECDSAPrivateKey(), // ETH private key
		ctx.node.ValidatorAddress(),         // validator address generated from validator key
		ctx.lockScriptStore,
		ctx.ethTrackers.WithState(ctx.deliver),
		ctx.proposalMaster.WithState(ctx.deliver),
		log.NewLoggerWithPrefix(ctx.logWriter, "internal_jobs").WithLevel(log.Level(ctx.cfg.Node.LogLevel)))
}

func (ctx *context) AddExternalTx(t action.Type, h action.Tx) error {
	return ctx.actionRouter.AddHandler(t, h)
}

func (ctx *context) AddExternalStore(storeType data.Type, storeObj interface{}) error {
	return ctx.extStores.Add(storeType, storeObj)
}

func (ctx *context) GetChainState() *storage.ChainState {
	return ctx.chainstate
}
