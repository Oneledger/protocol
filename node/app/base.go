package app

import (
	//"fmt"
	"encoding/json"
	"github.com/tendermint/abci/types"
	//"github.com/tendermint/tmlibs/common"
)

type ApplicationContext struct {
	types.BaseApplication
}

func NewApplicationContext() *ApplicationContext {
	return &ApplicationContext{}
}

type responseInfo struct {
	Hashes int `json:"hashes"`
	Txs    int `json:"txs"`
}

func (app *ApplicationContext) Info(req types.RequestInfo) types.ResponseInfo {
	hashes := 0
	txs := 0

	bytes, _ := json.Marshal(&responseInfo{Hashes: hashes, Txs: txs})

	// TODO: need to cast bytes into this type...
	info := string(bytes)

	return types.ResponseInfo{Data: info}
}

func (app *ApplicationContext) Query() types.ResponseQuery {
	return types.ResponseQuery{}
}

func (app *ApplicationContext) SetOption(req types.RequestSetOption) types.ResponseSetOption {
	return types.ResponseSetOption{}
}

func (app *ApplicationContext) DeliverTx(tx []byte) types.ResponseDeliverTx {
	return types.ResponseDeliverTx{}
}

func (app *ApplicationContext) CheckTx(tx []byte) types.ResponseCheckTx {
	return types.ResponseCheckTx{}
}

func (app *ApplicationContext) Commit() types.ResponseCommit {
	return types.ResponseCommit{}
}
