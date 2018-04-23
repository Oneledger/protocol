/*
	Copyright 2017-2018 OneLedger

	Encapsulate any terminal handling, to allow scripting later.
*/
package main

import "fmt"

type Terminal interface {
	Print(text string)
	Read() string
}

type Tty struct {
}

var Console Terminal

func init() {
	Console = NewTty()
}

func NewTty() *Tty {
	return &Tty{}
}

// TODO: Add varargs, pretty formatting, logging and detect disconnected terminals
func (tty *Tty) Print(text string) {
	fmt.Println(text)
}

// TODO: Catch a disconnected terminal, maybe read input from files?
func (tty *Tty) Read() string {
	return "missing"
}
