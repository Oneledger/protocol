/*
	Copyright 2017 - 2018 OneLedger

	Register this identity with the other nodes. As an externl identity
*/
package action

import (
	"github.com/Oneledger/protocol/node/status"
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
func (transaction Register) Validate() status.Code {
	log.Debug("Validating Register Transaction")

	if transaction.Identity == "" {
		return status.MISSING_DATA
	}

	if transaction.NodeName == "" {
		return status.MISSING_DATA
	}

	return status.SUCCESS
}

// Test to see if the identity already exists
func (transaction Register) ProcessCheck(app interface{}) status.Code {
	log.Debug("Processing Register Transaction for CheckTx")

	identities := GetIdentities(app)
	id, status := identities.FindName(transaction.Identity)

	if status != status.SUCCESS {
		return status
	}

	if id == nil {
		log.Debug("Success, it is a new Identity", "id", transaction.Identity)
		return status.SUCCESS
	}

	// Not necessarily a failure, since this identity might be local
	log.Debug("Identity already exists", "id", id)
	return status.SUCCESS
}

func (transaction Register) ShouldProcess(app interface{}) bool {
	return true
}

// Add the identity into the database as external, don't overwrite a local identity
func (transaction Register) ProcessDeliver(app interface{}) status.Code {
	log.Debug("Processing Register Transaction for DeliverTx")

	identities := GetIdentities(app)
	entry, status := identities.FindName(transaction.Identity)

	if status != status.SUCCESS {
		return status
	}

	if entry != nil {
		log.Debug("Ignoring Existin Identity")
	} else {
		identities.Add(id.NewIdentity(transaction.Identity, "Contact Information",
			true, global.Current.NodeName, transaction.AccountKey))
		log.Info("Updated External Identity", "id", transaction.Identity, "key", transaction.AccountKey)
	}

	return status.SUCCESS
}

func (transaction *Register) Resolve(app interface{}) Commands {
	return []Command{}
}
