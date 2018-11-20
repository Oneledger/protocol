/*
	Copyright 2017 - 2018 OneLedger

	Test the features of the underlying code.
	These tests need to be run as the only 'node', or they will fail on db access.
*/
package app

import (
	//"github.com/Oneledger/protocol/node/data"
	//"github.com/Oneledger/protocol/node/id"

	"testing"
)

type Object interface{}

// Control the execution
/*
func TestMain(m *testing.M) {
	flag.Parse()

	// Set the debug flags according to whether the -v flag is set in go test
	if testing.Verbose() {
		log.Debug("DEBUG TURNED ON")
		global.Current.Debug = true
	} else {
		log.Debug("DEBUG TURNED OFF")
		global.Current.Debug = false
	}

	global.Current.RootDir = "./test-db"
	// Run it all.
	code := m.Run()

	err := os.RemoveAll(global.Current.RootDir)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Removed")
	}

	os.Exit(code)
}
*/

// TODO: Remove code that causes this block to panic on duplicate key
//func TestRegister(t *testing.T) {
//	app := NewApplication()
//	fullset := []string{"Alice", "Bob", "Admin", "Alex", "Enrico"}
//
//	for _, current := range fullset {
//		register(app, current, t)
//	}
//}
////
//// Test the local storage of the database.
//func register(app *Application, idName string, t *testing.T) {
//	chains := [][]Object{
//		[]Object{"OneLedger", data.ONELEDGER},
//		[]Object{"Bitcoin", data.BITCOIN},
//		[]Object{"Ethereum", data.ETHEREUM},
//	}
//
//	// Create all of the accounts
//	identity := id.NewIdentity(idName, "Contact Info", false, "Alice-Node", id.AccountKey(nil))
//
//	log.Debug("Adding", "name", idName)
//
//	app.Identities.Add(identity)
//
//	for _, set := range chains {
//		//var key action.PublicKey
//
//		name := set[0]
//		chain := set[1]
//		log.Debug("Adding", "name", name, "chain", chain)
//
//		account := id.NewAccount(chain.(data.ChainType), idName+"-"+name.(string),
//			id.NilPublicKey(), id.NilPrivateKey())
//
//		app.Accounts.Add(account)
//	}
//
//	app.Identities.Close()
//	app.Identities = id.NewIdentities("identities")
//
//	// TODO: Test that they all exist
//	log.Info("Output")
//	app.Identities.Dump()
//	app.Accounts.Dump()
//}

func TestSendData(t *testing.T) {
}

func TestChangeState(t *testing.T) {
}

func TestSendToChains(t *testing.T) {
}

func TestPollChains(t *testing.T) {
}
