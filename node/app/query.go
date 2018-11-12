/*
	Copyright 2017-2018 OneLedger

	Implement all of the query mechanics for the node and the chain
*/
package app

import (
	"encoding/hex"
	"strings"

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
func HandleQuery(app Application, path string, message []byte) (buffer []byte) {

	var result interface{}

	switch path {
	case "/nodeName":
		result = HandleNodeNameQuery(app, message)

	case "/identity":
		result = HandleIdentityQuery(app, message)

	case "/accountKey":
		result = HandleAccountKeyQuery(app, message)

	case "/account":
		result = HandleAccountQuery(app, message)

	case "/balance":
		result = HandleBalanceQuery(app, message)

	case "/version":
		result = HandleVersionQuery(app, message)

	case "/swapAddress":
		result = HandleSwapAddressQuery(app, message)

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
	balance := app.Balances.Get(accountKey)
	if balance != nil {
		return balance
	}
	result := data.NewBalance(0, "OLT")
	return &result
}

func HandleSwapAddressQuery(app Application, message []byte) interface{} {
	log.Debug("SwapAddressQuery", "message", message)

	text := string(message)
	conv := convert.NewConvert()
	var chain data.ChainType
	parts := strings.Split(text, "=")
	if len(parts) > 1 {
		chain = conv.GetChain(parts[1])
	}

	//todo: make it general
	if chain == data.ONELEDGER {
		account, e := app.Accounts.FindName(global.Current.NodeAccountName)
		if e == status.SUCCESS {
			return account.AccountKey()
		}
	}

	return SwapAddress(chain)
}

func SwapAddress(chain data.ChainType) interface{} {
	return common.GetSwapAddress(chain)
}

// Return a nicely formatted error message
func HandleError(text string, path string, message []byte) interface{} {
	return "Unknown Query " + text + " " + path + " " + string(message)
}
