package app

import (
	"net/url"
	"os"
	"sync"

	"github.com/Oneledger/protocol/data/network_delegation"

	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/service"
	"github.com/tendermint/tendermint/store"

	"github.com/Oneledger/protocol/app/node"
	"github.com/Oneledger/protocol/config"
	"github.com/Oneledger/protocol/consensus"
	"github.com/Oneledger/protocol/data/accounts"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/bitcoin"
	"github.com/Oneledger/protocol/data/chain"
	"github.com/Oneledger/protocol/data/ethereum"
	"github.com/Oneledger/protocol/data/ons"
	"github.com/Oneledger/protocol/event"
	"github.com/Oneledger/protocol/identity"
	"github.com/Oneledger/protocol/log"
	"github.com/Oneledger/protocol/serialize"
	"github.com/Oneledger/protocol/storage"
)

// Ensure this App struct can control the underlying ABCI app
var _ abciController = &App{}

type App struct {
	Context context

	name     string
	nodeName string
	logger   *log.Logger
	sdk      service.Service // Probably needs to be changed

	header Header // Tendermint last header info

	abci *ABCI

	node       *consensus.Node
	genesisDoc *config.GenesisDoc
}

// New returns new app fresh and ready to start
func NewApp(cfg *config.Server, nodeContext *node.Context) (*App, error) {
	if cfg == nil || nodeContext == nil {
		return nil, errors.New("got nil argument")
	}

	// TODO: Determine the final logWriter in the configuration file
	w := os.Stdout

	app := &App{
		name:   "OneLedger",
		logger: log.NewLoggerWithPrefix(w, "app").WithLevel(log.Level(cfg.Node.LogLevel)),
	}
	app.nodeName = cfg.Node.NodeName

	ctx, err := newContext(w, *cfg, nodeContext)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create new app context")
	}

	app.Context = ctx
	app.setNewABCI()
	return app, nil
}

// ABCI returns an ABCI-ready Application used to initialize the new Node
func (app *App) ABCI() *ABCI {
	return app.abci
}

// Header returns this node's header
func (app *App) Header() Header {
	return app.header
}

// Node returns the consensus.Node, use this value to communicate with the internal consensus engine
func (app *App) Node() *consensus.Node {
	return app.node
}

// setNewABCI returns a new ABCI struct with the current context-values set in App
func (app *App) setNewABCI() {
	app.abci = &ABCI{
		infoServer:       app.infoServer(),
		optionSetter:     app.optionSetter(),
		queryer:          app.queryer(),
		txChecker:        app.txChecker(),
		chainInitializer: app.chainInitializer(),
		blockBeginner:    app.blockBeginner(),
		txDeliverer:      app.txDeliverer(),
		blockEnder:       app.blockEnder(),
		commitor:         app.commitor(),
	}
}

// setupState reads the AppState portion of the genesis file and uses that to set the app to its initial state
func (app *App) setupState(stateBytes []byte) error {
	app.logger.Info("Setting up state...")
	var initial consensus.AppState
	// Deserialize and get the proper app state
	err := serialize.GetSerializer(serialize.JSON).Deserialize(stateBytes, &initial)
	if err != nil {
		return errors.Wrap(err, "setupState deserialization")
	}

	err = app.Context.govern.WithHeight(app.header.Height).SetStakingOptions(initial.Governance.StakingOptions)
	if err != nil {
		return errors.Wrap(err, "Setup Staking Options")
	}
	err = app.Context.govern.WithHeight(app.header.Height).SetEvidenceOptions(initial.Governance.EvidenceOptions)
	if err != nil {
		return errors.Wrap(err, "Setup Evidence Options")
	}
	// commit the initial currencies to the governance db
	err = app.Context.govern.WithHeight(app.header.Height).SetCurrencies(initial.Currencies)
	if err != nil {
		return errors.Wrap(err, "Setup State")
	}

	err = app.Context.govern.WithHeight(app.header.Height).SetProposalOptions(initial.Governance.PropOptions)
	if err != nil {
		return errors.Wrap(err, "Setup Proposal Options")
	}
	app.Context.proposalMaster.Proposal.SetOptions(&initial.Governance.PropOptions)

	err = app.Context.govern.WithHeight(app.header.Height).SetETHChainDriverOption(initial.Governance.ETHCDOption)
	if err != nil {
		return errors.Wrap(err, "Setup Ethereum Options")
	}
	app.Context.ethTrackers.SetupOption(&initial.Governance.ETHCDOption)

	err = app.Context.govern.WithHeight(app.header.Height).SetBTCChainDriverOption(initial.Governance.BTCCDOption)
	if err != nil {
		return errors.Wrap(err, "Setup BTC Options")
	}
	balanceCtx := app.Context.Balances()

	app.Context.btcTrackers.SetConfig(bitcoin.NewBTCConfig(app.Context.cfg.ChainDriver, initial.Governance.BTCCDOption.ChainType))
	app.Context.btcTrackers.SetOption(initial.Governance.BTCCDOption)

	err = app.Context.govern.WithHeight(app.header.Height).SetONSOptions(initial.Governance.ONSOptions)
	if err != nil {
		return errors.Wrap(err, "Error in setting up ONS options")
	}
	app.Context.domains.SetOptions(&initial.Governance.ONSOptions)

	err = app.Context.govern.WithHeight(app.header.Height).SetRewardOptions(initial.Governance.RewardOptions)
	if err != nil {
		return errors.Wrap(err, "Error in setting up Reward options")
	}
	app.Context.rewardMaster.SetOptions(&initial.Governance.RewardOptions)

	// (1) Register all the currencies and fee
	for _, currency := range initial.Currencies {
		err := balanceCtx.Currencies().Register(currency)
		if err != nil {
			return errors.Wrapf(err, "failed to register currency %s", currency.Name)
		}
	}

	err = app.Context.govern.WithHeight(app.header.Height).SetFeeOption(initial.Governance.FeeOption)
	if err != nil {
		return errors.Wrap(err, "Setup FeeOptions Options")
	}
	app.Context.feePool.SetupOpt(&initial.Governance.FeeOption)

	//TODO change back to genesis in future, right now network delegation option is hardcoded to avoid genesis deployment
	hardCodedOption := network_delegation.Options{
		RewardsMaturityTime: network_delegation.RewardsMaturityTime,
	}
	err = app.Context.govern.WithHeight(app.header.Height).SetNetworkDelegOptions(hardCodedOption)

	//err = app.Context.govern.WithHeight(app.header.Height).SetNetworkDelegOptions(initial.Governance.DelegOptions)
	if err != nil {
		return errors.Wrap(err, "Setup NetworkDelegation Options")
	}

	err = app.Context.govern.WithHeight(app.header.Height).SetAllLUH()
	if err != nil {
		return errors.Wrap(err, "Unable to set last Update height ")
	}
	// (2) Set balances to all those mentioned
	for _, bal := range initial.Balances {
		key := storage.StoreKey(bal.Address)
		c, ok := balanceCtx.Currencies().GetCurrencyByName(bal.Currency)
		if !ok {
			return errors.New("currency for initial balance not support")
		}
		coin := c.NewCoinFromAmount(bal.Amount)
		err = balanceCtx.Store().WithState(app.Context.deliver).AddToAddress([]byte(key), coin)
		if err != nil {
			return errors.Wrap(err, "failed to set balance")
		}
	}
	for _, stake := range initial.Witness {
		err = app.Context.witnesses.WithState(app.Context.deliver).AddWitness(chain.ETHEREUM, identity.Stake(stake))
		if err != nil {
			return errors.Wrap(err, "failed to add initial ethereum witness")
		}
	}
	for _, stake := range initial.Staking {
		err := app.Context.delegators.WithState(app.Context.deliver).Stake(stake.ValidatorAddress, stake.StakeAddress, identity.Stake(stake).Amount)
		if err != nil {
			return errors.Wrap(err, "failed to handle delegators staking")
		}
		err = app.Context.validators.WithState(app.Context.deliver).HandleStake(identity.Stake(stake), false, 0)
		if err != nil {
			return errors.Wrap(err, "failed to handle initial staking")
		}
	}

	if !app.Context.delegators.WithState(app.Context.deliver).LoadState(initial.Delegation) {
		return errors.Wrap(err, "failed to setup initial delegation")
	}

	if !app.Context.rewardMaster.WithState(app.Context.deliver).LoadState(initial.Rewards) {
		return errors.Wrap(err, "failed to setup initial rewards")
	}

	for _, domain := range initial.Domains {
		if ons.GetNameFromString(domain.Name).IsValid() && app.Context.domains.GetOptions().IsNameAllowed(ons.Name(domain.Name)) {
			d, err := ons.NewDomain(domain.Owner, domain.Beneficiary, domain.Name, 0, domain.URI, domain.ExpireHeight, domain.ActiveFlag)
			if err != nil {
				return errors.Wrap(err, "failed to create initial domain")
			}
			err = app.Context.domains.WithState(app.Context.deliver).Set(d)
			if err != nil {
				return errors.Wrap(err, "failed to setup initial domain")
			}
		}
	}

	for _, fee := range initial.Fees {
		c, ok := app.Context.currencies.GetCurrencyByName(fee.Currency)
		if !ok {
			return errors.New("currency for initial balance not support")
		}
		err := app.Context.feePool.WithState(app.Context.deliver).Set(fee.Address, c.NewCoinFromAmount(fee.Amount))
		if err != nil {
			return errors.Wrap(err, "failed to setup initial fee")
		}
	}
	//TODO: Initialize BTC Trackers in the future.
	for _, tracker := range initial.Trackers {
		if tracker.State == ethereum.BusyBroadcasting {
			tracker.State = ethereum.New
		}
		tr := &ethereum.Tracker{
			Type:          tracker.Type,
			State:         tracker.State,
			TrackerName:   tracker.TrackerName,
			SignedETHTx:   tracker.SignedETHTx,
			Witnesses:     tracker.Witnesses,
			ProcessOwner:  tracker.ProcessOwner,
			FinalityVotes: make([]ethereum.Vote, len(tracker.Witnesses)),
			To:            tracker.To,
		}
		switch tracker.State {
		case ethereum.Released:
			err = app.Context.ethTrackers.WithState(app.Context.deliver).WithPrefixType(ethereum.PrefixPassed).Set(tr)
		case ethereum.Failed:
			err = app.Context.ethTrackers.WithState(app.Context.deliver).WithPrefixType(ethereum.PrefixFailed).Set(tr)
		default:
			err = app.Context.ethTrackers.WithState(app.Context.deliver).WithPrefixType(ethereum.PrefixOngoing).Set(tr)
		}

		if err != nil {
			return errors.Wrap(err, "failed to setup initial Trackers")
		}
	}

	//Setup Proposals
	err = app.Context.proposalMaster.WithState(app.Context.deliver).LoadProposals(initial.Proposals)
	if err != nil {
		return errors.Wrap(err, "error setup proposal data")
	}

	//Setup Delegators
	err = app.Context.netwkDelegators.Deleg.WithState(app.Context.deliver).LoadDelegators(initial.NetDelegators)
	if err != nil {
		return errors.Wrap(err, "error setup network delegation data")
	}

	//Setup Delegators
	err = app.Context.netwkDelegators.Rewards.WithState(app.Context.deliver).LoadState(&initial.DelegatorRew)
	if err != nil {
		return errors.Wrap(err, "error setting up network delegation reward data")
	}

	app.Context.deliver.Write()
	return nil
}

func (app *App) setupValidators(req RequestInitChain, currencies *balance.CurrencySet) (types.ValidatorUpdates, error) {

	vu, err := app.Context.validators.WithState(app.Context.deliver).Init(req, currencies)

	//btcCfg := app.Context.btcTrackers.GetConfig()

	//vals, err := app.Context.validators.WithState(app.Context.deliver).GetBitcoinKeys(btcCfg.BTCParams)
	//threshold := (len(vals) * 2 / 3) + 1
	//for i := 0; i < 6; i++ {
	//	// appHash := app.genesisDoc.AppHash.Bytes()
	//
	//	randBytes := []byte("XOLT")
	//
	//	script, address, addressList, err := bitcoin2.CreateMultiSigAddress(threshold, vals, randBytes, btcCfg.BTCParams)
	//	if err != nil {
	//		return nil, err
	//	}
	//
	//	signers := make([]keys.Address, len(addressList))
	//	for i := range addressList {
	//		addr := base58.Decode(addressList[i])
	//		signers[i] = keys.Address(addr)
	//	}
	//
	//	tracker, err := bitcoin.NewTracker(address, threshold, signers)
	//	if err != nil {
	//		return nil, err
	//	}
	//
	//	name := fmt.Sprintf("tracker_%d", i)
	//	err = app.Context.btcTrackers.SetTracker(name, tracker)
	//	if err != nil {
	//		return nil, err
	//	}
	//
	//	err = app.Context.lockScriptStore.SaveLockScript(address, script)
	//	if err != nil {
	//		return nil, err
	//	}
	//}

	return vu, err
}

func (app *App) Prepare() error {
	testEnv := os.Getenv("OLTEST")

	//Register address for current node if in test environment.
	if app.Context.govern.InitialChain() && testEnv == "1" {
		app.logger.Debug("didn't get the currencies from db,  register self")
		nodeCtx := app.Context.Node()
		walletCtx := app.Context.Accounts()
		myPrivKey := nodeCtx.PrivKey()
		myPubKey := nodeCtx.PubKey()
		// Start registering myself
		app.logger.Info("Registering myself...")

		chainType := chain.Type(0)
		acct, err := accounts.NewAccount(
			chainType,
			nodeCtx.NodeName,
			&myPrivKey,
			&myPubKey)

		if err != nil {
			app.logger.Warn("Can't create a new account for myself", "err", err, "chainType", chainType)
		}

		if _, err := walletCtx.GetAccount(acct.Address()); err != nil {
			err = walletCtx.Add(acct)
			if err != nil {
				app.logger.Warn("Failed to register myself", "err", err)
			}
		}
		app.logger.Infof("Successfully registered myself: %s", acct.Address())
	}

	//get currencies from governance db
	if !app.Context.govern.InitialChain() {
		//TODO remove setting individual stores after all TX's directly use the Gov store
		currencies, err := app.Context.govern.WithHeight(app.header.Height).GetCurrencies()
		if err != nil {
			return err
		}
		for _, currency := range currencies {
			err := app.Context.currencies.Register(currency)
			if err != nil {
				return errors.Wrapf(err, "failed to register currency %s", currency.Name)
			}
		}

		app.logger.Infof("Read currencies from db %#v", currencies)

		feeOpt, err := app.Context.govern.WithHeight(app.header.Height).GetFeeOption()
		if err != nil {
			return err
		}

		app.Context.feePool.SetupOpt(feeOpt)

		cdOpt, err := app.Context.govern.WithHeight(app.header.Height).GetETHChainDriverOption()
		if err != nil {
			return err
		}
		app.Context.ethTrackers.SetupOption(cdOpt)

		btcOption, err := app.Context.govern.WithHeight(app.header.Height).GetBTCChainDriverOption()
		if err != nil {
			return err
		}
		btcConfig := bitcoin.NewBTCConfig(app.Context.cfg.ChainDriver, btcOption.ChainType)

		app.Context.btcTrackers.SetConfig(btcConfig)

		propOpt, err := app.Context.govern.WithHeight(app.header.Height).GetProposalOptions()
		if err != nil {
			return err
		}
		app.Context.proposalMaster.Proposal.SetOptions(propOpt)
		rewardsOpt, err := app.Context.govern.GetRewardOptions()
		if err != nil {
			return err
		}
		app.Context.rewardMaster.SetOptions(rewardsOpt)
	}

	nodecfg, err := consensus.ParseConfig(&app.Context.cfg)
	if err != nil {
		return errors.Wrap(err, "failed parse NodeConfig")
	}
	genesisDoc, err := nodecfg.GetGenesisDoc()
	if err != nil {
		return errors.Wrap(err, "failed get genesisDoc")
	}
	app.genesisDoc = genesisDoc

	blockStoreChan := make(chan *store.BlockStore)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		app.node, err = consensus.NewNode(app.ABCI(), nodecfg, blockStoreChan)
	}()

	// Init witness store after genesis witnesses loaded in above NewNode
	app.Context.witnesses.Init(chain.ETHEREUM, app.Context.node.ValidatorAddress())

	// set up dirty block store when ready to prevent issues during node sync
	app.logger.Debug("awaiting to set up dirty block store for context, rewards and state db")
	app.Context.SetBlockStore(<-blockStoreChan)

	app.logger.Debug("waiting node preparation...")
	wg.Wait()
	app.logger.Debug("node loaded")
	if err != nil {
		app.logger.Error("Failed to create consensus.Node")
		return errors.Wrap(err, "failed to create new consensus.Node")
	}

	// set up loaded block store for context and state db
	app.logger.Debug("set up loaded block store for context, rewards and state db")
	app.Context.SetBlockStore(app.node.BlockStore())

	// Initialize internal Services
	app.Context.internalService = event.NewService(app.Context.node,
		log.NewLoggerWithPrefix(app.Context.logWriter, "internal_service"), app.Context.internalRouter, app.node)

	return nil
}

// Start initializes the state
func (app *App) Start() error {

	err := app.Prepare()
	if err != nil {
		return err
	}

	// Starting App
	err = app.node.Start()
	if err != nil {
		app.logger.Error("Failed to start consensus.Node")
		return errors.Wrap(err, "failed to start new consensus.Node")
	}
	//Start Jobbus
	_ = app.Context.jobBus.Start(app.Context.JobContext())
	// Starting Legacy RPC
	startRPC, err := app.rpcStarter()
	if err != nil {
		return errors.Wrap(err, "failed to prepare legacy rpc service")
	}

	err = startRPC()
	if err != nil {
		app.logger.Error("Failed to start legacy rpc")
		return err
	}

	web3Services := app.Context.Web3Services(app.node)

	// Starting new (web3) RPC
	err = app.Context.web3.StartHTTP(web3Services)
	if err != nil {
		app.logger.Error("Failed to start web3 http rpc, details", err)
		return err
	}

	err = app.Context.web3.StartWS(web3Services)
	if err != nil {
		app.logger.Error("Failed to start web3 ws rpc, details", err)
		return err
	}

	//"btc" service temporarily disabled
	//err = btc.EnableBTCInternalTx(internalRouter)
	//if err != nil {
	//	app.logger.Error("Failed to register btc internal transactions")
	//	return err
	//}

	return nil
}

// Close closes the application
func (app *App) Close() {
	app.logger.Info("Closing App...")
	if app.node == nil {
		app.logger.Info("node is nil!")
	} else {
		app.node.OnStop()
	}
	app.Context.Close()
}

func (app *App) rpcStarter() (func() error, error) {
	noop := func() error { return nil }

	u, err := url.Parse(app.Context.cfg.Network.SDKAddress)
	if err != nil {
		return noop, err
	}

	services, err := app.Context.Services()
	if err != nil {
		return noop, err
	}
	for name, svc := range services {
		err := app.Context.rpc.Register(name, svc)
		if err != nil {
			app.logger.Errorf("failed to register service %s", name)
		}
	}

	restfulRouter, err := app.Context.Restful()
	if err != nil {
		return noop, err
	}
	app.Context.rpc.RegisterRestfulMap(restfulRouter)

	err = app.Context.rpc.Prepare(u)
	if err != nil {
		return noop, err
	}

	srv := app.Context.rpc

	return srv.Start, nil
}

type closer interface {
	Close() error
}
