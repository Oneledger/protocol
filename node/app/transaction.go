/*
	Copyright 2017-2018 OneLedger

	An incoming transaction
*/
package app

var ChainId string

type TransactionType byte

const (
	SEND_TRANSACTION TransactionType = iota
	SWAP_TRANSACTION TransactionType = iota
	VERIFY_PREPARE   TransactionType = iota
	VERIFY_COMMIT    TransactionType = iota
)

func init() {
	ChainId = "OneLedger-Root"
}

type Transaction interface {
}

// Base Data for each type
type TransactionBase struct {
	Type    TransactionType `json:"type"`
	ChainId string          `json:"chain_id"`
	Signers []PublicKey     `json:"signers"`
}

// Synchronize a swap between two users
type SwapTransaction struct {
	TransactionBase

	Party1       Address      `json:"party1"`
	Party2       Address      `json:"party2"`
	ExchangeRate ExchangeRate `json:"exchangeRate"`
	Amount       Coin         `json:"amount"`
	Fee          Coin         `json:"fee"`
}

// Synchronize a swap between two users
type SendTransaction struct {
	TransactionBase

	Gas     Coin         `json:"gas"`
	Fee     Coin         `json:"fee"`
	Inputs  []SendInput  `json:"inputs"`
	Outputs []SendOutput `json:"outputs"`
}

func Validate(transaction Transaction) Error {
	switch transaction.(type) {

	case SwapTransaction:
		return ValidateSwap(transaction)

	case SendTransaction:
		return ValidateSend(transaction)

	}
	return MISSING_VALUE
}

func ValidateSwap(transaction Transaction) Error {
	return SUCCESS
}

func ValidateSend(transaction Transaction) Error {
	return SUCCESS
}
