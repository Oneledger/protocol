package app

import (
	"encoding/hex"
	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/identity"
	"io"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"path/filepath"


	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/common"

	"github.com/Oneledger/protocol/client"
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

	name   string
	nodeName string
	logger *log.Logger
	sdk    common.Service // Probably needs to be changed

	header     Header      // Tendermint last header info
	abci       *ABCI
}

// NewApp returns new app fresh and ready to start, returns an error if
func NewApp(cfg config.Server, rootDir string) (*App, error) {
	app := &App{
		name:   "OneLedger",
		logger: log.NewLoggerWithPrefix(os.Stdout, "app"),
	}
	app.nodeName = cfg.Node.NodeName

	ctx, err := newContext(cfg, rootDir)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create new app context")
	}

	app.Context = ctx
	app.setNewABCI()

	go app.startRPCServer()
	return app, nil
}

func (app *App) Header() Header {
	return app.header
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

// ABCI returns an ABCI-ready Application used to initialize the new Node
func (app *App) ABCI() *ABCI {
	return app.abci
}

// setupState reads the AppState portion of the genesis file and uses that to set the app to its initial state
func (app *App) setupState(stateBytes []byte) error {
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
		err = balanceCtx.Store().Set(key, &si.Balance)
		if err != nil {
			return errors.Wrap(err, "failed to set balance")
		}
	}
	return nil

}

// The base context for the application, holds databases and other stateful information contained by the app.
// Used to derive other package-level Contexts
type context struct {
	chainID string
	cfg     config.Server
	rootDir string

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
	accounts accounts.Wallet

	validators *identity.Validators // Set of validators currently active

	currencies      map[string]balance.Currency

	logWriter io.Writer
}

func newContext(cfg config.Server, rootDir string) (context, error) {
	ctx := context{
		rootDir:   rootDir,
		cfg:       cfg,
		chainID:   cfg.ChainID(),
		logWriter: os.Stdout, // TODO: This should be driven by configuration
	}

	ctx.balances = balance.NewStore("balance", ctx.dbDir(), ctx.cfg.Node.DB, storage.PERSISTENT)

	ctx.accounts = accounts.NewWallet(cfg, cfg.Node.DBDir)
	ctx.validators = identity.NewValidators()

	ctx.actionRouter = action.NewRouter("action")

	return ctx, nil
}

func (ctx context) dbDir() string {
	return filepath.Join(ctx.rootDir, ctx.cfg.Node.DBDir)
}

func (ctx *context) Action() *action.Context  {
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

func (app *App) startRPCServer() {

	handlers := NewClientHandler(app.nodeName, app.Context.balances, app.Context.accounts)
	err := rpc.Register(handlers)
	if err != nil {
		app.logger.Fatal("error registering rpc handlers", "err", err)
	}

	rpc.HandleHTTP()

	l, e := net.Listen("tcp", client.RPC_ADDRESS)
	if e != nil {
		app.logger.Fatal("listen error:", e)
	}

	err =  http.Serve(l, nil)
	if err != nil {
		app.logger.Fatal("error while starting the RPC server", "err", err)
	}
}
