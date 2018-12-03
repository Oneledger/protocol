/*
	Copyright 2017-2018 OneLedger

	Implement all of the query mechanics for the node and the chain
*/
package app

import (
	"encoding/hex"
	"github.com/Oneledger/protocol/node/comm"
	"strings"

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

	case "/createSendRequest":
		result = HandleCreateSendRequest(app, arguments)

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

	identities := IdentityInfo(application, idName)
	identity := identities.([]id.Identity)[0]

	if amount == "" {
		log.Error("Missing an amount argument")
		return result
	}

	balance := Balance(application, identity.AccountKey).(*data.Balance).GetAmountByName("VT")

	log.Debug("HandleApplyValidatorQuery", "balance", balance)

	if &balance == nil {
		log.Error("Missing Balance")
		return result
	}

	stake := conv.GetCoin(amount, "VT")
	sequence := SequenceNumber(application, identity.AccountKey.Bytes())

	validator := &action.ApplyValidator{
		Base: action.Base{
			Type:     action.APPLY_VALIDATOR,
			ChainId:  ChainId,
			Owner:    identity.AccountKey,
			Signers:  GetSigners(identity.AccountKey, application),
			Sequence: sequence.(SequenceRecord).Sequence,
		},

		AccountKey:        identity.AccountKey,
		TendermintAddress: identity.TendermintAddress,
		TendermintPubKey:  identity.TendermintPubKey,
		Stake:             stake,
	}

	signed := SignTransaction(action.Transaction(validator), application)
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
		log.Error("Could not deserialize ApplyValidatorArguments", "error", err)
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
	payment := AccountKey(application, "Payment")

	if party == nil || counterParty == nil {
		log.Fatal("System doesn't recognize the parties", "args", args,
			"party", party, "counterParty", counterParty)
	}

	if args.Currency == "" || args.Amount == "" {
		log.Error("Missing an amount argument")
		return result
	}

	amount := conv.GetCoin(args.Amount, args.Currency)

	partyBalance := Balance(application, party).(*data.Balance).GetAmountByName(args.Currency)
	counterPartyBalance := Balance(application, counterParty).(*data.Balance).GetAmountByName(args.Currency)
	paymentBalance := Balance(application, payment).(*data.Balance).GetAmountByName(args.Currency)

	if &partyBalance == nil || &counterPartyBalance == nil {
		log.Error("Missing Balance", "party", partyBalance, "counterParty", counterPartyBalance)
		return result
	}

	fee := conv.GetCoin(args.Fee, args.Currency)
	gas := conv.GetCoin(args.Gas, args.Currency)

	inputs := make([]action.SendInput, 0)
	inputs = append(inputs,
		action.NewSendInput(party, partyBalance),
		action.NewSendInput(counterParty, counterPartyBalance),
		action.NewSendInput(payment, paymentBalance))

	// Build up the outputs
	outputs := make([]action.SendOutput, 0)
	outputs = append(outputs,
		action.NewSendOutput(party, partyBalance.Minus(amount).Minus(fee)),
		action.NewSendOutput(counterParty, counterPartyBalance.Plus(amount)),
		action.NewSendOutput(payment, paymentBalance.Plus(fee)))

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
			Sequence: sequence.(SequenceRecord).Sequence,
		},
		Inputs:  inputs,
		Outputs: outputs,
		Fee:     fee,
		Gas:     gas,
	}

	signed := SignTransaction(action.Transaction(send), application)
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
		return balance
	}
	result := data.NewBalance()
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

func SequenceNumber(app Application, accountKey []byte) interface{} {
	sequenceRecord := NextSequence(&app, accountKey)
	return sequenceRecord
}
