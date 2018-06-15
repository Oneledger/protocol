/*
	Copyright 2017-2018 OneLedger

	Implement all of the query mechanics for the node and the chain
*/
package app

import (
	"fmt"
	"strings"

	"github.com/Oneledger/protocol/node/data"
	"github.com/Oneledger/protocol/node/id"
	"github.com/Oneledger/protocol/node/log"
	"github.com/Oneledger/protocol/node/version"
)

// Top-level list of all query types
func HandleQuery(app Application, path string, message []byte) []byte {

	switch path {
	case "/identity":
		return HandleIdentityQuery(app, message)

	case "/account":
		return HandleAccountQuery(app, message)

	case "/utxo":
		return HandleUtxoQuery(app, message)

	case "/version":
		return HandleVersionQuery(app, message)

	case "/accountKey":
		return HandleAccountKeyQuery(app, message)
	}

	return HandleError("Unknown Path", path, message)
}

// Get the account information for a given user
func HandleAccountKeyQuery(app Application, message []byte) []byte {
	log.Debug("AccountKeyQuery", "message", message)

	text := string(message)

	name := ""
	parts := strings.Split(text, "=")
	if len(parts) > 1 {
		name = parts[1]
	}
	return AccountKey(app, name)
}

func AccountKey(app Application, name string) []byte {
	identity, _ := app.Identities.FindName(name)

	if identity != nil {
		return identity.AccountKey
	}
	return []byte("Not Found")
}

// Get the account information for a given user
func HandleIdentityQuery(app Application, message []byte) []byte {
	log.Debug("IdentityQuery", "message", message)

	text := string(message)

	name := ""
	parts := strings.Split(text, "=")
	if len(parts) > 1 {
		name = parts[1]
	}
	return IdentityInfo(app, name)
}

func IdentityInfo(app Application, name string) []byte {
	if name == "" {
		identities := app.Identities.FindAll()

		count := fmt.Sprintf("%d", len(identities))
		buffer := "Answer: " + count + " "

		for _, curr := range identities {
			buffer += curr.AsString() + ", "
		}
		return []byte(buffer)
	}
	identity, _ := app.Identities.FindName(name)

	return []byte(identity.AsString())
}

// Get the account information for a given user
func HandleAccountQuery(app Application, message []byte) []byte {
	log.Debug("AccountQuery", "message", message)

	text := string(message)

	name := ""
	parts := strings.Split(text, "=")
	if len(parts) > 1 {
		name = parts[1]
	}
	return AccountInfo(app, name)
}

// Return the information for a given account
func AccountInfo(app Application, name string) []byte {

	var buffer string
	if name == "" {
		accounts := app.Accounts.FindAll()

		count := fmt.Sprintf("%d", len(accounts))
		buffer = "Answer[" + count + "]: "

		for _, curr := range accounts {
			buffer += curr.AsString()
			if curr.Chain() == data.ONELEDGER {
				buffer += " " + GetBalance(app, curr)
			}
			buffer += ", "
		}
		return []byte(buffer)
	}

	account, _ := app.Accounts.FindName(name)
	log.Debug("account", "account", account)

	if account != nil {
		buffer = "Answer[1]: " + account.AsString()
		if account.Chain() == data.ONELEDGER {
			buffer += " " + GetBalance(app, account)
		}
	} else {
		buffer = "Answer[0]: "
	}
	return []byte(buffer)
}

func HandleUtxoQuery(app Application, message []byte) []byte {
	log.Debug("UtxoQuery", "message", message)

	text := string(message)

	name := ""
	parts := strings.Split(text, "=")
	if len(parts) > 1 {
		name = parts[1]
	}
	result := UtxoInfo(app, name)
	log.Debug("Returning", "result", string(result))
	return result
}

func UtxoInfo(app Application, name string) []byte {
	buffer := ""
	if name == "" {
		entries := app.Utxo.FindAll()
		for key, value := range entries {
			account, _ := app.Accounts.FindKey([]byte(key))
			buffer += account.Name() + ":" + value.AsString() + ", "
		}

	} else {
		value := app.Utxo.Find(data.DatabaseKey(name))
		buffer += name + ":" + value.AsString()

	}
	return []byte(buffer)
}

// Get the balancd for an account
func GetBalance(app Application, account id.Account) string {
	result := app.Utxo.Find(account.AccountKey())
	if result == nil {
		log.Debug("Balance Not Found", "key", account.AccountKey())
		return " [nil]"
	}

	return result.AsString()
}

// Return a nicely formatted error message
func HandleError(text string, path string, massage []byte) []byte {
	return []byte("Invalid Query")
}

func HandleVersionQuery(app Application, message []byte) []byte {
	return []byte(version.Current.String())
}
