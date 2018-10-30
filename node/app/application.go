/*
	Copyright 2017-2018 OneLedger

	AN ABCi application node to process transactions from Tendermint Consensus
*/
package app

import (
	"bytes"

	"github.com/Oneledger/protocol/node/abci"
	"github.com/Oneledger/protocol/node/action"
	"github.com/Oneledger/protocol/node/data"
	"github.com/Oneledger/protocol/node/err"
	"github.com/Oneledger/protocol/node/id"
	"github.com/Oneledger/protocol/node/log"
	"github.com/Oneledger/protocol/node/serial"
	"github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/common"
)

var ChainId string

func init() {
	// TODO: Should be driven from config
	ChainId = "OneLedger-Root"
}

// ApplicationContext keeps all of the upper level global values.
type Application struct {
	types.BaseApplication

	Admin      data.Datastore   // any administrative parameters
	Status     data.Datastore   // current state of any composite transactions (pending, verified, etc.)
	Identities *id.Identities   // Keep a higher-level identity for a given user
	Accounts   *id.Accounts     // Keep all of the user accounts locally for their node (identity management)
	Utxo       *data.ChainState // unspent transction output (for each type of coin)
	Event      data.Datastore   // Event for any action that need to be tracked
	Contract   data.Datastore   // contract for reuse.

	LastHeader types.Header // Tendermint last header info
}

// NewApplicationContext initializes a new application
func NewApplication() *Application {
	return &Application{
		Admin:      data.NewDatastore("admin", data.PERSISTENT),
		Status:     data.NewDatastore("status", data.PERSISTENT),
		Identities: id.NewIdentities("identities"),
		Accounts:   id.NewAccounts("accounts"),
		Utxo:       data.NewChainState("utxo", data.PERSISTENT),
		Event:      data.NewDatastore("event", data.PERSISTENT),
		Contract:   data.NewDatastore("contract", data.PERSISTENT),
	}
}

// Initial the state of the application from persistent data
func (app Application) Initialize() {
	param := app.Admin.Get(data.DatabaseKey("NodeAccountName"))
	if param != nil {
		log.Dump("The parameter is", param)
		//global.Current.NodeAccountName = param.(string)
	}
}

type BasicState struct {
	Account string `json:"account"`
	Amount  int64  `json:"coins"` // TODO: Should be corrected as Amount, not coins
}

// Use the Genesis block to initialze the system
func (app Application) SetupState(stateBytes []byte) {
	log.Debug("SetupState", "state", string(stateBytes))

	var base BasicState

	des, errx := serial.Deserialize(stateBytes, &base, serial.JSON)
	if errx != nil {
		log.Fatal("Failed to deserialize stateBytes during SetupState")
	}

	state := des.(*BasicState)
	log.Debug("Deserialized State", "state", state, "state.Account", state.Account)

	// TODO: Can't generate a different key for each node. Needs to be in the genesis? Or ignored?
	//publicKey, privateKey := id.GenerateKeys([]byte(state.Account)) // TODO: switch with passphrase
	publicKey, privateKey := id.NilPublicKey(), id.NilPrivateKey()

	// TODO: This should probably only occur on the Admin node, for other nodes how do I know the key?
	// Register the identity and account first
	RegisterLocally(&app, state.Account, "OneLedger", data.ONELEDGER, publicKey, privateKey)
	account, status := app.Accounts.FindName(state.Account + "-OneLedger")
	if status != err.SUCCESS {
		log.Fatal("Recently Added Account is missing", "name", state.Account, "status", status)
	}

	// Use the account key in the database.
	balance := data.NewBalance(state.Amount, "OLT")
	app.Utxo.Set(account.AccountKey(), balance)

	// TODO: Until a block is commited, this data is not persistent
	//app.Utxo.Commit()

	log.Info("Genesis State UTXO database", "balance", balance)
}

// InitChain is called when a new chain is getting created
func (app Application) InitChain(req RequestInitChain) ResponseInitChain {
	log.Debug("ABCI: InitChain", "req", req)

	app.SetupState(req.AppStateBytes)

	return ResponseInitChain{}
}

// SetOption changes the underlying options for the ABCi app
func (app Application) SetOption(req RequestSetOption) ResponseSetOption {
	log.Debug("ABCI: SetOption", "key", req.Key, "value", req.Value)

	SetOption(&app, req.Key, req.Value)

	return ResponseSetOption{
		Code: types.CodeTypeOK,
		Log:  "Log Data",
		Info: "Info Data",
	}
}

// Info returns the current block information
func (app Application) Info(req RequestInfo) ResponseInfo {

	info := abci.NewResponseInfo(0, 0, 0)

	// TODO: Get the correct height from the last committed tree
	// lastHeight := app.Utxo.Commit.Height()

	log.Debug("ABCI: Info", "req", req, "info", info)

	result := ResponseInfo{
		Data: info.JSON(),
		//Version: convert.GetString64(app.Utxo.Version),

		// The version of the tree, needs to match the height of the chain
		//LastBlockHeight: int64(0),
		LastBlockHeight: int64(app.Utxo.Version),

		// TODO: Should return a valid AppHash
		LastBlockAppHash: app.Utxo.Hash,
	}

	log.Dump("Info Response is", result)
	return result
}

// Query returns a transaction or a proof
func (app Application) Query(req RequestQuery) ResponseQuery {
	log.Debug("ABCI: Query", "req", req, "path", req.Path, "data", req.Data)

	result := HandleQuery(app, req.Path, req.Data)

	return ResponseQuery{
		Code:   2,
		Log:    "Log Information",
		Info:   "Info Information",
		Index:  0,
		Key:    action.Message("result"),
		Value:  result,
		Proof:  nil,
		Height: int64(app.Utxo.Version),
	}
}

// CheckTx tests to see if a transaction is valid
func (app Application) CheckTx(tx []byte) ResponseCheckTx {
	log.Debug("ABCI: CheckTx", "tx", tx)

	if tx == nil {
		log.Warn("Empty Transaction, Ignoring", "tx", tx)
		return ResponseCheckTx{Code: err.PARSE_ERROR}
	}

	result, err := action.Parse(action.Message(tx))
	if err != 0 || result == nil {
		return ResponseCheckTx{Code: err}
	}

	// Check that this is a valid transaction
	if err = result.Validate(); err != 0 {
		return ResponseCheckTx{Code: err}
	}

	// Check that this transaction works in the context
	if err = result.ProcessCheck(&app); err != 0 {
		return ResponseCheckTx{Code: err}
	}

	return ResponseCheckTx{
		Code:      types.CodeTypeOK,
		Data:      []byte("Data"),
		Log:       "Log Data",
		Info:      "Info Data",
		GasWanted: 1000,
		GasUsed:   1000,
		Tags:      []common.KVPair(nil),
	}
}

var chainKey data.DatabaseKey = data.DatabaseKey("chainId")

// BeginBlock is called when a new block is started
func (app Application) BeginBlock(req RequestBeginBlock) ResponseBeginBlock {
	//log.Debug("ABCI: BeginBlock", "req", req)
	app.LastHeader = req.Header

	newChainId := action.Message(req.Header.ChainID)

	chainId := app.Admin.Get(chainKey)

	if chainId == nil {
		session := app.Admin.Begin()
		session.Set(chainKey, newChainId)
		session.Commit()

		// TODO: This is questionable?
		chainId = newChainId

	} else if bytes.Compare(chainId.([]byte), newChainId) != 0 {
		log.Warn("Mismatching chains", "chainId", chainId, "newChainId", newChainId)
	}

	return ResponseBeginBlock{
		Tags: []common.KVPair(nil),
	}
}

// DeliverTx accepts a transaction and updates all relevant data
func (app Application) DeliverTx(tx []byte) ResponseDeliverTx {
	log.Debug("ABCI: DeliverTx", "tx", tx)

	result, err := action.Parse(action.Message(tx))
	if err != 0 || result == nil {
		return ResponseDeliverTx{Code: err}
	}

	if err = result.Validate(); err != 0 {
		return ResponseDeliverTx{Code: err}
	}

	if result.ShouldProcess(app) {
		ttype, _ := action.UnpackMessage(action.Message(tx))

		if ttype == action.SWAP || ttype == action.PUBLISH || ttype == action.VERIFY {
			go result.ProcessDeliver(&app)
		} else {
			if err = result.ProcessDeliver(&app); err != 0 {
				log.Warn("Processing Failed", "err", err)
				return ResponseDeliverTx{Code: err}
			}
		}
	}

	return ResponseDeliverTx{
		Code:      types.CodeTypeOK,
		Data:      []byte("Data"),
		Log:       "Log Data",
		Info:      "Info Data",
		GasWanted: 1000,
		GasUsed:   1000,
		Tags:      []common.KVPair(nil),
	}
}

// EndBlock is called at the end of all of the transactions
func (app Application) EndBlock(req RequestEndBlock) ResponseEndBlock {
	log.Debug("ABCI: EndBlock", "req", req)

	return ResponseEndBlock{
		Tags: []common.KVPair(nil),
	}
}

// Commit tells the app to make everything persistent
func (app Application) Commit() ResponseCommit {
	log.Debug("ABCI: Commit")

	// Commit any pending changes.
	hash, version := app.Utxo.Commit()

	log.Debug("-- Committed New Block", "hash", hash, "version", version)

	return ResponseCommit{
		Data: hash,
	}
}
