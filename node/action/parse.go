/*
	Copyright 2017-2018 OneLedger

	Parse the incoming transactions

	TODO: switch from individual wire calls, to reading/writing directly to structs
*/
package action

import (
	"github.com/Oneledger/protocol/node/err"
	"github.com/Oneledger/protocol/node/log"
	"github.com/Oneledger/protocol/node/serial"
)

// Parse a message into the appropriate transaction
func Parse(message Message) (Transaction, err.Code) {
	var tx Transaction

	transaction, transactionErr := serial.Deserialize(message, tx, serial.CLIENT)

	if transactionErr == nil {
		return transaction.(Transaction), err.SUCCESS
	}

	log.Error("Could not deserialize a transaction", "error",  transactionErr)

	return nil, err.PARSE_ERROR
}
