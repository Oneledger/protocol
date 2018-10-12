/*
	Copyright 2017 - 2018 OneLedger
*/
package data

import "github.com/Oneledger/protocol/node/serial"

type ChainType int

// TODO: These should be in a database
const (
	UNKNOWN ChainType = iota
	ONELEDGER
	BITCOIN
	ETHEREUM
)

type Chain struct {
	ChainType   ChainType
	Description string
	Features    []string
}

type ChainNode struct {
	ChainType ChainType
	Location  string

	// TODO: Causing cycle...
	//Owner     id.Identity
}

func init() {
	serial.Register(ChainType(0))
	serial.Register(Chain{})
	serial.Register(ChainNode{})
}
