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

	"math/big"

	"crypto/sha256"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"time"
	"crypto/rand"
	"github.com/Oneledger/protocol/node/chains/common"
    "reflect"
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
	Key      id.AccountKey				`json:"key"`
	Accounts map[data.ChainType][]byte	`json:"accounts"`
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
		    them := GetParty(commands[i].Data[THEM_ACCOUNT])
            event := Event{Type: SWAP, Key: them.Key , Nonce: matchedSwap.Nonce }
            if commands[i].Function == FINISH {
                SaveEvent(app, event, true)
                return err.SUCCESS
            }

			status, result := Execute(app, commands[i], lastResult)
			if status != err.SUCCESS {
				log.Error("Failed to Execute", "command", commands[i])
				SaveEvent(app, event, false)
				return err.EXPAND_ERROR
			}

			lastResult = result
		}
	} else {
		log.Debug("Not Involved or Not Ready")
	}

	return err.SUCCESS
}

// TODO: Change to return Role as INITIATOR or PARTICIPANT
func FindMatchingSwap(status *data.Datastore, accountKey id.AccountKey, transaction *Swap, isParty bool) (matched *Swap) {

	result := FindSwap(status, accountKey)
	if result != nil {
		entry := result.(*Swap)
		if matching := MatchSwap(entry, transaction); matching {
			log.Debug("MatchSwap", "matching", matching, "transaction", transaction, "entry", entry, "isParty", isParty)
			var base Swap
			matched = &base
			if isParty {
                matched.Base = entry.Base //put them as base for easy access the key to store
				matched.Party = transaction.Party
				matched.CounterParty = entry.Party
				matched.Amount = transaction.Amount
				matched.Exchange = transaction.Exchange
			} else {
                matched.Base = transaction.Base //put them as base for easy access the key to store
				matched.Party = entry.Party
				matched.CounterParty = transaction.Party
				matched.Amount = transaction.Exchange
				matched.Exchange = transaction.Amount
			}

			matched.Fee = transaction.Fee
			matched.Nonce = transaction.Nonce

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
	account := GetNodeAccount(app)

	isParty := transaction.IsParty(account)

	if isParty == nil {
		log.Debug("No Account", "account", account)
		return nil
	}

	if *isParty {
		matchedSwap := FindMatchingSwap(status, transaction.CounterParty.Key, transaction, *isParty)
		if matchedSwap != nil {
		    log.Debug("Swap is ready", "swap", matchedSwap)
		    if matchedSwap.getRole(*isParty) == PARTICIPANT {
		        SaveSwap(status,  transaction.CounterParty.Key, matchedSwap)
            }
			return matchedSwap
		} else {
			SaveSwap(status, transaction.CounterParty.Key, transaction)
			log.Debug("Not Ready", "account", account)
			return nil
		}
	} else {
		matchedSwap := FindMatchingSwap(status, transaction.Party.Key, transaction, *isParty)
		if matchedSwap != nil {
            if matchedSwap.getRole(*isParty) == PARTICIPANT {
                SaveSwap(status,  matchedSwap.Party.Key, matchedSwap)
            }
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


// Is this node one of the partipants in the swap
func (transaction *Swap) ShouldProcess(app interface{}) bool {
	account := GetNodeAccount(app)

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

// Get the correct chains order for this action
func (swap *Swap) getChains() []data.ChainType {

	var first, second data.ChainType
	if swap.Amount.Currency.Id < swap.Exchange.Currency.Id {
		first = data.Currencies[swap.Amount.Currency.Name].Chain
		second = data.Currencies[swap.Exchange.Currency.Name].Chain
	} else {
		first = data.Currencies[swap.Exchange.Currency.Name].Chain
		second = data.Currencies[swap.Amount.Currency.Name].Chain
	}

	return []data.ChainType{data.ONELEDGER, first, second}
}

func (swap *Swap) getRole(isParty bool) Role {

	if swap.Amount.Currency.Id < swap.Exchange.Currency.Id {
		if isParty {
			return INITIATOR
		} else {
			return PARTICIPANT
		}
	} else {
		if isParty {
			return PARTICIPANT
		} else {
			return INITIATOR
		}
	}
}

// Given a transaction, expand it into a list of Commands to execute against various chains.
func (transaction *Swap) Expand(app interface{}) Commands {
	chains := transaction.getChains()

	account := GetNodeAccount(app)
	isParty := transaction.IsParty(account)

	role := transaction.getRole(*isParty)

	return GetCommands(SWAP, role, chains)
}

// Plug in data from the rest of a system into a set of commands
func (swap *Swap) Resolve(app interface{}, commands Commands) {
	account := GetNodeAccount(app)

	identities := GetIdentities(app)
	_ = identities
	name := global.Current.NodeIdentity
	_ = name

	utxo := GetUtxo(app)
	_ = utxo

	chains := swap.getChains()
	isParty := swap.IsParty(account)
	role := swap.getRole(*isParty)

	for i := 0; i < len(commands); i++ {

		side := commands[i].Order
		log.Info("side", "s", side)
		if &side == nil {
			commands[i].Chain = data.ONELEDGER
		} else {
            commands[i].Chain = chains[side]
        }

		var key id.AccountKey

		if *isParty {
		    commands[i].Data[MY_ACCOUNT] = swap.Party
		    commands[i].Data[THEM_ACCOUNT] = swap.CounterParty
			commands[i].Data[AMOUNT] = swap.Amount
			commands[i].Data[EXCHANGE] = swap.Exchange
            key = swap.CounterParty.Key
		} else {
		    commands[i].Data[MY_ACCOUNT] = swap.CounterParty
		    commands[i].Data[THEM_ACCOUNT] = swap.Party
			commands[i].Data[AMOUNT] = swap.Exchange
			commands[i].Data[EXCHANGE] = swap.Amount
            key = swap.Party.Key
		}

		commands[i].Data[ROLE] = role
		commands[i].Data[NONCE] = swap.Nonce

		if role == INITIATOR {
		    if commands[i].Function == INITIATE {

                var secret [32]byte
                _, err := rand.Read(secret[:])
                if err != nil {
                    log.Error("failed to get random secret with 32 length", "err", err)
                }
                secretHash := sha256.Sum256(secret[:])
                commands[i].Data[PASSWORD] = secret
                commands[i].Data[PREIMAGE] = secretHash
                tokens[key.String()] = secret
            } else {
                secret := tokens[key.String()]
                commands[i].Data[PASSWORD] = secret
                commands[i].Data[PREIMAGE] = sha256.Sum256(secret[:])
            }
		}
		commands[i].Data[EVENTTYPE] = SWAP
	}
	return
}

// TODO: Needs to be configurable
var lockPeriod = 5 * time.Minute
// todo: need to store this in db
var tokens = make(map[string][32]byte)

func CreateContractBTC(app interface{}, context map[Parameter]FunctionValue) (bool, map[Parameter]FunctionValue) {

	timeout := time.Now().Add(2 * lockPeriod).Unix()

    value := GetCoin(context[AMOUNT]).Amount

    receiverParty := GetParty(context[THEM_ACCOUNT])
    receiver := common.GetBTCAddressFromByteArray(data.BITCOIN, receiverParty.Accounts[data.BITCOIN])
    if receiver == nil {
        log.Error("Failed to get btc address from string", "address", receiverParty.Accounts[data.BITCOIN], "target", reflect.TypeOf(receiver))
        return false, nil
    }

    preimage := GetByte32(context[PREIMAGE])
    if context[PASSWORD] != nil {
        scr := GetByte32(context[PASSWORD])
        scrHash := sha256.Sum256(scr[:])
        if !bytes.Equal(preimage[:], scrHash[:]) {
            log.Error("Secret and Secret Hash doesn't match", "preimage", preimage, "scrHash", scrHash)
            return false, nil
        }
    }

	cli := bitcoin.GetBtcClient(global.Current.BTCAddress)

	amount := bitcoin.GetAmount(value.String())

	initCmd := htlc.NewInitiateCmd(receiver, amount, timeout, preimage)

	_, err := initCmd.RunCommand(cli)
	if err != nil {
		log.Error("Bitcoin Initiate", "err", err, "context", context)
		return false, nil
	}

    contract := &bitcoin.HTLContract{Contract: initCmd.Contract, ContractTx: *initCmd.ContractTx}

	context[BTCCONTRACT] = contract

	nonce := GetInt64(context[NONCE])
	SaveContract(app, receiverParty.Key.Bytes(), nonce , contract)
	log.Debug("btc contract","contract", context[BTCCONTRACT])
	return true, context
}

func CreateContractETH(app interface{}, context map[Parameter]FunctionValue) (bool, map[Parameter]FunctionValue) {
    me := GetParty(context[MY_ACCOUNT])

    contractMessage := FindContract(app, me.Key.Bytes(), 0)
    var contract *ethereum.HTLContract
    if contractMessage == nil {
        contract = ethereum.CreateHtlContract()
        SaveContract(app, me.Key.Bytes(), 0, contract)
    } else {
        contract = ethereum.GetHTLCFromMessage(contractMessage)
    }

	var receiverParty Party
    preimage := GetByte32(context[PREIMAGE])

    value := GetCoin(context[AMOUNT]).Amount

    receiverParty = GetParty(context[THEM_ACCOUNT])
	receiver := common.GetETHAddressFromByteArray(data.ETHEREUM,receiverParty.Accounts[data.ETHEREUM])
    if receiver == nil {
        log.Error("Failed to get eth address from string", "address", receiverParty.Accounts[data.ETHEREUM], "target", reflect.TypeOf(receiver))
    }

	timeoutSecond := int64(lockPeriod.Seconds())
	log.Debug("Create ETH HTLC", "value", value, "receiver", receiver, "preimage", preimage)
	err := contract.Funds(value)
	if err != nil {
		return false, nil
	}

	err = contract.Setup(big.NewInt(timeoutSecond), *receiver, preimage)
	if err != nil {
		return false, nil
	}

	context[ETHCONTRACT] = contract
	return true, context
}


func AuditContractBTC(app interface{}, context map[Parameter]FunctionValue) (bool, map[Parameter]FunctionValue) {
	contract := GetBTCContract(context[BTCCONTRACT])

	them := GetParty(context[THEM_ACCOUNT])
	cmd := htlc.NewAuditContractCmd(contract.Contract, &contract.ContractTx)
	cli := bitcoin.GetBtcClient(global.Current.BTCAddress)
	e := cmd.RunCommand(cli)
	if e != nil {
		log.Error("Bitcoin Audit", "err", e)
		return false, nil
	}

	nonce := GetInt64(context[NONCE])
	SaveContract(app, them.Key.Bytes(), nonce, contract)
	context[PREIMAGE] = cmd.SecretHash
	return true, context
}

func AuditContractETH(app interface{}, context map[Parameter]FunctionValue) (bool, map[Parameter]FunctionValue) {
	contract := GetETHContract(context[ETHCONTRACT])

	scrHash := GetByte32(context[PREIMAGE])
	address := ethereum.GetAddress()

	receiver, e := contract.HTLContractObject().Receiver(&bind.CallOpts{Pending: true})
	if e != nil {
		log.Error("can't get the receiver", "err", e)
		return false, nil
	}
	if !bytes.Equal(address.Bytes(), receiver.Bytes()) {
        log.Error("receiver not correct",  "contract", contract.Address, "receiver", receiver, "myAddress", address)
        return false, nil
    }

    value := GetCoin(context[EXCHANGE]).Amount

    setVale := contract.Balance()
    setScrhash := contract.ScrHash()
    if !bytes.Equal(scrHash[:], setScrhash[:]) {
        log.Error("Secret Hash doesn't match", "sh", scrHash, "setSh", setScrhash)
        return false, nil
    }

    if value.Cmp(setVale) != 0 {
        log.Error("Value doesn't match", "value", value, "setValue", setVale)
        return false, nil
    }

    log.Debug("Auditing ETH Contract","receiver", address, "value", value, "scrHash", scrHash)


    log.Debug("Set ETH Contract","receiver", receiver, "value", setVale, "scrHash", setScrhash)
	//e = contract.Audit(address, value ,scrHash)
	//if e != nil {
	//	log.Error("Failed to audit the contract with correct input", "err", e)
	//	return false, nil
	//}

	context[PREIMAGE] = scrHash
	them := GetParty(context[THEM_ACCOUNT])
	nonce := GetInt64(context[NONCE])
	SaveContract(app, them.Key.Bytes(), nonce, contract)
	return true, context
}

func ParticipateBTC(app interface{}, context map[Parameter]FunctionValue) (bool, map[Parameter]FunctionValue) {
    success, result := CreateContractBTC(app, context)
    if success != false {
        log.Error("failed to participate because can't create contract")
        return false, nil
    }
    return true, result
}

func ParticipateETH(app interface{}, context map[Parameter]FunctionValue) (bool, map[Parameter]FunctionValue) {
	success, result := CreateContractETH(app, context)
	if success == false {
		log.Error("failed to participate because can't create contract")
		return false, nil
	}
	return true, result
}

func RedeemBTC(app interface{}, context map[Parameter]FunctionValue) (bool, map[Parameter]FunctionValue) {
    them := GetParty(context[THEM_ACCOUNT])
    nonce := GetInt64(context[NONCE])
    contractMessage := FindContract(app, them.Key.Bytes(), nonce)
    if contractMessage == nil {
        return false,nil
    }
    contract := bitcoin.GetHTLCFromMessage(contractMessage)
    if contract == nil {
        log.Error("BTC Htlc contract not found")
        return false, nil
    }

    scr := GetByte32(context[PASSWORD])

    cmd := htlc.NewRedeemCmd(contract.Contract, &contract.ContractTx, scr[:])
    cli := bitcoin.GetBtcClient(global.Current.BTCAddress)
    _, e := cmd.RunCommand(cli)
    if e != nil {
        log.Error("Bitcoin redeem htlc", "err", e)
        return false, nil
    }
    context[BTCCONTRACT] = &bitcoin.HTLContract{Contract: contract.Contract, ContractTx: *cmd.RedeemContractTx}
    return true, context
}

func RedeemETH(app interface{}, context map[Parameter]FunctionValue) (bool, map[Parameter]FunctionValue) {
    them := GetParty(context[THEM_ACCOUNT])
    nonce := GetInt64(context[NONCE])
    contractMessage := FindContract(app, them.Key.Bytes(), nonce)
    if contractMessage == nil {
        return false, nil
    }
    contract := ethereum.GetHTLCFromMessage(contractMessage)
    if contract == nil {
        return false, nil
    }

	scr := GetByte32(context[PASSWORD])
	err := contract.Redeem(scr[:])
	if err != nil {
	    log.Error("Ethereum redeem htlc", "err", err)
	    return false, nil
    }
	context[ETHCONTRACT] = contract
	return true, context
}

func RefundBTC(app interface{}, context map[Parameter]FunctionValue) (bool, map[Parameter]FunctionValue) {

    them := GetParty(context[THEM_ACCOUNT])
    nonce := GetInt64(context[NONCE])
    contractMessage := FindContract(app, them.Key.Bytes(), nonce)
    if contractMessage == nil {
        log.Error("BTC Htlc contract not found")
        return false,nil
    }
    contract := bitcoin.GetHTLCFromMessage(contractMessage)
    if contract == nil {
        log.Error("BTC Htlc contract can't be parsed")
        return false, nil
    }

    cmd := htlc.NewRefundCmd(contract.Contract, &contract.ContractTx)
    cli := bitcoin.GetBtcClient(global.Current.BTCAddress)
    _, e := cmd.RunCommand(cli)
    if e != nil {
        log.Error("Bitcoin refund htlc", "err", e)
        return false, nil
    }
    return true, context
}

func RefundETH(app interface{}, context map[Parameter]FunctionValue) (bool, map[Parameter]FunctionValue) {
    me := GetParty(context[MY_ACCOUNT])
    contractMessage := FindContract(app, me.Key.Bytes(), 0)
    if contractMessage == nil {
        log.Error("ETH Htlc contract not found")
        return false, nil
    }
    contract := ethereum.GetHTLCFromMessage(contractMessage)
    if contract == nil {
        log.Error("ETH Htlc contract can't be parsed")
        return false, nil
    }

	err := contract.Refund()
	if err != nil {
	    return false, nil
    }
	context[ETHCONTRACT] = contract
	return true, context
}

func ExtractSecretBTC(app interface{}, context map[Parameter]FunctionValue) (bool, map[Parameter]FunctionValue) {
    contract := GetBTCContract(context[BTCCONTRACT])
    scrHash := GetByte32(context[PREIMAGE])
    cmd := htlc.NewExtractSecretCmd(&contract.ContractTx, scrHash)
    cli := bitcoin.GetBtcClient(global.Current.BTCAddress)
    e := cmd.RunCommand(cli)
    if e != nil {
        log.Error("Bitcoin extract hltc", "err", e)
        return false, nil
    }
    var tmpScr [32]byte
    copy(tmpScr[:], string(cmd.Secret))
    context[PASSWORD] = tmpScr
    log.Debug("extracted secret", "secretbytearray", cmd.Secret, "secretbyte32", tmpScr)
    return true, context

}

func ExtractSecretETH(app interface{}, context map[Parameter]FunctionValue) (bool, map[Parameter]FunctionValue) {
    contract := GetETHContract(context[ETHCONTRACT])
    //todo: make it correct scr, by extract or from local storage

    scr := contract.Extract()
    if scr == nil {
        return false, nil
    }
    var tmpScr [32]byte
    copy(tmpScr[:], string(scr))
    context[PASSWORD] = tmpScr
    log.Debug("extracted secret", "secret", scr, "r", tmpScr)
    return true, context
}

func CreateContractOLT(app interface{}, context map[Parameter]FunctionValue) (bool, map[Parameter]FunctionValue) {
    log.Warn("Not supported")
    return true, context
}

func ParticipateOLT(app interface{}, context map[Parameter]FunctionValue) (bool, map[Parameter]FunctionValue) {
    log.Warn("Not supported")
    return true, context
}

func AuditContractOLT(app interface{}, context map[Parameter]FunctionValue) (bool, map[Parameter]FunctionValue) {
    log.Warn("Not supported")
    return true, context
}


func RedeemOLT(app interface{}, context map[Parameter]FunctionValue) (bool, map[Parameter]FunctionValue) {
    log.Warn("Not supported")
    return true, context
}


func RefundOLT(app interface{}, context map[Parameter]FunctionValue) (bool, map[Parameter]FunctionValue) {
    log.Warn("Not supported")
    return true, context
}


func ExtractSecretOLT(app interface{}, context map[Parameter]FunctionValue) (bool, map[Parameter]FunctionValue) {
    log.Warn("Not supported")
    return true, context
}