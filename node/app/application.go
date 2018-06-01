/*
	Copyright 2017-2018 OneLedger

	AN ABCi application node to process transactions from Tendermint Consensus
*/
package app

import (
	//"fmt"
	"bytes"

	"github.com/Oneledger/protocol/node/abci"
	"github.com/Oneledger/protocol/node/action"
	"github.com/Oneledger/protocol/node/data"
	"github.com/Oneledger/protocol/node/id"
	"github.com/Oneledger/protocol/node/log"
	"github.com/tendermint/abci/types"
)

var ChainId string

func init() {
	// TODO: Should be driven from config
	ChainId = "OneLedger-Root"
}

// ApplicationContext keeps all of the upper level global values.
type Application struct {
	types.BaseApplication

	Admin      *data.Datastore  // any administrative parameters
	Status     *data.Datastore  // current state of any composite transactions (pending, verified, etc.)
	Identities *id.Identities   // Keep a higher-level identity for a given user
	Accounts   *id.Accounts     // Keep all of the user accounts locally for their node (identity management)
	Utxo       *data.ChainState // unspent transction output (for each type of coin)
}

// NewApplicationContext initializes a new application
func NewApplication() *Application {
	return &Application{
		Admin:      data.NewDatastore("admin", data.PERSISTENT),
		Status:     data.NewDatastore("status", data.PERSISTENT),
		Identities: id.NewIdentities("identities"),
		Accounts:   id.NewAccounts("accounts"),
		Utxo:       data.NewChainState("utxo", data.PERSISTENT),
	}
}

// Access to the local persistent databases
func (app Application) GetAdmin() interface{} {
	return app.Admin
}

// Access to the local persistent databases
func (app Application) GetStatus() interface{} {
	return app.Status
}

// Access to the local persistent databases
func (app Application) GetIdentities() interface{} {
	return app.Identities
}

// Access to the local persistent databases
func (app Application) GetAccounts() interface{} {
	return app.Accounts
}

// InitChain is called when a new chain is getting created
func (app Application) InitChain(req RequestInitChain) ResponseInitChain {
	log.Debug("Message: InitChain", "req", req)

	// TODO: Insure that all of the databases and shared resources are reset here

	return ResponseInitChain{}
}

// SetOption changes the underlying options for the ABCi app
func (app Application) SetOption(req RequestSetOption) ResponseSetOption {
	log.Debug("Message: SetOption")

	SetOption(&app, req.Key, req.Value)

	return ResponseSetOption{Code: types.CodeTypeOK}
}

// Info returns the current block information
func (app Application) Info(req RequestInfo) ResponseInfo {
	info := abci.NewResponseInfo(0, 0, 0)

	// lastHeight := app.Utxo.Commit.Height()

	log.Debug("Message: Info", "req", req, "info", info)

	return ResponseInfo{
		Data: info.JSON(),
		// Version: version,
		// LastBlockHeight: lastHeight,
		// LastBlockAppHash: lastAppHash,
	}
}

// Query returns a transaction or a proof
func (app Application) Query(req RequestQuery) ResponseQuery {
	log.Debug("Message: Query", "req", req, "path", req.Path, "data", req.Data)

	result := HandleQuery(req.Path, req.Data)

	return ResponseQuery{Key: action.Message("result"), Value: result}
}

// CheckTx tests to see if a transaction is valid
func (app Application) CheckTx(tx []byte) ResponseCheckTx {
	log.Debug("Message: CheckTx", "tx", tx)

	result, err := action.Parse(action.Message(tx))
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

var chainKey data.DatabaseKey = data.DatabaseKey("chainId")

// BeginBlock is called when a new block is started
func (app Application) BeginBlock(req RequestBeginBlock) ResponseBeginBlock {
	log.Debug("Message: BeginBlock", "req", req)

	newChainId := action.Message(req.Header.ChainID)

	chainId := app.Admin.Load(chainKey)

	if chainId == nil {
		chainId = app.Admin.Store(chainKey, newChainId)

	} else if bytes.Compare(chainId, newChainId) != 0 {
		log.Error("Mismatching chains", "chainId", chainId, "newChainId", newChainId)
	}

	return ResponseBeginBlock{}
}

// DeliverTx accepts a transaction and updates all relevant data
func (app Application) DeliverTx(tx []byte) ResponseDeliverTx {
	log.Debug("Message: DeliverTx", "tx", tx)

	result, err := action.Parse(action.Message(tx))
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
	log.Debug("Message: EndBlock", "req", req)

	return ResponseEndBlock{}
}

// Commit tells the app to make everything persistent
func (app Application) Commit() ResponseCommit {
	log.Debug("Message: Commit")

	// Commit any pending changes.
	//app.Utxo.Commit()

	return ResponseCommit{}
}
