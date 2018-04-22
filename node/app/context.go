/*
	Copyright 2017-2018 OneLedger

	The overall running context. Initialized right away, but is mutable.

	Precedence:
		- Default values
	 	- Environment variables (like $OLROOT)
		- Configuration files
		- Command line arguments
		- Overrides
*/
package app

import "os"

var Current *Context

type Context struct {
	Debug     bool   // DEBUG flag
	Name      string // Name of this instance
	RootDir   string // Working directory for this instance
	Transport string // socket vs grpc
	Address   string // address
}

func init() {
	Current = NewContext("Global")
}

// Set the default values for any context variables here (and no where else)
func NewContext(name string) *Context {
	return &Context{
		Name:    name,
		Debug:   false,
		RootDir: os.Getenv("OLDATA") + "/fullnode",
	}
}
