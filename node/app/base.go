package app

import (
	//"fmt"
	"github.com/tendermint/abci/types"
	//"github.com/tendermint/tmlibs/common"
)

type Context struct {
	types.BaseApplication
}

func NewContext() *Context {
	return &Context{}
}

func (app *Context) Info(req types.RequestInfo) types.ResponseInfo {
	return types.ResponseInfo{
		Data: "{ \"hashes\": 0, \"txs\": 0 }",
	}
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
