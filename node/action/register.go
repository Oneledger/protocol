/*
	Copyright 2017 - 2018 OneLedger

	Register this identity with the other nodes. As an externl identity
*/
package action

import (
	"github.com/Oneledger/protocol/node/err"
	"github.com/Oneledger/protocol/node/global"
	"github.com/Oneledger/protocol/node/id"
	"github.com/Oneledger/protocol/node/log"
)

// Register an identity with the chain
type Register struct {
	Base

	Identity   string
	NodeName   string
	AccountKey id.AccountKey
}

// Check the fields to make sure they have valid values.
func (transaction Register) Validate() err.Code {
	log.Debug("Validating Register Transaction")

	if transaction.Identity == "" {
		return err.MISSING_DATA
	}

	if transaction.NodeName == "" {
		return err.MISSING_DATA
	}

	// TODO: Make sure all of the parameters are there
	// TODO: Check all signatures and keys
	// TODO: Vet that the sender has the values
	return err.SUCCESS
}

// Test to see if the identity already exists
func (transaction Register) ProcessCheck(app interface{}) err.Code {
	log.Debug("Processing Register Transaction for CheckTx")

	identities := GetIdentities(app)
	id, errs := identities.FindName(transaction.Identity)

	if errs != err.SUCCESS {
		return errs
	}

	if id == nil {
		log.Debug("Success, it is a new Identity", "id", transaction.Identity)
		return err.SUCCESS
	}

	log.Debug("Identity already exists", "id", id)

	// TODO: Not necessarily a failure, since this identity might be local
	return err.SUCCESS
}

func (transaction Register) ShouldProcess(app interface{}) bool {
	return true
}

// Add the identity into the database as external, don't overwrite a local identity
func (transaction Register) ProcessDeliver(app interface{}) err.Code {
	log.Debug("Processing Register Transaction for DeliverTx")

	identities := GetIdentities(app)
	entry, status := identities.FindName(transaction.Identity)

	if status != err.SUCCESS {
		return status
	}

	if entry != nil {
		log.Debug("Ignoring Existin Identity")
	} else {
		identities.Add(id.NewIdentity(transaction.Identity, "Contact Information",
			true, global.Current.NodeName, transaction.AccountKey))
		log.Info("Updated External Identity", "id", transaction.Identity, "key", transaction.AccountKey)
	}

	return err.SUCCESS
}

// Given a transaction, expand it into a list of Commands to execute against various chains.
func (transaction Register) Expand(app interface{}) Commands {
	// TODO: Table-driven mechanics, probably elsewhere
	return []Command{}
}
