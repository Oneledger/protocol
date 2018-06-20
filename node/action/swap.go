/*
	Copyright 2017-2018 OneLedger

	An incoming transaction, send, swap, ready, verification, etc.
*/
package action

import (
	"bytes"

	"github.com/Oneledger/protocol/node/chains/bitcoin"
	"github.com/Oneledger/protocol/node/comm"
	"github.com/Oneledger/protocol/node/data"
	"github.com/Oneledger/protocol/node/err"
	"github.com/Oneledger/protocol/node/global"
	"github.com/Oneledger/protocol/node/id"
	"github.com/Oneledger/protocol/node/log"
	"github.com/Oneledger/protocol/node/chains/ethereum"
	"github.com/Oneledger/protocol/node/chains/ethereum/htlc"
	"github.com/ethereum/go-ethereum/common"
)

// Synchronize a swap between two users
type Swap struct {
	Base

	Party        id.AccountKey `json:"party"`
	CounterParty id.AccountKey `json:"counter_party"`
	Amount       data.Coin     `json:"amount"`
	Exchange     data.Coin     `json:"exchange"`
	Fee          data.Coin     `json:"fee"`
	Gas          data.Coin     `json:"fee"`
	Nonce        int64         `json:"nonce"`
	Preimage     []byte        `json:"preimage"`
}

// Ensure that all of the base values are at least reasonable.
func (transaction *Swap) Validate() err.Code {
	log.Debug("Validating Swap Transaction")

	if transaction.Party == nil {
		log.Debug("Missing Party")
		return err.MISSING_DATA
	}

	if transaction.CounterParty == nil {
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

	if ProcessSwap(app, transaction) {
		log.Debug("Expanding the Transaction into Functions")
		commands := transaction.Expand(app)

		transaction.Resolve(app, commands)

		for i := 0; i < commands.Count(); i++ {
			status := Execute(app, commands[i])
			if status != err.SUCCESS {
				log.Error("Failed to Execute", "command", commands[i])
				return err.EXPAND_ERROR
			}
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
func FindMatchingSwap(status *data.Datastore, accountKey id.AccountKey, transaction *Swap) *Swap {

	result := FindSwap(status, accountKey)
	if result != nil {
		entry := result.(*Swap)
		if MatchSwap(entry, transaction) {
			return entry
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
	if bytes.Compare(left.Party, right.CounterParty) != 0 {
		log.Debug("Party/CounterParty is wrong")
		return false
	}
	if bytes.Compare(left.CounterParty, right.Party) != 0 {
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

func ProcessSwap(app interface{}, transaction *Swap) bool {
	status := GetStatus(app)
	account := transaction.GetNodeAccount(app)

	isParty := transaction.IsParty(account)

	if isParty == nil {
		log.Debug("No Account", "account", account)
		return false
	}

	if *isParty {
		otherSide := FindMatchingSwap(status, transaction.CounterParty, transaction)
		if otherSide != nil {
			return true
		} else {
			SaveSwap(status, transaction.CounterParty, transaction)
			log.Debug("Not Ready", "account", account)
			return false
		}
	} else {
		otherSide := FindMatchingSwap(status, transaction.Party, transaction)
		if otherSide != nil {
			return true

		} else {
			SaveSwap(status, transaction.Party, transaction)
			log.Debug("Not Ready", "account", account)
			return false
		}
	}

	log.Debug("Not Involved", "account", account)
	return false
}

func SaveSwap(status *data.Datastore, accountKey id.AccountKey, transaction *Swap) {
	log.Debug("SaveSwap", "key", accountKey)
	buffer, _ := comm.Serialize(transaction)
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
	if bytes.Compare(transaction.Party, account.AccountKey()) == 0 {
		isParty = true
		return &isParty
	}

	if bytes.Compare(transaction.CounterParty, account.AccountKey()) == 0 {
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

	utxo := GetUtxo(app)
	_ = utxo

	var iindex, pindex int

	chains := GetChains(transaction)
	for i := 0; i < len(commands); i++ {
		isParty := swap.IsParty(account)
		if *isParty {
			commands[i].Chain = chains[0]
			iindex = 0
			pindex = 1
		} else {
			commands[i].Chain = chains[1]
			iindex = 1
			pindex = 0
		}

		role := PARTICIPANT
		if *isParty {
			role = INITIATOR
		}

		commands[i].Data[ROLE] = role
		commands[i].Data[INITIATOR_ACCOUNT] = chains[iindex]
		commands[i].Data[PARTICIPANT_ACCOUNT] = chains[pindex]

		commands[i].Data[AMOUNT] = swap.Amount
		commands[i].Data[EXCHANGE] = swap.Exchange
		commands[i].Data[NONCE] = swap.Nonce
		commands[i].Data[PREIMAGE] = swap.Preimage

		commands[i].Data[PASSWORD] = "password" // TODO: Needs to be corrected
	}
	return
}

// Execute the function
func Execute(app interface{}, command Command) err.Code {
	if command.Execute() {
		return err.SUCCESS
	}
	return err.NOT_IMPLEMENTED
}

func CreateContractBTC(context map[Parameter]FunctionValue) bool {
	address := global.Current.BTCAddress

	role := GetRole(context[ROLE])
	password := GetString(context[PASSWORD])

	_ = role
	_ = password

	cli := bitcoin.GetBtcClient(address)
	_ = cli
	//todo: runCommand(initCmd,cli)

	return true
}

func CreateContractETH(context map[Parameter]FunctionValue) (bool, *htlc.Htlc, common.Address) {
	cli := ethereum.GetEthClient()

	_ = htlc.DeployHtlc(,cli,)


	return true,
}

func CreateContractOLT(context map[Parameter]FunctionValue) bool {
	return true
}


func ParticipateETH(context map[Parameter]FunctionValue) bool {
	return true
}