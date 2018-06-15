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
	Name string // A unique name for the identity

	NodeName    string // The origin of this account
	ContactInfo string

	AccountKey AccountKey // A key

	External   bool
	Additional []Account

	Nodes map[string]data.ChainNode
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
	if err != nil {
		log.Error("Serialize Failed", "err", err)
		return
	}

	key := identity.Key()
	ids.data.Store(key, buffer)
	ids.data.Commit()
}

func (ids *Identities) Close() {
	ids.data.Close()
}

func (ids *Identities) Delete() {
}

func (ids *Identities) Exists(name string) bool {
	id := NewIdentity(name, "", true, "", nil)

	value := ids.data.Load(id.Key())
	if value != nil {
		return true
	}

	return false
}

func (ids *Identities) FindName(name string) (*Identity, err.Code) {
	id := NewIdentity(name, "", true, "", nil)

	value := ids.data.Load(id.Key())
	if value != nil {
		identity := &Identity{}
		base, _ := comm.Deserialize(value, identity)

		return base.(*Identity), err.SUCCESS
	}

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
		log.Info("Identity", "Name", identity.Name, "NodeName", identity.NodeName, "AccountKey", identity.AccountKey)
	}
}

func NewIdentity(name string, contactInfo string, external bool, nodeName string, accountKey AccountKey) *Identity {
	return &Identity{
		Name:        name,
		ContactInfo: contactInfo,
		External:    external,
		NodeName:    nodeName,
		AccountKey:  accountKey,
	}
}

func (id *Identity) IsExternal() bool {
	return id.External
}

func (id *Identity) Key() data.DatabaseKey {
	return data.DatabaseKey(id.Name)
}

func (id *Identity) AsString() string {
	buffer := ""
	buffer += id.Name
	if id.External {
		buffer += "(External)"
	} else {
		buffer += "(Local) " + id.ContactInfo
	}
	return buffer
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
