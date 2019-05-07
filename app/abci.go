package app

import (
	"github.com/Oneledger/protocol/version"
	abci "github.com/tendermint/tendermint/abci/types"
)

// The following set of functions will be passed to the abciController

// query connection: for querying the application state; only uses Query and Info
func (app *App) infoServer() infoServer {
	return func(info RequestInfo) ResponseInfo {
		return ResponseInfo{
			Data:             app.name,
			Version:          version.Fullnode.String(),
			LastBlockHeight:  app.header.Height,
			LastBlockAppHash: app.header.AppHash,
		}
	}
}

func (app *App) queryer() queryer {
	return func(RequestQuery) ResponseQuery {
		// Do stuff
		return ResponseQuery{}
	}
}

func (app *App) optionSetter() optionSetter {
	return func(RequestSetOption) ResponseSetOption {
		// TODO
		return ResponseSetOption{
			Code: CodeOK,
		}
	}
}

// consensus methods: for executing transactions that have been committed. Message sequence is -for every block

func (app *App) chainInitializer() chainInitializer {
	return func(req RequestInitChain) ResponseInitChain {
		err := app.setupState(req.AppStateBytes)
		// This should cause consensus to halt
		if err != nil {
			return ResponseInitChain{}
		}
		return ResponseInitChain{}
	}
}

func (app *App) blockBeginner() blockBeginner {
	return func(RequestBeginBlock) ResponseBeginBlock {
		// TODO: Do stuff
		return ResponseBeginBlock{}
	}
}

func (app *App) txDeliverer() txDeliverer {
	return func([]byte) ResponseDeliverTx {
		return abci.ResponseDeliverTx{
			Code: CodeNotOK,
			Info: "Unimplemented",
		}
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

// ABCI methods
type infoServer func(RequestInfo) ResponseInfo
type optionSetter func(RequestSetOption) ResponseSetOption
type queryer func(RequestQuery) ResponseQuery
type txChecker func([]byte) ResponseCheckTx
type chainInitializer func(RequestInitChain) ResponseInitChain
type blockBeginner func(RequestBeginBlock) ResponseBeginBlock
type txDeliverer func([]byte) ResponseDeliverTx
type blockEnder func(RequestEndBlock) ResponseEndBlock
type commitor func() ResponseCommit

// abciController ensures that the implementing type can control an underlying ABCI app
type abciController interface {
	infoServer() infoServer
	optionSetter() optionSetter
	queryer() queryer
	txChecker() txChecker
	chainInitializer() chainInitializer
	blockBeginner() blockBeginner
	txDeliverer() txDeliverer
	blockEnder() blockEnder
	commitor() commitor
}

var _ ABCIApp = &ABCI{}

// ABCI is used as an input for creating a new node
type ABCI struct {
	infoServer       infoServer
	optionSetter     optionSetter
	queryer          queryer
	txChecker        txChecker
	chainInitializer chainInitializer
	blockBeginner    blockBeginner
	txDeliverer      txDeliverer
	blockEnder       blockEnder
	commitor         commitor
}

func (app *ABCI) Info(request RequestInfo) ResponseInfo {
	return app.infoServer(request)
}

func (app *ABCI) SetOption(request RequestSetOption) ResponseSetOption {
	return ResponseSetOption{}
}

func (app *ABCI) Query(request RequestQuery) ResponseQuery {
	return app.queryer(request)
}

func (app *ABCI) CheckTx(tx []byte) ResponseCheckTx {
	return app.txChecker(tx)
}

func (app *ABCI) InitChain(request RequestInitChain) ResponseInitChain {
	return app.chainInitializer(request)
}

func (app *ABCI) BeginBlock(request RequestBeginBlock) ResponseBeginBlock {
	return app.blockBeginner(request)
}

func (app *ABCI) DeliverTx(tx []byte) ResponseDeliverTx {
	return app.txDeliverer(tx)
}

func (app *ABCI) EndBlock(request RequestEndBlock) ResponseEndBlock {
	return app.blockEnder(request)
}

func (app *ABCI) Commit() ResponseCommit {
	return app.commitor()
}
