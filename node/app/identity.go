/*
	Copyright 2017-2018 OneLedger

	Identities for any of the chains
*/
package app

import (
	crypto "github.com/tendermint/go-crypto"
	"github.com/tendermint/go-wire/data"
)

// Aliases to hide some of the basic underlying types.
type Address = data.Bytes
type Signature = crypto.Signature
type PublicKey = crypto.PubKey
type PrivateKey = crypto.PrivKey

// ENUM for type
type IdentityType int

const (
	ONELEDGER IdentityType = iota
	BITCOIN
	ETHEREUM
)

// Polymorphism
type Identity interface {
}

type IdentityBase struct {
	Type IdentityType
	Name string
}

// Information we need about our fullnode identities
type IdentityOneLedger struct {
	IdentityBase
	PublicKey PublicKey
	PrivteKey PrivateKey
}

// Information we need for the installed Bitcoin node
type IdentityBitcoin struct {
	IdentityBase
	Address   Address
	PublicKey PublicKey
	PrivteKey PrivateKey
}

// Information we need for the installed Ethereum node
type IdentityEthereum struct {
	IdentityBase
	Address   Address
	PublicKey PublicKey
	PrivteKey PrivateKey
}
