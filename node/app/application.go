/*
	Copyright 2017-2018 OneLedger

	AN ABCi application node to process transactions from Tendermint Consensus
*/
package app

import (
	"bytes"
	"encoding/hex"
	"math/big"
	"strconv"

	"github.com/Oneledger/protocol/node/abci"
	"github.com/Oneledger/protocol/node/action"
	"github.com/Oneledger/protocol/node/comm"
	"github.com/Oneledger/protocol/node/data"
	"github.com/Oneledger/protocol/node/global"
	"github.com/Oneledger/protocol/node/id"
	"github.com/Oneledger/protocol/node/log"
	"github.com/Oneledger/protocol/node/serial"
	"github.com/Oneledger/protocol/node/status"
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

	// Global Chain state (data is identical on all nodes in the chain)
	Utxo       *data.ChainState // unspent transction output (for each type of coin)
	SDK        common.Service
	Identities *id.Identities // Keep a higher-level identity for a given user

	// Local Node state (data is different for each node)
	Admin    data.Datastore // any administrative parameters
	Status   data.Datastore // current state of any composite transactions (pending, verified, etc.)
	Accounts *id.Accounts   // Keep all of the user accounts locally for their node (identity management)
	Event    data.Datastore // Event for any action that need to be tracked
	Contract data.Datastore // contract for reuse.

	// Tendermint's last block information
	LastHeader types.Header // Tendermint last header info
}

// NewApplicationContext initializes a new application, reconnects to the databases.
func NewApplication() *Application {
	return &Application{
		Identities: id.NewIdentities("identities"),
		Utxo:       data.NewChainState("utxo", data.PERSISTENT),

		Admin:    data.NewDatastore("admin", data.PERSISTENT),
		Status:   data.NewDatastore("status", data.PERSISTENT),
		Accounts: id.NewAccounts("accounts"),
		Event:    data.NewDatastore("event", data.PERSISTENT),
		Contract: data.NewDatastore("contract", data.PERSISTENT),
	}
}

type AdminParameters struct {
	NodeAccountName string
	NodeName        string
}

func init() {
	serial.Register(AdminParameters{})
}

// Initial the state of the application from persistent data
func (app Application) Initialize() {
	raw := app.Admin.Get(data.DatabaseKey("NodeAccountName"))
	if raw != nil {
		params := raw.(AdminParameters)
		global.Current.NodeAccountName = params.NodeAccountName
	} else {
		log.Debug("NodeAccountName not currently set")
	}

	// SDK Server should start when the --sdkrpc argument is passed to fullnode
	sdkPort := global.Current.SDKAddress
	if sdkPort == 0 {
		return
	}

	s, err := NewSDKServer(&app, sdkPort)
	if err != nil {
		panic(err)
	} else {
		app.SDK = s
		app.SDK.Start()
	}
}

type BasicState struct {
	Account string `json:"account"`
	Amount  string `json:"coins"` // TODO: Should be corrected as Amount, not coins
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
	privateKey, publicKey := id.GenerateKeys([]byte(state.Account), false) // TODO: switch with passphrase

	CreateAccount(app, state.Account, state.Amount, publicKey, privateKey)

	publicKey, privateKey = id.OnePublicKey(), id.OnePrivateKey()
	CreateAccount(app, "Payment", "0", publicKey, privateKey)
}

func CreateAccount(app Application, stateAccount string, stateAmount string, publicKey id.PublicKeyED25519, privateKey id.PrivateKeyED25519) {

	// TODO: This should probably only occur on the Admin node, for other nodes how do I know the key?
	// Register the identity and account first
	RegisterLocally(&app, stateAccount, "OneLedger", data.ONELEDGER, publicKey, privateKey)
	account, ok := app.Accounts.FindName(stateAccount + "-OneLedger")
	if ok != status.SUCCESS {
		log.Fatal("Recently Added Account is missing", "name", stateAccount, "status", ok)
	}

	// Use the account key in the database
	balance := NewBalanceFromString(stateAmount, "OLT")
	app.Utxo.Set(account.AccountKey(), balance)

	// TODO: Until a block is commited, this data is not persistent
	//app.Utxo.Commit()

	log.Info("Genesis State UTXO database", "balance", balance)
}

func NewBalanceFromString(amount string, currency string) data.Balance {
	value := big.NewInt(0)
	value.SetString(amount, 10)
	coin := data.Coin{
		Currency: data.NewCurrency(currency),
		Amount:   value,
	}
	if !coin.IsValid() {
		log.Fatal("Create Invalid Coin", "coin", coin)
	}
	return data.Balance{Amount: coin}
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
		return ResponseCheckTx{Code: status.PARSE_ERROR}
	}

	signedTransaction, err := action.Parse(action.Message(tx))
	if err != 0 {
		return ResponseCheckTx{Code: err}
	}

	if action.ValidateSignature(signedTransaction) == false {
		return ResponseCheckTx{Code: status.INVALID_SIGNATURE}
	}

	transaction := signedTransaction.Transaction

	// Check that this is a valid transaction
	if err = transaction.Validate(); err != 0 {
		return ResponseCheckTx{Code: err}
	}

	// Check that this transaction works in the context
	if err = transaction.ProcessCheck(&app); err != 0 {
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

	app.MakePayment(req)

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

// EndBlock is called at the end of all of the transactions
func (app Application) MakePayment(req RequestBeginBlock) {
	account, err := app.Accounts.FindName("Payment-OneLedger")
	if err != status.SUCCESS {
		log.Fatal("ABCI: BeginBlock Fatal Status", "status", err)
	}
	balance := app.Utxo.Get(account.AccountKey())
	if balance == nil {
		interimBalance := data.NewBalance(0, "OLT")
		balance = &interimBalance
	}

	if !balance.Amount.LessThanEqual(0) {
		list := req.LastCommitInfo.GetValidators()
		badValidators := req.ByzantineValidators

		numberValidators := data.NewCoin(int64(len(list)), "OLT")
		quotient := balance.Amount.Quotient(numberValidators)

		var identities []id.Identity

		if int(quotient.Amount.Int64()) > 0 {
			for _, entry := range list {
				entryIsBad := IsByzantine(entry.Validator, badValidators)
				if !entryIsBad {
					formatted := hex.EncodeToString(entry.Validator.Address)
					identity := app.Identities.FindTendermint(formatted)
					identities = append(identities, identity)
				}
			}

			result := CreatePaymentRequest(app, identities, quotient)
			if result != nil {
				// TODO: check this later
				comm.BroadcastAsync(result)
			}
		}
	}

}

func IsByzantine(validator types.Validator, badValidators []types.Evidence) (result bool) {
	for _, entry := range badValidators {
		if bytes.Equal(validator.Address, entry.Validator.Address) {
			return true
		}
	}
	return false
}

// DeliverTx accepts a transaction and updates all relevant data
func (app Application) DeliverTx(tx []byte) ResponseDeliverTx {
	log.Debug("ABCI: DeliverTx", "tx", tx)

	signedTransaction, err := action.Parse(action.Message(tx))
	if err != 0 {
		return ResponseDeliverTx{Code: err}
	}

	if action.ValidateSignature(signedTransaction) == false {
		return ResponseDeliverTx{Code: status.INVALID_SIGNATURE}
	}

	transaction := signedTransaction.Transaction

	log.Debug("Validating")
	if err = transaction.Validate(); err != 0 {
		return ResponseDeliverTx{Code: err}
	}

	log.Debug("Starting processing")
	if transaction.ShouldProcess(app) {
		if err = transaction.ProcessDeliver(&app); err != 0 {
			return ResponseDeliverTx{Code: err}
		}
	}
	tagType := strconv.FormatInt(int64(transaction.TransactionType()), 10)
	tags := make([]common.KVPair, 1)
	tag := common.KVPair{
		Key:   []byte("tx.type"),
		Value: []byte(tagType),
	}
	tags = append(tags, tag)

	return ResponseDeliverTx{
		Code:      types.CodeTypeOK,
		Data:      []byte("Data"),
		Log:       "Log Data",
		Info:      "Info Data",
		GasWanted: 1000,
		GasUsed:   1000,
		Tags:      tags,
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

// Close closes every datastore in app
func (app Application) Close() {
	app.Admin.Close()
	app.Status.Close()
	app.Identities.Close()
	app.Accounts.Close()
	app.Utxo.Close()
	app.Event.Close()
	app.Contract.Close()

	if app.SDK != nil {
		app.SDK.Stop()
	}
}
