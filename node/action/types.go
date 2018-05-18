/*
	Copyright 2017-2018 OneLedger

	Declare basic types used by the Application

	If a type requires functions or a few types are intertwinded, then should be in their own file.
*/
package action

import (
	"github.com/Oneledger/protocol/node/data"
	"github.com/Oneledger/protocol/node/id"
)

// inputs into a send transaction (similar to Bitcoin)
type SendInput struct {
	Address   id.Address   `json:"address"`
	Coins     data.Coins   `json:"coins"`
	Sequence  int          `json:"sequence"`
	Signature id.Signature `json:"signature"`
	PubKey    PublicKey    `json:"pub_key"`
}

// outputs for a send transaction (similar to Bitcoin)
type SendOutput struct {
	Address id.Address `json:"address"`
	Coins   data.Coins `json:"coins"`
}
