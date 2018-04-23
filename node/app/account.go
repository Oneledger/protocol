/*
	Copyright 2017-2018 OneLedger

	Current state of a given user, assembled from persistence
*/
package app

// The persistent collection of all accounts known by this node
type Accounts struct {
	accounts *Datastore
}

// Initialize or reconnect to the database
func NewAccounts(name string) *Accounts {
	accounts := NewDatastore(name, PERSISTENT)
	return &Accounts{accounts: accounts}
}

// Given an identity, get the account
func GetAccount(identity Identity) (string, Error) {
	return identity.Name(), 0
}
