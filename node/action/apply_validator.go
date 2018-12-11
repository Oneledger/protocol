/*
	Copyright 2017-2018 OneLedger

	An incoming transaction, send, swap, ready, verification, etc.
*/
package action

import (
	"github.com/Oneledger/protocol/node/data"
	"github.com/Oneledger/protocol/node/id"
	"github.com/Oneledger/protocol/node/log"
	"github.com/Oneledger/protocol/node/serial"
	"github.com/Oneledger/protocol/node/status"
)

// Apply a dynamic validator
type ApplyValidator struct {
	Base

	AccountKey        id.AccountKey
	Identity          string
	NodeName          string
	TendermintAddress string
	TendermintPubKey  string

	Stake data.Coin
}

func init() {
	serial.Register(ApplyValidator{})
}

func (transaction *ApplyValidator) Validate() status.Code {
	log.Debug("Validating ApplyValidator Transaction")

	baseValidate := transaction.Base.Validate()

	if baseValidate != status.SUCCESS {
		return baseValidate
	}

	if transaction.AccountKey == nil || len(transaction.AccountKey) == 0 {
		log.Debug("Missing AccountKey", "transaction", transaction)
		return status.MISSING_DATA
	}

	if transaction.TendermintAddress == "" {
		log.Debug("Missing TendermintAddress", "transaction", transaction)
		return status.MISSING_DATA
	}

	if transaction.TendermintPubKey == "" {
		log.Debug("Missing TendermintPubKey", "transaction", transaction)
		return status.MISSING_DATA
	}

	if transaction.Stake.LessThan(0) {
		log.Debug("Missing Stake", "ApplyValidator", transaction)
		return status.MISSING_DATA
	}

	return status.SUCCESS
}

func (transaction *ApplyValidator) ProcessCheck(app interface{}) status.Code {
	log.Debug("Processing ApplyValidator Transaction for CheckTx")

	return status.SUCCESS
}

func (transaction *ApplyValidator) ShouldProcess(app interface{}) bool {
	return true
}

func (transaction *ApplyValidator) ProcessDeliver(app interface{}) status.Code {
	log.Debug("Processing ApplyValidator Transaction for DeliverTx")
	identities := GetIdentities(app)
	entry, ok := identities.FindName(transaction.Identity)

	if ok != status.SUCCESS && ok != status.MISSING_DATA {
		log.Warn("Can't process Registration", "ok", ok)
		return ok
	}

	if entry.Name != "" {
		log.Debug("Ignoring Existing Identity", "identity", transaction.Identity)
	} else {
		identity := id.NewIdentity(transaction.Identity, "Contact Information",
			true, transaction.NodeName, transaction.AccountKey, transaction.TendermintAddress, transaction.TendermintPubKey)

		identities.Add(*identity)
		log.Info("Updated External Identity", "id", transaction.Identity, "key", transaction.AccountKey)
	}

	validators := GetValidators(app)
	validator := id.GetTendermintValidator(transaction.TendermintAddress, transaction.TendermintPubKey, 1)
	if validator == nil {
		return status.EXECUTE_ERROR
	}
	validators.NewValidators = append(validators.NewValidators, *validator)

	return status.SUCCESS
}

func (transaction *ApplyValidator) Resolve(app interface{}) Commands {
	return []Command{}
}
