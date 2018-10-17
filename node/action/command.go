/*
	Copyright 2017 - 2018 OneLedger
*/

package action

import (
	"github.com/Oneledger/protocol/node/data"
	"github.com/Oneledger/protocol/node/status"
)

type CommandType int

// Set of possible commands that can be driven from a transaction
const (
	NOOP CommandType = iota
)

type FunctionValue interface{}

type FunctionValues map[Parameter]FunctionValue

// A command to execute again a chain, needs to be polymorphic
type Command struct {
	opfunc func(app interface{}, chain data.ChainType, data FunctionValues) (bool, FunctionValues)
	chain  data.ChainType
	data   FunctionValues
}

func (command Command) Execute(app interface{}) (bool, FunctionValues) {
	return command.opfunc(app, command.chain, command.data)
}

type Commands []Command

func (commands Commands) Count() int {
	return len(commands)
}

func (cs Commands) Execute(app interface{}) status.Code {
	var lastResult FunctionValues
	var ok bool

	for i := 0; i < cs.Count(); i++ {
		ok, lastResult = cs[i].Execute(app)
		if !ok {
			return status.EXECUTE_ERROR
		}

		if len(lastResult) > 0 {
			cs[i+1].data = lastResult
		}
	}
	return status.SUCCESS
}
