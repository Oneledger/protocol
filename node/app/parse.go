/*
	Copyright 2017-2018 OneLedger

	Parse the incoming transactions

	TODO: switch from individual wire calls, to reading/writing directly to structs
*/
package app

import (
	"github.com/Oneledger/protocol/node/log"
	wire "github.com/tendermint/go-wire"
)

type Error = uint32 // Matches Tendermint status

const (
	SUCCESS         Error = 0
	PARSE_ERROR     Error = 101
	NOT_IMPLEMENTED Error = 201
	MISSING_VALUE   Error = 301
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
func Parse(message Message) (Transaction, Error) {
	log.Debug("Parsing a Transaction")

	command, body := UnpackMessage(message)

	switch command {

	case SEND_TRANSACTION:
		transaction := ParseSend(body)

		return transaction, SUCCESS

	case SWAP_TRANSACTION:
		transaction := ParseSwap(body)

		return transaction, SUCCESS

	case READY_TRANSACTION:
		transaction := ParseReady(body)

		return transaction, SUCCESS

	case VERIFY_TRANSACTION:
		transaction := ParseVerify(body)

		return transaction, SUCCESS

	default:
		log.Error("Unknown type", "command", command)
	}

	return nil, PARSE_ERROR
}

// Parse a send request
func ParseSend(message Message) *SendTransaction {
	log.Debug("Have a Send")

	return &SendTransaction{
		TransactionBase: TransactionBase{Type: SEND_TRANSACTION},
	}
}

// Parse a swap request
func ParseSwap(message Message) *SwapTransaction {
	log.Debug("Have a Swap")

	//return &SwapTransaction{Type: SWAP_TRANSACTION}
	return &SwapTransaction{
		TransactionBase: TransactionBase{Type: SWAP_TRANSACTION},
	}
}

// Parse a ready request
func ParseReady(message Message) *ReadyTransaction {
	log.Debug("Have a Ready")

	//return &SwapTransaction{Type: SWAP_TRANSACTION}
	return &ReadyTransaction{
		TransactionBase: TransactionBase{Type: READY_TRANSACTION},
	}
}

// Parse a ready request
func ParseVerify(message Message) *VerifyTransaction {
	log.Debug("Have a Ready")

	//return &SwapTransaction{Type: SWAP_TRANSACTION}
	return &VerifyTransaction{
		TransactionBase: TransactionBase{Type: VERIFY_TRANSACTION},
	}
}
