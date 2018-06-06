/*
	Copyright 2017-2018 OneLedger

	Implement all of the query mechanics for the node and the chain
*/
package app

import (
	"fmt"
	"strings"

	"github.com/Oneledger/protocol/node/log"
)

// Top-level list of all query types
func HandleQuery(app Application, path string, message []byte) []byte {

	switch path {
	case "/identity":
		return HandleIdentityQuery(app, message)

	case "/account":
		return HandleAccountQuery(app, message)
	}

	return HandleError("Unknown Path", path, message)
}

// Get the account information for a given user
func HandleIdentityQuery(app Application, message []byte) []byte {
	log.Debug("IdentityQuery", "message", message)

	text := string(message)

	parts := strings.Split(text, "=")
	if len(parts) == 0 {
		ids := app.Identities.FindAll()
		buffer := ""
		for _, curr := range ids {
			buffer += curr.AsString()
		}
		return []byte(buffer)

	} else if len(parts) == 2 {
		if parts[0] == "Identity" {
			return []byte("Identity Information")
		}
	}
	return []byte("Unknown Identity Query")
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

func AccountInfo(app Application, name string) []byte {
	if name == "" || name == "undefined" {
		accounts := app.Accounts.FindAll()

		count := fmt.Sprintf("%d", len(accounts))
		buffer := "Answer: " + count + " "

		for _, curr := range accounts {
			buffer += curr.AsString() + ", "
		}
		return []byte(buffer)
	}
	account, _ := app.Accounts.Find(name)

	return []byte(account.AsString())
}

// Return a nicely formatted error message
func HandleError(text string, path string, massage []byte) []byte {
	return []byte("Invalid Query")
}
