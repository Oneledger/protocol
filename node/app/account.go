/*
	Copyright 2017-2018 OneLedger

	Current state of a given user, assembled from persistence
*/
package app

// The current
type Account struct {
	Name     string
	Identity Identity
	Balance  Coins
	Chains   []Identity
}

// TODO: Set defaults here
func NewAccount() *Account {
	return &Account{}
}

func GetAccount(identity Identity) (*Account, Error) {
	// TODO: Build up the Account information
	return &Account{}, 28
}
