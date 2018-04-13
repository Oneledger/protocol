/*
	OneLedger Copywrite 2017-2018

	ABCi based node to process transactions from Tendermint
*/
package app

import (
	//"fmt"

	"github.com/tendermint/abci/types"
	"github.com/tendermint/tmlibs/db"
	"github.com/tendermint/tmlibs/log"
)

var logger log.Logger

// NewApplicationContext initializes a new application
func NewApplicationContext(logger log.Logger) *ApplicationContext {
	return &ApplicationContext{
		log: logger,
		db:  db.NewMemDB(),
	}
}

func (app ApplicationContext) Store(key Key, value Message) {
	app.db.Set(key, value)
}

func (app ApplicationContext) Load(key Key) (value Message) {
	return app.db.Get(key)
}

// Response arguments
type responseInfo struct {
	Hashes int `json:"hashes"`
	Txs    int `json:"txs"`
}

func (app ApplicationContext) Info(req types.RequestInfo) types.ResponseInfo {
	app.log.Debug("Message: Info")

	info := NewResponseInfo(10, 10)

	//json := info.Json()
	_ = info.Json()

	app.log.Debug("Message: Info", "info", info)

	//return types.ResponseInfo{Data: json}
	return types.ResponseInfo{}
}

func (app ApplicationContext) InitChain(req types.RequestInitChain) types.ResponseInitChain {
	app.log.Debug("Message: InitChain")

	return types.ResponseInitChain{}
}

func (app ApplicationContext) Query(req types.RequestQuery) types.ResponseQuery {
	app.log.Debug("Message: Query")

	return types.ResponseQuery{}
}

func (app ApplicationContext) SetOption(req types.RequestSetOption) types.ResponseSetOption {
	app.log.Debug("Message: SetOption")

	return types.ResponseSetOption{}
}

func (app ApplicationContext) CheckTx(tx []byte) types.ResponseCheckTx {
	app.log.Debug("Message: CheckTx")

	return types.ResponseCheckTx{Code: types.CodeTypeOK}
}

func (app ApplicationContext) BeginBlock(req types.RequestBeginBlock) types.ResponseBeginBlock {
	app.log.Debug("Message: BeginBlock")

	return types.ResponseBeginBlock{}
}

func (app ApplicationContext) DeliverTx(tx []byte) types.ResponseDeliverTx {
	app.log.Debug("Message: DeliverTx")

	return types.ResponseDeliverTx{Code: types.CodeTypeOK}
}

func (app ApplicationContext) EndBlock(req types.RequestEndBlock) types.ResponseEndBlock {
	app.log.Debug("Message: EndBlock")

	return types.ResponseEndBlock{}
}

func (app ApplicationContext) Commit() types.ResponseCommit {
	app.log.Debug("Message: Commit")

	return types.ResponseCommit{}
}
