package app

import (
	"github.com/Oneledger/protocol/data"
	"github.com/tendermint/tendermint/libs/common"
)

// Ensure this App struct can control the underlying ABCI app
var _ abciController = &App{}

type App struct {
	balances         *data.ChainState
	identities       data.Store
	smartContract    data.Store
	executionContext data.Store
	admin            data.Store
	accounts         data.Store
	sequence         data.Store
	status           data.Store
	contract         data.Store
	event            data.Store

	sdk           common.Service

	header     Header   // Tendermint last header info
	validators interface{} // Set of validators currently active
	abci *ABCI
}

// Getters
func (app *App) Balances() *data.ChainState {
	return app.balances
}
func (app *App) Identities() data.Store {
	return app.identities
}
func (app *App) SmartContract() data.Store {
	return app.smartContract
}
func (app *App) ExecutionContext() data.Store {
	return app.executionContext
}
func (app *App) Admin() data.Store {
	return app.admin
}
func (app *App) Accounts() data.Store {
	return app.accounts
}
func (app *App) Sequence() data.Store {
	return app.sequence
}
func (app *App) Status() data.Store {
	return app.status
}
func (app *App) Contract() data.Store {
	return app.contract
}
func (app *App) Event() data.Store {
	return app.event
}

// TODO: Add proper types
func (app *App) Validators() interface{} {
	return app.validators
}

// setNewABCI returns a new ABCI struct with the current context-values set in App
func (app *App) setNewABCI() {
	app.abci = &ABCI{
		infoServer: app.infoServer(),
		optionSetter: app.optionSetter(),
		queryer: app.queryer(),
		txChecker: app.txChecker(),
		chainInitializer: app.chainInitializer(),
		blockBeginner: app.blockBeginner(),
		txDeliverer: app.txDeliverer(),
		blockEnder: app.blockEnder(),
		commitor: app.commitor(),
	}
}

// ABCI returns an ABCI-ready Application used to initialize the new Node
func (app *App) ABCI() *ABCI {
	return app.abci
}

func (app *App) infoServer() infoServer {
	return func(info RequestInfo) ResponseInfo {
		return ResponseInfo{}
	}
}

func(app *App) optionSetter() optionSetter {
	return func(RequestSetOption) ResponseSetOption {
		// Do stuff
		return ResponseSetOption{}
	}
}
func(app *App) queryer() queryer {
	return func(RequestQuery) RequestQuery {
		// Do stuff
		return RequestQuery{}
	}
}
func(app *App) txChecker() txChecker {
	return func([]byte) ResponseCheckTx {
		// Do stuff
		return ResponseCheckTx{}
	}
}
func(app *App) chainInitializer() chainInitializer {
	return func(RequestInitChain) ResponseInitChain {
		// Do stuff
		return ResponseInitChain{}
	}
}
func(app *App) blockBeginner() blockBeginner {
	return func(RequestBeginBlock) ResponseBeginBlock {
		// Do stuff
		return ResponseBeginBlock{}
	}
}
func(app *App) txDeliverer() txDeliverer {
	return func([]byte) ResponseDeliverTx {
		// Do stuff
		return ResponseDeliverTx{}
	}
}
func(app *App) blockEnder() blockEnder {
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
