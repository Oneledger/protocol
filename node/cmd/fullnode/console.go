/*
	Copyright 2017-2018 OneLedger

	Encapsulate any terminal handling, to allow scripting later.
*/
package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"

	"github.com/Oneledger/protocol/node/log"
	isatty "github.com/mattn/go-isatty"
)

// Simple Terminal operations
type Terminal interface {
	Print(text ...interface{})
	Read() string
}

// Enforce a type underneath
type Tty struct {
}

// Declare a global Console that is usable from all of the cmds
var Console Terminal

func init() {
	Console = NewTty()
}

func NewTty() *Tty {
	return &Tty{}
}

// Print to a terminal, add a layer of presentation on top.
func (tty *Tty) Print(text ...interface{}) {
	fmt.Println(text...)
}

// Test to see if the process is still connected to a terminal
func inputIsTty(buf *bufio.Reader) bool {
	if isatty.IsTerminal(os.Stdin.Fd()) {
		return true
	}

	// Windows portability
	if isatty.IsCygwinTerminal(os.Stdin.Fd()) {
		return true
	}

	return false
}

// Read in user input from a terminal, if we are connected to one.
func (tty *Tty) Read() string {
	var buffer []byte
	var size int
	var err error

	input := bufio.NewReader(os.Stdin)
	if inputIsTty(input) {
		size, err = os.Stdin.Read(buffer)
		if err != nil {
			log.Error("Input Error", "status", err)
			// Go down hard, input is broken.
			panic(err)
		}
		if size == 0 {
			log.Error("Empty Input")
			panic(errors.New("Missing Input"))
		}
	} else {
		return "missing input"
	}

	return string(buffer)
}
