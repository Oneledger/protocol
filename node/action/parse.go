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
func UnpackMessage(message Message) (TransactionType, Message) {
	value, size, err := wire.GetVarint(message)
	if err != nil {
		log.Debug("Wire returned an error", "err", err)
		panic("Wire Error")
	}
	if size != 2 {
		log.Debug("Wire returned a bad size", "size", size)
		panic("Sizing Error")
	}
	return TransactionType(value), message[1:]

}

// Parse a message into the appropriate transaction
func Parse(message Message) (Transaction, err.Code) {
	log.Debug("Parsing a Transaction")

	command, body := UnpackMessage(message)

	switch command {

	case SEND:
		transaction := ParseSend(body)

		return transaction, err.SUCCESS

	case SWAP:
		transaction := ParseSwap(body)

		return transaction, err.SUCCESS

	case EXTERNAL_SEND:
		transaction := ParseExternalSend(body)

		return transaction, err.SUCCESS

	case EXTERNAL_LOCK:
		transaction := ParseExternalLock(body)

		return transaction, err.SUCCESS

	case PREPARE:
		transaction := ParsePrepare(body)

		return transaction, err.SUCCESS

	case COMMIT:
		transaction := ParseCommit(body)

		return transaction, err.SUCCESS

	case FORGET:
		transaction := ParseForget(body)

		return transaction, err.SUCCESS

	default:
		log.Error("Unknown type", "command", command)
	}

	return nil, err.PARSE_ERROR
}

// Parse a send request
func ParseSend(message Message) *Send {
	log.Debug("Have a Send")

	return &Send{
		TransactionBase: TransactionBase{Type: SEND},
	}
}

// Parse a swap request
func ParseSwap(message Message) *Swap {
	log.Debug("Have a Swap")

	//return &SwapTransaction{Type: SWAP_TRANSACTION}
	return &Swap{
		TransactionBase: TransactionBase{Type: SWAP},
	}
}

// Parse a send request
func ParseExternalSend(message Message) *ExternalSend {
	log.Debug("Have an ExternalSend")

	return &ExternalSend{
		TransactionBase: TransactionBase{Type: EXTERNAL_SEND},
	}
}

// Parse a send request
func ParseExternalLock(message Message) *ExternalLock {
	log.Debug("Have an ExternalLock")

	return &ExternalLock{
		TransactionBase: TransactionBase{Type: EXTERNAL_LOCK},
	}
}

// Parse a ready request
func ParsePrepare(message Message) *Prepare {
	log.Debug("Have a Ready")

	//return &SwapTransaction{Type: SWAP_TRANSACTION}
	return &Prepare{
		TransactionBase: TransactionBase{Type: PREPARE},
	}
}

// Parse a ready request
func ParseCommit(message Message) *Commit {
	log.Debug("Have a Ready")

	//return &SwapTransaction{Type: SWAP_TRANSACTION}
	return &Commit{
		TransactionBase: TransactionBase{Type: COMMIT},
	}
}

// Forget the transaction
func ParseForget(message Message) *Forget {
	log.Debug("Have a Ready")

	//return &SwapTransaction{Type: SWAP_TRANSACTION}
	return &Forget{
		TransactionBase: TransactionBase{Type: FORGET},
	}
}
