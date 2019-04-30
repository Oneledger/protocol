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
	"fmt"
	"github.com/Oneledger/protocol/node/log"
)

type handler func(Request, *Response)

type Router interface {
	AddHandler(query string, h handler) error
	Handle(req Request, resp *Response)
}

type router struct {
	routes map[string]handler
}

func NewRouter() Router {
	return &router{map[string]handler{}}
}

// AddHandler adds a new path to the router alongwith its Handler function
func (r *router) AddHandler(path string, h handler) error {

	if _, ok := r.routes[path]; ok {
		return errors.New("duplicate path")
	}

	r.routes[path] = h
	return nil
}

func (r *router) Handle(req Request, resp *Response) {

	h, ok := r.routes[req.Query]
	if  !ok {
		resp.Data = []byte{}
		resp.ErrorMsg = "path not found"
		resp.Success = false

		log.Error("path not found", )
	}

	h(req, resp)
}


/*

	Example

 */
func RunExample() {

	r := NewRouter()

	err := r.AddHandler("exec", exec)
	if err != nil {
		log.Fatal("router error", err)
	}

	p := map[string]interface{}{
		"name": "exec",
		"number": 1,
	}
	req := NewRequest("/exec", p)
	resp := &Response{}
	r.Handle(*req, resp)
}

func exec(req Request, resp Response) {
	fmt.Println(req.Query)

	fmt.Println("name", req.GetString("name"))
	fmt.Println("not set", req.GetString("not_set"))

	fmt.Println("number", req.GetInt("number"))

	err := resp.JSON("function response")
	if err != nil {
		log.Error("error marshalling response", err)
	}
}

