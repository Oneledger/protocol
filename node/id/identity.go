/*
	Copyright 2017-2018 OneLedger

	Current state of a given user, assembled from persistence
*/
package id

import (
	"strings"

	"github.com/Oneledger/protocol/node/data"
	"github.com/Oneledger/protocol/node/log"
	"github.com/Oneledger/protocol/node/serial"
	"github.com/Oneledger/protocol/node/status"
)

// The persistent collection of all accounts known by this node
type Identities struct {
	store data.Datastore
}

// A user of a OneLedger node, but not necessarily the chain itself.
type Identity struct {
	Name string // A unique name for the identity

	NodeName    string // The origin of this account
	ContactInfo string

	AccountKey AccountKey // A key

	External bool
	Chain    map[data.ChainType]AccountKey // TODO: Should be more than one account per chain

	Nodes map[string]data.ChainNode

	TendermintAddress string
	TendermintPubKey  string
}

func init() {
	serial.Register(Identity{})
}

// Initialize or reconnect to the database
func NewIdentities(name string) *Identities {
	store := data.NewDatastore(name, data.PERSISTENT)

	return &Identities{
		store: store,
	}
}

func (ids *Identities) Add(identity Identity) {
	key := identity.Key()
	session := ids.store.Begin()
	session.Set(key, identity)
	session.Commit()
}

func (ids *Identities) Close() {
	ids.store.Close()
}

func (ids *Identities) Delete() {
}

func (ids *Identities) Exists(name string) bool {
	id := NewIdentity(name, "", true, "", nil, "", "")

	value := ids.store.Get(id.Key())
	if value != nil {
		return true
	}

	return false
}

func (ids *Identities) FindName(name string) (Identity, status.Code) {
	value := ids.store.Get(data.DatabaseKey(name))
	if value != nil {
		return value.(Identity), status.SUCCESS
	}
	return Identity{}, status.MISSING_DATA
}

func (ids *Identities) FindAll() []Identity {
	keys := ids.store.FindAll()
	size := len(keys)
	results := make([]Identity, size, size)
	for i := 0; i < size; i++ {
		result := ids.store.Get(keys[i])
		results[i] = result.(Identity)
	}
	return results
}

func (ids *Identities) FindTendermint(tendermintAddress string) Identity {
	keys := ids.store.FindAll()
	size := len(keys)
	for i := 0; i < size; i++ {
		identity := ids.store.Get(keys[i]).(Identity)
		if strings.ToLower(tendermintAddress) == strings.ToLower(identity.TendermintAddress) {
			return identity
		}
	}
	return Identity{}
}

func (ids *Identities) Dump() {
	list := ids.FindAll()
	size := len(list)
	for i := 0; i < size; i++ {
		identity := list[i]
		log.Info("Identity", "Name", identity.Name, "NodeName", identity.NodeName, "AccountKey", identity.AccountKey,
			"TendermintAddress", identity.TendermintAddress, "TendermintPubKey", identity.TendermintPubKey)
	}
}

func NewIdentity(name string, contactInfo string, external bool, nodeName string, accountKey AccountKey, tendermintAddress string, tendermintPubKey string) *Identity {
	return &Identity{
		Name:              name,
		ContactInfo:       contactInfo,
		External:          external,
		NodeName:          nodeName,
		AccountKey:        accountKey,
		Chain:             make(map[data.ChainType]AccountKey, 2),
		TendermintAddress: tendermintAddress,
		TendermintPubKey:  tendermintPubKey,
	}
}

func (id *Identity) SetAccount(chain data.ChainType, account Account) {
	id.Chain[chain] = account.AccountKey()
}

func (id *Identity) IsExternal() bool {
	return id.External
}

func (id *Identity) Key() data.DatabaseKey {
	return data.DatabaseKey(id.Name)
}

//String used in fmt and Dump
func (id *Identity) String() string {
	buffer := ""
	buffer += id.Name
	if id.External {
		buffer += "(External)"
	} else {
		buffer += "(Local) " + id.ContactInfo
	}
	return buffer
}
