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

package app

import (
	"errors"
	"github.com/Oneledger/protocol/log"
	"os"
)

type handler func(Request, *Response)

// Router interface supplies functionality to add a handler function and
// Handle a request.
type Router interface {
	AddHandler(query string, h handler) error
	Handle(req Request, resp *Response)
}

// router is an implementation of a Router interface, currently all routes are stored in a map
type router struct {
	name string
	routes map[string]handler
	logger *log.Logger
}
// router implements Router
var _ Router = &router{}

// NewRouter creates a new router object with given name.
func NewRouter(name string) Router {
	return &router{name, map[string]handler{}, log.NewLoggerWithPrefix(os.Stdout, "app/router")}
}

// AddHandler adds a new path to the router alongwith its Handler function
func (r *router) AddHandler(path string, h handler) error {

	if _, ok := r.routes[path]; ok {
		return errors.New("duplicate path")
	}

	r.routes[path] = h
	return nil
}

// Handle
func (r *router) Handle(req Request, resp *Response) {

	h, ok := r.routes[req.Query]
	if !ok {
		resp.Data = []byte{}
		resp.ErrorMsg = "path not found"
		resp.Success = false

		r.logger.Error("path not found", "path", req.Query)
	}

	h(req, resp)
}

