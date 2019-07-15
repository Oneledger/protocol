/*
   ____             _              _                      _____           _                  _
  / __ \           | |            | |                    |  __ \         | |                | |
 | |  | |_ __   ___| |     ___  __| | __ _  ___ _ __     | |__) | __ ___ | |_ ___   ___ ___ | |
 | |  | | '_ \ / _ \ |    / _ \/ _` |/ _` |/ _ \ '__|    |  ___/ '__/ _ \| __/ _ \ / __/ _ \| |
 | |__| | | | |  __/ |___|  __/ (_| | (_| |  __/ |       | |   | | | (_) | || (_) | (_| (_) | |
  \____/|_| |_|\___|______\___|\__,_|\__, |\___|_|       |_|   |_|  \___/ \__\___/ \___\___/|_|
                                      __/ |
                                     |___/


Copyright 2017 - 2019 OneLedger
*/

package action

import (
	"errors"
	"os"

	"github.com/Oneledger/protocol/log"
)

// Router interface supplies functionality to add a handler function and
// Handle a request.
type Router interface {
	AddHandler(Type, Tx) error
	Handler(Type) Tx
}

// router is an implementation of a Router interface, currently all routes are stored in a map
type router struct {
	name   string
	routes map[Type]Tx
	logger *log.Logger
}

// router implements Router
var _ Router = &router{}

// NewRouter creates a new router object with given name.
func NewRouter(name string) Router {
	router := &router{name, map[Type]Tx{}, log.NewLoggerWithPrefix(os.Stdout, "action/router")}

	return router
}

// AddHandler adds a new path to the router alongwith its Handler function
func (r *router) AddHandler(t Type, h Tx) error {

	if _, ok := r.routes[t]; ok {
		return errors.New("duplicate path")
	}

	r.routes[t] = h
	return nil
}

// Handle
func (r *router) Handler(t Type) Tx {

	h, ok := r.routes[t]
	if !ok {
		r.logger.Error("handler not found", t)
	}

	return h
}
