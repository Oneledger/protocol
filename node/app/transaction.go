/*
	Copyright 2017-2018 OneLedger

	An incoming transaction
*/
package app

var ChainId string

type Message []byte // Contents of a transaction

// ENUM for type
type TransactionType byte

const (
	SEND_TRANSACTION TransactionType = iota
	SWAP_TRANSACTION
	VERIFY_PREPARE
	VERIFY_COMMIT
)

// Polymorphism and Serializable
type Transaction interface {
	Validate() Error
	ProcessTransaction() Error
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

func init() {
	ChainId = "OneLedger-Root"
}

func (transaction *SendTransaction) Validate() Error {
	return SUCCESS
}

func (transaction *SendTransaction) ProcessTransaction() Error {
	return SUCCESS
}

func (transaction *SwapTransaction) Validate() Error {
	return SUCCESS
}

func (transaction *SwapTransaction) ProcessTransaction() Error {
	return SUCCESS
}
