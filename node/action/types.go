/*
	Copyright 2017-2018 OneLedger

	Declare basic types used by the Application

	If a type requires functions or a few types are intertwinded, then should be in their own file.
*/
package action

import "github.com/Oneledger/protocol/node/id"

// Coin is the basic amount, specified in integers, at the smallest increment (i.e. a satoshi, not a bitcoin)
type Coin struct {
	Currency string `json:"currency"`
	Amount   int64  `json:"amount"`
}

type Coins []Coin

// inputs into a send transaction (similar to Bitcoin)
type SendInput struct {
	Address   id.Address   `json:"address"`
	Coins     Coins        `json:"coins"`
	Sequence  int          `json:"sequence"`
	Signature id.Signature `json:"signature"`
	PubKey    PublicKey    `json:"pub_key"`
}

// outputs for a send transaction (similar to Bitcoin)
type SendOutput struct {
	Address id.Address `json:"address"`
	Coins   Coins      `json:"coins"`
}
