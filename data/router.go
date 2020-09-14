package data

import (
	"errors"
	"github.com/Oneledger/protocol/storage"
)

type Type string

var (
	errInvalidInput = errors.New("invalid store type or interface")
	errKeyNotFound  = errors.New("store does not exist for key")
)

// Router interface supplies functionality to add a store to the Data package.
type Router interface {
	Add(Type, ExtStore) error
	Get(Type) (ExtStore, error)
}

var _ Router = &StorageRouter{}

type StorageRouter struct {
	router map[Type]ExtStore
}

// Add a new store to the router
func (s StorageRouter) Add(storeType Type, storeObj ExtStore) error {
	if storeType == "" || storeObj == nil {
		return errInvalidInput
	}
	s.router[storeType] = storeObj
	return nil
}

// Get the structure of store ,using the Type
func (s StorageRouter) Get(storeType Type) (ExtStore, error) {
	if store, ok := s.router[storeType]; ok {
		return store, nil
	}
	return nil, errKeyNotFound
}

func NewStorageRouter() StorageRouter {
	return StorageRouter{
		router: make(map[Type]ExtStore),
	}
}

func (s *StorageRouter) WithState(state *storage.State) *StorageRouter {
	for _, v := range s.router {
		v.WithState(state)
	}
	return s
}

type ExtStore interface{
	WithState(*storage.State) ExtStore
}
