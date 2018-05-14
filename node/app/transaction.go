/*
	Copyright 2017-2018 OneLedger

	An incoming transaction, send, swap, ready, verification, etc.
*/
package app

import "github.com/Oneledger/protocol/node/log"

var ChainId string

type Message = []byte // Contents of a transaction

// ENUM for type
type TransactionType byte

const (
	SEND_TRANSACTION TransactionType = iota
	SWAP_TRANSACTION
	READY_TRANSACTION
	VERIFY_TRANSACTION
)

func init() {
	ChainId = "OneLedger-Root"
}

// Polymorphism and Serializable
type Transaction interface {
	Validate() Error
	ProcessCheck(*Application) Error
	ProcessDeliver(*Application) Error
}

// Base Data for each type
type TransactionBase struct {
	Type    TransactionType `json:"type"`
	ChainId string          `json:"chain_id"`
	Signers []PublicKey     `json:"signers"`

	// TODO: Should these be for all transactions or just driving ones?
	Sequence int     `json:"sequence"`
	Owner    Address `json:"owner"`
}

// Synchronize a swap between two users
type SendTransaction struct {
	TransactionBase

	Gas     Coin         `json:"gas"`
	Fee     Coin         `json:"fee"`
	Inputs  []SendInput  `json:"inputs"`
	Outputs []SendOutput `json:"outputs"`
}

func (transaction *SendTransaction) Validate() Error {
	log.Debug("Validating Send Transaction")

	// TODO: Make sure all of the parameters are there
	// TODO: Check all signatures and keys
	// TODO: Vet that the sender has the values
	return SUCCESS
}

func (transaction *SendTransaction) ProcessCheck(app *Application) Error {
	log.Debug("Processing Send Transaction for CheckTx")

	// TODO: // Update in memory copy of Merkle Tree
	return SUCCESS
}

func (transaction *SendTransaction) ProcessDeliver(app *Application) Error {
	log.Debug("Processing Send Transaction for DeliverTx")

	// TODO: // Update in final copy of Merkle Tree
	return SUCCESS
}

// Synchronize a swap between two users
type SwapTransaction struct {
	TransactionBase

	Party1   Address `json:"party1"`
	Party2   Address `json:"party2"`
	Fee      Coin    `json:"fee"`
	Gas      Coin    `json:"fee"`
	Amount   Coin    `json:"amount"`
	Exchange Coin    `json:"exchange"`
}

// Issue swaps across other chains, make sure fees are collected
func (transaction *SwapTransaction) Validate() Error {
	log.Debug("Validating Swap Transaction")
	return SUCCESS
}

func (transaction *SwapTransaction) ProcessCheck(app *Application) Error {
	log.Debug("Processing Swap Transaction for CheckTx")
	return SUCCESS
}

func (transaction *SwapTransaction) ProcessDeliver(app *Application) Error {
	log.Debug("Processing Swap Transaction for DeliverTx")
	return SUCCESS
}

// TODO: This needs to be filled out properly, need a model for other-chain actions...
type ReadyTransaction struct {
	TransactionBase

	Target string `json:"target"`
}

func (transaction *ReadyTransaction) Validate() Error {
	log.Debug("Validating Ready Transaction")
	return SUCCESS
}

func (transaction *ReadyTransaction) ProcessCheck(app *Application) Error {
	log.Debug("Processing Ready Transaction for CheckTx")
	return SUCCESS
}

func (transaction *ReadyTransaction) ProcessDeliver(app *Application) Error {
	log.Debug("Processing Ready Transaction for DeliverTx")
	return SUCCESS
}

// TODO: This needs to be filled out properly, need a model for other-chain actions...
type VerifyTransaction struct {
	TransactionBase

	Target string `json:"target"`
}

func (transaction *VerifyTransaction) Validate() Error {
	log.Debug("Validating Verify Transaction")
	return SUCCESS
}

func (transaction *VerifyTransaction) ProcessCheck(app *Application) Error {
	log.Debug("Processing Verify Transaction for CheckTx")
	return SUCCESS
}

func (transaction *VerifyTransaction) ProcessDeliver(app *Application) Error {
	log.Debug("Processing Verify Transaction for DeliverTx")
	return SUCCESS
}
