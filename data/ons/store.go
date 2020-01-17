/*

 */

package ons

import (
	"github.com/pkg/errors"

	"github.com/Oneledger/protocol/serialize"
	"github.com/Oneledger/protocol/storage"
)

// DomainStore wraps the persistent storage and the serializer giving
// handy methods to access Domain objects
type DomainStore struct {
	State  *storage.State
	opt    *Options
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

func (ds *DomainStore) SetOptions(opt *Options) {
	ds.opt = opt
}

func (ds *DomainStore) GetOptions() *Options {
	return ds.opt
}

// Get is used to retrieve the domain object from the domain name
func (ds *DomainStore) Get(name Name) (*Domain, error) {
	key := name.toKey()
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
	key := d.Name.toKey()

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

func (ds *DomainStore) Exists(name Name) bool {
	key := name.toKey()
	key = append(ds.prefix, key...)
	return ds.State.Exists(key)
}

func (ds *DomainStore) Iterate(fn func(name Name, domain *Domain) bool) (stopped bool) {
	return ds.State.IterateRange(
		ds.prefix,
		storage.Rangefix(string(ds.prefix)),
		true,
		func(key, value []byte) bool {
			nameKey := string(key[len(ds.prefix):])
			domain := &Domain{}
			err := ds.szlr.Deserialize(value, domain)
			if err != nil {
				return false
			}
			return fn(Name(reverse(nameKey)), domain)
		},
	)
}

func (ds *DomainStore) IterateSubDomain(parentName Name, fn func(name Name, domain *Domain) bool) (stopped bool) {
	start := append(ds.prefix, ("." + parentName).toKey()...)
	end := storage.Rangefix(string(start))
	return ds.State.IterateRange(
		start,
		end,
		true,
		func(key, value []byte) bool {
			name := string(key[len(ds.prefix):])
			domain := &Domain{}
			err := ds.szlr.Deserialize(value, domain)
			if err != nil {
				return false
			}
			return fn(Name(reverse(name)), domain)
		},
	)
}

func (ds *DomainStore) DeleteAllSubdomains(name Name) error {

	ds.IterateSubDomain(name, func(name Name, domain *Domain) bool {

		prefixed := append(ds.prefix, name.toKey()...)
		_, err := ds.State.Delete(prefixed)
		if err != nil {
			return false
		}
		return false
	})
	return nil
}

func (ds *DomainStore) DeleteASubdomain(subdomainName Name) error {
	domain, err := ds.Get(subdomainName)
	if err != nil {
		return err
	}

	if !domain.Name.IsSub() {
		return errors.New("not a subdomain")
	}

	prefixed := append(ds.prefix, subdomainName.toKey()...)
	_, err = ds.State.Delete(prefixed)

	return err
}
