/*
	Copyright 2017-2018 OneLedger

	ABCi application node to process transactions from Tendermint
*/
package app

import (
	//"fmt"

	"github.com/tendermint/abci/types"
	"github.com/tendermint/tmlibs/log"
)

// ApplicationContext keeps all of the upper level global values.
type Application struct {
	types.BaseApplication

	log log.Logger // inherited logger
	db  Datastore  // key/value database in memory
}

// NewApplicationContext initializes a new application
func NewApplication(logger log.Logger) *Application {
	return &Application{
		log: logger,
		db:  *NewDatastore(),
	}
}

var logger log.Logger
var key Key = Key("x")

func (app Application) Info(req types.RequestInfo) types.ResponseInfo {
	app.log.Debug("Message: Info")

	info := NewResponseInfo(0, 0, 0)

	//_ = info.Json()
	json := info.JsonMessage()
	app.db.Store(key, json)

	app.log.Debug("Message: Info", "info", info)

	return types.ResponseInfo{Data: string(json)}
	//return types.ResponseInfo{}
}

func (app Application) InitChain(req types.RequestInitChain) types.ResponseInitChain {
	app.log.Debug("Message: InitChain")

	return types.ResponseInitChain{}
}

func (app Application) Query(req types.RequestQuery) types.ResponseQuery {
	app.log.Debug("Message: Query")

	return types.ResponseQuery{}
}

func (app Application) SetOption(req types.RequestSetOption) types.ResponseSetOption {
	app.log.Debug("Message: SetOption")

	return types.ResponseSetOption{}
}

func (app Application) CheckTx(tx []byte) types.ResponseCheckTx {
	app.log.Debug("Message: CheckTx")

	return types.ResponseCheckTx{Code: types.CodeTypeOK}
}

func (app Application) BeginBlock(req types.RequestBeginBlock) types.ResponseBeginBlock {
	app.log.Debug("Message: BeginBlock")

	return types.ResponseBeginBlock{}
}

func (app Application) DeliverTx(tx []byte) types.ResponseDeliverTx {
	app.log.Debug("Message: DeliverTx")

	return types.ResponseDeliverTx{Code: types.CodeTypeOK}
}

func (app Application) EndBlock(req types.RequestEndBlock) types.ResponseEndBlock {
	app.log.Debug("Message: EndBlock")

	return types.ResponseEndBlock{}
}

func (app Application) Commit() types.ResponseCommit {
	app.log.Debug("Message: Commit")

	return types.ResponseCommit{}
}
