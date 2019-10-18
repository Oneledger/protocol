/*

 */

package ons

import (
	"strings"

	"github.com/Oneledger/protocol/serialize"
	"github.com/Oneledger/protocol/storage"
	"github.com/pkg/errors"
)

// DomainStore wraps the persistent storage and the serializer giving
// handy methods to access Domain objects
type DomainStore struct {
	State  *storage.State
	szlr   serialize.Serializer
	prefix []byte
}

// NewDomainStore creates a new storage object from filepath and other configurations
func NewDomainStore(prefix string, state *storage.State) *DomainStore {

	return &DomainStore{
		State:  state,
		szlr:   serialize.GetSerializer(serialize.PERSISTENT),
		prefix: storage.Prefix(prefix),
	}
}

func (ds *DomainStore) WithState(state *storage.State) *DomainStore {
	ds.State = state
	return ds
}

// Get is used to retrieve the domain object from the domain name
func (ds *DomainStore) Get(name string) (*Domain, error) {
	key := keyFromName(name)
	key = append(ds.prefix, key...)
	exists := ds.State.Exists(key)
	if !exists {
		return nil, ErrDomainNotFound
	}

	data, _ := ds.State.Get(key)

	d := &Domain{}
	err := ds.szlr.Deserialize(data, d)
	if err != nil {
		return nil, errors.Wrap(err, "error de-serializing domain")
	}

	return d, nil
}

func (ds *DomainStore) Set(d *Domain) error {
	key := keyFromName(d.Name)

	data, err := ds.szlr.Serialize(d)
	if err != nil {
		return err
	}

	key = append(ds.prefix, key...)
	err = ds.State.Set(key, data)
	if err != nil {
		return err
	}

	return nil
}

func (ds *DomainStore) Exists(name string) bool {
	key := keyFromName(name)
	key = append(ds.prefix, key...)
	return ds.State.Exists(key)
}

func keyFromName(name string) []byte {

	return []byte(strings.ToLower(name))
}

func (ds *DomainStore) Iterate(fn func(name string, domain *Domain) bool) (stopped bool) {
	return ds.State.IterateRange(
		ds.prefix,
		storage.Rangefix(string(ds.prefix)),
		true,
		func(key, value []byte) bool {
			name := string(key[len(ds.prefix):])
			domain := &Domain{}
			err := ds.szlr.Deserialize(value, domain)
			if err != nil {
				return false
			}
			return fn(name, domain)
		},
	)
}
