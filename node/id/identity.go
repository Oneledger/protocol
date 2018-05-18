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
type Identities struct {
	data *data.Datastore
}

type Identity struct {
	UserId      string
	UserName    string
	ContactInfo string
	Primary     Account
	Secondary   []Account
}

// Initialize or reconnect to the database
func NewIdentities(name string) *Identities {
	data := data.NewDatastore(name, data.PERSISTENT)

	return &Identities{
		data: data,
	}
}

func NewIdentity(userId string, userName string, contactInfo string) *Identity {
	return &Identity{
		UserId:      userId,
		UserName:    userName,
		ContactInfo: contactInfo,
	}
}

func (ids *Identities) AddIdentity(identity Identity) {
	buffer, err := comm.Serialize(identity)
	if err != nil {
		log.Error("Serialize Failed", "err", err)
		return
	}
	ids.data.Store(identity.Key(), buffer)
}

func (ids *Identities) DeleteAccount() {
}

func (ids *Identities) FindIdentity(name string) (*Identity, err.Code) {
	return nil, err.SUCCESS
}

func (ids *Identities) AllIdentities() []Identity {
	return nil
}

func (id *Identity) Key() []byte {
	return []byte(id.UserId)
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
