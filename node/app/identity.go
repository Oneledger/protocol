/*
	Copyright 2017-2018 OneLedger

	Identities for any of the chains
*/
package app

type IdentityType int

// TODO: should be dynamic?

const (
	ONELEDGER IdentityType = iota
	BITCOIN   IdentityType = iota
	ETHEREUM  IdentityType = iota
)

type Identity struct {
	ttype IdentityType
}

type IdentityOneLedger struct {
	base Identity

	name string
}

type IdentityBitcoin struct {
	base Identity

	name string
}

type IdentityEthereum struct {
	base Identity

	name string
}
