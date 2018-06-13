/*
	Copyright 2017-2018 OneLedger

	The overall running context. Initialized right away, but is mutable.

	Contains the main variables.

	Precedence:
		- Default values
	 	- Environment variables (like $OLROOT)
		- Configuration files
		- Command line arguments
		- Overrides
*/
package global

import (
	"os"

	"github.com/Oneledger/protocol/node/persist"
)

var Current *Context

type Context struct {
	Application persist.Access // Global Access to the application when it is running

	Debug bool // DEBUG flag

	NodeName   string // Name of this instance
	RootDir    string // Working directory for this instance
	AppAddress string // app address
	RpcAddress string // rpc address
	Transport  string // socket vs grpc

	Sequence int // replay protection
	BTCRpcPort int
	ETHRpcPort int
}

func init() {
	Current = NewContext("OneLedger")
}

// Set the default values for any context variables here (and no where else)
func NewContext(name string) *Context {
	return &Context{
		NodeName: name,
		Debug:    false,
		RootDir:  os.Getenv("OLDATA") + "/" + name + "/fullnode",
		Sequence: 101,
	}
}

func (context *Context) SetApplication(app persist.Access) persist.Access {
	context.Application = app
	return app
}

func (context *Context) GetApplication() persist.Access {
	return context.Application
}
