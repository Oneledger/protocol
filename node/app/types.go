/*
	Copyright 2017-2018 OneLedger

	Declare all of the types used by the Application
*/
package app

import (
	crypto "github.com/tendermint/go-crypto"
	"github.com/tendermint/go-wire/data"
)

// TODO: Should we be using type aliases here?
type Message []byte // Contents of a transaction

type DatabaseKey []byte // Database key

type Address = data.Bytes
type Signature = crypto.Signature
type PublicKey = crypto.PubKey
type PrivateKey = crypto.PrivKey

type Coin struct {
	Currency string `json:"denom"`
	Amount   int64  `json:"amount"`
}

type Coins []Coin

type SendInput struct {
	Address   Address   `json:"address"`   // Hash of the PubKey
	Coins     Coins     `json:"coins"`     //
	Sequence  int       `json:"sequence"`  // Must be 1 greater than the last committed TxInput
	Signature Signature `json:"signature"` // Depends on the PubKey type and the whole Tx
	PubKey    PublicKey `json:"pub_key"`   // Is present iff Sequence == 0
}

type SendOutput struct {
	Address Address `json:"address"` // Hash of the PubKey
	Coins   Coins   `json:"coins"`   //
}

type SendTransaction struct {
	Type    TransactionType `json:"type"`
	Gas     Coin            `json:"gas"`
	Fee     Coin            `json:"fee"`
	Inputs  []SendInput     `json:"inputs"`
	Outputs []SendOutput    `json:"outputs"`
}

type FullSendTransaction struct {
	ChainId     string           `json:"chain_id"`
	Signers     []PublicKey      `json:"signers"`
	Transaction *SendTransaction `json:"transaction"`
}
