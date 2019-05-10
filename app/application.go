package app

import (
	"encoding/hex"
	"io"
	"net"
	"net/http"
	"net/rpc"
	"net/url"
	"os"
	"path/filepath"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/identity"

	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/common"

	"github.com/Oneledger/protocol/config"
	"github.com/Oneledger/protocol/consensus"
	"github.com/Oneledger/protocol/data/accounts"
	"github.com/Oneledger/protocol/data/balance"
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

	go app.startRPCServer()
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

// The base context for the application, holds databases and other stateful information contained by the app.
// Used to derive other package-level Contexts
type context struct {
	chainID string
	cfg     config.Server

	// identities       data.Store
	// smartContract    data.Store
	// executionContext data.Store
	// admin            data.Store

	// sequence         data.Store
	// status           data.Store
	// contract         data.Store
	// event            data.Store
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
		chainID:    cfg.ChainID(),
		logWriter:  logWriter,
		currencies: make(map[string]balance.Currency),
	}

	ctx.validators = identity.NewValidators()

	ctx.actionRouter = action.NewRouter("action")

	closers := make([]closer, 0)

	ctx.balances = balance.NewStore("balances", ctx.dbDir(), ctx.cfg.Node.DB, storage.PERSISTENT)
	closers = append(closers, ctx.balances)

	ctx.accounts = accounts.NewWallet(cfg, ctx.dbDir())
	closers = append(closers, ctx.accounts)

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

// Close all things that need to be closed
func (ctx *context) Close() {
	closers := []closer{ctx.balances, ctx.accounts}
	for _, closer := range closers {
		closer.Close()
	}
}

func (app *App) startRPCServer() {
	handlers := NewClientHandler(app.Context.cfg.Node.NodeName, app.Context.balances, app.Context.accounts, app.Context.currencies)

	err := rpc.Register(handlers)
	if err != nil {
		app.logger.Fatal("error registering rpc handlers", "err", err)
	}

	u, err := url.Parse(app.Context.cfg.Network.SDKAddress)
	if err != nil {
		app.logger.Error("Failed to parse sdk address")
	}
	rpc.HandleHTTP()

	// TODO: Thunk this function for better error handling
	// Change: startRPCServer: func()
	// To:     startRPCServer: func() -> (func(), err)
	l, e := net.Listen("tcp", u.Host)
	if e != nil {
		app.Close()
		app.logger.Fatal("listen error:", e)
	}

	err = http.Serve(l, nil)
	if err != nil {
		app.logger.Fatal("error while starting the RPC server", "err", err)
	}
}
