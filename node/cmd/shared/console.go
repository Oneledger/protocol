/*
	Copyright 2017-2018 OneLedger

	Encapsulate any reads/writes to a terminal, to allow scripting later.

	Should be separate from logging...
*/
package shared

import (
	"fmt"

	"github.com/Oneledger/protocol/node/global"
	"github.com/Oneledger/protocol/node/log"
	"github.com/bgentry/speakeasy"
)

type Tty struct {
}

type Terminal interface {
	// Output
	//Print(text ...interface{})
	Question(text ...interface{})
	Info(text ...interface{})
	Warning(text ...interface{})
	Error(text ...interface{})

	// Input
	Read(string) string
	Password() string
}

// A globally accessable terminal called Console
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

// Get a password from the console, needs to be attached to work correctly
func (tty *Tty) Password() string {

	// Debugging option to make like easier.
	if global.Current.DisablePasswords {
		return "password"
	}

	input := ""
	isValid := false

	for isValid {
		input = tty.Read("Enter a passpharse")
		if len(input) > 7 {
			isValid = true
		}
	}
	return input
}

// TODO: Catch a disconnected terminal, maybe read input from files?
// TODO: Will need this to handled getting secure passwords at the client...
func (tty *Tty) Read(prompt string) string {
	input, err := speakeasy.Ask(prompt)
	if err != nil {
		log.Fatal("Console Read", "err", err)
	}
	return input
}
