/*
	Copyright 2017 - 2018 OneLedger
*/

package action

import (
	"github.com/Oneledger/protocol/node/data"
	"github.com/Oneledger/protocol/node/serial"
)

type CommandType int

// Set of possible commands that can be driven from a transaction
const (
	NOOP CommandType = iota
	PREPARE_TRANSACTION
	SUBMIT_TRANSACTION
	INITIATE
	PARTICIPATE
	REDEEM
	REFUND
	EXTRACTSECRET
	AUDITCONTRACT
	WAIT_FOR_CHAIN
	FINISH
)

type FunctionValue interface{}

// A command to execute again a chain, needs to be polymorphic
type Command struct {
	Function CommandType
	Chain    data.ChainType
	Data     map[Parameter]FunctionValue
	Order    int
}

func init() {
	serial.Register(Command{})
}

func (command Command) Execute(app interface{}) (bool, map[Parameter]FunctionValue) {
	switch command.Function {
	case NOOP:
		return Noop(app, command.Chain, command.Data)

	case PREPARE_TRANSACTION:
		return PrepareTransaction(app, command.Chain, command.Data)

	case SUBMIT_TRANSACTION:
		return SubmitTransaction(app, command.Chain, command.Data)

	case INITIATE:
		return Initiate(app, command.Chain, command.Data)

	case PARTICIPATE:
		return Participate(app, command.Chain, command.Data)

	case REDEEM:
		return Redeem(app, command.Chain, command.Data)

	case REFUND:
		return Refund(app, command.Chain, command.Data)

	case EXTRACTSECRET:
		return ExtractSecret(app, command.Chain, command.Data)

	case AUDITCONTRACT:
		return AuditContract(app, command.Chain, command.Data)

	case WAIT_FOR_CHAIN:
		return WaitForChain(app, command.Chain, command.Data)
	}

	return true, nil
}

type Commands []Command

func (commands Commands) Count() int {
	return len(commands)
}
