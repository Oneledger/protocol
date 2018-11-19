/*
	Copyright 2017-2018 OneLedger

	An incoming transaction, send, swap, ready, verification, etc.
*/
package action

import (
	"github.com/Oneledger/protocol/node/id"
	"github.com/Oneledger/protocol/node/log"
	"github.com/Oneledger/protocol/node/serial"
	"github.com/Oneledger/protocol/node/status"
)

type Message = []byte // Contents of a transaction
// ENUM for type
type Type int
type Role int

func init() {
	serial.Register(Type(0))
	serial.Register(Role(0))
	serial.Register(Message(""))
}

const (
	INVALID       Type = iota
	REGISTER           // Register a new identity with the chain
	SEND               // Do a normal send transaction on local chain
	EXTERNAL_SEND      // Do send on external chain
	EXTERNAL_LOCK      // Lock some data on external chain
	SWAP               // Start a swap between chains
	VERIFY             // Verify if a transaction finished
	PUBLISH            // Exchange data on a chain
	READ               // Read a specific transaction on a chain
	PREPARE            // Do everything, except commit
	COMMIT             // Commit to doing the work
	FORGET             // Rollback and forget that this happened
)

const (
	ALL         Role = iota
	INITIATOR        // Register a new identity with the chain
	PARTICIPANT      // Do a normal send transaction on local chain
	NONE
)

type PublicKey = id.PublicKey

// Polymorphism and Serializable
type Transaction interface {
	TransactionType() Type
	Validate() status.Code
	ProcessCheck(interface{}) status.Code
	ShouldProcess(interface{}) bool
	ProcessDeliver(interface{}) status.Code
	Resolve(interface{}) Commands
}

type TransactionSignature struct {
	Signature []byte
}

type SignedTransaction struct {
	Transaction
	Signatures []TransactionSignature
}

// Base Data for each type
type Base struct {
	Type    Type   `json:"type"`
	ChainId string `json:"chain_id"`

	Owner  id.AccountKey `json:"owner"`
	Target id.AccountKey `json:"target"`

	Signers []PublicKey `json:"signers"`

	Sequence int64 `json:"sequence"`
	Delay    int64 `json:"delay"` // Pause the transaction in the mempool
}

func ValidateSignature(transaction SignedTransaction) bool {
	log.Debug("Signature validation", "transaction", transaction)
	var signers []id.PublicKey

	// TODO need to simplify it
	switch v := transaction.Transaction.(type) {
	case *Swap: signers = v.Base.Signers
	case *Send: signers = v.Base.Signers
	case *Register: signers = v.Base.Signers
	default: log.Warn("Signature validation (unknown transaction type)", "transaction", transaction)
	}

	if signers == nil {
		log.Warn("Signature validation (no signers)", "transaction", transaction)
		return false
	}

	if transaction.Signatures == nil {
		log.Warn("Signature validation (no signatures)", "transaction", transaction)
		return false
	}

	if len(signers) == 0 {
		log.Warn("Signature validation (no signers)", "transaction", transaction)
		return false
	}

	if len(signers) != len(transaction.Signatures) {
		log.Warn("Signature validation (wrong number of signatures)", "transaction", transaction)
		return false
	}

	message, err := serial.Serialize(transaction.Transaction, serial.CLIENT)

	if err != nil {
		log.Error("Signature validation (failed to serialize)", "error", err, "transaction", transaction)
		return false
	}

	for i := 0; i < len(signers); i++ {
		if signers[i].VerifyBytes(message, transaction.Signatures[i].Signature) == false {
			log.Warn("Signature validation (invalid signature)", "index", i, "transaction", transaction)
			return false
		}
	}

	log.Debug("Signature validation", "success", true)

	return true
}

func init() {
	serial.Register(Base{})
	serial.Register(TransactionSignature{})
	serial.Register(SignedTransaction{})
}
