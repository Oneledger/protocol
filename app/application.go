package app

import (
	"encoding/hex"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	action_ons "github.com/Oneledger/protocol/action/ons"
	"github.com/Oneledger/protocol/action/staking"
	"github.com/Oneledger/protocol/action/transfer"
	"github.com/Oneledger/protocol/data/ons"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/app/node"
	"github.com/Oneledger/protocol/client"
	"github.com/Oneledger/protocol/config"
	"github.com/Oneledger/protocol/consensus"
	"github.com/Oneledger/protocol/data/accounts"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/identity"
	"github.com/Oneledger/protocol/log"
	"github.com/Oneledger/protocol/rpc"
	"github.com/Oneledger/protocol/serialize"
	"github.com/Oneledger/protocol/service"
	"github.com/Oneledger/protocol/storage"
	"github.com/tendermint/tendermint/abci/types"

	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/common"
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
	abci   *ABCI

	node *consensus.Node
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
		logger: log.NewLoggerWithPrefix(w, "app"),
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

	// commit the initial currencies to the admin db
	session := app.Context.admin.BeginSession()
	_ = session.Set([]byte(ADMIN_CURRENCY_KEY), stateBytes)
	session.Commit()

	nodeCtx := app.Context.Node()
	balanceCtx := app.Context.Balances()
	walletCtx := app.Context.Accounts()

	// (1) Register all the currencies
	for _, currency := range initial.Currencies {
		err := balanceCtx.Currencies().Register(currency)
		if err != nil {
			return errors.Wrapf(err, "failed to register currency %s", currency.Name)
		}
	}

	// (2) Set balances to all those mentioned
	for _, state := range initial.States {
		si := state.StateInput()
		addrBytes, err := hex.DecodeString(si.Address)
		if err != nil {
			return errors.Wrapf(err, "failed to decode address %s", si.Address)
		}

		key := storage.StoreKey(addrBytes)
		err = balanceCtx.Store().Set([]byte(key), si.Balance)
		if err != nil {
			return errors.Wrap(err, "failed to set balance")
		}

		app.logger.Debug(strings.ToUpper(hex.EncodeToString(key)))
	}

	myPrivKey := nodeCtx.PrivKey()
	myPubKey := nodeCtx.PubKey()

	// Start registering myself
	app.logger.Info("Registering myself...")
	for _, currency := range initial.Currencies {
		chainType := currency.Chain
		acct, err := accounts.NewAccount(
			chainType,
			nodeCtx.NodeName,
			&myPrivKey,
			&myPubKey)

		if err != nil {
			app.logger.Warn("Can't create a new account for myself", "err", err, "chainType", chainType)
			continue
		}

		if _, err := walletCtx.GetAccount(acct.Address()); err != nil {
			err = walletCtx.Add(acct)
			if err != nil {
				app.logger.Warn("Failed to register myself", "err", err)
				continue
			}
		}
		app.logger.Info("Successfully registered myself!")
	}
	return nil
}

func (app *App) setupValidators(req RequestInitChain, currencies *balance.CurrencyList) (types.ValidatorUpdates, error) {
	return app.Context.validators.Init(req, currencies)
}

// Start initializes the state
func (app *App) Start() error {
	app.logger.Info("Starting node...")

	//get currencies from admin db
	result, err := app.Context.admin.Get([]byte(ADMIN_CURRENCY_KEY))
	if err != nil {
		app.logger.Debug("didn't get the currencies from db")
	} else {

		as := &consensus.AppState{}
		err = serialize.GetSerializer(serialize.PERSISTENT).Deserialize(result, as)
		if err != nil {
			app.logger.Error("failed to deserialize the currencies from db")
			return errors.Wrap(err, "failed to get the currencies")
		}

		for _, currency := range as.Currencies {
			err := app.Context.currencies.Register(currency)
			if err != nil {
				return errors.Wrapf(err, "failed to register currency %s", currency.Name)
			}
		}

		app.logger.Infof("Read currencies from db %#v", app.Context.currencies)
	}

	node, err := consensus.NewNode(app.ABCI(), &app.Context.cfg)
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

	app.node = node
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

	err = app.Context.rpc.Prepare(u)
	if err != nil {
		return noop, err
	}

	srv := app.Context.rpc

	return srv.Start, nil
}

// The base context for the application, holds databases and other stateful information contained by the app.
// Used to derive other package-level Contexts
type context struct {
	node node.Context
	cfg  config.Server

	rpc          *rpc.Server
	actionRouter action.Router

	balances *balance.Store

	domains *ons.DomainStore

	validators *identity.ValidatorStore // Set of validators currently active
	accounts   accounts.Wallet
	currencies *balance.CurrencyList
	admin      storage.SessionedStorage

	logWriter io.Writer
}

type closer interface {
	Close()
}

func newContext(logWriter io.Writer, cfg config.Server, nodeCtx *node.Context) (context, error) {
	ctx := context{
		cfg:        cfg,
		logWriter:  logWriter,
		currencies: balance.NewCurrencyList(),
		node:       *nodeCtx,
	}

	ctx.rpc = rpc.NewServer(logWriter)
	ctx.validators = identity.NewValidatorStore(cfg, ctx.dbDir(), ctx.cfg.Node.DB)
	ctx.actionRouter = action.NewRouter("action")
	ctx.balances = balance.NewStore("balances", ctx.dbDir(), ctx.cfg.Node.DB, storage.PERSISTENT)
	ctx.accounts = accounts.NewWallet(cfg, ctx.dbDir())
	ctx.domains = ons.NewDomainStore("domains", ctx.dbDir(), ctx.cfg.Node.DB, storage.PERSISTENT)
	ctx.admin = storage.NewStorageDB(storage.KEYVALUE, "admin", ctx.dbDir(), ctx.cfg.Node.DB)

	_ = transfer.EnableSend(ctx.actionRouter)
	_ = staking.EnableApplyValidator(ctx.actionRouter)
	_ = action_ons.EnableONS(ctx.actionRouter)
	return ctx, nil
}

func (ctx context) dbDir() string {
	return filepath.Join(ctx.cfg.RootDir(), ctx.cfg.Node.DBDir)
}

func (ctx *context) Action(header *Header) *action.Context {
	actionCtx := action.NewContext(
		ctx.actionRouter,
		header,
		ctx.accounts,
		ctx.balances,
		ctx.currencies,
		ctx.validators,
		ctx.domains,
		log.NewLoggerWithPrefix(ctx.logWriter, "action"))

	return actionCtx
}

func (ctx *context) ID() {}
func (ctx *context) Accounts() accounts.Wallet {
	return ctx.accounts
}

func (ctx *context) ValidatorCtx() *identity.ValidatorContext {
	return identity.NewValidatorContext(ctx.balances)
}

// Returns a balance.Context
func (ctx *context) Balances() *balance.Context {
	return balance.NewContext(
		log.NewLoggerWithPrefix(ctx.logWriter, "balances"),
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
		Cfg:          ctx.cfg,
		NodeContext:  ctx.node,
		ValidatorSet: ctx.validators,
		Domains:      ctx.domains,
		Router:       ctx.actionRouter,
		Logger:       log.NewLoggerWithPrefix(ctx.logWriter, "rpc"),
		Services:     extSvcs,
	}

	return service.NewMap(svcCtx), nil
}

// Close all things that need to be closed
func (ctx *context) Close() {
	closers := []closer{ctx.balances, ctx.accounts, ctx.rpc}
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
