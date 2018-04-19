/*
	Copyright 2017-2018 OneLedger
*/
package app

import (
	wire "github.com/tendermint/go-wire"
)

type Error uint32 // Matches Tendermint status

const (
	SUCCESS         Error = 0
	PARSE_ERROR     Error = 101
	NOT_IMPLEMENTED Error = 201
)

func UnpackMessage(message Message) (TransactionType, Message) {
	value, size, err := wire.GetVarint(message)
	if err != nil {
		Log.Debug("Wire returned an error", "err", err)
		panic("Wire Error")
	}
	if size != 2 {
		Log.Debug("Wire returned a bad size", "size", size)
		panic("Sizing Error")
	}
	return TransactionType(value), message[1:]

}

// Parse a message into the appropriate transaction
func Parse(message Message) (Transaction, Error) {
	Log.Debug("Parsing a Transaction")

	command, body := UnpackMessage(message)

	switch command {

	case SEND_TRANSACTION:
		transaction := ParseSend(body)

		return transaction, SUCCESS

	case SWAP_TRANSACTION:
		transaction := ParseSwap(body)

		return transaction, SUCCESS

	case VERIFY_PREPARE:
		Log.Error("Have Prepare, not implemented yet")

		return nil, NOT_IMPLEMENTED

	case VERIFY_COMMIT:
		Log.Error("Have Commit, not implemented yet")

		return nil, NOT_IMPLEMENTED

	default:
		Log.Error("Unknown type", "command", command)
	}

	return nil, PARSE_ERROR
}

func ParseSend(message Message) *SendTransactionBase {
	Log.Debug("Have a Send")

	return &SendTransactionBase{ttype: SEND_TRANSACTION}
}

func ParseSwap(message Message) *SwapTransactionBase {
	Log.Debug("Have a Swap")

	return &SwapTransactionBase{ttype: SWAP_TRANSACTION}
}
