/*
	Copyright 2017-2018 OneLedger

	Declare basic types used by the Application

	If a type requires functions or a few types are intertwinded, then should be in their own file.
*/
package app

import (
	crypto "github.com/tendermint/go-crypto"
	"github.com/tendermint/go-wire/data"
)

// TODO: Should we be using type aliases here?
type Message []byte     // Contents of a transaction
type DatabaseKey []byte // Database key

// Hide some of the underlying types.
type Address = data.Bytes
type Signature = crypto.Signature
type PublicKey = crypto.PubKey
type PrivateKey = crypto.PrivKey

// Coin is the basic amount, specified in integers, at the smallest increment (i.e. a satoshi, not a bitcoin)
type Coin struct {
	Currency string `json:"denom"`
	Amount   int64  `json:"amount"`
}

type Coins []Coin

// A rate of exchange, agreeed upon between two parties
type ExchangeRate struct {
	//Rate float64 // TODO: should this actually be a rational pair? wire is really unhappy about floats...
	Numerator   int64 `json:"numerator"`
	Denominator int64 `json:"denominator"`
}

// inputs into a transaction
type SendInput struct {
	Address   Address   `json:"address"`   // Hash of the PubKey
	Coins     Coins     `json:"coins"`     //
	Sequence  int       `json:"sequence"`  // Must be 1 greater than the last committed Input (reply protection?)
	Signature Signature `json:"signature"` // Depends on the PubKey type and the whole Tx (?)
	PubKey    PublicKey `json:"pub_key"`   // Is present iff Sequence == 0 (?)
}

// outputs for a transaction
type SendOutput struct {
	Address Address `json:"address"` // Hash of the PubKey
	Coins   Coins   `json:"coins"`   //
}
