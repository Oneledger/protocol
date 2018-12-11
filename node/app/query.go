/*
	Copyright 2017-2018 OneLedger

	Implement all of the query mechanics for the node and the chain
*/
package app

import (
	"encoding/hex"
	"strings"

	"github.com/Oneledger/protocol/node/comm"

	"github.com/Oneledger/protocol/node/action"
	"github.com/Oneledger/protocol/node/chains/common"

	"github.com/Oneledger/protocol/node/convert"
	"github.com/Oneledger/protocol/node/data"

	"github.com/Oneledger/protocol/node/global"
	"github.com/Oneledger/protocol/node/status"

	"github.com/Oneledger/protocol/node/id"
	"github.com/Oneledger/protocol/node/log"
	"github.com/Oneledger/protocol/node/serial"
	"github.com/Oneledger/protocol/node/version"
)

// Top-level list of all query types
func HandleQuery(app Application, path string, arguments map[string]string) []byte {

	var result interface{}

	switch path {
	case "/applyValidators":
		result = HandleApplyValidatorQuery(app, arguments)

	case "/createExSendRequest":
		result = HandleCreateExSendRequest(app, arguments)

	case "/createSendRequest":
		result = HandleCreateSendRequest(app, arguments)

	case "/createMintRequest":
		result = HandleCreateMintRequest(app, arguments)

	case "/createSwapRequest":
		result = HandleSwapRequest(app, arguments)

	case "/nodeName":
		result = HandleNodeNameQuery(app, arguments)

	case "/identity":
		result = HandleIdentityQuery(app, arguments)

	case "/signTransaction":
		result = HandleSignTransaction(app, arguments)

	case "/accountPublicKey":
		result = HandleAccountPublicKeyQuery(app, arguments)

	case "/accountKey":
		result = HandleAccountKeyQuery(app, arguments)

	case "/account":
		result = HandleAccountQuery(app, arguments)

	case "/balance":
		result = HandleBalanceQuery(app, arguments)

	case "/version":
		result = HandleVersionQuery(app, arguments)

	case "/currencyAddress":
		result = HandleCurrencyAddressQuery(app, arguments)

	case "/sequenceNumber":
		result = HandleSequenceNumberQuery(app, arguments)

	case "/testScript":
		result = HandleTestScript(app, arguments)

	default:
		result = HandleError("Unknown Query", path, arguments)
	}

	buffer, err := serial.Serialize(result, serial.CLIENT)
	if err != nil {
		log.Debug("Failed to serialize query")
	}

	return buffer
}

func HandleApplyValidatorQuery(application Application, arguments map[string]string) interface{} {
	log.Debug("HandleApplyValidatorQuery", "arguments", arguments)

	conv := convert.NewConvert()

	result := make([]byte, 0)

	var argsHolder comm.ApplyValidatorArguments

	text := arguments["parameters"]

	args, err := serial.Deserialize([]byte(text), argsHolder, serial.CLIENT)

	if err != nil {
		log.Error("Could not deserialize ApplyValidatorArguments", "error", err)
		return result
	}

	amount := args.(*comm.ApplyValidatorArguments).Amount
	idName := args.(*comm.ApplyValidatorArguments).Id

	identity, ok := application.Identities.FindName(idName)
	if ok != status.SUCCESS {
		log.Error("The identity is not found", "id", idName)
		return result
	}

	if amount == "" {
		log.Error("Missing an amount argument")
		return result
	}

	balance := application.Balances.Get(data.DatabaseKey(identity.AccountKey))

	if &balance == nil {
		log.Error("Missing Balance")
		return result
	}

	stake := conv.GetCoin(amount, "VT")

	if balance.GetAmountByName("VT").LessThanCoin(stake) {
		log.Error("Validator token is not enough")
		return result
	}

	sequence := SequenceNumber(application, identity.AccountKey.Bytes())

	validator := &action.ApplyValidator{
		Base: action.Base{
			Type:     action.APPLY_VALIDATOR,
			ChainId:  ChainId,
			Owner:    identity.AccountKey,
			Signers:  GetSigners(identity.AccountKey, application),
			Sequence: sequence.Sequence,
		},

		AccountKey:        identity.AccountKey,
		Identity:          identity.Name,
		NodeName:          global.Current.NodeName,
		TendermintAddress: identity.TendermintAddress,
		TendermintPubKey:  identity.TendermintPubKey,
		Stake:             stake,
	}

	signed := SignTransaction(action.Transaction(validator), application)
	packet := action.PackRequest(signed)

	return packet
}

func HandleCreateExSendRequest(application Application, arguments map[string]string) interface{} {
	conv := convert.NewConvert()

	result := make([]byte, 0)

	var argsHolder comm.ExSendArguments

	text := arguments["parameters"]

	argsDeserialized, err := serial.Deserialize([]byte(text), argsHolder, serial.CLIENT)

	if err != nil {
		log.Error("Could not deserialize ExSendArguments", "error", err)
		return result
	}

	args := argsDeserialized.(*comm.ExSendArguments)

	partyKey := AccountKey(application, args.SenderId)
	cpartyKey := AccountKey(application, args.ReceiverId)

	fee := conv.GetCoin(args.Fee, "OLT")
	gas := conv.GetCoin(args.Gas, "OLT")
	amount := conv.GetCoin(args.Amount, args.Currency)
	chain := conv.GetChainFromCurrency(args.Chain)

	sender := CurrencyAddress(application, conv.GetCurrency(args.Currency), args.SenderId)
	reciever := CurrencyAddress(application, conv.GetCurrency(args.Currency), args.ReceiverId)

	sequence := SequenceNumber(application, partyKey)

	exSend := &action.ExternalSend{
		Base: action.Base{
			Type:     action.EXTERNAL_SEND,
			ChainId:  ChainId,
			Signers:  GetSigners(sender, application),
			Owner:    partyKey,
			Target:   cpartyKey,
			Sequence: sequence.Sequence,
		},
		Gas:      gas,
		Fee:      fee,
		Chain:    chain,
		Sender:   string(sender),
		Receiver: string(reciever),
		Amount:   amount,
	}

	signed := SignTransaction(action.Transaction(exSend), application)
	packet := action.PackRequest(signed)

	return packet
}

func HandleCreateSendRequest(application Application, arguments map[string]string) interface{} {
	conv := convert.NewConvert()

	result := make([]byte, 0)

	var argsHolder comm.SendArguments

	text := arguments["parameters"]

	argsDeserialized, err := serial.Deserialize([]byte(text), argsHolder, serial.CLIENT)

	if err != nil {
		log.Error("Could not deserialize SendArguments", "error", err)
		return result
	}

	args := argsDeserialized.(*comm.SendArguments)

	if args.Party == "" {
		log.Error("Missing Party argument")
		return result
	}

	if args.CounterParty == "" {
		log.Error("Missing CounterParty argument")
		return result
	}

	// TODO: Can't convert identities to accounts, this way!
	party := AccountKey(application, args.Party)
	counterParty := AccountKey(application, args.CounterParty)

	if party == nil || counterParty == nil {
		log.Fatal("System doesn't recognize the parties", "args", args,
			"party", party, "counterParty", counterParty)
	}

	if args.Currency == "" || args.Amount == "" {
		log.Error("Missing an amount argument")
		return result
	}

	amount := conv.GetCoin(args.Amount, args.Currency)

	fee := conv.GetCoin(args.Fee, "OLT")
	gas := conv.GetCoin(args.Gas, "OLT")

	sendTo := action.SendTo{
		AccountKey: counterParty,
		Amount:     amount,
	}

	if conv.HasErrors() {
		log.Error("Conversion error", "error", conv.GetErrors())
		return result
	}

	sequence := SequenceNumber(application, party.Bytes())

	send := &action.Send{
		Base: action.Base{
			Type:     action.SEND,
			ChainId:  ChainId,
			Owner:    party,
			Signers:  GetSigners(party, application),
			Sequence: sequence.Sequence,
		},
		SendTo: sendTo,
		Fee:    fee,
		Gas:    gas,
	}

	signed := SignTransaction(action.Transaction(send), application)
	packet := action.PackRequest(signed)

	return packet
}

func HandleCreateMintRequest(application Application, arguments map[string]string) interface{} {
	result := make([]byte, 0)

	var argsHolder comm.SendArguments

	text := arguments["parameters"]

	argsDeserialized, err := serial.Deserialize([]byte(text), argsHolder, serial.CLIENT)

	if err != nil {
		log.Error("Could not deserialize SendArguments", "error", err)
		return result
	}

	args := argsDeserialized.(*comm.SendArguments)

	conv := convert.NewConvert()

	if args.Party == "" {
		log.Warn("Missing Party arguments", "args", args)
		return result
	}

	zero := AccountKey(application, "Zero")
	party := AccountKey(application, args.Party)

	if party == nil || zero == nil {
		log.Warn("Missing Party information", "args", args, "party", party, "zero", zero)
		return result
	}

	amount := conv.GetCoin(args.Amount, args.Currency)

	log.Debug("Getting TestMint Account Balances")
	zeroBalance := Balance(application, zero).(*data.Balance).GetAmountByName(args.Currency)

	if &zeroBalance == nil {
		log.Warn("Missing Balances", "zero", zero)
		return result
	}

	if zeroBalance.LessThanEqual(0) {
		log.Warn("No more money left...")
		return result
	}

	sendTo := action.SendTo{
		AccountKey: party,
		Amount:     amount,
	}

	gas := conv.GetCoin(args.Gas, "OLT")
	fee := conv.GetCoin(args.Fee, "OLT")

	if conv.HasErrors() {
		log.Error("Conversion error", "error", conv.GetErrors())
		return result
	}

	sequence := SequenceNumber(application, party)

	send := &action.Send{
		Base: action.Base{
			Type:     action.SEND,
			ChainId:  ChainId,
			Signers:  GetSigners(zero, application),
			Owner:    zero,
			Sequence: sequence.Sequence,
		},
		SendTo: sendTo,
		Fee:    fee,
		Gas:    gas,
	}

	signed := SignTransaction(action.Transaction(send), application)
	packet := action.PackRequest(signed)

	return packet
}

func HandleSwapRequest(application Application, arguments map[string]string) interface{} {
	result := make([]byte, 0)

	var argsHolder comm.SwapArguments

	text := arguments["parameters"]

	argsDeserialized, err := serial.Deserialize([]byte(text), argsHolder, serial.CLIENT)

	if err != nil {
		log.Error("Could not deserialize SwapArgumentsArguments", "error", err)
		return result
	}

	args := argsDeserialized.(*comm.SwapArguments)

	conv := convert.NewConvert()

	partyKey := AccountKey(application, args.Party)
	counterPartyKey := AccountKey(application, args.CounterParty)

	fee := conv.GetCoin(args.Fee, "OLT")
	gas := conv.GetCoin(args.Gas, "OLT")

	amount := conv.GetCoin(args.Amount, args.Currency)
	exchange := conv.GetCoin(args.Exchange, args.Excurrency)

	if conv.HasErrors() {
		log.Error("Conversion error", "error", conv.GetErrors())
		return result
	}

	account := make(map[data.ChainType][]byte)
	counterAccount := make(map[data.ChainType][]byte)

	account[conv.GetChainFromCurrency(args.Currency)] = CurrencyAddress(application, conv.GetCurrency(args.Currency), args.Party)
	account[conv.GetChainFromCurrency(args.Excurrency)] = CurrencyAddress(application, conv.GetCurrency(args.Excurrency), args.Party)

	party := action.Party{Key: partyKey, Accounts: account}
	counterParty := action.Party{Key: counterPartyKey, Accounts: counterAccount}

	swapInit := action.SwapInit{
		Party:        party,
		CounterParty: counterParty,
		Fee:          fee,
		Gas:          gas,
		Amount:       amount,
		Exchange:     exchange,
		Nonce:        args.Nonce,
	}

	sequence := SequenceNumber(application, partyKey.Bytes())

	swap := &action.Swap{
		Base: action.Base{
			Type:     action.SWAP,
			ChainId:  ChainId,
			Signers:  GetSigners(partyKey, application),
			Owner:    partyKey,
			Target:   counterPartyKey,
			Sequence: sequence.Sequence,
		},
		SwapMessage: swapInit,
		Stage:       action.SWAP_MATCHING,
	}

	signed := SignTransaction(action.Transaction(swap), application)
	packet := action.PackRequest(signed)

	return packet
}

func HandleNodeNameQuery(app Application, arguments map[string]string) interface{} {
	return global.Current.NodeName
}

func HandleSignTransaction(app Application, arguments map[string]string) interface{} {
	log.Debug("SignTransactionQuery", "arguments", arguments)

	var tx action.Transaction

	text := arguments["parameters"]
	transaction, transactionErr := serial.Deserialize([]byte(text), tx, serial.CLIENT)

	signatures := []action.TransactionSignature{}

	if transactionErr != nil {
		log.Error("Could not deserialize a transaction", "error", transactionErr)
		return signatures
	}

	accountKey := transaction.(action.Transaction).GetOwner()

	if accountKey == nil {
		log.Error("Account key is null", "transaction", transaction)
		return signatures
	}

	account, accountStatus := app.Accounts.FindKey(accountKey)

	if accountStatus != status.SUCCESS {
		log.Error("Could not find an account", "status", accountStatus)
		return signatures
	}

	privateKey := account.PrivateKey()

	signature, signatureError := privateKey.Sign([]byte(text))

	if signatureError != nil {
		log.Error("Could not sign a transaction", "error", signatureError)
		return signatures
	}

	signatures = append(signatures, action.TransactionSignature{signature})

	return signatures
}

func HandleAccountPublicKeyQuery(app Application, arguments map[string]string) interface{} {
	log.Debug("AccountPublicKeyQuery", "arguments", arguments)

	text := arguments["parameters"]

	account, accountStatus := app.Accounts.FindKey([]byte(text))

	if accountStatus != status.SUCCESS {
		log.Error("Could not find an account", "status", accountStatus)
		return []byte{}
	}

	privateKey := account.PrivateKey()
	publicKey := privateKey.PubKey()

	if publicKey.Equals(account.PublicKey()) == false {
		log.Warn("Public keys don't match", "derivedPublicKey",
			publicKey, "accountPublicKey", account.PublicKey())
	}
	return publicKey
}

// Get the account information for a given user
func HandleAccountKeyQuery(app Application, arguments map[string]string) interface{} {
	log.Debug("AccountKeyQuery", "arguments", arguments)

	text := arguments["parameters"]

	name := ""
	parts := strings.Split(text, "=")
	if len(parts) > 1 {
		name = parts[1]
	}
	return AccountKey(app, name)
}

func AccountKey(app Application, name string) id.AccountKey {

	// Check Itdentities First
	identity, ok := app.Identities.FindName(name)
	if ok == status.SUCCESS && identity.Name != "" {
		return identity.AccountKey
	}

	// TODO: This is a bit dangerous (can cause confusion)
	// Maybe this is an AccountName, not an identity
	account, ok := app.Accounts.FindName(name)
	if ok == status.SUCCESS && account.Name() != "" {
		return account.AccountKey()
	}
	return nil
}

// Get the account information for a given user
func HandleIdentityQuery(app Application, arguments map[string]string) interface{} {
	log.Debug("IdentityQuery", "arguments", arguments)

	text := arguments["parameters"]

	name := ""
	parts := strings.Split(text, "=")
	if len(parts) > 1 {
		name = parts[1]
	}
	return IdentityInfo(app, name)
}

func IdentityInfo(app Application, name string) interface{} {
	if name == "" {
		identities := app.Identities.FindAll()
		return identities
	}

	identity, ok := app.Identities.FindName(name)
	if ok == status.SUCCESS {
		return []id.Identity{identity}
	}
	return "Identity " + name + " Not Found" + global.Current.NodeName
}

// Get the account information for a given user
func HandleAccountQuery(app Application, arguments map[string]string) interface{} {
	log.Debug("AccountQuery", "arguments", arguments)

	text := arguments["parameters"]

	name := ""
	parts := strings.Split(text, "=")
	if len(parts) > 1 {
		name = parts[1]
	}
	return AccountInfo(app, name)
}

// AccountInfo returns the information for a given account
func AccountInfo(app Application, name string) interface{} {
	if name == "" {
		accounts := app.Accounts.FindAll()
		return accounts
	}

	account, ok := app.Accounts.FindName(name)
	if ok == status.SUCCESS {
		return account
	}
	return "Account " + name + " Not Found" + global.Current.NodeName
}

func HandleVersionQuery(app Application, arguments map[string]string) interface{} {
	return version.Current.String()
}

// Get the account information for a given user
func HandleBalanceQuery(app Application, arguments map[string]string) interface{} {
	log.Debug("BalanceQuery", "arguments", arguments)

	text := arguments["parameters"]

	var key []byte
	parts := strings.Split(text, "=")
	if len(parts) > 1 {

		// TODO: Encoded because it is dumped into a string
		key, _ = hex.DecodeString(parts[1])
	}
	return Balance(app, key)
}

func Balance(app Application, accountKey []byte) interface{} {
	balance := app.Balances.Get(accountKey)
	if balance != nil {
		log.Dump("###### FOUND BALANCE #########", balance)
		return balance
	}
	result := data.NewBalance()
	log.Dump("###### NEW BALANCE #########", result)

	return &result
}

func HandleCurrencyAddressQuery(app Application, arguments map[string]string) interface{} {
	log.Debug("CurrencyAddressQuery", "arguments", arguments)

	text := arguments["parameters"]
	conv := convert.NewConvert()
	var chain data.ChainType
	var identity id.Identity
	parts := strings.Split(text, "|")
	if len(parts) == 2 {
		currencyString := strings.Split(parts[0], "=")
		if len(currencyString) > 1 {
			chain = conv.GetChainFromCurrency(currencyString[1])
		} else {
			return nil
		}

		identityString := strings.Split(parts[1], "=")
		if len(identityString) > 1 {
			identity, ok := app.Identities.FindName(identityString[1])
			if ok == status.SUCCESS {
				if chain != data.ONELEDGER {
					//todo : right now, the idenetity.Chain do not contain the address on the other chain
					// need to fix register to remove this part
					_ = identity
					return ChainAddress(chain)
				}

			}
		}

	}
	return identity.Chain[chain]
}

func CurrencyAddress(application Application, currency string, identityName string) []byte {
	conv := convert.NewConvert()

	chain := conv.GetChainFromCurrency(currency)

	identity, ok := application.Identities.FindName(identityName)

	if ok == status.SUCCESS {
		if chain != data.ONELEDGER {
			//todo : right now, the idenetity.Chain do not contain the address on the other chain, need to fix register to remove this part
			_ = identity
			return ChainAddress(chain).([]byte)
		}
	}

	return identity.Chain[chain]
}

func ChainAddress(chain data.ChainType) interface{} {
	return common.GetChainAddress(chain)
}

// Return a nicely formatted error message
func HandleError(text string, path string, arguments map[string]string) interface{} {
	// TODO: Add in arguments to output
	return "Unknown Query " + text + " " + path
}

func SignTransaction(transaction action.Transaction, application Application) action.SignedTransaction {
	packet, err := serial.Serialize(transaction, serial.CLIENT)

	signed := action.SignedTransaction{transaction, nil}

	if err != nil {
		log.Error("Failed to Serialize packet: ", "error", err)
	} else {
		request := action.Message(packet)
		arguments := map[string]string{
			"parameters": string(request),
		}

		response := HandleSignTransaction(application, arguments)

		if response == nil {
			log.Warn("Query returned no signature", "request", request)
		} else {
			signed.Signatures = response.([]action.TransactionSignature)
		}
	}

	return signed
}

func GetSigners(owner []byte, application Application) []id.PublicKey {
	log.Debug("GetSigners", "owner", owner)

	arguments := map[string]string{
		"parameters": string(owner),
	}
	publicKey := HandleAccountPublicKeyQuery(application, arguments)

	if publicKey == nil {
		return nil
	}

	return []id.PublicKey{publicKey.(id.PublicKey)}
}

// Get the account information for a given user
func HandleSequenceNumberQuery(app Application, arguments map[string]string) interface{} {
	log.Debug("SequenceNumberQuery", "arguments", arguments)

	text := arguments["parameters"]

	var key []byte
	parts := strings.Split(text, "=")
	if len(parts) > 1 {

		// TODO: Encoded because it is dumped into a string
		key, _ = hex.DecodeString(parts[1])
	}
	return SequenceNumber(app, key)
}

func SequenceNumber(app Application, accountKey []byte) id.SequenceRecord {
	sequenceRecord := id.NextSequence(&app, accountKey)
	return sequenceRecord
}

func HandleTestScript(app Application, arguments map[string]string) interface{} {
	log.Debug("TestScript", "arguments", arguments)
	text := arguments["parameters"]

	results := RunTestScriptName(text)

	return results
}
