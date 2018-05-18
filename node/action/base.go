/*
	Copyright 2017-2018 OneLedger

	An incoming transaction, send, swap, ready, verification, etc.
*/
package action

import (
	"github.com/Oneledger/protocol/node/err"
	"github.com/Oneledger/protocol/node/id"
	crypto "github.com/tendermint/go-crypto"
)

type Message = []byte // Contents of a transaction
type PublicKey = crypto.PubKey

// ENUM for type
type TransactionType byte

const (
	SEND TransactionType = iota
	EXTERNAL_SEND
	EXTERNAL_LOCK
	SWAP
	PREPARE
	COMMIT
	FORGET
)

// Polymorphism and Serializable
type Transaction interface {
	Validate() err.Code
	ProcessCheck(interface{}) err.Code
	ProcessDeliver(interface{}) err.Code
}

// Base Data for each type
type TransactionBase struct {
	Type    TransactionType `json:"type"`
	ChainId string          `json:"chain_id"` // TODO: Not necessary?
	Signers []PublicKey     `json:"signers"`

	// TODO: Should these be for all transactions or just driving ones?
	Sequence int        `json:"sequence"`
	Owner    id.Address `json:"owner"`
}
