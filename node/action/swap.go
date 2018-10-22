/*
	Copyright 2017-2018 OneLedger

	An incoming transaction, send, swap, ready, verification, etc.
*/
package action

import (
	"bytes"
	"github.com/Oneledger/protocol/node/comm"
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
	SwapMessage SwapMessage   `json:"message"`
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

	if bytes.Equal(transaction.Base.Target, account.AccountKey()) {
		return true
	}

	if bytes.Equal(transaction.Base.Owner, account.AccountKey()) {
		return true
	}

	return false
}

// Start the swap
func (transaction *Swap) ProcessDeliver(app interface{}) status.Code {
	log.Debug("Processing Swap Transaction for DeliverTx")
	commands := transaction.Resolve(app)

	if commands.Count() == 0 {
		return status.EXPAND_ERROR
	}

	for i := 0; i < commands.Count(); i++ {
		//them := GetParty(commands[i].Data[THEM_ACCOUNT])
		//event := Event{Type: SWAP, Key: them.Key , Nonce: matchedSwap.Nonce }
		//if commands[i].Function == FINISH {
		//	SaveEvent(app, event, true)
		//	return status.SUCCESS
		//}

		ok, result := commands[i].Execute(app)
		if !ok {
			log.Error("Failed to Execute", "command", commands[i])
			//SaveEvent(app, event, false)
			return status.EXPAND_ERROR
		}
		if len(result) > 0 {
			commands[i+1].data = result
			chain := GetChain(result[NEXTCHAINNAME])
			commands[i+1].chain = chain
		}
	}

	return status.SUCCESS
}

// Plug in data from the rest of a system into a set of commands
func (swap *Swap) Resolve(app interface{}) Commands {
	return swap.SwapMessage.resolve(app, swap.Stage)
}

var swapStageFlow = map[swapStageType]swapStage{
	SWAP_MATCHING: {
		Stage:    SWAP_MATCHING,
		Commands: Commands{Command{opfunc: SaveUnmatchSwap}, Command{opfunc: NextStage}},
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

func (si SwapInit) resolve(app interface{}, stageType swapStageType) Commands {
	chains := si.getChains()
	account := GetNodeAccount(app)
	isParty := si.IsParty(account)
	role := si.getRole(*isParty)

	stage, ok := swapStageFlow[stageType]
	if !ok {
		log.Error("stage not found", "stagetype", stageType)
		return nil
	}

	commands := amino.DeepCopy(stage.Commands).(Commands)

	command := &commands[0]
	if stageType == SWAP_MATCHING {

		swap := MatchingSwap(app, &si)

		var storeInfo SwapInit
		if swap == nil {
			command.data[STAGE] = WAIT_FOR_CHAIN
			if *isParty {
				command.data[STOREKEY] = si.CounterParty.Key
			} else {
				command.data[STOREKEY] = si.Party.Key
			}
			storeInfo = si
		} else {
			command.data[STAGE] = SWAP_MATCHED
			command.data[STOREKEY] = _hash(swap)
			storeInfo = *swap
		}
		buffer, err := serial.Serialize(storeInfo, serial.NETWORK)
		if err != nil {
			log.Error("Serialize swapInit failed", "err", err)
		}
		command.data[STOREMESSAGE] = buffer
	} else {
		swapKey := _hash(si)
		swap := FindSwap(app, swapKey, false)
		if swap == nil {
			log.Error("SwapInit not find", "key", swapKey)
			return nil
		}
		command.data[STAGE] = SWAP_MATCHED
		if *isParty {
			command.data[MY_ACCOUNT] = swap.Party
			command.data[THEM_ACCOUNT] = swap.CounterParty
			command.data[AMOUNT] = swap.Amount
			command.data[EXCHANGE] = swap.Exchange
		} else {
			log.Dump("this should never show")
			command.data[MY_ACCOUNT] = swap.CounterParty
			command.data[THEM_ACCOUNT] = swap.Party
			command.data[AMOUNT] = swap.Exchange
			command.data[EXCHANGE] = swap.Amount
		}
		command.data[STOREKEY] = swapKey
		command.data[ROLE] = role
		command.data[NONCE] = swap.Nonce
		command.data[EVENTTYPE] = SWAP

		if stage.Stage == INITIATOR_INITIATE {
			command.chain = chains[1]

			var secret [32]byte
			_, err := rand.Read(secret[:])
			if err != nil {
				log.Error("failed to get random secret with 32 length", "status", err)
			}
			secretHash := sha256.Sum256(secret[:])
			command.data[PASSWORD] = secret
			command.data[PREIMAGE] = secretHash
		}
	}
	return commands
}

func (si SwapInit) IsParty(account id.Account) *bool {

	if account == nil {
		log.Debug("Getting Role for empty account")
		return nil
	}

	var isParty bool
	if bytes.Compare(si.Party.Key, account.AccountKey()) == 0 {
		isParty = true
		return &isParty
	}

	if bytes.Compare(si.CounterParty.Key, account.AccountKey()) == 0 {
		isParty = false
		return &isParty
	}

	// TODO: Shouldn't be in-band this way
	return nil
}

// Get the correct chains order for this action
func (si SwapInit) getChains() []data.ChainType {

	var first, second data.ChainType
	if si.Amount.Currency.Id < si.Exchange.Currency.Id {
		first = data.Currencies[si.Amount.Currency.Name].Chain
		second = data.Currencies[si.Exchange.Currency.Name].Chain
	} else {
		first = data.Currencies[si.Exchange.Currency.Name].Chain
		second = data.Currencies[si.Amount.Currency.Name].Chain
	}

	return []data.ChainType{data.ONELEDGER, first, second}
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

type SwapExchange struct {
	Contract   Message        `json:"message"` //message converted from HTLContract
	SecretHash [32]byte       `json:"secrethash"`
	Chain      data.ChainType `json:"chain"`
}

func (se SwapExchange) validate() status.Code {
	log.Debug("Validating SwapExchange")

	if se.Contract == nil {
		log.Debug("Missing Contract")
		return status.MISSING_DATA
	}

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

	command := &commands[0]

	if se.Chain == data.BITCOIN {
		contract := bitcoin.GetHTLCFromMessage(se.Contract)
		command.data[BTCCONTRACT] = contract

	} else if se.Chain == data.ETHEREUM {
		contract := ethereum.GetHTLCFromMessage(se.Contract)
		command.data[ETHCONTRACT] = contract
	}
	command.data[PREIMAGE] = se.SecretHash
	command.chain = se.Chain
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

	return nil
}

func SaveUnmatchSwap(app interface{}, chain data.ChainType, context FunctionValues) (bool, FunctionValues) {

	swap := GetBytes(context[STOREMESSAGE])
	key := GetAccountKey(context[STOREKEY])
	SaveSwap(app, key, swap)

	//After finished the save, then create the transaction will be create for next stage.
	them := GetParty(context[THEM_ACCOUNT])
	nonce := GetInt64(context[NONCE])
	event := Event{Type: SWAP, Key: them.Key, Nonce: nonce}
	swapVerify := &SwapVerify{
		Event: event,
	}

	buffer, err := serial.Serialize(swapVerify, serial.NETWORK)
	if err != nil {
		log.Error("Serialize SwapVeify failed", "err", err)
	}
	context[STOREMESSAGE] = buffer

	return true, context
}

//get the next stage from the current stage
func NextStage(app interface{}, chain data.ChainType, context FunctionValues) (bool, FunctionValues) {
	// Make sure it is pushed forward first...
	global.Current.Sequence += 32

	message := GetBytes(context[STOREMESSAGE])
	var proto SwapMessage
	swapMessage, err := serial.Deserialize(message, proto, serial.NETWORK)
	if err != nil {
		log.Error("Deserialize SwapMessage failed", "err", err)
	}
	stage := GetInt(context[STAGE])

	signers := make([]PublicKey, 0)
	owner := GetParty(context[MY_ACCOUNT])
	target := GetParty(context[THEM_ACCOUNT])
	chainId := GetString(context[CHAINID])
	//log.Debug("parsed contract", "contract", contract, "chain", chain, "context", context, "count", count)
	swap := &Swap{
		Base: Base{
			Type:     SWAP,
			ChainId:  chainId,
			Signers:  signers,
			Owner:    owner.Key,
			Target:   target.Key,
			Sequence: global.Current.Sequence,
		},
		SwapMessage: swapMessage.(SwapMessage),
		Stage:       swapStageType(stage),
	}

	packet := SignAndPack(SWAP, Transaction(swap))

	result := comm.Broadcast(packet)
	log.Debug("Submit Transaction to OLT successfully", "result", result)
	return true, nil
}

// Two matching swap requests from different parties
func isMatch(left *SwapInit, right *SwapInit) bool {

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

func MatchingSwap(app interface{}, transaction *SwapInit) *SwapInit {

	account := GetNodeAccount(app)

	isParty := transaction.IsParty(account)

	if isParty == nil {
		log.Debug("No Account", "account", account)
		return nil
	}
	var matchedSwap *SwapInit

	storage := GetStatus(app)
	matchedSwap.Fee = transaction.Fee
	matchedSwap.Nonce = transaction.Nonce

	if *isParty {
		result := FindSwap(storage, transaction.CounterParty.Key, true)
		if result != nil {
			if matching := isMatch(result, transaction); matching {
				matchedSwap.Party = transaction.Party
				matchedSwap.CounterParty = result.Party
				matchedSwap.Amount = transaction.Amount
				matchedSwap.Exchange = transaction.Exchange
				return matchedSwap
			}
		}
	} else {
		result := FindSwap(storage, transaction.Party.Key, true)
		if result != nil {
			if matching := isMatch(result, transaction); matching {
				matchedSwap.Party = result.Party
				matchedSwap.CounterParty = transaction.Party
				matchedSwap.Amount = transaction.Exchange
				matchedSwap.Exchange = transaction.Amount
				return matchedSwap
			}
		}
	}
	return nil
}

func SaveSwap(app interface{}, swapKey id.AccountKey, transaction interface{}) {
	log.Debug("SaveSwap", "key", swapKey)
	storage := GetStatus(app)
	buffer, err := serial.Serialize(transaction, serial.PERSISTENT)
	if err != nil {
		log.Error("Failed to Serialize SaveSwap transaction")
	}
	storage.Store(swapKey, buffer)
	storage.Commit()
}

func FindSwap(app interface{}, key id.AccountKey, delete bool) *SwapInit {
	storage := GetStatus(app)
	result := storage.Load(key)
	if result == nil {
		return nil
	}

	if delete {
		storage.Delete(key)
	}

	var transaction Swap
	buffer, err := serial.Deserialize(result, &transaction, serial.CLIENT)
	if err != nil {
		return nil
	}
	return buffer.(*SwapInit)
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

func CreateCheckEvent(app interface{}, chain data.ChainType, context FunctionValues) (bool, FunctionValues) {
	swapKey := GetBytes(context[STOREKEY])
	si := FindSwap(app, swapKey, false)
	if si == nil {
		log.Error("Saved swap not found", "key", swapKey)
	}

	event := Event{Type: SWAP, Key: si.CounterParty.Key, Nonce: si.Nonce}
	SaveEvent(app, event, false)

	return true, context
}

func FinalizeSwap(app interface{}, chain data.ChainType, context FunctionValues) (bool, FunctionValues) {
	//todo:
	return true, context
}

func VerifySwap(app interface{}, chain data.ChainType, context FunctionValues) (bool, FunctionValues) {
	//todo:
	return true, context
}

func ClearEvent(app interface{}, chain data.ChainType, context FunctionValues) (bool, FunctionValues) {
	//todo:
	return true, context
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
		log.Error("Bitcoin Initiate", "status", err, "context", context)
		return false, nil
	}

	contract := &bitcoin.HTLContract{Contract: initCmd.Contract, ContractTx: *initCmd.ContractTx}

	context[BTCCONTRACT] = contract

	nonce := GetInt64(context[NONCE])
	SaveContract(app, receiverParty.Key.Bytes(), nonce, contract)
	log.Debug("btc contract", "contract", context[BTCCONTRACT])
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
		log.Error("Bitcoin Audit", "status", e)
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
		log.Error("can't get the receiver", "status", e)
		return false, nil
	}
	if !bytes.Equal(address.Bytes(), receiver.Bytes()) {
		log.Error("receiver not correct", "contract", contract.Address, "receiver", receiver, "myAddress", address)
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

	log.Debug("Auditing ETH Contract", "receiver", address, "value", value, "scrHash", scrHash)

	log.Debug("Set ETH Contract", "receiver", receiver, "value", setVale, "scrHash", setScrhash)
	//e = contract.Audit(address, value ,scrHash)
	//if e != nil {
	//	log.Error("Failed to audit the contract with correct input", "status", e)
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
		return false, nil
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
		log.Error("Bitcoin redeem htlc", "status", e)
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
		log.Error("Ethereum redeem htlc", "status", err)
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
		return false, nil
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
		log.Error("Bitcoin refund htlc", "status", e)
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
		log.Error("Bitcoin extract hltc", "status", e)
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
	contract := NewMultiSigBox(1, 1, message)
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

func WaitForChain(app interface{}, chain data.ChainType, context FunctionValues) (bool, FunctionValues) {
	log.Info("Executing WaitForChain Command", "chain", chain, "context", context)
	//todo : make this to check finish status, and then rollback if necessary
	// Make sure it is pushed forward first...
	global.Current.Sequence += 32

	signers := []id.PublicKey(nil)
	owner := GetParty(context[MY_ACCOUNT])
	target := GetParty(context[THEM_ACCOUNT])
	eventType := GetType(context[EVENTTYPE])
	nonce := GetInt64(context[NONCE])
	stage := GetInt(context[STAGE])
	verify := SwapVerify{
		Event: Event{
			Type:  eventType,
			Key:   target.Key,
			Nonce: nonce,
		},
	}
	swap := &Swap{
		Base: Base{
			Type:     SWAP,
			ChainId:  "OneLedger-Root",
			Owner:    owner.Key,
			Signers:  signers,
			Sequence: global.Current.Sequence,
		},
		SwapMessage: verify,
		Stage:       swapStageType(stage),
	}
	DelayedTransaction(SWAP, Transaction(swap), 3*lockPeriod)

	return true, nil
}
