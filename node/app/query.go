/*
	Copyright 2017-2018 OneLedger

	Implement all of the query mechanics for the node and the chain
*/
package app

import (
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
func HandleQuery(app Application, path string, arguments map[string]string) []byte {

	var result interface{}

	switch path {
	case "/nodeName":
		result = HandleNodeNameQuery(app, arguments)

	case "/identity":
		result = HandleIdentityQuery(app, arguments)

	case "/accountKey":
		result = HandleAccountKeyQuery(app, arguments)

	case "/account":
		result = HandleAccountQuery(app, arguments)

	case "/balance":
		result = HandleBalanceQuery(app, arguments)

	case "/version":
		result = HandleVersionQuery(app, arguments)

	case "/swapAddress":
		result = HandleSwapAddressQuery(app, arguments)

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
		//key, _ = hex.DecodeString(parts[1])
		key = []byte(parts[1])
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

func HandleSwapAddressQuery(app Application, arguments map[string]string) interface{} {
	log.Debug("SwapAddressQuery", "arguments", arguments)

	text := arguments["parameter"]
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
func HandleError(text string, path string, arguments map[string]string) interface{} {
	// TODO: Add in arguments to output
	return "Unknown Query " + text + " " + path
}
