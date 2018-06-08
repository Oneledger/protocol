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
	Data     map[string]string
	Chain    data.ChainType
}

type Commands []Command

func (commands Commands) Count() int {
	return len(commands)
}
