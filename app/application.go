package app

import (
	"fmt"
	"net/url"
	"os"

	"github.com/btcsuite/btcutil/base58"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/common"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/action/btc"
	"github.com/Oneledger/protocol/action/eth"
	"github.com/Oneledger/protocol/app/node"
	bitcoin2 "github.com/Oneledger/protocol/chains/bitcoin"
	"github.com/Oneledger/protocol/config"
	"github.com/Oneledger/protocol/consensus"
	"github.com/Oneledger/protocol/data/accounts"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/bitcoin"
	"github.com/Oneledger/protocol/data/chain"
	"github.com/Oneledger/protocol/data/keys"
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
	sdk      common.Service // Probably needs to be changed

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

	// commit the initial currencies to the governance db

	err = app.Context.govern.SetCurrencies(initial.Currencies)
	if err != nil {
		return errors.Wrap(err, "Setup State")
	}

	err = app.Context.govern.SetETHChainDriverOption(initial.ETHCDOption)
	if err != nil {
		return errors.Wrap(err, "Setup State")
	}

	balanceCtx := app.Context.Balances()

	// (1) Register all the currencies and fee
	for _, currency := range initial.Currencies {
		err := balanceCtx.Currencies().Register(currency)
		if err != nil {
			return errors.Wrapf(err, "failed to register currency %s", currency.Name)
		}
	}
	app.Context.feeOption.FeeCurrency = initial.FeeOption.FeeCurrency
	app.Context.feeOption.MinFeeDecimal = initial.FeeOption.MinFeeDecimal
	app.Context.ethTrackers.SetupOption(&initial.ETHCDOption)
	err = app.Context.govern.SetFeeOption(initial.FeeOption)
	if err != nil {
		return errors.Wrap(err, "Setup State")
	}
	app.Context.feePool.SetupOpt(app.Context.feeOption)

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

	for _, stake := range initial.Staking {
		err := app.Context.validators.WithState(app.Context.deliver).HandleStake(identity.Stake(stake))
		if err != nil {
			return errors.Wrap(err, "failed to handle initial staking")
		}
	}

	for _, domain := range initial.Domains {
		d := ons.NewDomain(domain.OwnerAddress, domain.AccountAddress, domain.Name, 0)
		err := app.Context.domains.WithState(app.Context.deliver).Set(d)
		if err != nil {
			return errors.Wrap(err, "failed to setup initial domain")
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
	app.Context.deliver.Write()
	return nil
}

func (app *App) setupValidators(req RequestInitChain, currencies *balance.CurrencySet) (types.ValidatorUpdates, error) {

	vu, err := app.Context.validators.WithState(app.Context.deliver).Init(req, currencies)

	params := bitcoin2.GetChainParams(app.Context.cfg.ChainDriver.BitcoinChainType)

	vals, err := app.Context.validators.WithState(app.Context.deliver).GetBitcoinKeys(params)
	threshold := (len(vals) * 2 / 3) + 1
	for i := 0; i < 6; i++ {
		// appHash := app.genesisDoc.AppHash.Bytes()

		randBytes := []byte("XOLT")

		script, address, addressList, err := bitcoin2.CreateMultiSigAddress(threshold, vals, randBytes, params)
		if err != nil {
			return nil, err
		}

		signers := make([]keys.Address, len(addressList))
		for i := range addressList {
			addr := base58.Decode(addressList[i])
			signers[i] = keys.Address(addr)
		}

		tracker, err := bitcoin.NewTracker(address, threshold, signers)
		if err != nil {
			return nil, err
		}

		name := fmt.Sprintf("tracker_%d", i)
		err = app.Context.btcTrackers.SetTracker(name, tracker)
		if err != nil {
			return nil, err
		}

		err = app.Context.lockScriptStore.SaveLockScript(address, script)
		if err != nil {
			return nil, err
		}
	}

	return vu, err
}

func (app *App) Prepare() error {
	//get currencies from governance db

	if app.Context.govern.InitialChain() {
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

	} else {
		currencies, err := app.Context.govern.GetCurrencies()
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

		feeOpt, err := app.Context.govern.GetFeeOption()
		if err != nil {
			return err
		}
		app.Context.feeOption = feeOpt
		app.Context.feePool.SetupOpt(feeOpt)

		cdOpt, err := app.Context.govern.GetETHChainDriverOption()
		if err != nil {
			return err
		}
		app.Context.ethTrackers.SetupOption(cdOpt)
	}
	return nil
}

// Start initializes the state
func (app *App) Start() error {
	app.logger.Info("Starting node...")

	err := app.Prepare()
	if err != nil {
		return err
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

	node, err := consensus.NewNode(app.ABCI(), nodecfg)
	if err != nil {
		app.logger.Error("Failed to create consensus.Node")
		return errors.Wrap(err, "failed to create new consensus.Node")
	}
	err = node.Start()
	if err != nil {
		app.logger.Error("Failed to start consensus.Node")
		return errors.Wrap(err, "failed to start new consensus.Node")
	}

	startRPC, err := app.rpcStarter()
	if err != nil {
		return errors.Wrap(err, "failed to prepare rpc service")
	}

	err = startRPC()
	if err != nil {
		app.logger.Error("Failed to start rpc")
		return err
	}

	internalRouter := action.NewRouter("internal")
	err = btc.EnableBTCInternalTx(internalRouter)
	if err != nil {
		app.logger.Error("Failed to register btc internal transactions")
		return err
	}
	err = eth.EnableInternalETH(internalRouter)
	if err != nil {
		app.logger.Error("failed to register eth internal transaction")
	}

	app.node = node
	app.Context.internalService = event.NewService(app.Context.node,
		log.NewLoggerWithPrefix(app.Context.logWriter, "internal_service"), internalRouter, node)

	_ = app.Context.jobBus.Start(app.Context.JobContext())

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
	Close()
}
