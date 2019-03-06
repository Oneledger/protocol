/*
	Copyright 2017-2018 OneLedger

	An incoming transaction, send, swap, ready, verification, etc.
*/
package action

import (
	"bytes"
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
	Purge bool
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

	if transaction.Stake.Currency.Name != "VT" {
		log.Debug("Wrong token used for apply validator", "token", transaction.Stake)
		return status.INVALID
	}

	return status.SUCCESS
}

func (transaction *ApplyValidator) ProcessCheck(app interface{}) status.Code {
	log.Debug("Processing ApplyValidator Transaction for CheckTx")

	result := CheckBalances(app, transaction.Owner, transaction.AccountKey, transaction.Stake)
	if result == true {
		return status.SUCCESS
	} else {
		return status.INVALID
	}
}

func CheckBalances(app interface{}, owner id.AccountKey, identityAccountKey id.AccountKey, stake data.Coin) bool {

	balances := GetBalances(app)

	//check identity's VT is equal to the stake
	identityBalance := balances.Get(identityAccountKey, false)
	if identityBalance.GetAmountByName("VT").LessThanCoin(stake) {
		return false
	}

	//check administrator's VT is greater than 10
	if bytes.Compare(owner, identityAccountKey) != 0 {
		ownerBalance := balances.Get(owner, false)
		if ownerBalance.GetAmountByName("VT").LessThanCoin(data.NewCoinFromFloat(10.0, "VT")) {
			return false
		}
	}

	return true
}

func (transaction *ApplyValidator) ShouldProcess(app interface{}) bool {
	return true
}

func (transaction *ApplyValidator) ProcessDeliver(app interface{}) status.Code {
	log.Debug("Processing ApplyValidator Transaction for DeliverTx")

	result := CheckBalances(app, transaction.Owner, transaction.AccountKey, transaction.Stake)
	if result == false {
		return status.INVALID
	}

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
	validator := id.GetValidator(transaction.TendermintAddress, transaction.TendermintPubKey, 1)
	if validator == nil {
		return status.EXECUTE_ERROR
	}

	apply := id.ApplyValidator{
		Validator: *validator,
		Stake:     transaction.Stake,
	}

	if transaction.Purge == true {
		validators.ToBeRemoved = append(validators.ToBeRemoved, apply)
	} else {
		validators.NewValidators = append(validators.NewValidators, apply)
	}

	return status.SUCCESS
}

func (transaction *ApplyValidator) Resolve(app interface{}) Commands {
	return []Command{}
}
