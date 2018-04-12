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

func (app ApplicationContext) Info(req types.RequestInfo) types.ResponseInfo {
	hashes := 0
	txs := 0

	bytes, _ := json.Marshal(&responseInfo{Hashes: hashes, Txs: txs})

	return types.ResponseInfo{Data: string(bytes)}
}

func (app ApplicationContext) InitChain(req types.RequestInitChain) types.ResponseInitChain {
	return types.ResponseInitChain{}
}

func (app ApplicationContext) Query(req types.RequestQuery) types.ResponseQuery {
	return types.ResponseQuery{}
}

func (app ApplicationContext) SetOption(req types.RequestSetOption) types.ResponseSetOption {
	return types.ResponseSetOption{}
}

func (app ApplicationContext) CheckTx(tx []byte) types.ResponseCheckTx {
	return types.ResponseCheckTx{}
}

func (app ApplicationContext) BeginBlock(req types.RequestBeginBlock) types.ResponseBeginBlock {
	return types.ResponseBeginBlock{}
}

func (app ApplicationContext) DeliverTx(tx []byte) types.ResponseDeliverTx {
	return types.ResponseDeliverTx{}
}

func (app ApplicationContext) EndBlock(req types.RequestEndBlock) types.ResponseEndBlock {
	return types.ResponseEndBlock{}
}

func (app ApplicationContext) Commit() types.ResponseCommit {
	return types.ResponseCommit{}
}
