/*
	Copyright 2017 - 2018 OneLedger

	Test the features of the underlying code.
	These tests need to be run as the only 'node', or they will fail on db access.
*/
package app

import (
	"testing"

	"github.com/Oneledger/protocol/node/action"
	"github.com/Oneledger/protocol/node/id"
	"github.com/Oneledger/protocol/node/log"
)

type Object interface{}

func TestRegister(t *testing.T) {
	app := NewApplication()
	fullset := []string{"Alice", "Bob", "Admin", "Alex", "Enrico"}

	for _, current := range fullset {
		register(app, current, t)
	}
}

func register(app *Application, idName string, t *testing.T) {
	chains := [][]Object{
		[]Object{"OneLedger", id.ONELEDGER},
		[]Object{"Bitcoin", id.BITCOIN},
		[]Object{"Ethereum", id.ETHEREUM},
	}

	// Create all of the accounts
	identity := id.NewIdentity(idName, "Contact Info")

	log.Debug("Adding", "name", idName)

	app.Identities.Add(identity)

	for _, set := range chains {
		var key action.PublicKey

		name := set[0]
		chain := set[1]
		log.Debug("Adding", "name", name, "chain", chain)

		account := id.NewAccount(chain.(id.AccountType), idName+"-"+name.(string), key)
		app.Accounts.Add(account)
	}

	// TODO: Test that they all exist
	log.Info("Output")
	app.Identities.Dump()
	app.Accounts.Dump()
}

func TestSendData(t *testing.T) {
}

func TestChangeState(t *testing.T) {
}

func TestSendToChains(t *testing.T) {
}

func TestPollChains(t *testing.T) {
}
