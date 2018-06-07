/*
	Copyright 2017-2018 OneLedger

	An incoming transaction, send, swap, ready, verification, etc.
*/
package action

import (
	"github.com/Oneledger/protocol/node/data"
	"github.com/Oneledger/protocol/node/err"
	"github.com/Oneledger/protocol/node/id"
	crypto "github.com/tendermint/go-crypto"
)

type Message = []byte // Contents of a transaction
type PublicKey = crypto.PubKey

// ENUM for type
type Type byte

const (
	INVALID       Type = iota
	REGISTER           // Register a new identity with the chain
	SEND               // Do a normal send transaction on local chain
	EXTERNAL_SEND      // Do send on external chain
	EXTERNAL_LOCK      // Lock some data on external chain
	SWAP               // Start a swap between chains
	VERIFY             // Verify that a lockbox is correct
	PUBLISH            // Publish data (preimage) on a chain
	READ               // Read a specific transaction on a chain
	PREPARE            // Do everything, except commit
	COMMIT             // Commit to doing the work
	FORGET             // Rollback and forget that this happened
)

// Polymorphism and Serializable
type Transaction interface {
	Validate() err.Code
	ProcessCheck(interface{}) err.Code
	ProcessDeliver(interface{}) err.Code
	Expand(interface{}) Commands
}

// Base Data for each type
type Base struct {
	Type    Type        `json:"type"`
	ChainId string      `json:"chain_id"` // TODO: Not necessary?
	Signers []PublicKey `json:"signers"`

	// TODO: Should these be for all transactions or just driving ones?
	Sequence int        `json:"sequence"`
	Owner    id.Address `json:"owner"`
}

// Get the correct chain for an action
func GetChain(transaction interface{}) data.ChainType {
	// TODO: Need to fix this
	return data.BITCOIN
}
