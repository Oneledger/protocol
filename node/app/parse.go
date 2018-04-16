/*
	Copyright 2017-2018 OneLedger
*/
package app

type Error uint32 // Matches Tendermint

const (
	PARSE_ERROR     Error = 101
	NOT_IMPLEMENTED Error = 201
)

// Parse a message into the appropriate transaction
func Parse(transaction Message) (Transaction, Error) {
	Log.Debug("Parsing a Transaction")

	command := TransactionType(transaction[0])

	switch command {

	case SWAP_TRANSACTION:
		Log.Debug("Have a Swap")

		transaction := &SwapTransaction{ttype: command}
		return transaction, 0

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
