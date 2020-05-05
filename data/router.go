package data

import "errors"

type Type string

var (
	errInvalidInput = errors.New("invalid store type or interface")
	errKeyNotFound  = errors.New("store does not exist for key")
)

// Router interface supplies functionality to add a store to the Data package.
type Router interface {
	Add(Type, interface{}) error
	Get(Type) (interface{}, error)
}

var _ Router = &StorageRouter{}

type StorageRouter struct {
	router map[Type]interface{}
}

// Add a new store to the router
func (s StorageRouter) Add(storeType Type, storeObj interface{}) error {
	if storeType == "" || storeObj == nil {
		return errInvalidInput
	}
	s.router[storeType] = storeObj
	return nil
}

// Get the structure of store ,using the Type
func (s StorageRouter) Get(storeType Type) (interface{}, error) {
	if store, ok := s.router[storeType]; ok {
		return store, nil
	}
	return nil, errKeyNotFound
}

func NewStorageRouter() StorageRouter {
	return StorageRouter{
		router: make(map[Type]interface{}),
	}
}
