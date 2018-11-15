/*
	Copyright 2017-2018 OneLedger

	Parse the incoming transactions

	TODO: switch from individual wire calls, to reading/writing directly to structs
*/
package action

import (
	"github.com/Oneledger/protocol/node/log"
	"github.com/Oneledger/protocol/node/serial"
	"github.com/Oneledger/protocol/node/status"
)

// Parse a message into the appropriate transaction
func Parse(message Message) (SignedTransaction, status.Code) {
	var tx SignedTransaction

	transaction, transactionErr := serial.Deserialize(message, tx, serial.CLIENT)

	if transactionErr == nil {
		return transaction.(SignedTransaction), status.SUCCESS
	}

	log.Error("Could not deserialize a transaction", "error",  transactionErr)

	return SignedTransaction{}, status.PARSE_ERROR
}
