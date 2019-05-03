package app

import (
	"net"
	"net/http"
	"net/rpc"
	"os"

	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/common"

	"github.com/Oneledger/protocol/client"
	"github.com/Oneledger/protocol/config"
	"github.com/Oneledger/protocol/consensus"
	"github.com/Oneledger/protocol/data"
	"github.com/Oneledger/protocol/data/accounts"
	"github.com/Oneledger/protocol/log"
	"github.com/Oneledger/protocol/serialize"
	"github.com/Oneledger/protocol/storage"
)

// Ensure this App struct can control the underlying ABCI app
var _ abciController = &App{}

// The base context for the application, holds databases and other stateful information contained by the app.
// Used to derive other package-level Contexts
type context struct {
	chainID string

	balances         *storage.ChainState
	identities       data.Store
	smartContract    data.Store
	executionContext data.Store
	admin            data.Store
	accounts         data.Store
	sequence         data.Store
	status           data.Store
	contract         data.Store
	event            data.Store
	wallet           accounts.WalletStore
}

// TODO: Fill in these context creators, these should return package-specific Contexts
/* Example
func (ctx *context) Action() *action.Context {
	// If context members are private:
	ctx := &action.Context{}
	ctx.SetChainID(ctx.chainID)
	ctx.SetBalances(ctx.balances)
	ctx.SetAccounts(ctx)
	return ctx

	// If members are public
	return &action.Context{
		chainID,
		balances,
		identities,
		accounts,
	}
}
*/
func (ctx *context) Action()  {}
func (ctx *context) ID()      {}
func (ctx *context) Account() {}

type App struct {
	Context context

	logger *log.Logger
	sdk    common.Service // Probably needs to be changed

	header Header // Tendermint last header info
	// ? Should this be in Context?
	validators interface{} // Set of validators currently active
	abci       *ABCI
	chainID    string // Signed with every transaction

	rpcServer *rpc.Server
}

func NewApp(cfg config.Server) *App {
	app := &App{
		// sdk:
		logger: log.NewLoggerWithPrefix(os.Stdout, "app:"),
	}

	app.Context.wallet = accounts.NewWallet(cfg)
	// TODO add other data stores

	return app
}

func (app *App) Header() Header {
	return app.header
}

// Getters
func (app *App) Balances() *storage.ChainState {
	return app.Context.balances
}
func (app *App) Identities() data.Store {
	return app.Context.identities
}
func (app *App) SmartContract() data.Store {
	return app.Context.smartContract
}
func (app *App) ExecutionContext() data.Store {
	return app.Context.executionContext
}
func (app *App) Admin() data.Store {
	return app.Context.admin
}
func (app *App) Accounts() data.Store {
	return app.Context.accounts
}
func (app *App) WalletStore() accounts.WalletStore {
	return app.Context.wallet
}
func (app *App) Sequence() data.Store {
	return app.Context.sequence
}
func (app *App) Status() data.Store {
	return app.Context.status
}
func (app *App) Contract() data.Store {
	return app.Context.contract
}
func (app *App) Event() data.Store {
	return app.Context.event
}

// ChainID returns the chainID of this network & application
func (app *App) ChainID() string {
	return app.Context.chainID
}

// TODO: Add proper types
func (app *App) Validators() interface{} {
	return app.validators
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

// query connection: for querying the application state; only uses Query and Info
func (app *App) infoServer() infoServer {
	return func(info RequestInfo) ResponseInfo {
		return ResponseInfo{}
	}
}

func (app *App) queryer() queryer {
	return func(RequestQuery) ResponseQuery {
		// Do stuff
		return ResponseQuery{}
	}
}

// consensus methods: for executing transactions that have been committed. Message sequence is -for every block
func (app *App) optionSetter() optionSetter {
	return func(RequestSetOption) ResponseSetOption {
		// Do stuff
		return ResponseSetOption{}
	}
}

func (app *App) chainInitializer() chainInitializer {
	return func(req RequestInitChain) ResponseInitChain {
		// Do stuff
		err := app.setupState(req.AppStateBytes)
		if err != nil {
			//
		}
		return ResponseInitChain{}
	}
}

func (app *App) blockBeginner() blockBeginner {
	return func(RequestBeginBlock) ResponseBeginBlock {
		// Do stuff
		return ResponseBeginBlock{}
	}
}

func (app *App) txDeliverer() txDeliverer {
	return func([]byte) ResponseDeliverTx {
		// Do stuff
		return ResponseDeliverTx{}
	}
}

func (app *App) blockEnder() blockEnder {
	return func(RequestEndBlock) ResponseEndBlock {
		// Do stuff
		return ResponseEndBlock{}
	}
}

func (app *App) commitor() commitor {
	return func() ResponseCommit {
		return ResponseCommit{}
	}
}

// mempool connection: for checking if transactions should be relayed before they are committed
func (app *App) txChecker() txChecker {
	return func([]byte) ResponseCheckTx {
		// Do stuff
		return ResponseCheckTx{}
	}
}

func (app *App) setupState(stateBytes []byte) error {
	var initial consensus.AppState

	// 	(1) Return the appropriate errors with stateBytes
	err := serialize.JSONSzr.Deserialize(stateBytes, &initial)
	if err != nil {
		return errors.Wrap(err, "setupState deserialization")
	}

	// 	TODO: (2) Generate the node account locally
	// id.GenerateKeys([]byte(global.Current.PaymentAccount), false)

	// 	TODO: (3) Generate the Zero account (?)
	// createAccount(app, &consensus.AppState{global.Current.PaymentAccount, states}, publicKey, privateKey, nil)
	return nil

}

func (app *App) startRPCServer() {

	handlers := client.NewClientHandler(app.Balances(), app.Accounts(), app.WalletStore())
	err := rpc.Register(handlers)
	if err != nil {
		app.logger.Fatal("error registering rpc handlers", "err", err)
	}

	rpc.HandleHTTP()

	l, e := net.Listen("tcp", client.RPC_ADDRESS)
	if e != nil {
		app.logger.Fatal("listen error:", e)
	}

	go http.Serve(l, nil)
}
