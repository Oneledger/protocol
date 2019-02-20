/*
	Copyright 2017-2018 OneLedger

	Common Transaction utilities, helps to create them consistently
*/
package shared

import (
	"github.com/Oneledger/protocol/node/id"
)

// Given an Identity or Account, get the correct associated public key
func GetPublicKey() id.PublicKey {
	// TODO: Really not sure about this.
	return id.NilPublicKey()
}
