/*
	Copyright 2017 -2018 OneLedger
*/

package persist

type Access interface {
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
