package app

import (
	"encoding/hex"
	"io"
	"net/url"
	"os"
	"path/filepath"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/config"
	"github.com/Oneledger/protocol/consensus"
	"github.com/Oneledger/protocol/data/accounts"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/identity"
	"github.com/Oneledger/protocol/log"
	"github.com/Oneledger/protocol/rpc"
	"github.com/Oneledger/protocol/serialize"
	"github.com/Oneledger/protocol/storage"

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

// New returns new app fresh and ready to start, returns an error if
func NewApp(cfg *config.Server) (*App, error) {
	// TODO: Determine the final logWriter in the configuration file
	w := os.Stdout

	app := &App{
		name:   "OneLedger",
		logger: log.NewLoggerWithPrefix(w, "app"),
	}
	app.nodeName = cfg.Node.NodeName

	ctx, err := newContext(*cfg, w)
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

	balanceCtx := app.Context.Balances()

	// (1) Register all the currencies
	for _, currency := range initial.Currencies {
		err := balanceCtx.RegisterCurrency(currency)
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
		err = balanceCtx.Store().Set(key, si.Balance)
		if err != nil {
			return errors.Wrap(err, "failed to set balance")
		}

	}
	return nil
}

// Start initializes the state
func (app *App) Start() error {
	app.logger.Info("Starting node...")
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

	handlers := NewClientHandler(app.Context.cfg.Node.NodeName, app.Context.balances, app.Context.accounts, app.Context.currencies)

	u, err := url.Parse(app.Context.cfg.Network.SDKAddress)
	if err != nil {
		return noop, err
	}

	err = app.Context.rpc.Prepare(u, handlers)
	if err != nil {
		return noop, err
	}

	srv := app.Context.rpc

	return srv.Start, nil
}

// The base context for the application, holds databases and other stateful information contained by the app.
// Used to derive other package-level Contexts
type context struct {
	cfg config.Server

	rpc          *rpc.Server
	actionRouter action.Router

	balances *balance.Store

	validators *identity.Validators // Set of validators currently active
	accounts   accounts.Wallet
	currencies map[string]balance.Currency

	logWriter io.Writer
}

type closer interface {
	Close()
}

func newContext(cfg config.Server, logWriter io.Writer) (context, error) {
	ctx := context{
		cfg:        cfg,
		logWriter:  logWriter,
		currencies: make(map[string]balance.Currency),
	}

	ctx.rpc = rpc.NewServer(logWriter)
	ctx.validators = identity.NewValidators()
	ctx.actionRouter = action.NewRouter("action")
	ctx.balances = balance.NewStore("balances", ctx.dbDir(), ctx.cfg.Node.DB, storage.PERSISTENT)
	ctx.accounts = accounts.NewWallet(cfg, ctx.dbDir())

	return ctx, nil
}

func (ctx context) dbDir() string {
	return filepath.Join(ctx.cfg.RootDir(), ctx.cfg.Node.DBDir)
}

func (ctx *context) Action() *action.Context {
	return action.NewContext(
		ctx.actionRouter,
		ctx.accounts,
		ctx.balances,
		ctx.currencies,
		log.NewLoggerWithPrefix(ctx.logWriter, "action"))
}

func (ctx *context) ID()       {}
func (ctx *context) Accounts() {}

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

func (ctx *context) RPC() *RPCServerContext {
	return &RPCServerContext{
		nodeName:   ctx.cfg.Node.NodeName,
		balances:   ctx.balances,
		accounts:   ctx.accounts,
		currencies: ctx.currencies,
		logger:     log.NewLoggerWithPrefix(ctx.logWriter, "rpc"),
	}
}

// Close all things that need to be closed
func (ctx *context) Close() {
	closers := []closer{ctx.balances, ctx.accounts, ctx.rpc}
	for _, closer := range closers {
		closer.Close()
	}
}
