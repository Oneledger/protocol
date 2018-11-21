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
	Balances   *data.ChainState // unspent transction output (for each type of coin)
	Identities *id.Identities   // Keep a higher-level identity for a given user

	// Local Node state (data is different for each node)
	Accounts *id.Accounts   // Keep all of the user accounts locally for their node (identity management)
	Admin    data.Datastore // any administrative parameters
	Event    data.Datastore // Event for any action that need to be tracked
	Status   data.Datastore // current state of any composite transactions (pending, verified, etc.)
	Contract data.Datastore // contract for reuse.
	Sequence data.Datastore // Store sequence number per account

	SDK common.Service

	// Tendermint's last block information
	LastHeader types.Header // Tendermint last header info
	Validators ValidatorList
}

// NewApplicationContext initializes a new application, reconnects to the databases.
func NewApplication() *Application {
	return &Application{
		Identities: id.NewIdentities("identities"),
		Balances:   data.NewChainState("balances", data.PERSISTENT),

		Accounts: id.NewAccounts("accounts"),
		Admin:    data.NewDatastore("admin", data.PERSISTENT),
		Event:    data.NewDatastore("event", data.PERSISTENT),
		Status:   data.NewDatastore("status", data.PERSISTENT),
		Contract: data.NewDatastore("contract", data.PERSISTENT),
		Sequence: data.NewDatastore("sequence", data.PERSISTENT),
	}
}

type AdminParameters struct {
	NodeAccountName string
	NodeName        string
}

func init() {
	serial.Register(AdminParameters{})
}

type SequenceRecord struct {
	Sequence int64
}

func init() {
	serial.Register(SequenceRecord{})
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
	app.StartSDK()
}

// Start up a local server for direct connections from clients
func (app Application) StartSDK() {

	// SDK Server should start when the --sdkrpc argument is passed to fullnode
	sdkAddress := global.Current.SDKAddress

	sdk, err := NewSDKServer(&app, sdkAddress)
	if err != nil {
		log.Fatal("SDK Server Failed", "err", err)
	}

	app.SDK = sdk
	app.SDK.Start()
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

	// Tendermint serializes this data, so we have to use raw JSON serialization to read it.
	des, errx := serial.Deserialize(stateBytes, &base, serial.JSON)
	if errx != nil {
		log.Fatal("Failed to deserialize stateBytes during SetupState")
	}

	state := des.(*BasicState)
	log.Debug("Deserialized State", "state", state)

	// TODO: Can't generate a different key for each node. Needs to be in the genesis? Or ignored?
	privateKey, publicKey := id.GenerateKeys([]byte(state.Account), false) // TODO: switch with passphrase

	CreateAccount(app, state, publicKey, privateKey)

	privateKey, publicKey = id.GenerateKeys([]byte(global.Current.PaymentAccount), false) // TODO: make a user put a real key actually
	CreateAccount(app, &BasicState{global.Current.PaymentAccount, []State{State{"0", "OLT"}}}, publicKey, privateKey)
}

func CreateAccount(app Application, state *BasicState, publicKey id.PublicKeyED25519, privateKey id.PrivateKeyED25519) {

	// TODO: This should probably only occur on the Admin node, for other nodes how do I know the key?
	// Register the identity and account first
	AddAccount(&app, state.Account, data.ONELEDGER, publicKey, privateKey, false)
	//RegisterLocally(&app, stateAccount, "OneLedger", data.ONELEDGER, publicKey, privateKey)

	account, ok := app.Accounts.FindName(state.Account)

	if ok != status.SUCCESS {
		log.Fatal("Recently Added Account is missing", "name", state.Account, "status", ok)
	}

	// Use the account key in the database
	balance := NewBalanceFromStates(state.States)

	app.Balances.Set(account.AccountKey(), balance)

	// TODO: Until a block is commited, this data is not persistent
	//app.Balances.Commit()

	log.Info("Genesis State Balances database", "balance", balance)
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

	result := ResponseInitChain{}
	log.Debug("ABCI: InitChain Result", "result", result)
	return result
}

// SetOption changes the underlying options for the ABCi app
func (app Application) SetOption(req RequestSetOption) ResponseSetOption {
	log.Debug("ABCI: SetOption", "key", req.Key, "value", req.Value)

	SetOption(&app, req.Key, req.Value)

	result := ResponseSetOption{
		Code: types.CodeTypeOK,
		Log:  "Log Data",
		Info: "Info Data",
	}

	log.Debug("ABCI: SetOption Result", "result", result)
	return result
}

// Info returns the current block information
func (app Application) Info(req RequestInfo) ResponseInfo {

	// TODO: Check this...
	info := abci.NewResponseInfo(0, 0, 0)

	log.Debug("ABCI: Info", "req", req, "info", info)

	result := ResponseInfo{
		Data:             info.JSON(),
		LastBlockHeight:  int64(app.Balances.Version),
		LastBlockAppHash: app.Balances.Hash,
	}

	log.Debug("ABCI: Info Result", "result", result)
	return result
}

func ParseData(message []byte) map[string]string {
	result := map[string]string{
		"parameters": string(message),
	}
	return result
}

// Query comes from tendermint node, and returns data and/or a proof
func (app Application) Query(req RequestQuery) ResponseQuery {
	log.Debug("ABCI: Query", "req", req, "path", req.Path, "data", req.Data)

	arguments := ParseData(req.Data)
	response := HandleQuery(app, req.Path, arguments)

	result := ResponseQuery{
		Code:  types.CodeTypeOK,
		Index: 0, // TODO: What is this for?

		Log:  "Log Information",
		Info: "Info Information",

		Key:   action.Message("result"),
		Value: response,

		Proof:  nil,
		Height: int64(app.Balances.Version),
	}

	log.Debug("ABCI: Query Result", "result", result)
	return result
}

// CheckTx tests to see if a transaction is valid
func (app Application) CheckTx(tx []byte) ResponseCheckTx {
	log.Debug("ABCI: CheckTx", "tx", tx)

	errorCode := types.CodeTypeOK

	if tx == nil {
		log.Warn("Empty Transaction, Ignoring", "tx", tx)
		errorCode = status.PARSE_ERROR

	} else {
		signedTransaction, err := action.Parse(action.Message(tx))
		if err != status.SUCCESS {
			errorCode = err

		} else if action.ValidateSignature(signedTransaction) == false {
			errorCode = status.INVALID_SIGNATURE

		} else {
			transaction := signedTransaction.Transaction
			if err = transaction.Validate(); err != status.SUCCESS {
				errorCode = err

			} else if err = transaction.ProcessCheck(&app); err != status.SUCCESS {
				errorCode = err
			}
		}
	}

	result := ResponseCheckTx{
		Code: errorCode,

		Data: []byte("Data"),
		Log:  "Log Data",
		Info: "Info Data",

		GasWanted: 0,
		GasUsed:   0,
		Tags:      []common.KVPair(nil),
	}

	log.Debug("ABCI: CheckTx Result", "result", result)
	return result
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

	result := ResponseBeginBlock{
		Tags: []common.KVPair(nil),
	}

	log.Debug("ABCI: BeginBlock Result", "result", result)
	return result
}

// EndBlock is called at the end of all of the transactions
func (app Application) MakePayment(req RequestBeginBlock) {
	account, err := app.Accounts.FindName("Payment")
	if err != status.SUCCESS {
		log.Fatal("ABCI: BeginBlock Fatal Status", "status", err)
	}

	paymentBalance := app.Balances.Get(account.AccountKey())
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

	errorCode := types.CodeTypeOK
	var transaction action.Transaction

	signedTransaction, err := action.Parse(action.Message(tx))
	if err != status.SUCCESS {
		errorCode = err

	} else if action.ValidateSignature(signedTransaction) == false {
		return ResponseDeliverTx{Code: status.INVALID_SIGNATURE}

	} else {
		transaction = signedTransaction.Transaction
		if err = transaction.Validate(); err != status.SUCCESS {
			errorCode = err

		} else if transaction.ShouldProcess(app) {
			if err = transaction.ProcessDeliver(&app); err != status.SUCCESS {
				errorCode = err
			}
		}
	}

	tags := transaction.TransactionTags()

	result := ResponseDeliverTx{
		Code:      errorCode,
		Data:      []byte("Data"),
		Log:       "Log Data",
		Info:      "Info Data",
		GasWanted: 0,
		GasUsed:   0,
		Tags:      tags,
	}

	log.Debug("ABCI: DeliverTx Result", "result", result)
	return result
}

// EndBlock is called at the end of all of the transactions
func (app Application) EndBlock(req RequestEndBlock) ResponseEndBlock {
	log.Debug("ABCI: EndBlock", "req", req)

	result := ResponseEndBlock{
		Tags: []common.KVPair(nil),
	}

	log.Debug("ABCI: EndBlock Result", "result", result)
	return result
}

// Commit tells the app to make everything persistent
func (app Application) Commit() ResponseCommit {
	log.Debug("ABCI: Commit")

	// Commit any pending changes.
	hash, version := app.Balances.Commit()

	log.Debug("-- Committed New Block", "hash", hash, "version", version)

	result := ResponseCommit{
		Data: hash,
	}

	log.Debug("ABCI: EndBlock Result", "result", result)
	return result
}

func NextSequence(app *Application, accountkey id.AccountKey) SequenceRecord {
	sequence := int64(1)
	raw := app.Sequence.Get(accountkey)
	if raw != nil {
		interim := raw.(SequenceRecord)
		sequence = interim.Sequence + 1
	}

	sequenceRecord := SequenceRecord{
		Sequence: sequence,
	}

	session := app.Sequence.Begin()
	session.Set(accountkey, sequenceRecord)
	session.Commit()

	return sequenceRecord
}

// Close closes every datastore in app
func (app Application) Close() {
	app.Admin.Close()
	app.Status.Close()
	app.Identities.Close()
	app.Accounts.Close()
	app.Balances.Close()
	app.Event.Close()
	app.Contract.Close()
	app.Sequence.Close()

	if app.SDK != nil {
		app.SDK.Stop()
	}
}
