/*
	Copyright 2017-2018 OneLedger

	An incoming transaction, send, swap, ready, verification, etc.
*/
package action

import (
	"bytes"

	"github.com/tendermint/go-amino"

	"github.com/Oneledger/protocol/node/chains/bitcoin"
	"github.com/Oneledger/protocol/node/chains/bitcoin/htlc"
	"github.com/Oneledger/protocol/node/chains/ethereum"
	"github.com/Oneledger/protocol/node/data"
	"github.com/Oneledger/protocol/node/global"
	"github.com/Oneledger/protocol/node/id"
	"github.com/Oneledger/protocol/node/log"
	"github.com/Oneledger/protocol/node/serial"
	"github.com/Oneledger/protocol/node/status"

	"math/big"

	"crypto/rand"
	"crypto/sha256"
	"reflect"
	"time"

	"github.com/Oneledger/protocol/node/chains/common"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

func init() {
	serial.Register(Swap{})
	serial.Register(Party{})
	var SwapMessage SwapMessage
	serial.RegisterInterface(&SwapMessage)
	serial.Register(SwapInit{})
	serial.Register(SwapExchange{})
	serial.Register(SwapVerify{})
}

type swapStageType int

const (
	NOSTAGE swapStageType = iota
	SWAP_MATCHING
	SWAP_MATCHED
	INITIATOR_INITIATE
	PARTICIPANT_PARTICIPATE
	INITIATOR_REDEEM
	PARTICIPANT_REDEEM
	WAIT_FOR_CHAIN
	SWAP_REFUND
	SWAP_FINISH
)

// Synchronize a swap between two users
type Swap struct {
	Base
	SwapMessage SwapMessage   `json:"swapmessage"`
	Stage       swapStageType `json:"stage"`
}

// Ensure that all of the base values are at least reasonable.
func (transaction *Swap) Validate() status.Code {
	log.Debug("Validating Swap Transaction")

	if transaction.SwapMessage == nil {
		log.Error("swap don't contain message")
		return status.MISSING_DATA
	}

	if transaction.SwapMessage.validate() != status.SUCCESS {
		log.Debug("SwapMessage not validate")
		return status.INVALID
	}

	log.Debug("Swap is validated!")
	return status.SUCCESS
}

func (transaction *Swap) ProcessCheck(app interface{}) status.Code {
	log.Debug("Processing Swap Transaction for CheckTx")

	// TODO: Check all of the data to make sure it is valid.

	return status.SUCCESS
}

// Is this node one of the partipants in the swap
func (transaction *Swap) ShouldProcess(app interface{}) bool {
	account := GetNodeAccount(app)
	if account == nil {
		log.Warn("NodeAccount not setup for processing")
		return false
	}

	if bytes.Equal(transaction.Base.Target, account.AccountKey()) {
		log.Debug("Swap involved", "swap", transaction.SwapMessage, "stage", transaction.Stage)
		return true
	}

	if bytes.Equal(transaction.Base.Owner, account.AccountKey()) && transaction.Stage == SWAP_MATCHING {
		log.Debug("Swap involved", "swap", transaction.SwapMessage, "stage", transaction.Stage)
		return true
	}

	log.Debug("Swap not involved", "me", account, "owner", transaction.Base.Owner, "target", transaction.Base.Target, "stage", transaction.Stage)

	return false
}

// Start the swap
func (transaction *Swap) ProcessDeliver(app interface{}) status.Code {
	log.Debug("Processing Swap Transaction for DeliverTx")
	commands := transaction.Resolve(app)

	if commands.Count() == 0 {
		return status.EXPAND_ERROR
	}

	c := make(chan status.Code)
	go transaction.CommandExecute(app, commands, c)

	return status.SUCCESS
}

// Plug in data from the rest of a system into a set of commands
func (swap *Swap) Resolve(app interface{}) Commands {
	log.Debug("Resolve Swap", "stage", swap.Stage)
	var commands Commands
	if swap.Stage == SWAP_MATCHING {
		si := swap.SwapMessage.(SwapInit)
		stageType := matchingSwap(app, &si)
		commands = si.resolve(app, stageType)

	} else {
		// after the swap_matching, every swapmessage should be swapexchange or swapverify.
		commands = swap.SwapMessage.resolve(app, swap.Stage)
		if commands.Count() == 0 {
			log.Error("Swap resolve no commands", "swap", swap)
			return nil
		}
		//todo: the owner and signer should change when support wallet/light client.
		// the owner should be the wallet/light client sender not the operation node. so as the signer
		commands[0].data[PREVIOUS] = _hash(swap)
	}
	return commands
}

func (swap *Swap) CommandExecute(app interface{}, commands Commands, c chan status.Code) {
	defer close(c)
	for i := 0; i < commands.Count(); i++ {

		context := commands[i].data
		if v, _ := context[NEXTCHAINNAME]; v != nil {
			chain := GetChain(v)
			commands[i].chain = chain
		}
		ok, result := commands[i].Execute(app)
		if !ok {
			log.Error("Failed to Execute", "command", commands[i])
			//SaveEvent(app, event, false)
			c <- status.EXPAND_ERROR
		}
		if len(result) > 0 {
			commands[i+1].data = result
		}
	}
	c <- status.SUCCESS
}

var swapStageFlow = map[swapStageType]swapStage{
	SWAP_MATCHING: {
		Stage:    SWAP_MATCHING,
		Commands: Commands{Command{opfunc: NextStage}},
		InStage:  NOSTAGE,
		OutStage: WAIT_FOR_CHAIN,
	},
	SWAP_MATCHED: {
		Stage:    SWAP_MATCHED,
		Commands: Commands{Command{opfunc: CreateCheckEvent}, Command{opfunc: NextStage}},
		InStage:  SWAP_MATCHING,
		OutStage: INITIATOR_INITIATE,
	},
	INITIATOR_INITIATE: {
		Stage:    INITIATOR_INITIATE,
		Commands: Commands{Command{opfunc: Initiate}, Command{opfunc: NextStage}},
		InStage:  SWAP_MATCHED,
		OutStage: PARTICIPANT_PARTICIPATE,
	},
	PARTICIPANT_PARTICIPATE: {
		Stage:    PARTICIPANT_PARTICIPATE,
		Commands: Commands{Command{opfunc: AuditContract}, Command{opfunc: Participate}, Command{opfunc: NextStage}},
		InStage:  INITIATOR_INITIATE,
		OutStage: INITIATOR_REDEEM,
	},
	INITIATOR_REDEEM: {
		Stage:    INITIATOR_REDEEM,
		Commands: Commands{Command{opfunc: AuditContract}, Command{opfunc: Redeem}, Command{opfunc: NextStage}},
		InStage:  PARTICIPANT_PARTICIPATE,
		OutStage: PARTICIPANT_REDEEM,
	},
	PARTICIPANT_REDEEM: {
		Stage:    PARTICIPANT_REDEEM,
		Commands: Commands{Command{opfunc: ExtractSecret}, Command{opfunc: Redeem}, Command{opfunc: NextStage}},
		InStage:  INITIATOR_REDEEM,
		OutStage: SWAP_FINISH,
	},
	SWAP_FINISH: {
		Stage:    SWAP_FINISH,
		Commands: Commands{Command{opfunc: FinalizeSwap}},
		InStage:  PARTICIPANT_REDEEM,
		OutStage: NOSTAGE,
	},
	WAIT_FOR_CHAIN: {
		Stage:    WAIT_FOR_CHAIN,
		Commands: Commands{Command{opfunc: VerifySwap}, Command{opfunc: ClearEvent}},
		InStage:  SWAP_MATCHING,
		OutStage: NOSTAGE,
	},
	SWAP_REFUND: {
		Stage:    SWAP_REFUND,
		Commands: Commands{Command{opfunc: Refund}},
		InStage:  WAIT_FOR_CHAIN,
		OutStage: NOSTAGE,
	},
}

type swapStage struct {
	Stage    swapStageType
	Commands Commands

	InStage  swapStageType
	OutStage swapStageType
}

type SwapMessage interface {
	validate() status.Code
	resolve(interface{}, swapStageType) Commands
}

type Party struct {
	Key      id.AccountKey             `json:"key"`
	Accounts map[data.ChainType][]byte `json:"accounts"`
}

type SwapKey struct {
	Initiator   id.AccountKey `json:"initiator"`
	Participant id.AccountKey `json:"participant"`
	Amount      data.Coin     `json:"amount"`
	Exchange    data.Coin     `json:"exchange"`
	Nonce       int64         `json:"nonce"`
}

func (sk SwapKey) toHash() []byte {
	return _hash(sk)
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

func (si SwapInit) validate() status.Code {
	log.Debug("Validating SwapInit")
	if si.Party.Key == nil {
		log.Debug("Missing Party")
		return status.MISSING_DATA
	}

	if si.CounterParty.Key == nil {
		log.Debug("Missing CounterParty")
		return status.MISSING_DATA
	}

	if !si.Amount.IsCurrency("BTC", "ETH", "OLT") {
		log.Debug("Swap on Currency isn't implement yet")
		return status.NOT_IMPLEMENTED
	}

	if !si.Exchange.IsCurrency("BTC", "ETH", "OLT") {
		log.Debug("Swap on Currency isn't implement yet")
		return status.NOT_IMPLEMENTED
	}

	return status.SUCCESS
}

func (si SwapInit) store(app interface{}) []byte {
	sk := si.getKey()
	storeKey := sk.toHash()
	log.Debug("Store SwapInit", "key", storeKey, "si", si)
	SaveSwap(app, storeKey, si)
	return storeKey
}

func (si SwapInit) getKey() *SwapKey {
	sk := &SwapKey{
		Initiator:   si.Party.Key,
		Participant: si.CounterParty.Key,
		Amount:      si.Amount,
		Exchange:    si.Exchange,
		Nonce:       si.Nonce,
	}
	return sk
}

func (si SwapInit) resolve(app interface{}, stageType swapStageType) Commands {
	var sv SwapVerify
	stage, _ := swapStageFlow[SWAP_MATCHING]
	key := si.getKey().toHash()
	account := GetNodeAccount(app)

	commands := amino.DeepCopy(stage.Commands).(Commands)
	context := make(FunctionValues)

	context[STAGE] = stage.Stage
	switch stageType {
	case SWAP_MATCHING:
		context[NEXTSTAGE] = WAIT_FOR_CHAIN
		event := Event{Type: SWAP, SwapKeyHash: key, Step: 0}
		sv = SwapVerify{
			Event: event,
		}
		context[SWAPMESSAGE] = sv
	case SWAP_MATCHED:
		context[NEXTSTAGE] = SWAP_MATCHED
		chains := si.getChains()
		se := SwapExchange{
			Contract:    nil,
			SwapKeyHash: key,
			Chain:       chains[0],
		}
		context[SWAPMESSAGE] = se
	default:
		log.Warn("Unexpected stage for SwapInit", "stage", stageType)
	}
	context[OWNER] = account.AccountKey()
	context[TARGET] = account.AccountKey()
	commands[0].data = context
	return commands
}

func (si SwapInit) getParty(accountKey id.AccountKey) Party {
	var cp Party
	if accountKey == nil {
		log.Debug("Getting Role for empty account")
		return cp
	}

	if bytes.Compare(si.Party.Key, accountKey) == 0 {
		cp = si.Party

	} else if bytes.Compare(si.CounterParty.Key, accountKey) == 0 {
		cp = si.CounterParty

	}
	return cp
}

// Get the correct chains order for this action
func (si SwapInit) getChains() []data.ChainType {

	if si.Amount.Currency.Id < si.Exchange.Currency.Id {
		return []data.ChainType{si.Amount.Currency.Chain, si.Exchange.Currency.Chain}
	} else {
		return []data.ChainType{si.Exchange.Currency.Chain, si.Amount.Currency.Chain}
	}
}

func (si SwapInit) getRole(isParty bool) Role {

	if si.Amount.Currency.Id < si.Exchange.Currency.Id {
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

func (si *SwapInit) order() bool {
	chains := si.getChains()

	if si.Amount.Currency.Chain == chains[0] {
		// don't need to switch
		return true
	} else {
		si.Party, si.CounterParty = si.CounterParty, si.Party
		si.Amount, si.Exchange = si.Exchange, si.Amount
		return false
	}

}

type SwapExchange struct {
	Contract    common.Contract `json:"message"`
	SwapKeyHash []byte          `json:"swapkeyhash"`
	Chain       data.ChainType  `json:"chain"`
	PreviousTx  []byte          `json:"previoustx"`
}

func (se SwapExchange) validate() status.Code {
	log.Debug("Validating SwapExchange")

	//if se.Contract == nil {
	//	log.Debug("Missing Contract")
	//	return status.MISSING_DATA
	//}

	log.Debug("SwapExchange is validated!")
	return status.SUCCESS
}

func (se SwapExchange) resolve(app interface{}, stageType swapStageType) Commands {

	stage, ok := swapStageFlow[stageType]
	if !ok {
		log.Error("stage not found", "stageType", stageType)
		return nil
	}

	commands := amino.DeepCopy(stage.Commands).(Commands)

	context := make(FunctionValues)

	context[SWAPMESSAGE] = se
	context[STOREKEY] = se.SwapKeyHash
	context[STAGE] = stage.Stage
	context[NEXTSTAGE] = stage.OutStage

	si := FindSwap(app, se.SwapKeyHash)
	chains := si.getChains()
	account := GetNodeAccount(app)

	switch stage.Stage {
	case SWAP_MATCHED:
		if bytes.Equal(si.Party.Key, account.AccountKey()) {
			context[OWNER] = si.Party.Key
			context[TARGET] = si.Party.Key
		} else {
			context[OWNER] = si.CounterParty.Key
			context[TARGET] = si.CounterParty.Key
			context[STAGE] = NOSTAGE
		}

	case INITIATOR_INITIATE:

		var secret [32]byte
		_, err := rand.Read(secret[:])
		if err != nil {
			log.Error("failed to get random secret with 32 length", "status", err)
		}
		secretHash := sha256.Sum256(secret[:])
		context[PASSWORD] = secret
		context[PREIMAGE] = secretHash

		context[NEXTCHAINNAME] = chains[0]
		SaveContract(app, se.SwapKeyHash, 0, secret[:])

		context[OWNER] = si.Party.Key
		context[TARGET] = si.CounterParty.Key

	case PARTICIPANT_PARTICIPATE:
		context[NEXTCHAINNAME] = se.Chain
		context[OWNER] = si.CounterParty.Key
		context[TARGET] = si.Party.Key

	case INITIATOR_REDEEM:
		scr := FindContract(app, se.SwapKeyHash, 0)
		_, err := rand.Read(scr[:])
		if err != nil {
			log.Error("failed to get random secret with 32 length", "status", err)
		}
		secretHash := sha256.Sum256(scr[:])
		context[PASSWORD] = scr
		context[PREIMAGE] = secretHash
		context[NEXTCHAINNAME] = se.Chain
		context[OWNER] = si.Party.Key
		context[TARGET] = si.CounterParty.Key
	case PARTICIPANT_REDEEM:
		context[NEXTCHAINNAME] = se.Chain
		scrHash := FindContract(app, se.SwapKeyHash, 0)
		context[PREIMAGE] = scrHash
		context[OWNER] = si.CounterParty.Key
		context[TARGET] = si.Party.Key
	case SWAP_FINISH:

	default:
		log.Warn("Unexpected stage for SwapExchange", "stage", stageType)
	}
	commands[0].data = context
	return commands
}

type SwapVerify struct {
	Event Event `json:"event"`
}

func (sv SwapVerify) validate() status.Code {
	log.Debug("Validating SwapVerify")

	if &sv.Event == nil {
		log.Debug("Missing Event")
		return status.MISSING_DATA
	}

	log.Debug("SwapVerify is validated!")
	return status.SUCCESS
}

func (sv SwapVerify) resolve(app interface{}, stageType swapStageType) Commands {

	stage, ok := swapStageFlow[stageType]
	if !ok {
		log.Error("stage not found", "stageType", stageType)
		return nil
	}

	commands := amino.DeepCopy(stage.Commands).(Commands)
	context := make(FunctionValues)

	context[SWAPMESSAGE] = sv
	context[STOREKEY] = sv.Event.SwapKeyHash
	context[STAGE] = stage.Stage

	commands[0].data = context
	return nil
}

//get the next stage from the current stage
func NextStage(app interface{}, chain data.ChainType, context FunctionValues, tx Transaction) (bool, FunctionValues) {
	// Make sure it is pushed forward first...
	global.Current.Sequence += 32

	message := GetSwapMessage(context[SWAPMESSAGE])
	stage := getStageType(context[STAGE])
	nextStage := getStageType(context[NEXTSTAGE])

	switch stage {
	case SWAP_MATCHING:

	case SWAP_MATCHED:

	case INITIATOR_INITIATE:

	case PARTICIPANT_PARTICIPATE:

	case INITIATOR_REDEEM:

	case PARTICIPANT_REDEEM:
		storeKey := GetBytes(context[STOREKEY])
		event := Event{Type: SWAP, SwapKeyHash: storeKey, Step: 0}
		message = SwapVerify{
			Event: event,
		}
	default:
		return false, nil
	}

	owner := GetAccountKey(context[OWNER])
	target := GetAccountKey(context[TARGET])
	chainId := GetChainID(app)
	//log.Debug("parsed contract", "contract", contract, "chain", chain, "context", context, "count", count)
	swap := &Swap{
		Base: Base{
			Type:     SWAP,
			ChainId:  chainId,
			Signers:  GetSigners(owner),
			Owner:    owner,
			Target:   target,
			Sequence: global.Current.Sequence,
		},
		SwapMessage: message,
		Stage:       nextStage,
	}
	log.Debug("NextStage Swap", "swap", swap)
	if nextStage == WAIT_FOR_CHAIN {
		waitTime := 3 * lockPeriod
		DelayedTransaction(SWAP, swap, waitTime)
	} else {
		waitTime := 1 * time.Second
		DelayedTransaction(SWAP, swap, waitTime)
	}
	return true, nil
}

// Two matching swap requests from different parties
func isMatch(left *SwapInit, right *SwapInit) bool {

	if bytes.Compare(left.Party.Key, right.Party.Key) != 0 {
		log.Debug("Party/CounterParty is wrong")
		return false
	}
	if bytes.Compare(left.CounterParty.Key, right.CounterParty.Key) != 0 {
		log.Debug("CounterParty/Party is wrong")
		return false
	}
	if !left.Amount.Equals(right.Amount) {
		log.Debug("Amount/Exchange is wrong")
		return false
	}
	if !left.Exchange.Equals(right.Exchange) {
		log.Debug("Exchange/Amount is wrong")
		return false
	}
	if left.Nonce != right.Nonce {
		log.Debug("Nonce is wrong")
		return false
	}

	return true
}

func matchingSwap(app interface{}, si *SwapInit) swapStageType {
	ordered := si.order()
	key := si.getKey().toHash()

	result := FindSwap(app, key)
	if result != nil {
		if matching := isMatch(result, si); matching {
			if ordered {
				si.CounterParty = result.CounterParty
			} else {
				si.Party = result.Party
			}
			si.store(app)
			return SWAP_MATCHED
		} else {
			log.Warn("Swap Stored Not matched", "si", si, "result", result)
		}
	} else {
		log.Debug("Swap not found in storage", "key", key)
		si.store(app)
	}
	return SWAP_MATCHING
}

func SaveSwap(app interface{}, swapKey []byte, transaction SwapInit) {
	storage := GetStatus(app)
	session := storage.Begin()
	session.Set(swapKey, transaction)
	session.Commit()
}

func FindSwap(app interface{}, key []byte) *SwapInit {
	storage := GetStatus(app)
	buffer := storage.Get(key)
	if buffer == nil {
		return nil
	}
	result := buffer.(SwapInit)
	return &result
}

func DeleteSwap(app interface{}, key id.AccountKey) {
	storage := GetStatus(app)
	session := storage.Begin()
	ok := session.Delete(key)
	if !ok {
		log.Error("Delete swap failed", "key", key)
		session.Rollback()
	}
	session.Commit()
	return
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

func CreateCheckEvent(app interface{}, chain data.ChainType, context FunctionValues, tx Transaction) (bool, FunctionValues) {
	storeKey := GetBytes(context[STOREKEY])
	si := FindSwap(app, storeKey)
	if si == nil {
		log.Error("Saved swap not found", "key", storeKey)
		return false, nil
	}
	se := GetSwapMessage(context[SWAPMESSAGE]).(SwapExchange)
	previous := GetBytes(context[PREVIOUS])
	se.PreviousTx = previous
	context[SWAPMESSAGE] = se

	event := Event{Type: SWAP, SwapKeyHash: se.SwapKeyHash, Step: 0}
	SaveEvent(app, event, false)

	return true, context
}

func FinalizeSwap(app interface{}, chain data.ChainType, context FunctionValues, tx Transaction) (bool, FunctionValues) {

	sv := GetSwapMessage(context[SWAPMESSAGE]).(SwapVerify)
	SaveEvent(app, sv.Event, true)
	return true, context
}

func VerifySwap(app interface{}, chain data.ChainType, context FunctionValues, tx Transaction) (bool, FunctionValues) {

	sv := GetSwapMessage(context[SWAPMESSAGE]).(SwapVerify)

	finish := FindEvent(app, sv.Event)
	context[FINISHED] = finish
	return true, context
}

func ClearEvent(app interface{}, chain data.ChainType, context FunctionValues, tx Transaction) (bool, FunctionValues) {
	finish := GetBool(context[FINISHED])
	if finish == nil {
		log.Error("Failed to detected the event result")
		return false, nil
	} else if *finish == true {
		storeKey := GetBytes(context[STOREKEY])
		si := FindSwap(app, storeKey)
		chains := si.getChains()
		DeleteSwap(app, storeKey)
		DeleteContract(app, storeKey, 0)
		DeleteContract(app, storeKey, int64(chains[0]))
		DeleteContract(app, storeKey, int64(chains[1]))
		return true, nil
	} else {
		return false, nil
	}
}

// TODO: Needs to be configurable
var lockPeriod = 5 * time.Minute

func Initiate(app interface{}, chain data.ChainType, context FunctionValues, tx Transaction) (bool, FunctionValues) {
	log.Info("Executing Initiate Command", "chain", chain, "context", context)
	switch chain {

	case data.BITCOIN:
		return CreateContractBTC(app, context, tx)
	case data.ETHEREUM:
		return CreateContractETH(app, context, tx)
	case data.ONELEDGER:
		return CreateContractOLT(app, context, tx)
	default:
		log.Warn("Chain not support", "Chain", chain)
		return false, nil
	}
}

func Participate(app interface{}, chain data.ChainType, context FunctionValues, tx Transaction) (bool, FunctionValues) {
	log.Info("Executing Participate Command", "chain", chain, "context", context)
	switch chain {

	case data.BITCOIN:
		return ParticipateBTC(app, context, tx)
	case data.ETHEREUM:
		return ParticipateETH(app, context, tx)
	case data.ONELEDGER:
		return ParticipateOLT(app, context, tx)
	default:
		log.Warn("Chain not support", "Chain", chain)
		return false, nil
	}
}

func Redeem(app interface{}, chain data.ChainType, context FunctionValues, tx Transaction) (bool, FunctionValues) {
	log.Info("Executing Redeem Command", "chain", chain, "context", context)

	switch chain {

	case data.BITCOIN:
		return RedeemBTC(app, context, tx)
	case data.ETHEREUM:
		return RedeemETH(app, context, tx)
	case data.ONELEDGER:
		return RedeemOLT(app, context, tx)
	default:
		log.Warn("Chain not support", "Chain", chain)
		return false, nil
	}
}

func Refund(app interface{}, chain data.ChainType, context FunctionValues, tx Transaction) (bool, FunctionValues) {
	log.Info("Executing Refund Command", "chain", chain, "context", context)
	switch chain {

	case data.BITCOIN:
		return RefundBTC(app, context, tx)
	case data.ETHEREUM:
		return RefundETH(app, context, tx)
	case data.ONELEDGER:
		return RefundOLT(app, context, tx)
	default:
		log.Warn("Chain not support", "Chain", chain)
		return false, nil
	}
}

func ExtractSecret(app interface{}, chain data.ChainType, context FunctionValues, tx Transaction) (bool, FunctionValues) {
	log.Info("Executing ExtractSecret Command", "chain", chain, "context", context)
	switch chain {

	case data.BITCOIN:
		return ExtractSecretBTC(app, context, tx)
	case data.ETHEREUM:
		return ExtractSecretETH(app, context, tx)
	case data.ONELEDGER:
		return ExtractSecretOLT(app, context, tx)
	default:
		log.Warn("Chain not support", "Chain", chain)
		return false, nil
	}
}

func AuditContract(app interface{}, chain data.ChainType, context FunctionValues, tx Transaction) (bool, FunctionValues) {
	log.Info("Executing AuditContract Command", "chain", chain, "context", context)
	switch chain {

	case data.BITCOIN:
		return AuditContractBTC(app, context, tx)
	case data.ETHEREUM:
		return AuditContractETH(app, context, tx)
	case data.ONELEDGER:
		return AuditContractOLT(app, context, tx)
	default:
		log.Warn("Chain not support", "Chain", chain)
		return false, nil
	}
}

func CreateContractBTC(app interface{}, context FunctionValues, tx Transaction) (bool, FunctionValues) {

	timeout := time.Now().Add(2 * lockPeriod).Unix()

	storeKey := GetBytes(context[STOREKEY])
	si := FindSwap(app, storeKey)

	stage := getStageType(context[STAGE])

	var value *big.Int
	var receiverParty Party
	if stage == INITIATOR_INITIATE {
		value = si.Amount.Amount
		receiverParty = si.CounterParty
		context[OWNER] = si.Party.Key
		context[TARGET] = si.CounterParty.Key
	} else {
		value = si.Exchange.Amount
		receiverParty = si.Party
		context[OWNER] = si.CounterParty.Key
		context[TARGET] = si.Party.Key
	}

	receiver := common.GetBTCAddressFromByteArray(data.BITCOIN, receiverParty.Accounts[data.BITCOIN])
	if receiver == nil {
		log.Error("Failed to get btc address from string", "address", receiverParty.Accounts[data.BITCOIN], "target", reflect.TypeOf(receiver))
		return false, nil
	}

	preimage := GetByte32(context[PREIMAGE])
	//if context[PASSWORD] != nil {
	//	scr := GetByte32(context[PASSWORD])
	//	scrHash := sha256.Sum256(scr[:])
	//	if !bytes.Equal(preimage[:], scrHash[:]) {
	//		log.Error("Secret and Secret Hash doesn't match", "preimage", preimage, "scrHash", scrHash)
	//		return false, nil
	//	}
	//}

	cli := bitcoin.GetBtcClient(global.Current.BTCAddress)

	amount := bitcoin.GetAmount(value.String())

	initCmd := htlc.NewInitiateCmd(receiver, amount, timeout, preimage)

	_, err := initCmd.RunCommand(cli)
	if err != nil {
		log.Error("Bitcoin Initiate", "status", err, "context", context)
		return false, nil
	}

	contract := &bitcoin.HTLContract{}
	contract.FromMsgTx(initCmd.Contract, initCmd.ContractTx)

	previous := GetBytes(context[PREVIOUS])
	se := SwapExchange{
		Contract:    contract,
		SwapKeyHash: storeKey,
		Chain:       contract.Chain(),
		PreviousTx:  previous,
	}

	context[SWAPMESSAGE] = se

	SaveContract(app, storeKey, int64(data.BITCOIN), contract.ToBytes())
	log.Debug("btc contract", "contract", contract)

	return true, context
}

func CreateContractETH(app interface{}, context FunctionValues, tx Transaction) (bool, FunctionValues) {
	storeKey := GetBytes(context[STOREKEY])
	si := FindSwap(app, storeKey)

	stage := getStageType(context[STAGE])

	var value *big.Int
	var receiverParty Party
	if stage == INITIATOR_INITIATE {
		value = si.Amount.Amount
		receiverParty = si.CounterParty
	} else {
		value = si.Exchange.Amount
		receiverParty = si.Party
	}
	//todo : need to have a better key to store ethereum contract.
	me := GetNodeAccount(app)

	contractMessage := FindContract(app, me.AccountKey().Bytes(), int64(data.ETHEREUM))
	var contract *ethereum.HTLContract
	if contractMessage == nil {
		contract = ethereum.CreateHtlContract()
		if contract == nil {
			return false, nil
		}
		SaveContract(app, me.AccountKey().Bytes(), int64(data.ETHEREUM), contract.ToBytes())
	} else {
		buffer, err := serial.Deserialize(contractMessage, contract, serial.JSON)
		if err != nil {
			log.Error("Can't deserialze loaded ETH contract", "buffer", contractMessage)
		}
		contract = buffer.(*ethereum.HTLContract)
	}

	preimage := GetByte32(context[PREIMAGE])
	//if context[PASSWORD] != nil {
	//	scr := GetByte32(context[PASSWORD])
	//	scrHash := sha256.Sum256(scr[:])
	//	if !bytes.Equal(preimage[:], scrHash[:]) {
	//		log.Error("Secret and Secret Hash doesn't match", "preimage", preimage, "scrHash", scrHash)
	//		return false, nil
	//	}
	//}

	receiver := common.GetETHAddressFromByteArray(data.ETHEREUM, receiverParty.Accounts[data.ETHEREUM])
	if receiver == nil {
		log.Error("Failed to get eth address from string", "address", receiverParty.Accounts[data.ETHEREUM], "target", reflect.TypeOf(receiver))
	}

	timeoutSecond := int64(lockPeriod.Seconds())
	log.Debug("Create ETH HTLC", "value", value, "receiver", receiver, "preimage", preimage)
	err := contract.Funds(value, big.NewInt(timeoutSecond), *receiver, preimage)
	if err != nil {
		return false, nil
	}

	previous := GetBytes(context[PREVIOUS])
	se := SwapExchange{
		Contract:    contract,
		SwapKeyHash: storeKey,
		Chain:       contract.Chain(),
		PreviousTx:  previous,
	}

	context[SWAPMESSAGE] = se
	return true, context
}

func AuditContractBTC(app interface{}, context FunctionValues, tx Transaction) (bool, FunctionValues) {

	storeKey := GetBytes(context[STOREKEY])
	si := FindSwap(app, storeKey)

	se := GetSwapMessage(context[SWAPMESSAGE]).(SwapExchange)
	contract := se.Contract.(*bitcoin.HTLContract)
	log.Debug("Contract to audit", "contract", contract)

	stage := getStageType(context[STAGE])
	switch stage {
	case PARTICIPANT_PARTICIPATE:

	case INITIATOR_REDEEM:
		scr := FindContract(app, storeKey, 0)
		if scr == nil {
			log.Error("secret not found", "key", storeKey)
			return false, nil
		}
		var secret [32]byte
		copy(secret[:], scr)
		context[PASSWORD] = secret
	}

	// if audit first chain, then participate on second chain,
	// if audit second chain, then redeem on second chain.
	chains := si.getChains()
	context[NEXTCHAINNAME] = chains[1]

	msgTx := contract.GetMsgTx()
	cmd := htlc.NewAuditContractCmd(contract.Contract, msgTx)
	cli := bitcoin.GetBtcClient(global.Current.BTCAddress)
	e := cmd.RunCommand(cli)
	if e != nil {
		log.Error("Bitcoin Audit", "status", e)
		return false, nil
	}

	SaveContract(app, storeKey, int64(data.BITCOIN), contract.ToBytes())
	context[PREIMAGE] = cmd.SecretHash
	return true, context
}

func AuditContractETH(app interface{}, context FunctionValues, tx Transaction) (bool, FunctionValues) {
	se := GetSwapMessage(context[SWAPMESSAGE]).(SwapExchange)
	contract := se.Contract.(*ethereum.HTLContract)
	log.Debug("Contract to audit", "contract", contract.Address.Hash(), "tx", contract.TxHash.String())

	storeKey := GetBytes(context[STOREKEY])
	si := FindSwap(app, storeKey)

	stage := getStageType(context[STAGE])

	var amount data.Coin
	switch stage {
	case PARTICIPANT_PARTICIPATE:
		amount = si.Amount
	case INITIATOR_REDEEM:
		amount = si.Exchange
		scr := FindContract(app, storeKey, 0)
		if scr == nil {
			log.Error("secret not found", "key", storeKey)
			return false, nil
		}
		var secret [32]byte
		copy(secret[:], scr)
		context[PASSWORD] = secret
	}

	// if audit first chain, then participate on second chain,
	// if audit second chain, then redeem on second chain.
	chains := si.getChains()
	context[NEXTCHAINNAME] = chains[1]

	//todo : when support light client, need to get this address from swapinit
	address := ethereum.GetAddress()

	receiver, e := contract.HTLContractObject().Receiver(&bind.CallOpts{Pending: true})
	if e != nil {
		log.Error("can't get the receiver", "status", e)
		return false, nil
	}
	if !bytes.Equal(address.Bytes(), receiver.Bytes()) {
		log.Error("receiver not correct", "contract", contract.Address, "receiver", receiver, "myAddress", address)
		return false, nil
	}

	value := amount.Amount

	setVale := contract.Balance()
	setScrhash := contract.ScrHash()
	//if !bytes.Equal(scrHash[:], setScrhash[:]) {
	//	log.Error("Secret Hash doesn't match", "sh", scrHash, "setSh", setScrhash)
	//	return false, nil
	//}

	if value.Cmp(setVale) != 0 {
		log.Error("Value doesn't match", "value", value, "setValue", setVale)
		return false, nil
	}

	//log.Debug("Auditing ETH Contract", "receiver", address, "value", value, "scrHash", scrHash)
	//
	//log.Debug("Set ETH Contract", "receiver", receiver, "value", setVale, "scrHash", setScrhash)
	//e = contract.Audit(address, value ,scrHash)
	//if e != nil {
	//	log.Error("Failed to audit the contract with correct input", "status", e)
	//	return false, nil
	//}

	context[PREIMAGE] = setScrhash
	SaveContract(app, storeKey, int64(data.ETHEREUM), contract.ToBytes())
	return true, context
}

func ParticipateBTC(app interface{}, context FunctionValues, tx Transaction) (bool, FunctionValues) {
	success, result := CreateContractBTC(app, context, tx)
	if success != false {
		log.Error("failed to participate because can't create contract")
		return false, nil
	}
	return true, result
}

func ParticipateETH(app interface{}, context FunctionValues, tx Transaction) (bool, FunctionValues) {
	success, result := CreateContractETH(app, context, tx)
	if success == false {
		log.Error("failed to participate because can't create contract")
		return false, nil
	}
	return true, result
}

func RedeemBTC(app interface{}, context FunctionValues, tx Transaction) (bool, FunctionValues) {

	storeKey := GetBytes(context[STOREKEY])

	buffer := FindContract(app, storeKey, int64(data.BITCOIN))
	if buffer == nil {
		log.Error("Failed to load the contract to Redeem", "key", storeKey)
		return false, nil
	}
	contract := &bitcoin.HTLContract{}
	contract.FromBytes(buffer)

	scr := GetByte32(context[PASSWORD])

	cmd := htlc.NewRedeemCmd(contract.Contract, contract.GetMsgTx(), scr[:])
	cli := bitcoin.GetBtcClient(global.Current.BTCAddress)
	_, e := cmd.RunCommand(cli)
	if e != nil {
		log.Error("Bitcoin redeem htlc", "status", e)
		return false, nil
	}

	newcontract := &bitcoin.HTLContract{}
	newcontract.FromMsgTx(contract.Contract, cmd.RedeemContractTx)

	previous := GetBytes(context[PREVIOUS])

	se := SwapExchange{
		Contract:    newcontract,
		SwapKeyHash: storeKey,
		Chain:       contract.Chain(),
		PreviousTx:  previous,
	}

	context[SWAPMESSAGE] = se
	return true, context
}

func RedeemETH(app interface{}, context FunctionValues, tx Transaction) (bool, FunctionValues) {
	storeKey := GetBytes(context[STOREKEY])

	buffer := FindContract(app, storeKey, int64(data.ETHEREUM))
	if buffer == nil {
		log.Error("Failed to load the contract to Redeem", "key", storeKey)
		return false, nil
	}
	contract := &ethereum.HTLContract{}
	contract.FromBytes(buffer)

	scr := GetByte32(context[PASSWORD])
	err := contract.Redeem(scr[:])
	if err != nil {
		log.Error("Ethereum redeem htlc", "status", err)
		return false, nil
	}

	previous := GetBytes(context[PREVIOUS])
	se := SwapExchange{
		SwapKeyHash: storeKey,
		Chain:       contract.Chain(),
		PreviousTx:  previous,
	}

	context[SWAPMESSAGE] = se
	return true, context
}

func RefundBTC(app interface{}, context FunctionValues, tx Transaction) (bool, FunctionValues) {
	storeKey := GetBytes(context[STOREKEY])

	buffer := FindContract(app, storeKey, int64(data.BITCOIN))
	if buffer == nil {
		return false, nil
	}
	contract := &bitcoin.HTLContract{}
	contract.FromBytes(buffer)

	if contract == nil {
		log.Error("BTC Htlc contract can't be parsed")
		return false, nil
	}

	cmd := htlc.NewRefundCmd(contract.Contract, contract.GetMsgTx())
	cli := bitcoin.GetBtcClient(global.Current.BTCAddress)
	_, e := cmd.RunCommand(cli)
	if e != nil {
		log.Error("Bitcoin refund htlc", "status", e)
		return false, nil
	}
	return true, context
}

func RefundETH(app interface{}, context FunctionValues, tx Transaction) (bool, FunctionValues) {
	storeKey := GetBytes(context[STOREKEY])

	buffer := FindContract(app, storeKey, int64(data.ETHEREUM))
	if buffer == nil {
		return false, nil
	}
	contract := &ethereum.HTLContract{}
	contract.FromBytes(buffer)

	if contract == nil {
		log.Error("ETH Htlc contract can't be parsed")
		return false, nil
	}

	err := contract.Refund()
	if err != nil {
		return false, nil
	}
	return true, context
}

func ExtractSecretBTC(app interface{}, context FunctionValues, tx Transaction) (bool, FunctionValues) {
	se := GetSwapMessage(context[SWAPMESSAGE]).(SwapExchange)
	contract := se.Contract.(*bitcoin.HTLContract)
	storeKey := se.SwapKeyHash
	//
	//buffer := FindContract(app, storeKey, int64(data.BITCOIN))
	//contract := &bitcoin.HTLContract{}
	//contract.FromBytes(buffer)

	si := FindSwap(app, storeKey)

	chains := si.getChains()
	context[NEXTCHAINNAME] = chains[0]

	scrHash := GetByte32(context[PREIMAGE])

	cmd := htlc.NewExtractSecretCmd(contract.GetMsgTx(), scrHash)
	cli := bitcoin.GetBtcClient(global.Current.BTCAddress)
	e := cmd.RunCommand(cli)
	if e != nil {
		log.Error("Bitcoin extract hltc", "status", e)
		return false, nil
	}
	var tmpScr [32]byte
	copy(tmpScr[:], string(cmd.Secret))
	context[PASSWORD] = tmpScr
	log.Debug("extracted secret", "secretbytearray", cmd.Secret, "secretbyte32", tmpScr)
	return true, context

}

func ExtractSecretETH(app interface{}, context FunctionValues, tx Transaction) (bool, FunctionValues) {

	se := GetSwapMessage(context[SWAPMESSAGE]).(SwapExchange)
	storeKey := se.SwapKeyHash

	si := FindSwap(app, storeKey)

	chains := si.getChains()
	context[NEXTCHAINNAME] = chains[0]

	//todo : need to have a better key to store ethereum contract.
	me := GetNodeAccount(app)

	buffer := FindContract(app, me.AccountKey().Bytes(), int64(data.ETHEREUM))
	if buffer == nil {
		log.Error("Failed to find eth local contract")
		return false, nil
	}

	contract := &ethereum.HTLContract{}
	contract.FromBytes(buffer)

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

func CreateContractOLT(app interface{}, context FunctionValues, tx Transaction) (bool, FunctionValues) {
	log.Warn("Not supported")
	//party := GetParty(context[MY_ACCOUNT])
	//counterParty := GetParty(context[THEM_ACCOUNT])
	//partyBalance := GetUtxo(app).Find(party.Key).Amount
	//counterPartyBalance := GetUtxo(app).Find(counterParty.Key).Amount
	//
	//preimage := GetByte32(context[PREIMAGE])
	//if context[PASSWORD] != nil {
	//	scr := GetByte32(context[PASSWORD])
	//	scrHash := sha256.Sum256(scr[:])
	//	if !bytes.Equal(preimage[:], scrHash[:]) {
	//		log.Error("Secret and Secret Hash doesn't match", "preimage", preimage, "scrHash", scrHash)
	//		return false, nil
	//	}
	//}
	//
	//inputs := make([]SendInput, 0)
	//inputs = append(inputs,
	//	NewSendInput(party.Key, partyBalance),
	//	NewSendInput(counterParty.Key, counterPartyBalance))
	//amount := GetCoin(context[AMOUNT])
	//// Build up the outputs
	//outputs := make([]SendOutput, 0)
	//outputs = append(outputs,
	//	NewSendOutput(party.Key, partyBalance.Minus(amount)),
	//	NewSendOutput(counterParty.Key, counterPartyBalance.Plus(amount)))
	//send := &Send{
	//	Base: Base{
	//		Type:     SEND,
	//		ChainId:  GetChainID(app),
	//		Signers:  nil,
	//		Sequence: global.Current.Sequence,
	//	},
	//	Inputs:  inputs,
	//	Outputs: outputs,
	//	Fee:     data.NewCoin(0, "OLT"),
	//	Gas:     data.NewCoin(0, "OLT"),
	//}
	//message := SignAndPack(SEND, Transaction(send))
	//contract := NewMultiSigBox(1, 1, message)
	//_ = contract
	return false, nil
}

func ParticipateOLT(app interface{}, context FunctionValues, tx Transaction) (bool, FunctionValues) {
	log.Warn("Not supported")
	return true, context
}

func AuditContractOLT(app interface{}, context FunctionValues, tx Transaction) (bool, FunctionValues) {
	log.Warn("Not supported")
	return true, context
}

func RedeemOLT(app interface{}, context FunctionValues, tx Transaction) (bool, FunctionValues) {
	log.Warn("Not supported")
	return true, context
}

func RefundOLT(app interface{}, context FunctionValues, tx Transaction) (bool, FunctionValues) {
	log.Warn("Not supported")
	return true, context
}

func ExtractSecretOLT(app interface{}, context FunctionValues, tx Transaction) (bool, FunctionValues) {
	log.Warn("Not supported")
	return true, context
}
