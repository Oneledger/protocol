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

// A user of a OneLedger node, but not necessarily the chain itself.
type Identity struct {
	UserName    string
	ContactInfo string
	Primary     Account
	Secondary   []Account
	Nodes       map[string]data.ChainNode
}

// Initialize or reconnect to the database
func NewIdentities(name string) *Identities {
	data := data.NewDatastore(name, data.PERSISTENT)

	return &Identities{
		data: data,
	}
}

func (ids *Identities) Add(identity *Identity) {
	buffer, err := comm.Serialize(identity)
	key := identity.Key()

	if err != nil {
		log.Error("Serialize Failed", "err", err)
		return
	}
	ids.data.Store(key, buffer)
	ids.data.Commit()
}

func (ids *Identities) Delete() {
}

func (ids *Identities) Exists(name string) bool {
	id := NewIdentity(name, "")

	value := ids.data.Load(id.Key())
	if value != nil {
		log.Debug("Identity Exists", "value", value)
		return true
	}
	log.Debug("Identity Does not Exist", "name", name)
	return false
}

func (ids *Identities) Find(name string) (*Identity, err.Code) {
	return nil, err.SUCCESS
}

func (ids *Identities) FindAll() []*Identity {
	keys := ids.data.List()
	size := len(keys)
	results := make([]*Identity, size, size)
	for i := 0; i < size; i++ {
		identity := &Identity{}
		base, _ := comm.Deserialize(ids.data.Load(keys[i]), identity)
		results[i] = base.(*Identity)
	}
	return results
}

func (ids *Identities) Dump() {
	list := ids.FindAll()
	size := len(list)
	for i := 0; i < size; i++ {
		identity := list[i]
		log.Info("Entry", "UserName", identity.UserName)
	}
}

func NewIdentity(userName string, contactInfo string) *Identity {
	return &Identity{
		UserName:    userName,
		ContactInfo: contactInfo,
	}
}

func (id *Identity) Key() data.DatabaseKey {
	return data.DatabaseKey(id.UserName)
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
