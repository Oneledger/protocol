package app

import (
	"github.com/pkg/errors"
)

var (
	errInvalidInput = errors.New("invalid store type or interface")
)

type txblock int8

const (
	BlockBeginner txblock = 1
	BlockEnder    txblock = 2
)

type Cfunction struct {
	function      func(interface{})
	functionParam interface{}
}

// Router interface supplies functionality to add a function to the blockender and blockbeginner
type Router interface {
	Add(txblock, Cfunction) error
	Iterate(txblock) ([]Cfunction, error)
}

type ControllerRouter struct {
	functionlist map[txblock][]Cfunction
}

func (r ControllerRouter) Add(t txblock, i Cfunction) error {
	if t != BlockBeginner && t != BlockEnder {
		return errInvalidInput
	}
	r.functionlist[t] = append(r.functionlist[t], i)
	return nil
}

func (r ControllerRouter) Iterate(t txblock) ([]Cfunction, error) {
	if t != BlockBeginner && t != BlockEnder {
		return []Cfunction{}, errInvalidInput
	}
	return r.functionlist[t], nil
}

func NewRouter() ControllerRouter {
	return ControllerRouter{
		functionlist: make(map[txblock][]Cfunction),
	}
}

var _ Router = &ControllerRouter{}
