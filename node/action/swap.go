/*
	Copyright 2017-2018 OneLedger

	An incoming transaction, send, swap, ready, verification, etc.
*/
package action

import (
	"bytes"

	"github.com/Oneledger/protocol/node/chains/bitcoin"
	"github.com/Oneledger/protocol/node/chains/bitcoin/htlc"
	"github.com/Oneledger/protocol/node/chains/ethereum"
	"github.com/Oneledger/protocol/node/comm"
	"github.com/Oneledger/protocol/node/data"
	"github.com/Oneledger/protocol/node/err"
	"github.com/Oneledger/protocol/node/global"
	"github.com/Oneledger/protocol/node/id"
	"github.com/Oneledger/protocol/node/log"
	"github.com/btcsuite/btcd/chaincfg"
	bwire "github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"

	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"crypto/sha256"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

// Synchronize a swap between two users
type Swap struct {
	Base

	Party        Party     `json:"party"`
	CounterParty Party     `json:"counter_party"`
	Amount       data.Coin `json:"amount"`
	Exchange     data.Coin `json:"exchange"`
	Fee          data.Coin `json:"fee"`
	Gas          data.Coin `json:"fee"`
	Nonce        int64     `json:"nonce"`
	Preimage     []byte    `json:"preimage"`
}

type Party struct {
	Key      id.AccountKey
	Accounts map[data.ChainType]string
}

// Ensure that all of the base values are at least reasonable.
func (transaction *Swap) Validate() err.Code {
	log.Debug("Validating Swap Transaction")

	if transaction.Party.Key == nil {
		log.Debug("Missing Party")
		return err.MISSING_DATA
	}

	if transaction.CounterParty.Key == nil {
		log.Debug("Missing CounterParty")
		return err.MISSING_DATA
	}

	if !transaction.Amount.IsCurrency("BTC", "ETH") {
		log.Debug("Swap on Currency isn't implement yet")
		return err.NOT_IMPLEMENTED
	}

	if !transaction.Exchange.IsCurrency("BTC", "ETH") {
		log.Debug("Swap on Currency isn't implement yet")
		return err.NOT_IMPLEMENTED
	}

	log.Debug("Swap is validated!")
	return err.SUCCESS
}

func (transaction *Swap) ProcessCheck(app interface{}) err.Code {
	log.Debug("Processing Swap Transaction for CheckTx")

	// TODO: Check all of the data to make sure it is valid.

	return err.SUCCESS
}

// Start the swap
func (transaction *Swap) ProcessDeliver(app interface{}) err.Code {
	log.Debug("Processing Swap Transaction for DeliverTx")
	matchedSwap := ProcessSwap(app, transaction)
	if matchedSwap != nil {
		log.Debug("Expanding the Transaction into Functions")
		commands := matchedSwap.Expand(app)

		matchedSwap.Resolve(app, commands)

		//before loop of execute, lastResult is nil
		var lastResult map[Parameter]FunctionValue

		for i := 0; i < commands.Count(); i++ {
			status, result := Execute(app, commands[i], lastResult)
			if status != err.SUCCESS {
				log.Error("Failed to Execute", "command", commands[i])
				return err.EXPAND_ERROR
			}
			lastResult = result
		}
	} else {
		log.Debug("Not Involved or Not Ready")
	}

	return err.SUCCESS
}

func FindSwap(status *data.Datastore, key id.AccountKey) Transaction {
	result := status.Load(key)
	if result == nil {
		return nil
	}

	var transaction Swap
	buffer, err := comm.Deserialize(result, &transaction)
	if err != nil {
		return nil
	}
	return buffer.(Transaction)
}

// TODO: Change to return Role as INITIATOR or PARTICIPANT
func FindMatchingSwap(status *data.Datastore, accountKey id.AccountKey, transaction *Swap, isParty bool) (matched *Swap) {

	result := FindSwap(status, accountKey)
	if result != nil {
		entry := result.(*Swap)
		if MatchSwap(entry, transaction) {
			if isParty {
				matched.Party = transaction.Party
				matched.CounterParty = entry.Party
				matched.Amount = transaction.Amount
				matched.Exchange = transaction.Exchange
			} else {
				matched.Party = entry.Party
				matched.CounterParty = transaction.Party
				matched.Amount = transaction.Exchange
				matched.Exchange = transaction.Amount
			}
			matched.Base = transaction.Base
			matched.Fee = transaction.Fee
			matched.Nonce = transaction.Nonce
			//todo: get preimage for this swap
			matched.Preimage = []byte("super cool pre-image")

			return matched
		} else {
			log.Debug("Swap doesn't match", "key", accountKey, "transaction", transaction, "entry", entry)
		}
	} else {
		log.Debug("Swap not found", "key", accountKey)
	}

	return nil
}

// Two matching swap requests from different parties
func MatchSwap(left *Swap, right *Swap) bool {
	if left.Base.Type != right.Base.Type {
		log.Debug("Type is wrong")
		return false
	}
	if left.Base.ChainId != right.Base.ChainId {
		log.Debug("ChainId is wrong")
		return false
	}
	if bytes.Compare(left.Party.Key, right.CounterParty.Key) != 0 {
		log.Debug("Party/CounterParty is wrong")
		return false
	}
	if bytes.Compare(left.CounterParty.Key, right.Party.Key) != 0 {
		log.Debug("CounterParty/Party is wrong")
		return false
	}
	if !left.Amount.Equals(right.Exchange) {
		log.Debug("Amount/Exchange is wrong")
		return false
	}
	if !left.Exchange.Equals(right.Amount) {
		log.Debug("Exchange/Amount is wrong")
		return false
	}
	if left.Nonce != right.Nonce {
		log.Debug("Nonce is wrong")
		return false
	}

	return true
}

func ProcessSwap(app interface{}, transaction *Swap) *Swap {
	status := GetStatus(app)
	account := transaction.GetNodeAccount(app)

	isParty := transaction.IsParty(account)

	if isParty == nil {
		log.Debug("No Account", "account", account)
		return nil
	}

	if *isParty {
		matchedSwap := FindMatchingSwap(status, transaction.CounterParty.Key, transaction, *isParty)
		if matchedSwap != nil {
			return matchedSwap
		} else {
			SaveSwap(status, transaction.CounterParty.Key, transaction)
			log.Debug("Not Ready", "account", account)
			return nil
		}
	} else {
		matchedSwap := FindMatchingSwap(status, transaction.Party.Key, transaction, *isParty)
		if matchedSwap != nil {
			return matchedSwap

		} else {
			SaveSwap(status, transaction.Party.Key, transaction)
			log.Debug("Not Ready", "account", account)
			return nil
		}
	}

	log.Debug("Not Involved", "account", account)
	return nil
}

func SaveSwap(status *data.Datastore, accountKey id.AccountKey, transaction *Swap) {
	log.Debug("SaveSwap", "key", accountKey)
	buffer, err := comm.Serialize(transaction)
	if err != nil {
		log.Error("Failed to Serialize SaveSwap transaction")
	}
	status.Store(accountKey, buffer)
	status.Commit()
}

// Is this node one of the partipants in the swap
func (transaction *Swap) ShouldProcess(app interface{}) bool {
	account := transaction.GetNodeAccount(app)

	if transaction.IsParty(account) != nil {
		return true
	}

	return false
}

func GetAccount(app interface{}, accountKey id.AccountKey) id.Account {
	accounts := GetAccounts(app)
	account, _ := accounts.FindKey(accountKey)

	return account
}

// Map the identity to a specific account on a chain
func GetChainAccount(app interface{}, name string, chain data.ChainType) id.Account {
	identities := GetIdentities(app)
	accounts := GetAccounts(app)

	identity, _ := identities.FindName(name)
	account, _ := accounts.FindKey(identity.Chain[chain])

	return account
}

func (transaction *Swap) GetNodeAccount(app interface{}) id.Account {

	accounts := GetAccounts(app)
	account, _ := accounts.FindName(global.Current.NodeAccountName)
	if account == nil {
		log.Error("Node does not have account", "name", global.Current.NodeAccountName)
		accounts.Dump()
		return nil
	}

	return account
}

func (transaction *Swap) IsParty(account id.Account) *bool {

	if account == nil {
		log.Debug("Getting Role for empty account")
		return nil
	}

	var isParty bool
	if bytes.Compare(transaction.Party.Key, account.AccountKey()) == 0 {
		isParty = true
		return &isParty
	}

	if bytes.Compare(transaction.CounterParty.Key, account.AccountKey()) == 0 {
		isParty = false
		return &isParty
	}

	// TODO: Shouldn't be in-band this way
	return nil
}

// Given a transaction, expand it into a list of Commands to execute against various chains.
func (transaction *Swap) Expand(app interface{}) Commands {
	chains := GetChains(transaction)

	account := transaction.GetNodeAccount(app)
	isParty := transaction.IsParty(account)
	role := PARTICIPANT
	if *isParty {
		role = INITIATOR
	}

	return GetCommands(SWAP, role, chains)
}

// Plug in data from the rest of a system into a set of commands
func (swap *Swap) Resolve(app interface{}, commands Commands) {
	transaction := Transaction(swap)

	account := swap.GetNodeAccount(app)

	identities := GetIdentities(app)
	_ = identities
	name := global.Current.NodeIdentity
	_ = name

	utxo := GetUtxo(app)
	_ = utxo

	var iindex, pindex int

	chains := GetChains(transaction)
	isParty := swap.IsParty(account)
	role := swap.getRole(*isParty)

	for i := 0; i < len(commands); i++ {

		if *isParty {
			commands[i].Chain = chains[0]
			iindex = 0
			pindex = 1
		} else {
			commands[i].Chain = chains[1]
			iindex = 1
			pindex = 0
		}

		_ = iindex
		_ = pindex
		if *isParty {
			if role == INITIATOR {
				commands[i].Data[INITIATOR_ACCOUNT] = swap.Party.Accounts
				commands[i].Data[PARTICIPANT_ACCOUNT] = swap.CounterParty.Accounts
			} else {
				commands[i].Data[INITIATOR_ACCOUNT] = swap.CounterParty.Accounts
				commands[i].Data[PARTICIPANT_ACCOUNT] = swap.Party.Accounts
			}

		} else {
			if role == PARTICIPANT {
				commands[i].Data[INITIATOR_ACCOUNT] = swap.Party.Accounts
				commands[i].Data[PARTICIPANT_ACCOUNT] = swap.CounterParty.Accounts
			} else {
				commands[i].Data[INITIATOR_ACCOUNT] = swap.CounterParty.Accounts
				commands[i].Data[PARTICIPANT_ACCOUNT] = swap.Party.Accounts
			}

		}

		commands[i].Data[ROLE] = role

		commands[i].Data[AMOUNT] = swap.Amount
		commands[i].Data[EXCHANGE] = swap.Exchange
		commands[i].Data[NONCE] = swap.Nonce
		commands[i].Data[PREIMAGE] = swap.Preimage

		commands[i].Data[PASSWORD] = "password" // TODO: Needs to be corrected
	}
	return
}

func (swap *Swap) getRole(isParty bool) Role {

	if isParty {
		if data.Currencies[swap.Amount.Currency] < data.Currencies[swap.Exchange.Currency] {
			return INITIATOR
		} else {
			return PARTICIPANT
		}
	} else {
		if data.Currencies[swap.Exchange.Currency] < data.Currencies[swap.Amount.Currency] {
			return PARTICIPANT
		} else {
			return INITIATOR
		}
	}
}

// Execute the function
func Execute(app interface{}, command Command, lastResult map[Parameter]FunctionValue) (err.Code, map[Parameter]FunctionValue) {
	if status, lastResult := command.Execute(); status {
		return err.SUCCESS, lastResult
	}
	return err.NOT_IMPLEMENTED, lastResult
}

func GetPubKeyHash(address string) *btcutil.AddressPubKeyHash {

	// TODO: Needs to be configurable
	chainParams := &chaincfg.RegressionNetParams
	hash, _ := btcutil.NewAddressPubKeyHash([]byte(address), chainParams)

	return hash
}

// TODO: Needs to be configurable
var timeout int64 = 100000

func CreateContractBTC(context map[Parameter]FunctionValue) (bool, map[Parameter]FunctionValue) {
	btcAddress := global.Current.BTCAddress

	amount := GetAmount(context[AMOUNT])

	var accountKey id.AccountKey
	var client int
	role := GetRole(context[ROLE])
	if role == INITIATOR {
		client = 1
		accountKey = GetAccountKey(context[INITIATOR_ACCOUNT])
	} else {
		client = 2
		accountKey = GetAccountKey(context[PARTICIPANT_ACCOUNT])
	}
	_ = accountKey

	password := GetString(context[PASSWORD])
	_ = password

	config := chaincfg.RegressionNetParams // TODO: should be passed in

	cli := bitcoin.GetBtcClient(btcAddress, client, &config)

	if role == INITIATOR {

		address := GetPubKeyHash(btcAddress)

		_, err := htlc.NewInitiateCmd(address, amount, timeout).RunCommand(cli)
		if err != nil {
			log.Error("Bitcoin Initiate", "err", err)
			return false, nil
		}
	} else {
		contract := []byte(nil)
		contractTx := &bwire.MsgTx{}
		_ = htlc.NewAuditContractCmd(contract, contractTx).RunCommand(cli)
	}

	return true, nil
}

func CreateContractETH(context map[Parameter]FunctionValue) (bool, map[Parameter]FunctionValue) {

	contract := ethereum.GetHtlContract()
	role := GetRole(context[ROLE])
	var value = big.NewInt(0)
	var receiver common.Address
	if role == INITIATOR {
		value = GetCoin(context[AMOUNT]).Amount
		receiverParty := GetParty(context[PARTICIPANT_ACCOUNT])
		receiver = common.BytesToAddress([]byte(receiverParty.Accounts[data.ETHEREUM]))
	} else if role == PARTICIPANT {
		value = GetCoin(context[EXCHANGE]).Amount
		receiverParty := GetParty(context[INITIATOR_ACCOUNT])
		receiver = common.BytesToAddress([]byte(receiverParty.Accounts[data.ETHEREUM]))
	}
	scr := GetBytes(context[PASSWORD])
	scrHash := sha256.Sum256([]byte(scr))



	contract.Funds(value)
	contract.Setup(big.NewInt(25*3600), receiver, scrHash)


	context[ETHCONTRACT] = contract
	return true, context
}

func CreateContractOLT(context map[Parameter]FunctionValue) (bool, map[Parameter]FunctionValue) {
	log.Warn("Not supported")
	return true, nil
}

func ParticipateETH(context map[Parameter]FunctionValue) (bool, map[Parameter]FunctionValue) {
	address := ethereum.GetAddress()
	contract := GetETHContract(context[ETHCONTRACT])
	scrHash, err := contract.Contract.ScrHash(&bind.CallOpts{Pending: true})
	if err != nil {
		log.Error("can't get the secret Hash", "contract", contract.Address, "err", err)
	}

	locktime, err := contract.Contract.LockPeriod(&bind.CallOpts{Pending: true})
	if err != nil {
		log.Error("can't get the lock period", "contract", contract.Address, "err", err)
	}
	_ = scrHash
	_ = locktime
	receiver, err := contract.Contract.Receiver(&bind.CallOpts{Pending: true})
	if err != nil || receiver != address  {
		log.Error("can't get the receiver or receiver not correct", "err", err, "contract", contract.Address, "receiver", receiver, "my address", address)
	}

	success, result := CreateContractBTC(context)
	if success != false {
		log.Error("failed to participate because can't create contract")
	}

	return true, result
}


func RedeemETH(context map[Parameter]FunctionValue) (bool, map[Parameter]FunctionValue) {
	contract := GetETHContract(context[ETHCONTRACT])
	//todo: make it correct scr, by extract or from local storage
	scr := []byte("my cool secret")
	contract.Redeem(scr)
	context[ETHCONTRACT] = contract
	return true, context
}

func RefundETH(context map[Parameter]FunctionValue) (bool, map[Parameter]FunctionValue) {
	contract := GetETHContract(context[ETHCONTRACT])
	//todo: make it correct scr, by extract or from local storage
	scr := []byte("my cool secret")
	contract.Refund(scr)
	context[ETHCONTRACT] = contract
	return true, context
}