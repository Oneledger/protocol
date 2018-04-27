/*
	Copyright 2017-2018 OneLedger

	Common Transaction utilities
*/
package main

import (
	"github.com/Oneledger/prototype/node/app"
	wire "github.com/tendermint/go-wire"
)

func GetPublicKey() app.PublicKey {
	// TODO: Really not sure about this.
	return app.PublicKey{}
}

// GetSigners will return the public keys of the signers
func GetSigners() []app.PublicKey {
	return nil
}

// SignTransaction with the local keys
func SignTransaction(transaction app.Transaction) app.Transaction {
	return transaction
}

// Pack a request into a transferable format (wire)
func PackRequest(request app.Transaction) []byte {
	packet := wire.BinaryBytes(request)
	return packet
}
