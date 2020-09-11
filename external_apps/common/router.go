package common

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

// Router interface supplies functionality to add a function to the blockender and blockbeginner
type Router interface {
	Add(txblock, func(interface{})) error
	Iterate(txblock) ([]func(interface{}), error)
}

type ControllerRouter struct {
	functionlist map[txblock][]func(interface{})
}

func (r ControllerRouter) Add(t txblock, i func(interface{})) error {
	if t != BlockBeginner && t != BlockEnder || i == nil {
		return errInvalidInput
	}
	r.functionlist[t] = append(r.functionlist[t], i)
	return nil
}

func (r ControllerRouter) Iterate(t txblock) ([]func(interface{}), error) {
	if t != BlockBeginner && t != BlockEnder {
		return nil, errInvalidInput
	}
	return r.functionlist[t], nil
}

func NewRouter() ControllerRouter {
	return ControllerRouter{
		functionlist: make(map[txblock][]func(interface{})),
	}
}

var _ Router = &ControllerRouter{}
