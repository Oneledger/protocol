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
	"github.com/tendermint/tendermint/libs/common"
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
	INVALID         Type = iota
	REGISTER             // Register a new identity with the chain
	SEND                 // Do a normal send transaction on local chain
	PAYMENT              // Do a payment transaction on local chain
	EXTERNAL_SEND        // Do send on external chain
	SWAP                 // Start a swap between chains
	SMART_CONTRACT       // Install and Execute smart contracts
	APPLY_VALIDATOR      // Apply a dynamic validator
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
	GetSigners() []id.PublicKey
	GetOwner() id.AccountKey
	TransactionTags() Tags
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

func (b Base) GetOwner() id.AccountKey {
	return b.Owner
}

func (b Base) GetSigners() []id.PublicKey {
	return b.Signers
}

func ValidateSignature(transaction SignedTransaction) bool {
	log.Debug("Signature validation", "transaction", transaction)

	signers := transaction.Transaction.GetSigners()

	if signers == nil {
		log.Warn("Signature validation (no signers)", "transaction", transaction)
		log.Dump("Signed Transaction is", transaction)
		return false
	}

	if transaction.Signatures == nil {
		log.Warn("Signature validation (no signatures)", "transaction", transaction)
		log.Dump("Signed Transaction is", transaction)
		return false
	}

	if len(signers) == 0 {
		log.Warn("Signature validation (no signers)", "transaction", transaction)
		log.Dump("Signed Transaction is", transaction)
		return false
	}

	if len(signers) != len(transaction.Signatures) {
		log.Warn("Signature validation (wrong number of signatures)", "transaction", transaction)
		log.Dump("Signed Transaction is", transaction)
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
			log.Dump("Signed Transaction is", transaction)
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

func Parse(message Message) (SignedTransaction, status.Code) {
	var tx SignedTransaction

	transaction, transactionErr := serial.Deserialize(message, tx, serial.CLIENT)

	if transactionErr == nil {
		return transaction.(SignedTransaction), status.SUCCESS
	}

	log.Error("Could not deserialize a transaction", "error", transactionErr)

	return SignedTransaction{}, status.PARSE_ERROR
}

func (t Type) String() string {
	switch t {
	case REGISTER:
		return "REGISTER"
	case SEND:
		return "SEND"
	case PAYMENT:
		return "PAYMENT"
	case EXTERNAL_SEND:
		return "EXTERNAL_SEND"
	case SWAP:
		return "SWAP"
	case SMART_CONTRACT:
		return "SMART_CONTRACT"
	case APPLY_VALIDATOR:
		return "APPLY_VALIDATOR"
	default:
		return "INVALID"
	}
}

type Tags common.KVPairs

func (b Base) TransactionTags() Tags {

	//Add transaction type as a tag
	tagType := b.Type.String()
	tag1 := common.KVPair{
		Key:   []byte("tx.type"),
		Value: []byte(tagType),
	}

	//Add owner as a tag
	tagOwner := b.Owner.String()
	tag2 := common.KVPair{
		Key:   []byte("tx.owner"),
		Value: []byte(tagOwner),
	}
	return Tags{tag1, tag2}
}
