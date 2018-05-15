/*
	Copyright 2017-2018 OneLedger

	Current state of a given user, assembled from persistence
*/
package id

import (
	"github.com/Oneledger/protocol/node/comm"
	"github.com/Oneledger/protocol/node/data"
	"github.com/Oneledger/protocol/node/err"
	"github.com/Oneledger/protocol/node/log"
)

// The persistent collection of all accounts known by this node
type Accounts struct {
	data *data.Datastore
}

// Initialize or reconnect to the database
func NewAccounts(name string) *Accounts {
	data := data.NewDatastore(name, data.PERSISTENT)

	return &Accounts{data: data}
}

func (acc *Accounts) AddAccount(identity Identity) {
	buffer, err := comm.Serialize(identity)
	if err != nil {
		log.Error("Serialize Failed", "err", err)
		return
	}
	acc.data.Store([]byte(identity.Name()), buffer)
}

func (acc *Accounts) DeleteAccount() {
}

func (acc *Accounts) FindAccount(name string) (Identity, err.Code) {
	return nil, err.SUCCESS
}

func (acc *Accounts) AllAccounts() []Identity {
	return nil
}

/*
func (identity Identity) Format() (string, err.Code) {
	return identity.Format(), err.SUCCESS
}

// Given an identity, get the account
func (identity Identity) GetName() (string, err.Code) {
	return identity.Name(), err.SUCCESS
}
*/
