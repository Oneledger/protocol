/*
	Copyright 2017-2018 OneLedger

	Parse the incoming transactions

	TODO: switch from individual wire calls, to reading/writing directly to structs
*/
package action

import (
	"github.com/Oneledger/protocol/node/err"
	"github.com/Oneledger/protocol/node/log"
	wire "github.com/tendermint/go-wire"
)

// Unpack an encoded (wire) message
func UnpackMessage(message Message) (Type, Message) {
	value, size, err := wire.GetVarint(message)
	if err != nil {
		log.Debug("Wire returned an error", "err", err)
		panic("Wire Error")
	}
	if size != 2 {
		log.Debug("Wire returned a bad size", "size", size)
		panic("Sizing Error")
	}
	return Type(value), message[1:]

}

// Parse a message into the appropriate transaction
func Parse(message Message) (Transaction, err.Code) {
	log.Debug("Parsing a Transaction")

	command, body := UnpackMessage(message)

	// TODO: Can I do this with deserialize?
	switch command {

	case SEND:
		action := ParseSend(body)

		return action, err.SUCCESS

	case SWAP:
		action := ParseSwap(body)

		return action, err.SUCCESS

	case EXTERNAL_SEND:
		action := ParseExternalSend(body)

		return action, err.SUCCESS

	case EXTERNAL_LOCK:
		action := ParseExternalLock(body)

		return action, err.SUCCESS

	case PREPARE:
		action := ParsePrepare(body)

		return action, err.SUCCESS

	case COMMIT:
		action := ParseCommit(body)

		return action, err.SUCCESS

	case FORGET:
		action := ParseForget(body)

		return action, err.SUCCESS

	default:
		log.Error("Unknown transaction", "command", command)
	}

	return nil, err.PARSE_ERROR
}

// Parse a send request
func ParseSend(message Message) *Send {
	log.Debug("Have a Send")

	return &Send{
		Base: Base{Type: SEND},
	}
}

// Parse a swap request
func ParseSwap(message Message) *Swap {
	log.Debug("Have a Swap")

	//return &SwapTransaction{Type: SWAP_TRANSACTION}
	return &Swap{
		Base: Base{Type: SWAP},
	}
}

// Parse a send request
func ParseExternalSend(message Message) *ExternalSend {
	log.Debug("Have an ExternalSend")

	return &ExternalSend{
		Base: Base{Type: EXTERNAL_SEND},
	}
}

// Parse a send request
func ParseExternalLock(message Message) *ExternalLock {
	log.Debug("Have an ExternalLock")

	return &ExternalLock{
		Base: Base{Type: EXTERNAL_LOCK},
	}
}

// Parse a ready request
func ParsePrepare(message Message) *Prepare {
	log.Debug("Have a Ready")

	//return &SwapTransaction{Type: SWAP_TRANSACTION}
	return &Prepare{
		Base: Base{Type: PREPARE},
	}
}

// Parse a ready request
func ParseCommit(message Message) *Commit {
	log.Debug("Have a Ready")

	//return &SwapTransaction{Type: SWAP_TRANSACTION}
	return &Commit{
		Base: Base{Type: COMMIT},
	}
}

// Forget the transaction
func ParseForget(message Message) *Forget {
	log.Debug("Have a Ready")

	//return &SwapTransaction{Type: SWAP_TRANSACTION}
	return &Forget{
		Base: Base{Type: FORGET},
	}
}
