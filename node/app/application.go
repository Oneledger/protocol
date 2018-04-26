/*
	Copyright 2017-2018 OneLedger

	AN ABCi application node to process transactions from Tendermint Consensus
*/
package app

import (
	//"fmt"
	"bytes"

	"github.com/Oneledger/prototype/node/abci"
	"github.com/tendermint/abci/types"
)

// ApplicationContext keeps all of the upper level global values.
type Application struct {
	types.BaseApplication

	Admin    *Datastore  // any administrative parameters
	Status   *Datastore  // current state of any composite transactions (pending, verified, etc.)
	Accounts *Accounts   // Keep all of the user accounts locally for their node (identity management)
	Utxo     *ChainState // unspent transction output (for each type of coin)

	// TODO: basecoin has fees and staking too?
}

// NewApplicationContext initializes a new application
func NewApplication() *Application {
	return &Application{
		Admin:    NewDatastore("admin", PERSISTENT),
		Status:   NewDatastore("status", PERSISTENT),
		Accounts: NewAccounts("accounts"),
		Utxo:     NewChainState("utxo", PERSISTENT),
	}
}

// InitChain is called when a new chain is getting created
func (app Application) InitChain(req RequestInitChain) ResponseInitChain {
	Log.Debug("Message: InitChain", "req", req)

	// TODO: Insure that all of the databases and shared resources are reset here

	return ResponseInitChain{}
}

// Info returns the current block information
func (app Application) Info(req RequestInfo) ResponseInfo {
	info := abci.NewResponseInfo(0, 0, 0)

	Log.Debug("Message: Info", "req", req, "info", info)

	return ResponseInfo{
		Data: info.JSON(),
		// LastBlockHeight: lastHeight,
		// LastBlockAppHash: lastAppHash,
	}
}

// Query returns a transaction or a proof
func (app Application) Query(req RequestQuery) ResponseQuery {
	Log.Debug("Message: Query", "req", req, "path", req.Path, "data", req.Data)

	result := HandleQuery(req.Path, req.Data)

	return ResponseQuery{Key: Message("result"), Value: result}
}

// SetOption changes the underlying options for the ABCi app
func (app Application) SetOption(req RequestSetOption) ResponseSetOption {
	Log.Debug("Message: SetOption")

	return ResponseSetOption{}
}

// CheckTx tests to see if a transaction is valid
func (app Application) CheckTx(tx []byte) ResponseCheckTx {
	Log.Debug("Message: CheckTx", "tx", tx)

	result, err := Parse(Message(tx))
	if err != 0 {
		return ResponseCheckTx{Code: err}
	}

	// Check that this is a valid transaction
	if err = result.Validate(); err != 0 {
		return ResponseCheckTx{Code: err}
	}

	if err = result.ProcessCheck(&app); err != 0 {
		return ResponseCheckTx{Code: err}
	}

	return ResponseCheckTx{Code: types.CodeTypeOK}
}

var chainKey DatabaseKey = DatabaseKey("chainId")

// BeginBlock is called when a new block is started
func (app Application) BeginBlock(req RequestBeginBlock) ResponseBeginBlock {
	Log.Debug("Message: BeginBlock", "req", req)

	newChainId := Message(req.Header.ChainID)

	chainId := app.Admin.Load(chainKey)

	if chainId == nil {
		chainId = app.Admin.Store(chainKey, newChainId)

	} else if bytes.Compare(chainId, newChainId) != 0 {
		//panic("Mismatching chains")
	}

	Log.Debug("ChainID is", "id", chainId)

	return ResponseBeginBlock{}
}

// DeliverTx accepts a transaction and updates all relevant data
func (app Application) DeliverTx(tx []byte) ResponseDeliverTx {
	Log.Debug("Message: DeliverTx", "tx", tx)

	result, err := Parse(Message(tx))
	if err != 0 {
		return ResponseDeliverTx{Code: err}
	}

	if err = result.Validate(); err != 0 {
		return ResponseDeliverTx{Code: err}
	}

	if err = result.ProcessDeliver(&app); err != 0 {
		return ResponseDeliverTx{Code: err}
	}

	return ResponseDeliverTx{Code: types.CodeTypeOK}
}

// EndBlock is called at the end of all of the transactions
func (app Application) EndBlock(req RequestEndBlock) ResponseEndBlock {
	Log.Debug("Message: EndBlock", "req", req)

	return ResponseEndBlock{}
}

// Commit tells the app to make everything persistent
func (app Application) Commit() ResponseCommit {
	Log.Debug("Message: Commit")

	// TODO: Empty commit for now, but all transactional work should be queued, and
	// only persisted on commit.

	return ResponseCommit{}
}
