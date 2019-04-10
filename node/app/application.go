/*
	Copyright 2017-2018 OneLedger

	AN ABCi application node to process transactions from Tendermint Consensus
*/
package app

import (
	"bytes"
	"encoding/json"
	"time"

	"github.com/Oneledger/protocol/node/abci"
	"github.com/Oneledger/protocol/node/action"
	"github.com/Oneledger/protocol/node/data"
	"github.com/Oneledger/protocol/node/global"
	"github.com/Oneledger/protocol/node/id"
	"github.com/Oneledger/protocol/node/log"
	"github.com/Oneledger/protocol/node/serial"
	"github.com/Oneledger/protocol/node/status"
	"github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/common"
)

var _ types.Application = Application{}

var ChainId string

func init() {
	// TODO: Should be driven from config
	ChainId = "OneLedger-Root"
}

// ApplicationContext keeps all of the upper level global values.
type Application struct {
	types.BaseApplication

	// Global Chain state (data is identical on all nodes in the chain)
	Balances         *data.ChainState // unspent transction output (for each type of coin)
	Identities       *id.Identities   // Keep a higher-level identity for a given user
	SmartContract    data.Datastore   //Store olvm smart contracts
	ExecutionContext data.Datastore   //Store last olvm execution

	// Local Node state (data is different for each node)
	Admin    data.Datastore // any administrative parameters
	Accounts *id.Accounts   // Keep all of the user accounts locally for their node (identity management)
	Sequence data.Datastore // Store sequence number per account
	Status   data.Datastore // current state of any composite transactions (pending, verified, etc.)
	Contract data.Datastore // contract for reuse.
	Event    data.Datastore // Event for any action that need to be tracked

	SDK common.Service

	// Tendermint's last block information
	Header     types.Header   // Tendermint last header info
	Validators *id.Validators // List of validators for this block
}

// NewApplicationContext initializes a new application, reconnects to the databases.
func NewApplication() *Application {
	return &Application{
		Balances:         data.NewChainState("balances", data.PERSISTENT),
		Identities:       id.NewIdentities("identities"),
		SmartContract:    data.NewDatastore("smartContract", data.PERSISTENT),
		ExecutionContext: data.NewDatastore("executionContext", data.PERSISTENT),

		Admin:    data.NewDatastore("admin", data.PERSISTENT),
		Accounts: id.NewAccounts("accounts"),
		Sequence: data.NewDatastore("sequence", data.PERSISTENT),
		Status:   data.NewDatastore("status", data.PERSISTENT),
		Contract: data.NewDatastore("contract", data.PERSISTENT),
		Event:    data.NewDatastore("event", data.PERSISTENT),

		Validators: id.NewValidatorList(),
	}
}

type AdminParameters struct {
	NodeAccountName string
	NodeName        string
}

func init() {
	serial.Register(AdminParameters{})
}

func (app Application) CheckIfInitialized() bool {
	if app.GetPassword() == nil {
		return false
	}

	return true
}

func (app Application) GetPassword() interface{} {
	return app.Admin.Get(data.DatabaseKey("Password"))
}

// Initial the state of the application from persistent data
func (app Application) Initialize() {

	// This config parameter is driven from the database, not the file or cli
	raw := app.Admin.Get(data.DatabaseKey("NodeAccountName"))
	if raw != nil {
		params := raw.(AdminParameters)
		global.Current.NodeAccountName = params.NodeAccountName
	} else {
		log.Debug("NodeAccountName not currently set")
	}

	app.StartSDK()
	log.Debug("SDK is started")

	StartOLVM()
	log.Debug("OLVM is started")
}

// Start up a local server for direct connections from clients
func (app Application) StartSDK() {

	// SDK Server should start when the --sdkrpc argument is passed to fullnode
	sdkAddress := global.Current.Config.Network.SDKAddress

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

// TODO: Not used anymore
type State struct {
	Amount   string `json:"amount"`
	Currency string `json:"currency"`
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

	CreateAccount(app, state, publicKey, privateKey, nil)

	// TODO: Make a user put in a real key
	privateKey, publicKey = id.GenerateKeys([]byte(global.Current.PaymentAccount), false)

	states := []State{
		State{Amount: "0", Currency: "OLT"},
	}
	CreateAccount(app, &BasicState{global.Current.PaymentAccount, states}, publicKey, privateKey, nil)
}

// TODO: DEBUG
var ZeroAccountKey id.AccountKey

func CreateAccount(app Application, state *BasicState, publicKey id.PublicKeyED25519, privateKey id.PrivateKeyED25519, chainkey interface{}) {

	// TODO: This should probably only occur on the Admin node, for other nodes how do I know the key?
	// Register the identity and account first
	AddAccount(&app, state.Account, data.ONELEDGER, publicKey, privateKey, chainkey, false)

	account, ok := app.Accounts.FindName(state.Account)

	if ok != status.SUCCESS {
		log.Fatal("Recently Added Account is missing", "name", state.Account, "status", ok)
	}

	// Use the account key in the database
	balance := NewBalanceFromStates(state.States)

	app.Balances.Set(account.AccountKey(), balance)
	if account.Name() == "Zero" {
		ZeroAccountKey = account.AccountKey()
	}

	// TODO: Until a block is commited, this data should not be persistent
	//app.Balances.Commit()

	log.Info("Genesis State Balances database", "name", state.Account, "balance", balance)
}

func NewBalanceFromStates(states []State) *data.Balance {
	var balance *data.Balance
	for i, value := range states {
		if i == 0 {
			balance = data.NewBalanceFromString(value.Amount, value.Currency)
		} else {
			coin := data.NewCoinFromString(value.Amount, value.Currency)
			balance.AddAmount(coin)
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

func ParseData(message []byte) map[string]interface{} {
	result := map[string]interface{}{
		"parameters": message,
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
	var transaction action.Transaction

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
			transaction = signedTransaction.Transaction
			if err = transaction.Validate(); err != status.SUCCESS {
				errorCode = err

			} else if err = transaction.ProcessCheck(&app); err != status.SUCCESS {
				errorCode = err
			}
		}
	}

	outputData := ""

	if transaction != nil {
		data, error := json.Marshal(transaction.GetData())
		if error != nil {
			log.Warn("transaction get data error", "error", error)
		} else {
			outputData = string(data)
		}
	}

	result := ResponseCheckTx{
		Code: errorCode,

		Data: []byte(outputData),
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

	votes := req.LastCommitInfo.GetVotes()
	byzantineValidators := req.ByzantineValidators

	app.Validators.Set(app, votes, byzantineValidators, req.Header.LastBlockId.Hash)

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

	result := ResponseBeginBlock{
		Tags: []common.KVPair(nil),
	}

	log.Debug("ABCI: BeginBlock Result", "result", result)
	return result
}

// make payment to validators
func (app *Application) MakePayment(req RequestBeginBlock) {
	log.Debug("MakePayment")

	account, err := app.Accounts.FindName(global.Current.PaymentAccount)
	if err != status.SUCCESS {
		log.Fatal("ABCI: BeginBlock Fatal Status", "status", err)
	}

	paymentBalance := app.Balances.Get(account.AccountKey(), true)
	if paymentBalance == nil {
		paymentBalance = data.NewBalance()
	}

	paymentRecordBlockHeight := int64(-1)
	height := int64(req.Header.Height)

	raw := app.Admin.Get(data.DatabaseKey("PaymentRecord"))
	if raw != nil {
		params := raw.(action.PaymentRecord)
		paymentRecordBlockHeight = params.BlockHeight

		log.Debug("Checking for stall", "height", height, "recHeight", paymentRecordBlockHeight)
		if paymentRecordBlockHeight != -1 {
			numTrans := height - paymentRecordBlockHeight
			if numTrans > 10 {
				//store payment record in database (O OLT, -1) because delete doesn't work
				amount := data.NewCoinFromInt(0, "OLT")
				app.SetPaymentRecord(amount, -1)
				paymentRecordBlockHeight = -1
			}
		}
	} else {
		log.Debug("Database uninitialized", "height", height, "recHeight", paymentRecordBlockHeight)
	}

	if (!paymentBalance.GetCoinByName("OLT").LessThanEqual(0)) && paymentRecordBlockHeight == -1 {
		if len(app.Validators.Approved) < 1 || app.Validators.SelectedValidator.Name == "" {
			log.Debug("Missing Validator Information")
			return
		}

		approvedValidatorIdentities := app.Validators.Approved
		selectedValidatorIdentity := app.Validators.SelectedValidator

		numberValidators := len(approvedValidatorIdentities)
		quotient := paymentBalance.GetCoinByName("OLT").Divide(numberValidators)

		if int(quotient.Amount.Int64()) > 0 {
			//store payment record in database
			totalPayment := quotient.MultiplyInt(numberValidators)
			app.SetPaymentRecord(totalPayment, height)

			// if global.Current.NodeName == selectedValidatorIdentity.NodeName {
			nodeAccount, err := app.Accounts.FindName(global.Current.NodeAccountName)

			if err == status.SUCCESS {
				if bytes.Compare(nodeAccount.AccountKey(), selectedValidatorIdentity.AccountKey) == 0 {
					result := CreatePaymentRequest(*app, quotient, height)
					if result != nil {
						// TODO: check this later
						log.Debug("Issuing Payment", "result", result)
						action.DelayedTransaction(result, 0*time.Second)
					}
				} else {
					log.Debug("Payment happens on a different node", "node",
						selectedValidatorIdentity.Name, "validator", selectedValidatorIdentity)
				}
			} else {
				log.Debug("Missing Node Account")
			}
		} else {
			log.Debug("Nothing to Pay")
		}
	} else {
		log.Debug("Not ready for Payment")
	}
}

func (app *Application) SetPaymentRecord(amount data.Coin, blockHeight int64) {
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

	tags := transaction.TransactionTags(app)

	data, error := json.Marshal(transaction.GetData())

	outputData := ""

	if error != nil {
		log.Warn("transaction get data error", "error", error)
	} else {
		outputData = string(data)
	}

	log.Debug("transaction type", "type", transaction.GetType())
	log.Debug("transaction data", "data", data)
	log.Debug("transaction output data", "output", outputData)

	result := ResponseDeliverTx{
		Code:      errorCode,
		Data:      []byte(outputData),
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
	updates := make([]id.Validator, 0)
	if req.Height > 1 && (len(app.Validators.NewValidators) > 0 || len(app.Validators.ToBeRemoved) > 0) {

		for _, validator := range app.Validators.ApprovedValidators {
			found := false
			for _, validatorToBePurged := range app.Validators.ToBeRemoved {
				if bytes.Compare(validator.Address, validatorToBePurged.Validator.Address) == 0 {
					found = true
					err := action.TransferVT(app, validatorToBePurged)
					if err != status.SUCCESS {
						log.Info("Remove Validator - error in transfer of VT")
					}
					break
				}
			}
			if found == true {
				validator.Power = 0
			}
			updates = append(updates, validator)
		}

		for _, applyValidator := range app.Validators.NewValidators {
			if id.HasValidatorToken(app, applyValidator.Validator) {
				updates = append(updates, applyValidator.Validator)
				err := action.TransferVT(app, applyValidator)
				if err != status.SUCCESS {
					log.Info("New Validator - error in transfer of VT")
				}
			} else {
				log.Info("Reject validator", "validatorPubKey", applyValidator.Validator)
			}

		}

	}

	validatorFinalUpdates := make([]types.ValidatorUpdate, 0)
	for _, validator := range updates {
		validatorUpdate := types.ValidatorUpdate{
			PubKey: validator.PubKey,
			Power:  validator.Power,
		}
		validatorFinalUpdates = append(validatorFinalUpdates, validatorUpdate)
	}
	result := ResponseEndBlock{
		ValidatorUpdates: validatorFinalUpdates,
		Tags:             []common.KVPair(nil),
	}

	log.Debug("ABCI: EndBlock Result", "result", result)
	return result
}

// Commit tells the app to make everything persistent
func (app Application) Commit() ResponseCommit {
	log.Debug("ABCI: Commit")

	// Commit any pending changes.
	hash, version := app.Balances.Commit()
	//log.Dump("ZERO IS NOW", app.Balances.Get(ZeroAccountKey))

	log.Debug("-- Committed New Block", "hash", hash, "version", version)

	result := ResponseCommit{
		Data: hash,
	}

	log.Debug("ABCI: EndBlock Result", "result", result)
	return result
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
	app.SmartContract.Close()

	if app.SDK != nil {
		app.SDK.Stop()
	}
}
