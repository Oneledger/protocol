/*

 */

package ons

import (
	"github.com/Oneledger/protocol/serialize"
	"github.com/Oneledger/protocol/storage"
	"github.com/pkg/errors"
)

// DomainStore wraps the persistent storage and the serializer giving
// handy methods to access Domain objects
type DomainStore struct {
	*storage.ChainState
	szlr serialize.Serializer
}

// NewDomainStore creates a new storage object from filepath and other configurations
func NewDomainStore(name, dbDir, configDB string, typ storage.StorageType) *DomainStore {
	cs := storage.NewChainState(name, dbDir, configDB, typ)

	return &DomainStore{cs,
		serialize.GetSerializer(serialize.PERSISTENT)}
}

// Get is used to retrieve the domain object from the domain name
func (ds *DomainStore) Get(name string) (*Domain, error) {

	key := keyFromName(name)

	exists := ds.ChainState.Exists(key)
	if !exists {
		return nil, errors.New("Domain doesn't exist")
	}

	data := ds.ChainState.Get(key, true)

	d := &Domain{}
	err := ds.szlr.Deserialize(data, d)
	if err != nil {
		return nil, err
	}

	return d, nil
}

func (ds *DomainStore) Set(d *Domain) error {
	key := []byte(d.Name)

	data, err := ds.szlr.Serialize(d)
	if err != nil {
		return err
	}

	err = ds.ChainState.Set(key, data)
	if err != nil {
		return err
	}

	return nil
}

func (ds *DomainStore) Exists(name string) bool {
	key := keyFromName(name)
	return ds.ChainState.Exists(key)
}

func keyFromName(name string) []byte {

	return []byte(name)
}
