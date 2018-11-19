/*
	Copyright 2017-2018 OneLedger

	Implement all of the query mechanics for the node and the chain
*/
package app

import (
	"encoding/hex"
	"github.com/Oneledger/protocol/node/action"
	"github.com/Oneledger/protocol/node/chains/common"
	"strings"

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
func HandleQuery(app Application, path string, message []byte) (buffer []byte) {

	var result interface{}

	switch path {
	case "/nodeName":
		result = HandleNodeNameQuery(app, message)

	case "/identity":
		result = HandleIdentityQuery(app, message)

	case "/signTransaction":
		result = HandleSignTransaction(app, message)

	case "/accountPublicKey":
		result = HandleAccountPublicKeyQuery(app, message)

	case "/accountKey":
		result = HandleAccountKeyQuery(app, message)

	case "/account":
		result = HandleAccountQuery(app, message)

	case "/balance":
		result = HandleBalanceQuery(app, message)

	case "/utxo":
		result = HandleUtxoQuery(app, message)

	case "/version":
		result = HandleVersionQuery(app, message)

	case "/currencyAddress":
		result = HandleCurrencyAddressQuery(app, message)

	default:
		result = HandleError("Unknown Query", path, message)
	}

	buffer, err := serial.Serialize(result, serial.CLIENT)
	if err != nil {
		log.Debug("Failed to serialize query")
	}

	return
}

func HandleNodeNameQuery(app Application, message []byte) interface{} {
	return global.Current.NodeName
}

func HandleSignTransaction(app Application, message []byte) interface{} {
	log.Debug("SignTransactionQuery", "message", message)

	var tx action.Transaction

	transaction, transactionErr := serial.Deserialize(message, tx, serial.CLIENT)

	signatures := []action.TransactionSignature{}

	if transactionErr != nil {
		log.Error("Could not deserialize a transaction", "error", transactionErr)
		return signatures
	}

	var accountKey id.AccountKey

	switch v := transaction.(type) {
	case *action.Swap:
		accountKey = v.Base.Owner
	case *action.Send:
		accountKey = v.Base.Owner
	case *action.Register:
		accountKey = v.Base.Owner
	case *action.Payment:
		accountKey = v.Base.Owner
	default:
		log.Error("Unknown transaction type", "transaction", transaction)
	}

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

	signature, signatureError := privateKey.Sign(message)

	if signatureError != nil {
		log.Error("Could not sign a transaction", "error", signatureError)
		return signatures
	}

	signatures = append(signatures, action.TransactionSignature{signature})

	return signatures
}

func HandleAccountPublicKeyQuery(app Application, message []byte) interface{} {
	log.Debug("AccountPublicKeyQuery", "message", message)

	account, accountStatus := app.Accounts.FindKey(message)

	if accountStatus != status.SUCCESS {
		log.Error("Could not find an account", "status", accountStatus)
		return []byte{}
	}

	privateKey := account.PrivateKey()
	publicKey := privateKey.PubKey()

	if publicKey.Equals(account.PublicKey()) == false {
		log.Warn("Public keys don't match", "derivedPublicKey", publicKey, "accountPublicKey", account.PublicKey())
	}

	return publicKey
}

// Get the account information for a given user
func HandleAccountKeyQuery(app Application, message []byte) interface{} {
	log.Debug("AccountKeyQuery", "message", message)

	text := string(message)

	name := ""
	parts := strings.Split(text, "=")
	if len(parts) > 1 {
		name = parts[1]
	}
	return AccountKey(app, name)
}

func AccountKey(app Application, name string) interface{} {
	identity, ok := app.Identities.FindName(name)

	if ok == status.SUCCESS && identity.Name != "" {
		return identity.AccountKey
	}

	// Maybe this is an AccountName, not an identity
	account, ok := app.Accounts.FindName(name)
	if ok == status.SUCCESS && identity.Name != "" {
		return account.AccountKey()
	}

	return "Account " + name + " Not Found"
}

// Get the account information for a given user
func HandleIdentityQuery(app Application, message []byte) interface{} {
	log.Debug("IdentityQuery", "message", message)

	text := string(message)

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

	return "Identity " + name + " Not Found"
}

// Get the account information for a given user
func HandleAccountQuery(app Application, message []byte) interface{} {
	log.Debug("AccountQuery", "message", message)

	text := string(message)

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

	return "Account " + name + " Not Found"
}

func HandleUtxoQuery(app Application, message []byte) interface{} {
	log.Debug("UtxoQuery", "message", message)

	text := string(message)

	name := ""
	parts := strings.Split(text, "=")
	if len(parts) > 1 {
		name = parts[1]
	}
	return UtxoInfo(app, name)
}

func UtxoInfo(app Application, name string) interface{} {
	if name == "" {
		entries := app.Utxo.FindAll()
		return entries
	}
	value := app.Utxo.Get(data.DatabaseKey(name))
	return value
}

func HandleVersionQuery(app Application, message []byte) interface{} {
	return version.Current.String()
}

// Get the account information for a given user
func HandleBalanceQuery(app Application, message []byte) interface{} {
	log.Debug("BalanceQuery", "message", message)

	text := string(message)

	var key []byte
	parts := strings.Split(text, "=")
	if len(parts) > 1 {
		key, _ = hex.DecodeString(parts[1])
	}
	return Balance(app, key)
}

func Balance(app Application, accountKey []byte) interface{} {

	balance := app.Utxo.Get(accountKey)
	if balance != nil {
		return balance
	}
	result := data.NewBalance()
	return &result
}

func HandleCurrencyAddressQuery(app Application, message []byte) interface{} {
	log.Debug("CurrencyAddressQuery", "message", message)

	text := string(message)
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
func HandleError(text string, path string, message []byte) interface{} {
	return "Unknown Query " + text + " " + path + " " + string(message)
}

func SignTransaction(transaction action.Transaction, applicaiton Application) action.SignedTransaction {
	packet, err := serial.Serialize(transaction, serial.CLIENT)

	signed := action.SignedTransaction{transaction, nil}

	if err != nil {
		log.Error("Failed to Serialize packet: ", "error", err)
	} else {
		request := action.Message(packet)

		response := HandleSignTransaction(applicaiton, request)

		if response == nil {
			log.Warn("Query returned no signature", "request", request)
		} else {
			signed.Signatures = response.([]action.TransactionSignature)
		}
	}

	return signed
}

func GetSigners(owner []byte, applicaiton Application) []id.PublicKey {
	log.Debug("GetSigners", "owner", owner)

	publicKey := HandleAccountPublicKeyQuery(applicaiton, owner)

	if publicKey == nil {
		return nil
	}

	return []id.PublicKey{publicKey.(id.PublicKey)}
}
