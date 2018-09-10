/*
	Copyright 2017-2018 OneLedger

	AN ABCi application node to process transactions from Tendermint Consensus
*/
package app

import (
	"bytes"
	"github.com/Oneledger/protocol/node/abci"
	"github.com/Oneledger/protocol/node/action"
	"github.com/Oneledger/protocol/node/comm"
	"github.com/Oneledger/protocol/node/convert"
	"github.com/Oneledger/protocol/node/data"
	"github.com/Oneledger/protocol/node/global"
	"github.com/Oneledger/protocol/node/id"
	"github.com/Oneledger/protocol/node/log"
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

	Admin      *data.Datastore  // any administrative parameters
	Status     *data.Datastore  // current state of any composite transactions (pending, verified, etc.)
	Identities *id.Identities   // Keep a higher-level identity for a given user
	Accounts   *id.Accounts     // Keep all of the user accounts locally for their node (identity management)
	Utxo       *data.ChainState // unspent transction output (for each type of coin)
	Event      *data.Datastore  // Event for any action that need to be tracked
	Contract	   *data.Datastore  // contract for reuse.
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
		Contract: 	data.NewDatastore("contract", data.PERSISTENT),
	}
}

// Initial the state of the application from persistent data
func (app Application) Initialize() {
	param := app.Admin.Load(data.DatabaseKey("NodeAccountName"))
	if param != nil {
		var name string
		buffer, err := comm.Deserialize(param, &name)
		if err != nil {
			log.Error("Failed to deserialize persistent data")
		}
		global.Current.NodeAccountName = *(buffer.(*string))
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

// Access to the local persistent databases
func (app Application) GetUtxo() interface{} {
	return app.Utxo
}

func (app Application) GetChainID() interface{} {
    return ChainId
}

func (app Application) GetEvent() interface{} {
    return app.Event
}

func (app Application) GetContract() interface{} {
	return app.Contract
}

type BasicState struct {
	Account string `json:"account"`
	Amount  int64  `json:"coins"` // TODO: Should be corrected as Amount, not coins
}

// Use the Genesis block to initialze the system
func (app Application) SetupState(stateBytes []byte) {
	log.Debug("SetupState", "state", string(stateBytes))

	var base BasicState
	des, err := comm.Deserialize(stateBytes, &base)
	if err != nil {
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
	account, _ := app.Accounts.FindName(state.Account + "-OneLedger")

	// TODO: Should be more flexible to match genesis block
	balance := data.Balance{
		Amount: data.NewCoin(state.Amount, "OLT"),
	}

	buffer, err := comm.Serialize(balance)
	if err != nil {
		log.Error("Failed to Serialize balance")
	}

	// Use the account key in the database.
	app.Utxo.Delivered.Set(account.AccountKey(), buffer)
	app.Utxo.Delivered.SaveVersion()
	app.Utxo.Commit()

	log.Info("Genesis State UTXO database", "balance", balance)
}

// InitChain is called when a new chain is getting created
func (app Application) InitChain(req RequestInitChain) ResponseInitChain {
	log.Debug("Contract: InitChain", "req", req)

	app.SetupState(req.AppStateBytes)

	return ResponseInitChain{}
}

// SetOption changes the underlying options for the ABCi app
func (app Application) SetOption(req RequestSetOption) ResponseSetOption {
	log.Debug("Contract: SetOption", "key", req.Key, "value", req.Value)

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

	log.Debug("Contract: Info", "req", req, "info", info)

	return ResponseInfo{
		Data:            info.JSON(),
		Version:         convert.GetString64(app.Utxo.Version),
		LastBlockHeight: int64(app.Utxo.Height),
		// TODO: Should return a valid AppHash
		//LastBlockAppHash: app.Utxo.Hash,
	}
}

// Query returns a transaction or a proof
func (app Application) Query(req RequestQuery) ResponseQuery {
	log.Debug("Contract: Query", "req", req, "path", req.Path, "data", req.Data)

	result := HandleQuery(app, req.Path, req.Data)

	return ResponseQuery{
		Code:   2,
		Log:    "Log Information",
		Info:   "Info Information",
		Index:  0,
		Key:    action.Message("result"),
		Value:  result,
		Proof:  nil,
		Height: int64(app.Utxo.Height),
	}
}

// CheckTx tests to see if a transaction is valid
func (app Application) CheckTx(tx []byte) ResponseCheckTx {
	log.Debug("Contract: CheckTx", "tx", tx)

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
		GasUsed: 1000,
		Tags:      []common.KVPair(nil),
	}
}

var chainKey data.DatabaseKey = data.DatabaseKey("chainId")

// BeginBlock is called when a new block is started
func (app Application) BeginBlock(req RequestBeginBlock) ResponseBeginBlock {
	//log.Debug("Contract: BeginBlock", "req", req)

	newChainId := action.Message(req.Header.ChainID)

	chainId := app.Admin.Load(chainKey)

	if chainId == nil {
		chainId = app.Admin.Store(chainKey, newChainId)

	} else if bytes.Compare(chainId, newChainId) != 0 {
		log.Warn("Mismatching chains", "chainId", chainId, "newChainId", newChainId)
	}

	return ResponseBeginBlock{}
}

// DeliverTx accepts a transaction and updates all relevant data
func (app Application) DeliverTx(tx []byte) ResponseDeliverTx {
	log.Debug("Contract: DeliverTx", "tx", tx)

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
	log.Debug("Contract: EndBlock", "req", req)

	return ResponseEndBlock{}
}

// Commit tells the app to make everything persistent
func (app Application) Commit() ResponseCommit {
	log.Debug("Contract: Commit")

	// Commit any pending changes.
	hash, version := app.Utxo.Commit()

	log.Debug("Committed", "hash", hash, "version", version)

	return ResponseCommit{}
}
