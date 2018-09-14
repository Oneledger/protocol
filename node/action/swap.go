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

	"crypto/rand"
	"crypto/sha256"
	"github.com/Oneledger/protocol/node/chains/common"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"reflect"
	"time"
)
type swapStageType int32

const (
    MATCHING  swapStageType = iota
    INITIATOR_INITIATE
    PARTICIPANT_PARTICIPATE
    INITIATOR_REDEEM
    PARTICIPANT_REDEEM
    WAIT_FOR_CHAIN
    REFUND_SWAP
    FINISH_SWAP
)

// Synchronize a swap between two users
type Swap struct {
	Base
    Message Message `json:"message"`
	Stage   string   `json:"stage"`
}

var swapStageFlow = map[swapStageType]SwapStage{
    MATCHING: {
        Stage:MATCHING,
        Commands: Commands{Command{opfunc: SaveUnmatchSwap}, Command{opfunc: NextStage}},
        InStage: nil,
        OutStage: WAIT_FOR_CHAIN,
    },
    INITIATOR_INITIATE: {
        Stage: INITIATOR_INITIATE,
        Commands: Commands{Command{opfunc: Initiate}, Command{opfunc: NextStage}},
        InStage: nil,
        OutStage: PARTICIPANT_PARTICIPATE,
    },
    PARTICIPANT_PARTICIPATE: {
        Stage: PARTICIPANT_PARTICIPATE,
        Commands: Commands{Command{opfunc: AuditContract}, Command{opfunc: Participate}, Command{opfunc: NextStage}},
        InStage: INITIATOR_INITIATE,
        OutStage: INITIATOR_REDEEM,
    },
    INITIATOR_REDEEM: {
        Stage: INITIATOR_REDEEM,
        Commands: Commands{Command{opfunc: AuditContract}, Command{opfunc: Redeem}, Command{opfunc: NextStage}},
        InStage: PARTICIPANT_PARTICIPATE,
        OutStage: PARTICIPANT_REDEEM,
    },
    PARTICIPANT_REDEEM: {
        Stage: PARTICIPANT_REDEEM,
        Commands: Commands{Command{opfunc: ExtractSecret}, Command{opfunc: Redeem}, Command{opfunc: NextStage}},
        InStage: INITIATOR_REDEEM,
        OutStage: FINISH_SWAP,
    },
    FINISH_SWAP: {
        Stage: FINISH_SWAP,
        Commands: Commands{Command{opfunc: FinalizeSwap}},
        InStage: PARTICIPANT_REDEEM,
        OutStage: nil,

    },
    WAIT_FOR_CHAIN: {
      Stage: WAIT_FOR_CHAIN,
      Commands: Commands{Command{opfunc: VerifySwap}},
      InStage: MATCHING,
      OutStage: nil,
    },
    REFUND_SWAP: {
        Stage: REFUND_SWAP,
        Commands: Commands{Command{opfunc: Refund}},
        InStage: WAIT_FOR_CHAIN,
        OutStage: nil,
    },
}


type SwapStage struct {
    Stage       swapStageType
    Commands    Commands

    InStage     swapStageType
    OutStage    swapStageType
}


type SwapInit struct {
    Party        Party     `json:"party"`
    CounterParty Party     `json:"counter_party"`
    Amount       data.Coin `json:"amount"`
    Exchange     data.Coin `json:"exchange"`
    Fee          data.Coin `json:"fee"`
    Gas          data.Coin `json:"fee"`
    Nonce        int64     `json:"nonce"`
    Preimage     []byte    `json:"preimage"`
}

func (si *SwapInit) Marshal() Message {

    return nil
}

func (si *SwapInit) UnMarshal(message Message) {

}

type SwapPublish struct {
    Contract   Message       `json:"message"` //message converted from HTLContract
    SecretHash [32]byte      `json:"secrethash"`
    Count      int           `json:"count"`
}

type SwapVerify struct {
    Event   Event           `json:"event"`
    Message Message         `json:"Message"`
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

	if !transaction.Amount.IsCurrency("BTC", "ETH", "OLT") {
		log.Debug("Swap on Currency isn't implement yet")
		return err.NOT_IMPLEMENTED
	}

	if !transaction.Exchange.IsCurrency("BTC", "ETH", "OLT") {
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
		var lastResult FunctionValues

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
		    SaveSwap(status,  transaction.CounterParty.Key, matchedSwap)
			return matchedSwap
		} else {
			SaveSwap(status, transaction.CounterParty.Key, transaction)
			log.Debug("Not Ready", "account", account)
			return nil
		}
	} else {
		matchedSwap := FindMatchingSwap(status, transaction.Party.Key, transaction, *isParty)
		if matchedSwap != nil {
		    SaveSwap(status,  matchedSwap.Party.Key, matchedSwap)
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

// Plug in data from the rest of a system into a set of commands
func (swap *Swap) Resolve(app interface{}) Commands {
	chains := swap.getChains()

	account := GetNodeAccount(app)
	isParty := swap.IsParty(account)

	role := swap.getRole(*isParty)

	commands := make(Commands, 0)

	append(commands, )
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
		    commands[i].data[MY_ACCOUNT] = swap.Party
		    commands[i].data[THEM_ACCOUNT] = swap.CounterParty
			commands[i].data[AMOUNT] = swap.Amount
			commands[i].data[EXCHANGE] = swap.Exchange
            key = swap.CounterParty.Key
		} else {
		    commands[i].data[MY_ACCOUNT] = swap.CounterParty
		    commands[i].data[THEM_ACCOUNT] = swap.Party
			commands[i].data[AMOUNT] = swap.Exchange
			commands[i].data[EXCHANGE] = swap.Amount
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


func Initiate(app interface{}, chain data.ChainType, context FunctionValues) (bool, FunctionValues) {
    log.Info("Executing Initiate Command", "chain", chain, "context", context)
    switch chain {

    case data.BITCOIN:
        return CreateContractBTC(app, context)
    case data.ETHEREUM:
        return CreateContractETH(app, context)
    case data.ONELEDGER:
        return CreateContractOLT(app, context)
    default:
        return false, nil
    }
}

func Participate(app interface{}, chain data.ChainType, context FunctionValues) (bool, FunctionValues) {
    log.Info("Executing Participate Command", "chain", chain, "context", context)
    switch chain {

    case data.BITCOIN:
        return ParticipateBTC(app, context)
    case data.ETHEREUM:
        return ParticipateETH(app, context)
    case data.ONELEDGER:
        return ParticipateOLT(app, context)
    default:
        return false, nil
    }
}

func Redeem(app interface{}, chain data.ChainType, context FunctionValues) (bool, FunctionValues) {
    log.Info("Executing Redeem Command", "chain", chain, "context", context)

    switch chain {

    case data.BITCOIN:
        return RedeemBTC(app, context)
    case data.ETHEREUM:
        return RedeemETH(app, context)
    case data.ONELEDGER:
        return RedeemOLT(app, context)
    default:
        return false, nil
    }
}

func Refund(app interface{}, chain data.ChainType, context FunctionValues) (bool, FunctionValues) {
    log.Info("Executing Refund Command", "chain", chain, "context", context)
    switch chain {

    case data.BITCOIN:
        return RefundBTC(app, context)
    case data.ETHEREUM:
        return RefundETH(app, context)
    case data.ONELEDGER:
        return RefundOLT(app, context)
    default:
        return false, nil
    }
}

func ExtractSecret(app interface{}, chain data.ChainType, context FunctionValues) (bool, FunctionValues) {
    log.Info("Executing ExtractSecret Command", "chain", chain, "context", context)
    switch chain {

    case data.BITCOIN:
        return ExtractSecretBTC(app, context)
    case data.ETHEREUM:
        return ExtractSecretETH(app, context)
    case data.ONELEDGER:
        return ExtractSecretOLT(app, context)
    default:
        return false, nil
    }
}

func AuditContract(app interface{}, chain data.ChainType, context FunctionValues) (bool, FunctionValues) {
    log.Info("Executing AuditContract Command", "chain", chain, "context", context)
    switch chain {

    case data.BITCOIN:
        return AuditContractBTC(app, context)
    case data.ETHEREUM:
        return AuditContractETH(app, context)
    case data.ONELEDGER:
        return AuditContractOLT(app, context)
    default:
        return false, nil
    }
}

func CreateContractBTC(app interface{}, context FunctionValues) (bool, FunctionValues) {

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

func CreateContractETH(app interface{}, context FunctionValues) (bool, FunctionValues) {
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
	if context[PASSWORD] != nil {
		scr := GetByte32(context[PASSWORD])
		scrHash := sha256.Sum256(scr[:])
		if !bytes.Equal(preimage[:], scrHash[:]) {
			log.Error("Secret and Secret Hash doesn't match", "preimage", preimage, "scrHash", scrHash)
			return false, nil
		}
	}

    value := GetCoin(context[AMOUNT]).Amount

    receiverParty = GetParty(context[THEM_ACCOUNT])
	receiver := common.GetETHAddressFromByteArray(data.ETHEREUM,receiverParty.Accounts[data.ETHEREUM])
    if receiver == nil {
        log.Error("Failed to get eth address from string", "address", receiverParty.Accounts[data.ETHEREUM], "target", reflect.TypeOf(receiver))
    }

	timeoutSecond := int64(lockPeriod.Seconds())
	log.Debug("Create ETH HTLC", "value", value, "receiver", receiver, "preimage", preimage)
	err := contract.Funds(value, big.NewInt(timeoutSecond), *receiver, preimage)
	if err != nil {
		return false, nil
	}

	context[ETHCONTRACT] = contract
	return true, context
}

func AuditContractBTC(app interface{}, context FunctionValues) (bool, FunctionValues) {
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

func AuditContractETH(app interface{}, context FunctionValues) (bool, FunctionValues) {
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

func ParticipateBTC(app interface{}, context FunctionValues) (bool, FunctionValues) {
    success, result := CreateContractBTC(app, context)
    if success != false {
        log.Error("failed to participate because can't create contract")
        return false, nil
    }
    return true, result
}

func ParticipateETH(app interface{}, context FunctionValues) (bool, FunctionValues) {
	success, result := CreateContractETH(app, context)
	if success == false {
		log.Error("failed to participate because can't create contract")
		return false, nil
	}
	return true, result
}

func RedeemBTC(app interface{}, context FunctionValues) (bool, FunctionValues) {
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

func RedeemETH(app interface{}, context FunctionValues) (bool, FunctionValues) {
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

func RefundBTC(app interface{}, context FunctionValues) (bool, FunctionValues) {

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

func RefundETH(app interface{}, context FunctionValues) (bool, FunctionValues) {
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

func ExtractSecretBTC(app interface{}, context FunctionValues) (bool, FunctionValues) {
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

func ExtractSecretETH(app interface{}, context FunctionValues) (bool, FunctionValues) {
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

func CreateContractOLT(app interface{}, context FunctionValues) (bool, FunctionValues) {
    log.Warn("Not supported")
    party := GetParty(context[MY_ACCOUNT])
    counterParty := GetParty(context[THEM_ACCOUNT])
    partyBalance := GetUtxo(app).Find(party.Key).Amount
    counterPartyBalance := GetUtxo(app).Find(counterParty.Key).Amount

	preimage := GetByte32(context[PREIMAGE])
	if context[PASSWORD] != nil {
		scr := GetByte32(context[PASSWORD])
		scrHash := sha256.Sum256(scr[:])
		if !bytes.Equal(preimage[:], scrHash[:]) {
			log.Error("Secret and Secret Hash doesn't match", "preimage", preimage, "scrHash", scrHash)
			return false, nil
		}
	}

	inputs := make([]SendInput, 0)
    inputs = append(inputs,
        NewSendInput(party.Key, partyBalance),
        NewSendInput(counterParty.Key, counterPartyBalance))
    amount := GetCoin(context[AMOUNT])
    // Build up the outputs
    outputs := make([]SendOutput, 0)
    outputs = append(outputs,
        NewSendOutput(party.Key, partyBalance.Minus(amount)),
        NewSendOutput(counterParty.Key, counterPartyBalance.Plus(amount)))
    send := &Send{
        Base: Base{
            Type:     SEND,
            ChainId:  GetChainID(app),
            Signers:  nil,
            Sequence: global.Current.Sequence,
        },
        Inputs:  inputs,
        Outputs: outputs,
        Fee:     data.NewCoin(0, "OLT"),
        Gas:     data.NewCoin(0, "OLT"),
    }
    message := SignAndPack(SEND, Transaction(send))
    contract := NewMultiSigBox(1,1, message)
    _ = contract
    return true, context
}

func ParticipateOLT(app interface{}, context FunctionValues) (bool, FunctionValues) {
    log.Warn("Not supported")
    return true, context
}

func AuditContractOLT(app interface{}, context FunctionValues) (bool, FunctionValues) {
    log.Warn("Not supported")
    return true, context
}


func RedeemOLT(app interface{}, context FunctionValues) (bool, FunctionValues) {
    log.Warn("Not supported")
    return true, context
}


func RefundOLT(app interface{}, context FunctionValues) (bool, FunctionValues) {
    log.Warn("Not supported")
    return true, context
}


func ExtractSecretOLT(app interface{}, context FunctionValues) (bool, FunctionValues) {
    log.Warn("Not supported")
    return true, context
}