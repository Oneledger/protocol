/*
	Copyright 2017-2018 OneLedger

	The overall running context. Initialized right away, but is mutable.

	Precedence:
		- Default values
		- Configuration files
		- Command line arguments
		- Overrides
*/
package app

type Context struct {
	name string
}

func NewContext(name string) *Context {
	return &Context{
		name: name,
	}
}
