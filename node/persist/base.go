/*
	Copyright 2017 -2018 OneLedger
*/

package persist

type Access interface {
	// Run a smart contract
	RunScript(script string) interface{}

	// Access the databases
	GetAdmin() interface{}
	GetStatus() interface{}
	GetIdentities() interface{}
	GetAccounts() interface{}
	GetBalances() interface{}
	GetChainID() interface{}
	GetEvent() interface{}
	GetContract() interface{}
	GetSmartContract() interface{}
	GetValidators() interface{}
	GetSequence() interface{}
}
