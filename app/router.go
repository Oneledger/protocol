package app

import (
	"os"

	"github.com/pkg/errors"

	"github.com/Oneledger/protocol/log"
)

var (
	errInvalidInput = errors.New("invalid store type or interface")
	errKeyNotFound  = errors.New("store does not exist for key")
)

type Type string

// Router interface supplies functionality to add a function to the blockender and blockbeginner
type Router interface {
	AddBlockBeginner(Type, func(app *App)) error
	GetBlockBeginner(Type) (func(app *App), error)
	AddBlockEnder(Type, func(app *App)) error
	GetBlockEnder(Type) (func(app *App), error)
	IterateBlockBeginner() []func(app *App)
}

// router is an implementation of a Router interface
type router struct {
	blockbeginnerFn map[Type]func(app *App)
	blockenderFn    map[Type]func(app *App)
	logger          *log.Logger
}

func (r router) AddBlockBeginner(t Type, i func(app *App)) error {
	if t == "" || i == nil {
		return errInvalidInput
	}
	r.blockbeginnerFn[t] = i
	return nil
}

func (r router) GetBlockBeginner(t Type) (func(app *App), error) {
	if store, ok := r.blockbeginnerFn[t]; ok {
		return store, nil
	}
	return nil, errKeyNotFound
}

func (r router) AddBlockEnder(t Type, i func(app *App)) error {
	if t == "" || i == nil {
		return errInvalidInput
	}
	r.blockenderFn[t] = i
	return nil
}

func (r router) GetBlockEnder(t Type) (func(app *App), error) {
	if store, ok := r.blockbeginnerFn[t]; ok {
		return store, nil
	}
	return nil, errKeyNotFound
}

func (r router) IterateBlockBeginner() []func(app *App) {
	var functionlist []func(app *App)
	for _, v := range r.blockbeginnerFn {
		functionlist = append(functionlist, v)
	}
	return functionlist
}

func NewRouter() router {
	return router{
		blockbeginnerFn: make(map[Type]func(app *App)),
		blockenderFn:    make(map[Type]func(app *App)),
		logger:          log.NewLoggerWithPrefix(os.Stdout, "app/router"),
	}
}

// router implements Router
var _ Router = &router{}
