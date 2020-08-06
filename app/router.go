package app

import (
	"os"

	"github.com/pkg/errors"

	"github.com/Oneledger/protocol/log"
)

var (
	errInvalidInput = errors.New("invalid store type or interface")
)

type txblock int8

const (
	BlockBeginner = 1
	BlockEnder    = 2
)

// Router interface supplies functionality to add a function to the blockender and blockbeginner
type Router interface {
	Add(txblock, func(interface{})) error
	Iterate(txblock) []func(interface{})
}

type ControllerRouter struct {
	functionlist map[txblock][]func(interface{})
	logger       *log.Logger
}

func (r ControllerRouter) Add(t txblock, i func(interface{})) error {
	if t != 1 && t != 2 || i == nil {
		return errInvalidInput
	}
	r.functionlist[t] = append(r.functionlist[t], i)
	return nil
}

func (r ControllerRouter) Iterate(t txblock) []func(interface{}) {
	return r.functionlist[t]
}

func NewRouter() ControllerRouter {
	return ControllerRouter{
		functionlist: make(map[txblock][]func(interface{})),
		logger:       log.NewLoggerWithPrefix(os.Stdout, "app/router"),
	}
}

var _ Router = &ControllerRouter{}
