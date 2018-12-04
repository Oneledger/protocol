/*
	Copyright 2017 - 2018 OneLedger

	Register this identity with the other nodes. As an externl identity
*/
package action

import (
	"github.com/Oneledger/protocol/node/id"
	"github.com/Oneledger/protocol/node/log"
	"github.com/Oneledger/protocol/node/serial"
	"github.com/Oneledger/protocol/node/status"
)

// Register an identity with the chain
type Register struct {
	Base

	Identity          string
	NodeName          string
	AccountKey        id.AccountKey
	TendermintAddress string
	TendermintPubKey  string
}

func init() {
	serial.Register(Register{})
}

// Check the fields to make sure they have valid values.
func (transaction Register) Validate() status.Code {
	log.Debug("Validating Register Transaction")

	baseValidate := transaction.Base.Validate()

	if baseValidate != status.SUCCESS {
		return baseValidate
	}

	if transaction.Identity == "" {
		log.Debug("Missing Identity", "transaction", transaction)
		return status.MISSING_DATA
	}

	if transaction.NodeName == "" {
		log.Debug("Missing NodeName", "transaction", transaction)
		return status.MISSING_DATA
	}

	if transaction.AccountKey == nil || len(transaction.AccountKey) == 0 {
		log.Debug("Missing AccountKey", "transaction", transaction)
		return status.MISSING_DATA
	}

	if transaction.TendermintAddress == "" {
		log.Debug("Missing TendermintAddress", "transaction", transaction)
		return status.MISSING_DATA
	}

	if transaction.TendermintAddress == "" {
		log.Debug("Missing TendermintPubKey", "transaction", transaction)
		return status.MISSING_DATA
	}

	return status.SUCCESS
}

// Test to see if the identity already exists
func (transaction Register) ProcessCheck(app interface{}) status.Code {
	log.Debug("Processing Register Transaction for CheckTx")
	/*
		identities := GetIdentities(app)
		id, ok := identities.FindName(transaction.Identity)

		if ok != status.SUCCESS {
			return ok
		}
	*/

	/*
		if id == nil {
			log.Debug("Success, it is a new Identity", "id", transaction.Identity)
			return err.SUCCESS
		}
	*/

	// Not necessarily a failure, since this identity might be local
	//log.Debug("Identity already exists", "id", id)
	return status.SUCCESS
}

func (transaction Register) ShouldProcess(app interface{}) bool {
	return true
}

// Add the identity into the database as external, don't overwrite a local identity
func (transaction Register) ProcessDeliver(app interface{}) status.Code {
	log.Debug("Processing Register Transaction for DeliverTx")

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

	return status.SUCCESS
}

func (transaction *Register) Resolve(app interface{}) Commands {
	return []Command{}
}

func CreateRegisterRequest(identity string, chainId string, sequence int64, nodeName string, signers []PublicKey, accountKey id.AccountKey) *Register {
	return &Register{
		Base: Base{
			Type:     REGISTER,
			ChainId:  chainId,
			Signers:  signers,
			Sequence: sequence,
		},
		Identity:   identity,
		NodeName:   nodeName,
		AccountKey: accountKey,
	}
}
