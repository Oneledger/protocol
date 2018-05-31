/*
	Copyright 2017 - 2018 OneLedger

	Register this identity with the other nodes. As an externl identity
*/
package action

import (
	"github.com/Oneledger/protocol/node/err"
	"github.com/Oneledger/protocol/node/id"
	"github.com/Oneledger/protocol/node/log"
	"github.com/Oneledger/protocol/node/persist"
)

// Register an identity with the chain
type Register struct {
	Base

	Identity string
}

func (transaction *Register) Validate() err.Code {
	log.Debug("Validating Send Transaction")

	// TODO: Make sure all of the parameters are there
	// TODO: Check all signatures and keys
	// TODO: Vet that the sender has the values
	return err.SUCCESS
}

func (transaction *Register) ProcessCheck(app persist.Access) err.Code {
	log.Debug("Processing Register Transaction for CheckTx")
	//xapp := global.Current.GetApplication()
	accounts := app.GetAccounts().(*id.Accounts)
	id, errs := accounts.Find(transaction.Identity)

	if errs != err.SUCCESS {
		return errs
	}

	if id == nil {
		return err.SUCCESS
	}

	// TODO: // Update in memory copy of Merkle Tree
	return err.DUPLICATE
}

func (transaction *Register) ProcessDeliver(app persist.Access) err.Code {
	log.Debug("Processing Register Transaction for DeliverTx")

	// TODO: // Update in final copy of Merkle Tree
	return err.SUCCESS
}

// Given a transaction, expand it into a list of Commands to execute against various chains.
func (transaction *Register) Expand(app persist.Access) Commands {
	// TODO: Table-driven mechanics, probably elsewhere
	return []Command{}
}
