/*
	Copyright 2017 -2018 OneLedger
*/

package persist

type Access interface {
	GetAdmin() interface{}
	GetStatus() interface{}
	GetIdentities() interface{}
	GetAccounts() interface{}
	GetUtxo() interface{}
	GetChainID() interface{}
}
