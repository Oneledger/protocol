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
	"github.com/Oneledger/protocol/log"
	"github.com/Oneledger/protocol/serialize"
	"os"
)

// Router interface supplies functionality to add a handler function and
// Handle a request.
type Router interface {
	AddHandler(Type, Tx) error
	Handler([]byte) Tx
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
	return &router{name, map[Type]Tx{}, log.NewLoggerWithPrefix(os.Stdout, "action/router")}
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
func (r *router) Handler(msg []byte) Tx {
	var tx BaseTx

	err := serialize.GetSerializer(serialize.PERSISTENT).Deserialize(msg, tx)
	if err != nil {
		r.logger.Errorf("failed to deserialize msg: %s, error: %s ", msg, err)
	}

	data := tx.Data

	h, ok := r.routes[data.Type()]
	if !ok {
		r.logger.Error("handler not found", tx)
	}

	return h
}
