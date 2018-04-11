package app

import (
	//"fmt"
	"encoding/json"
	"github.com/tendermint/abci/types"
	//"github.com/tendermint/tmlibs/common"
)

type Context struct {
	types.BaseApplication
}

func NewContext() *Context {
	return &Context{}
}

type responseInfo struct {
	Hashes int `json:"hashes"`
	Txs    int `json:"txs"`
}

func (app *Context) Info(req types.RequestInfo) types.ResponseInfo {
	hashes := 0
	txs := 0

	bytes, _ := json.Marshal(&responseInfo{Hashes: hashes, Txs: txs})

	return types.ResponseInfo{Data: info}
}

func (app *Context) Query() types.ResponseQuery {
	return types.ResponseQuery{}
}

func (app *Context) SetOption(req types.RequestSetOption) types.ResponseSetOption {
	return types.ResponseSetOption{}
}

func (app *Context) DeliverTx(tx []byte) types.ResponseDeliverTx {
	return types.ResponseDeliverTx{}
}

func (app *Context) CheckTx(tx []byte) types.ResponseCheckTx {
	return types.ResponseCheckTx{}
}

func (app *Context) Commit() types.ResponseCommit {
	return types.ResponseCommit{}
}
