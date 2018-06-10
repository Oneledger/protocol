/*
	Copyright 2017 - 2018 OneLedger
*/

package action

import "github.com/Oneledger/protocol/node/data"

type CommandType int

// Set of possible commands that can be driven from a transaction
const (
	NOOP CommandType = iota
	SUBMIT_TRANSACTION
	CREATE_LOCKBOX
	SIGN_LOCKBOX
	VERIFY_LOCKBOX
	SEND_KEY
	READ_CHAIN
	OPEN_LOCKBOX
	WAIT_FOR_CHAIN
)

// A command to execute again a chain, needs to be polymorphic
type Command struct {
	Function CommandType
	Chain    data.ChainType
	Data     map[string]string
}

func (command Command) Execute() bool {
	switch command.Function {
	case NOOP:
		return Noop(command.Chain, command.Data)

	case SUBMIT_TRANSACTION:
		return SubmitTransaction(command.Chain, command.Data)

	case CREATE_LOCKBOX:
		return CreateLockbox(command.Chain, command.Data)

	case SIGN_LOCKBOX:
		return SignLockbox(command.Chain, command.Data)

	case VERIFY_LOCKBOX:
		return VerifyLockbox(command.Chain, command.Data)

	case SEND_KEY:
		return SendKey(command.Chain, command.Data)

	case READ_CHAIN:
		return ReadChain(command.Chain, command.Data)

	case OPEN_LOCKBOX:
		return OpenLockbox(command.Chain, command.Data)

	case WAIT_FOR_CHAIN:
		return WaitForChain(command.Chain, command.Data)
	}

	return true
}

type Commands []Command

func (commands Commands) Count() int {
	return len(commands)
}
