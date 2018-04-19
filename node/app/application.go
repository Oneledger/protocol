/*
	Copyright 2017-2018 OneLedger

	AN ABCi application node to process transactions from Tendermint Consensus
*/
package app

import (
	//"fmt"

	"github.com/tendermint/abci/types"
)

// ApplicationContext keeps all of the upper level global values.
type Application struct {
	types.BaseApplication

	status   *Datastore // current state of any composite transactions
	accounts *Datastore // identity management
	utxo     *Datastore // unspent transctions

	// TODO: basecoin has fees and staking too?
}

// NewApplicationContext initializes a new application
func NewApplication() *Application {
	return &Application{
		status:   NewDatastore("status", MEMORY),
		accounts: NewDatastore("accounts", MEMORY),
		utxo:     NewDatastore("utxo", MEMORY),
	}
}

// Type aliases
type BeginRequest = types.RequestBeginBlock
type BeginResponse = types.ResponseBeginBlock

// InitChain is called when a new chain is getting created
func (app Application) InitChain(req types.RequestInitChain) types.ResponseInitChain {
	Log.Debug("Message: InitChain", "req", req)

	return types.ResponseInitChain{}
}

// Info returns the current block information
func (app Application) Info(req types.RequestInfo) types.ResponseInfo {
	info := NewResponseInfo(0, 0, 0)

	Log.Debug("Message: Info", "req", req, "info", info)

	return types.ResponseInfo{
		Data: info.JSON(),
	}
}

// Query returns a transaction or a proof
func (app Application) Query(req types.RequestQuery) types.ResponseQuery {
	Log.Debug("Message: Query", "req", req)

	return types.ResponseQuery{}
}

// SetOption changes the underlying options for the ABCi app
func (app Application) SetOption(req types.RequestSetOption) types.ResponseSetOption {
	Log.Debug("Message: SetOption")

	return types.ResponseSetOption{}
}

// CheckTx tests to see if a transaction is valid
func (app Application) CheckTx(tx []byte) types.ResponseCheckTx {
	Log.Debug("Message: CheckTx", "tx", tx)

	result, err := Parse(Message(tx))
	if err != 0 {
		//return types.ResponseCheckTx{Code: uint32(err)}
		Log.Error("Parse Failed, invalid type")
		return types.ResponseCheckTx{Code: types.CodeTypeOK}
	}

	// TODO: Do something real here
	_ = result

	return types.ResponseCheckTx{Code: types.CodeTypeOK}
}

// BeginBlock is called when a new block is started
func (app Application) BeginBlock(req BeginRequest) BeginResponse {
	Log.Debug("Message: BeginBlock", "req", req)

	return BeginResponse{}
}

// DeliverTx accepts a transaction and updates all relevant data
func (app Application) DeliverTx(tx []byte) types.ResponseDeliverTx {
	Log.Debug("Message: DeliverTx", "tx", tx)

	result, err := Parse(Message(tx))
	if err != 0 {
		return types.ResponseDeliverTx{Code: uint32(err)}
	}

	// TODO: Do something real here
	_ = result

	return types.ResponseDeliverTx{Code: types.CodeTypeOK}
}

// EndBlock is called at the end of all of the transactions
func (app Application) EndBlock(req types.RequestEndBlock) types.ResponseEndBlock {
	Log.Debug("Message: EndBlock", "req", req)

	return types.ResponseEndBlock{}
}

// Commit tells the app to make everything persistent
func (app Application) Commit() types.ResponseCommit {
	Log.Debug("Message: Commit")

	return types.ResponseCommit{}
}
