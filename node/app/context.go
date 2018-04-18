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
	Debug     bool
	Name      string
	RootDir   string
	Transport string
	Address   string
}

func init() {
	Current = NewContext("Global")
}

func NewContext(name string) *Context {
	return &Context{
		Name:    name,
		Debug:   false,
		RootDir: os.Getenv("OLROOT"),
	}
}
