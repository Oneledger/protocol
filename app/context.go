package app

import (
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/Oneledger/protocol/data/evm"
	"github.com/Oneledger/protocol/external_apps"
	"github.com/Oneledger/protocol/external_apps/common"
	"github.com/Oneledger/protocol/web3"

	tmdb "github.com/tendermint/tm-db"

	"github.com/Oneledger/protocol/data"
	"github.com/Oneledger/protocol/data/network_delegation"
	netwkDeleg "github.com/Oneledger/protocol/data/network_delegation"
	"github.com/Oneledger/protocol/data/rewards"
	"github.com/Oneledger/protocol/data/transactions"

	"github.com/pkg/errors"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/action/eth"
	action_pen "github.com/Oneledger/protocol/action/evidence"
	action_gov "github.com/Oneledger/protocol/action/governance"
	action_netwkdeleg "github.com/Oneledger/protocol/action/network_delegation"
	action_ons "github.com/Oneledger/protocol/action/ons"
	action_rewards "github.com/Oneledger/protocol/action/rewards"
	action_sc "github.com/Oneledger/protocol/action/smart_contract"
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
	"github.com/Oneledger/protocol/data/evidence"
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
	web3         *web3.Server
	actionRouter action.Router

	//db for chain state storage
	db         tmdb.DB
	chainstate *storage.ChainState
	check      *storage.State
	deliver    *storage.State

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
	netwkDelegators *netwkDeleg.MasterStore
	evidenceStore   *evidence.EvidenceStore
	rewardMaster    *rewards.RewardMasterStore
	transaction     *transactions.TransactionStore
	logWriter       io.Writer
	govupdate       *action.GovernaceUpdateAndValidate
	extApp          *common.ExtAppData
	extStores       data.StorageRouter
	extServiceMap   common.ExtServiceMap
	extFunctions    common.ControllerRouter

	// evm integration
	contracts     *evm.ContractStore
	accountKeeper balance.AccountKeeper
	stateDB       *action.CommitStateDB
}

func newContext(logWriter io.Writer, cfg config.Server, nodeCtx *node.Context) (context, error) {
	ctx := context{
		cfg:        cfg,
		logWriter:  logWriter,
		currencies: balance.NewCurrencySet(),
		node:       *nodeCtx,
	}

	ctx.rpc = rpc.NewServer(logWriter, &cfg)
	// new rpc service
	web3, err := web3.NewServer(logWriter, &cfg)
	if err != nil {
		return ctx, errors.Wrap(err, "web3 api failed")
	}
	ctx.web3 = web3

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

	ctx.validators = identity.NewValidatorStore("v", "purged", storage.NewState(ctx.chainstate))
	ctx.witnesses = identity.NewWitnessStore("w", storage.NewState(ctx.chainstate))
	ctx.balances = balance.NewStore("b", storage.NewState(ctx.chainstate))
	ctx.domains = ons.NewDomainStore("d", storage.NewState(ctx.chainstate))
	ctx.feePool = fees.NewStore("f", storage.NewState(ctx.chainstate))
	ctx.govern = governance.NewStore("g", storage.NewState(ctx.chainstate))
	ctx.proposalMaster = NewProposalMasterStore(ctx.chainstate)
	ctx.delegators = delegation.NewDelegationStore("st", storage.NewState(ctx.chainstate))
	ctx.netwkDelegators = netwkDeleg.NewMasterStore("deleg", "delegRwz", storage.NewState(ctx.chainstate))
	ctx.evidenceStore = evidence.NewEvidenceStore("es", storage.NewState(ctx.chainstate))
	ctx.rewardMaster = NewRewardMasterStore(ctx.chainstate)
	ctx.btcTrackers = bitcoin.NewTrackerStore("btct", storage.NewState(ctx.chainstate))
	//Separate DB and chainstate
	newDB := tmdb.NewDB("internaltxdb", tmdb.MemDBBackend, "")
	cs := storage.NewState(storage.NewChainState("chainstateTX", newDB))
	ctx.transaction = transactions.NewTransactionStore("intx", cs)

	ctx.ethTrackers = ethereum.NewTrackerStore("etht", "ethfailed", "ethsuccess", storage.NewState(ctx.chainstate))
	ctx.accounts = accounts.NewWallet(cfg, ctx.dbDir())

	// TODO check if validator
	// valAddr := ctx.node.ValidatorAddress()

	ctx.jobStore = jobs.NewJobStore(cfg, ctx.dbDir())
	ctx.lockScriptStore = bitcoin.NewLockScriptStore(cfg, ctx.dbDir())
	ctx.actionRouter = action.NewRouter("action")
	ctx.internalRouter = action.NewRouter("internal")
	ctx.extStores = data.NewStorageRouter()
	ctx.extServiceMap = common.NewExtServiceMap()
	ctx.extFunctions = common.NewFunctionRouter()

	// evm
	ctx.contracts = evm.NewContractStore(storage.NewState(ctx.chainstate))
	ctx.accountKeeper = balance.NewNesterAccountKeeper(
		storage.NewState(ctx.chainstate),
		ctx.balances,
		ctx.currencies,
	)
	logger := log.NewLoggerWithPrefix(ctx.logWriter, "stateDB").WithLevel(log.Level(ctx.cfg.Node.LogLevel))
	ctx.stateDB = action.NewCommitStateDB(ctx.contracts, ctx.accountKeeper, logger)

	err = external_apps.RegisterExtApp(ctx.chainstate, ctx.actionRouter, ctx.extStores, ctx.extServiceMap, ctx.extFunctions)
	if err != nil {
		return ctx, errors.Wrap(err, "error in registering external apps")
	}
	ctx.govupdate = action.NewGovUpdate()
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
	_ = action_sc.EnableSmartContract(ctx.actionRouter)
	_ = action_ons.EnableONS(ctx.actionRouter)

	//"btc" service temporarily disabled
	//_ = btc.EnableBTC(ctx.actionRouter)

	_ = eth.EnableETH(ctx.actionRouter)
	_ = eth.EnableInternalETH(ctx.internalRouter)

	_ = action_rewards.EnableRewards(ctx.actionRouter)
	_ = action_netwkdeleg.EnableNetworkDelegation(ctx.actionRouter)
	_ = action_gov.EnableGovernance(ctx.actionRouter)
	_ = action_gov.EnableInternalGovernance(ctx.internalRouter)
	_ = staking.EnableStaking(ctx.actionRouter)
	_ = action_pen.EnablePenalization(ctx.actionRouter)

	return ctx, nil
}

func NewProposalMasterStore(chainstate *storage.ChainState) *governance.ProposalMasterStore {
	proposals := governance.NewProposalStore("propActive", "propPassed", "propFailed", "propFinalized", "propFinalizeFailed", storage.NewState(chainstate))
	proposalFunds := governance.NewProposalFundStore("propFunds", storage.NewState(chainstate))
	proposalVotes := governance.NewProposalVoteStore("propVotes", storage.NewState(chainstate))
	return governance.NewProposalMasterStore(proposals, proposalFunds, proposalVotes)
}

func NewRewardMasterStore(chainstate *storage.ChainState) *rewards.RewardMasterStore {
	reward := rewards.NewRewardStore("rwz", "ri", "rwaddr", storage.NewState(chainstate))
	rewardCumula := rewards.NewRewardCumulativeStore("rwcum", storage.NewState(chainstate))
	return rewards.NewRewardMasterStore(reward, rewardCumula)
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
		ctx.delegators.WithState(state),
		ctx.netwkDelegators.WithState(state),
		ctx.evidenceStore.WithState(state),
		ctx.btcTrackers.WithState(state),
		ctx.ethTrackers.WithState(state),
		ctx.jobStore,
		ctx.lockScriptStore,
		log.NewLoggerWithPrefix(ctx.logWriter, "action").WithLevel(log.Level(ctx.cfg.Node.LogLevel)),
		ctx.proposalMaster.WithState(state),
		ctx.rewardMaster.WithState(state),
		ctx.govern.WithState(state),
		ctx.extStores.WithState(state),
		ctx.govupdate,
		ctx.stateDB.WithState(state),
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
		ctx.evidenceStore.WithState(ctx.deliver),
		ctx.govern.WithState(ctx.deliver),
		ctx.currencies,
		ctx.validators.WithState(ctx.deliver),
	)
}

// Returns a balance.Context
func (ctx *context) Balances() *balance.Context {
	return balance.NewContext(
		log.NewLoggerWithPrefix(ctx.logWriter, "balances").WithLevel(log.Level(ctx.cfg.Node.LogLevel)),
		ctx.balances,
		ctx.currencies)
}

func (ctx *context) Web3Services() (map[string]interface{}, error) {
	extSvcs, err := client.NewExtServiceContext(ctx.cfg.Network.RPCAddress, ctx.cfg.Network.SDKAddress)
	if err != nil {
		return nil, errors.Wrap(err, "failed to start service context")
	}
	web3Ctx := web3.NewContext(
		log.NewLoggerWithPrefix(ctx.logWriter, "web3").WithLevel(log.Level(ctx.cfg.Node.LogLevel)),
		&extSvcs,
		ctx.validators,
		ctx.contracts,
		ctx.accountKeeper,
		ctx.node,
	)
	// registering services
	web3Ctx.DefaultRegisterForAll()
	return web3Ctx.ServiceList(), nil
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

	onsStore := ons.NewDomainStore("d", storage.NewState(ctx.chainstate))

	proposalMaster := NewProposalMasterStore(ctx.chainstate)
	proposalMaster.Proposal.SetOptions(ctx.proposalMaster.Proposal.GetOptions())

	rewardMaster := NewRewardMasterStore(ctx.chainstate)
	rewardMaster.SetOptions(ctx.rewardMaster.GetOptions())

	netwkDelegators := netwkDeleg.NewMasterStore("deleg", "delegRwz", storage.NewState(ctx.chainstate))

	svcCtx := &service.Context{
		Balances:        balance.NewStore("b", storage.NewState(ctx.chainstate)),
		Accounts:        ctx.accounts,
		Currencies:      ctx.currencies,
		FeePool:         feePool,
		Cfg:             ctx.cfg,
		NodeContext:     ctx.node,
		ValidatorSet:    identity.NewValidatorStore("v", "purged", storage.NewState(ctx.chainstate)),
		WitnessSet:      identity.NewWitnessStore("w", storage.NewState(ctx.chainstate)),
		Domains:         onsStore,
		Delegators:      delegation.NewDelegationStore("st", storage.NewState(ctx.chainstate)),
		NetwkDelegators: netwkDelegators,
		ProposalMaster:  proposalMaster,
		EvidenceStore:   evidence.NewEvidenceStore("es", storage.NewState(ctx.chainstate)),
		RewardMaster:    rewardMaster,
		ExtStores:       ctx.extStores,
		ExtServiceMap:   ctx.extServiceMap,
		Router:          ctx.actionRouter,
		Logger:          log.NewLoggerWithPrefix(ctx.logWriter, "rpc").WithLevel(log.Level(ctx.cfg.Node.LogLevel)),
		Services:        extSvcs,
		EthTrackers:     ethTracker,
		Trackers:        btcTrackers,
		Govern:          governance.NewStore("g", storage.NewState(ctx.chainstate)),
		GovUpdate:       ctx.govupdate,
		Contracts:       ctx.contracts,
		AccountKeeper:   ctx.accountKeeper,
		StateDB:         ctx.stateDB,
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

		Trackers:      ctx.btcTrackers,
		Contracts:     ctx.contracts,
		AccountKeeper: ctx.accountKeeper,
		StateDB:       ctx.stateDB,
	}
	return service.NewRestfulService(svcCtx).Router(), nil
}

type StorageCtx struct {
	Balances        *balance.Store
	Domains         *ons.DomainStore
	Validators      *identity.ValidatorStore // Set of validators currently active
	Delegators      *delegation.DelegationStore
	RewardMaster    *rewards.RewardMasterStore
	ProposalMaster  *governance.ProposalMasterStore
	NetwkDelegators *network_delegation.MasterStore
	FeePool         *fees.Store
	Govern          *governance.Store
	Trackers        *ethereum.TrackerStore //TODO: Create struct to contain all tracker types including Bitcoin.

	Currencies *balance.CurrencySet
	FeeOption  *fees.FeeOption
	Hash       []byte
	Version    int64
	Chainstate *storage.ChainState
}

func (ctx *context) Storage() StorageCtx {
	return StorageCtx{
		Version:         ctx.chainstate.Version,
		Hash:            ctx.chainstate.Hash,
		Chainstate:      ctx.chainstate,
		Balances:        ctx.balances,
		Domains:         ctx.domains,
		Validators:      ctx.validators,
		Delegators:      ctx.delegators,
		RewardMaster:    ctx.rewardMaster,
		ProposalMaster:  ctx.proposalMaster,
		NetwkDelegators: ctx.netwkDelegators,
		FeePool:         ctx.feePool,
		Govern:          ctx.govern,
		Currencies:      ctx.currencies,
		FeeOption:       ctx.feePool.GetOpt(),
		Trackers:        ctx.ethTrackers,
	}
}

// Close all things that need to be closed
func (ctx *context) Close() {
	closers := []closer{ctx.db, ctx.accounts, ctx.rpc, ctx.jobBus}
	for _, closer := range closers {
		err := closer.Close()
		if err != nil {
			panic(err)
		}
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
