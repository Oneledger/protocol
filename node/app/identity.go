/*
	Copyright 2017-2018 OneLedger

	Identities management for any of the associated chains
*/
package app

import (
	crypto "github.com/tendermint/go-crypto"
	"github.com/tendermint/go-wire/data"
)

// Aliases to hide some of the basic underlying types.

type Address = data.Bytes // OneLedger address, like Tendermint the hash of the associated PubKey

type PublicKey = crypto.PubKey
type PrivateKey = crypto.PrivKey

type Signature = crypto.Signature

// ENUM for type
type IdentityType int

const (
	ONELEDGER IdentityType = iota
	BITCOIN
	ETHEREUM
)

// Polymorphism
type Identity interface {
	AddPublicKey(PublicKey)
	AddPrivateKey(PrivateKey)
}

type IdentityBase struct {
	Type IdentityType
	Name string
}

// Information we need about our own fullnode identities
type IdentityOneLedger struct {
	IdentityBase
	Address Address

	PublicKey PublicKey
	PrivteKey PrivateKey

	NodeId      string
	ExternalIds []string
}

// Information we need for a Bitcoin account
type IdentityBitcoin struct {
	IdentityBase
	Address Address

	PublicKey PublicKey
	PrivteKey PrivateKey
}

// Information we need for an Ethereum account
type IdentityEthereum struct {
	IdentityBase
	Address Address

	PublicKey PublicKey
	PrivteKey PrivateKey
}

func NewIdentity(name string, newType IdentityType) Identity {
	switch newType {

	case ONELEDGER:
		return &IdentityOneLedger{}

	case BITCOIN:
		return &IdentityBitcoin{}

	case ETHEREUM:
		return &IdentityEthereum{}

	default:
		panic("Unknown Type")
	}
}

func FindIdentityType(typeName string) (IdentityType, Error) {
	switch typeName {
	case "OneLedger":
		return ONELEDGER, 0

	case "Ethereum":
		return ETHEREUM, 0

	case "Bitcoin":
		return BITCOIN, 0
	}
	return 0, 42
}

func FindIdentity(name string) (Identity, Error) {
	// TODO: Lookup the identity in the node's database
	return &IdentityOneLedger{}, 42
}

func (identity *IdentityOneLedger) AddPublicKey(key PublicKey) {
}

func (identity *IdentityBitcoin) AddPublicKey(key PublicKey) {
}

func (identity *IdentityEthereum) AddPublicKey(key PublicKey) {
}

func (identity *IdentityOneLedger) AddPrivateKey(key PrivateKey) {
}

func (identity *IdentityBitcoin) AddPrivateKey(key PrivateKey) {
}

func (identity *IdentityEthereum) AddPrivateKey(key PrivateKey) {
}
