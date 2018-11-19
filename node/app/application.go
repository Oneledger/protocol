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
	"github.com/Oneledger/protocol/node/data"
	"github.com/Oneledger/protocol/node/global"
	"github.com/Oneledger/protocol/node/id"
	"github.com/Oneledger/protocol/node/log"
	"github.com/Oneledger/protocol/node/serial"
	"github.com/Oneledger/protocol/node/status"
	"github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/common"
	"math/big"
	"strconv"
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
	Validators ValidatorList
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
	Account string  `json:"account"`
	States  []State `json:"states"`
}

type State struct {
	Amount string `json:"amount"`
	Coin   string `json:"coin"`
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
	log.Debug("Deserialized State", "state", state)

	// TODO: Can't generate a different key for each node. Needs to be in the genesis? Or ignored?
	privateKey, publicKey := id.GenerateKeys([]byte(state.Account), false) // TODO: switch with passphrase

	CreateAccount(app, state, publicKey, privateKey)

	privateKey, publicKey = id.GenerateKeys([]byte("Payment"), false) // TODO: make a user put a real key actually
	CreateAccount(app, &BasicState{"Payment", []State{State{"0", "OLT"}}}, publicKey, privateKey)
}

func CreateAccount(app Application, state *BasicState, publicKey id.PublicKeyED25519, privateKey id.PrivateKeyED25519) {

	// TODO: This should probably only occur on the Admin node, for other nodes how do I know the key?
	// Register the identity and account first
	RegisterLocally(&app, state.Account, "OneLedger", data.ONELEDGER, publicKey, privateKey)
	account, ok := app.Accounts.FindName(state.Account + "-OneLedger")
	if ok != status.SUCCESS {
		log.Fatal("Recently Added Account is missing", "name", state.Account, "status", ok)
	}

	// Use the account key in the database
	balance := NewBalanceFromStates(state.States)
	app.Utxo.Set(account.AccountKey(), balance)

	// TODO: Until a block is commited, this data is not persistent
	//app.Utxo.Commit()

	log.Info("Genesis State UTXO database", "balance", balance)
}

func NewBalanceFromStates(states []State) data.Balance {
	var balance data.Balance
	for i, v := range states {
		if i == 0 {
			value := big.NewInt(0)
			value.SetString(v.Amount, 10)
			balance = data.NewBalanceFromString(value.Int64(), v.Coin)
		} else {
			value := big.NewInt(0)
			value.SetString(v.Amount, 10)
			coin := data.NewCoin(value.Int64(), v.Coin)
			balance.AddAmmount(coin)
		}
	}

	return balance
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
	log.Debug("ABCI: BeginBlock", "req", req)

	validators := req.LastCommitInfo.GetValidators()
	byzantineValidators := req.ByzantineValidators

	app.Validators.Set(validators, byzantineValidators)

	raw := app.Admin.Get(data.DatabaseKey("PaymentRecord"))
	if raw == nil {
		app.MakePayment(req)
	} else {
		params := raw.(action.PaymentRecord)
		if params.BlockHeight == -1 {
			app.MakePayment(req)
		}
	}

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
	//log.Debug("MakePayment", "req", req)
	account, err := app.Accounts.FindName("Payment-OneLedger")
	if err != status.SUCCESS {
		log.Fatal("ABCI: BeginBlock Fatal Status", "status", err)
	}

	paymentBalance := app.Utxo.Get(account.AccountKey())
	if paymentBalance == nil {
		interimBalance := data.NewBalance()
		paymentBalance = &interimBalance
	}

	paymentRecordBlockHeight := int64(-1)
	height := int64(req.Header.Height)

	raw := app.Admin.Get(data.DatabaseKey("PaymentRecord"))
	if raw != nil {
		params := raw.(action.PaymentRecord)
		paymentRecordBlockHeight = params.BlockHeight

		if paymentRecordBlockHeight != -1 {
			numTrans := height - paymentRecordBlockHeight
			if numTrans > 10 {
				//store payment record in database (O OLT, -1) because delete doesn't work
				amount := data.NewCoin(0, "OLT")
				SetPaymentRecord(amount, -1, app)
				paymentRecordBlockHeight = -1
			}
		}
	}

	if (!paymentBalance.GetAmountByName("OLT").LessThanEqual(0)) && paymentRecordBlockHeight == -1 {
		goodValidatorIdentities := app.Validators.FindGood(app)
		selectedValidatorIdentity := app.Validators.FindSelectedValidator(app, req.Header.LastBlockHash)

		numberValidators := data.NewCoin(int64(len(goodValidatorIdentities)), "OLT")
		quotient := paymentBalance.GetAmountByName("OLT").Quotient(numberValidators)

		if int(quotient.Amount.Int64()) > 0 {
			//store payment record in database
			totalPayment := quotient.Multiply(numberValidators)
			SetPaymentRecord(totalPayment, height, app)

			if global.Current.NodeName == selectedValidatorIdentity.NodeName {
				result := CreatePaymentRequest(app, goodValidatorIdentities, quotient, height)
				if result != nil {
					// TODO: check this later
					comm.BroadcastAsync(result)
				}
			}
		}
	}
}

func SetPaymentRecord(amount data.Coin, blockHeight int64, app Application) {
	var paymentRecordKey data.DatabaseKey = data.DatabaseKey("PaymentRecord")
	var paymentRecord action.PaymentRecord
	paymentRecord.Amount = amount
	paymentRecord.BlockHeight = blockHeight
	session := app.Admin.Begin()
	session.Set(paymentRecordKey, paymentRecord)
	session.Commit()
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
