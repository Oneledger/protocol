/*
	Copyright 2017-2018 OneLedger

	An incoming transaction
*/
package app

type TransactionType byte

const (
	SEND_TRANSACTION TransactionType = iota
	SWAP_TRANSACTION TransactionType = iota
	VERIFY_PREPARE   TransactionType = iota
	VERIFY_COMMIT    TransactionType = iota
)

type Transaction interface {
}

// Base Data for each type
type TransactionBase struct {
	ttype TransactionType
}

// Synchronize a swap between two users
type SwapTransactionBase struct {
	// TODO: Fix this to embed it properly.
	//TransactionBase
	ttype TransactionType

	party1       string // TODO: put in addresses here
	party2       string
	exchangeRate int
	amount       int
	fee          int
}

// Synchronize a swap between two users
type SendTransactionBase struct {
	// TODO: Fix this to embed it properly.
	//TransactionBase
	ttype TransactionType

	from         string // TODO: put in addresses here
	to           string
	exchangeRate int
	amount       int
	fee          int
}

// TODO: roughed in...
type CoinTransactionBase struct {
	TransactionBase

	inputs  []string
	outputs []string
	// What else is standard?
}
