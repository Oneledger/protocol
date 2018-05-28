/*
	Copyright 2017-2018 OneLedger

	Common Transaction utilities
*/
package main

import (
	"github.com/Oneledger/protocol/node/action"
	"github.com/Oneledger/protocol/node/id"
	wire "github.com/tendermint/go-wire"
)

func GetPublicKey() id.PublicKey {
	// TODO: Really not sure about this.
	return id.PublicKey{}
}

// GetSigners will return the public keys of the signers
func GetSigners() []id.PublicKey {
	return nil
}

// SignTransaction with the local keys
func SignTransaction(transaction action.Transaction) action.Transaction {
	return transaction
}

// Pack a request into a transferable format (wire)
func PackRequest(request action.Transaction) []byte {
	packet := wire.BinaryBytes(request)
	return packet
}
