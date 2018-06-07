/*
	Copyright 2017-2018 OneLedger

	Encapsulate any reads/writes to a terminal, to allow scripting later.

	Should be separate from logging...
*/
package shared

import "fmt"

type Terminal interface {
	//Print(text ...interface{})
	Question(text ...interface{})
	Info(text ...interface{})
	Warning(text ...interface{})
	Error(text ...interface{})

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

// TODO: Depreciate
func (tty *Tty) Print(text ...interface{}) {
	fmt.Println(text...)
}

func (tty *Tty) Question(text ...interface{}) {
	fmt.Println(text...)
}

func (tty *Tty) Info(text ...interface{}) {
	fmt.Println(text...)
}

func (tty *Tty) Warning(text ...interface{}) {
	fmt.Println(text...)
}

func (tty *Tty) Error(text ...interface{}) {
	fmt.Println(text...)
}

// TODO: Catch a disconnected terminal, maybe read input from files?
func (tty *Tty) Read() string {
	return "missing input"
}
