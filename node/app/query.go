/*
	Copyright 2017-2018 OneLedger

	Implement all of the query mechanics for the node and the chain
*/
package app

import (
	"encoding/hex"
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

func AccountKey(app Application, name string) interface{} {

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
	return "AccountKey: Identity " + name + " not Found on " + global.Current.NodeName
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
