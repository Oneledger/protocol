package app

import (
	"bytes"
	"encoding/hex"
	"github.com/Oneledger/protocol/node/id"
	"github.com/Oneledger/protocol/node/serial"
	"github.com/tendermint/tendermint/abci/types"
	"math/big"
)

type ValidatorList struct {
	Signers   []types.SigningValidator
	Byzantine []types.Evidence
}

func init() {
	serial.Register(ValidatorList{})
}

func NewValidatorList() *ValidatorList {
	return &ValidatorList{}
}

func (list *ValidatorList) Set(validators []types.SigningValidator, badValidators []types.Evidence) {
	list.Signers = validators
	list.Byzantine = badValidators
}

func (list *ValidatorList) FindSelectedValidator(app Application, hash []byte) id.Identity {
	goodList := list.FindGood(app)
	countBigInt := big.NewInt(int64(len(goodList)))
	hashBigInt := new(big.Int).SetBytes(hash)
	indexBigInt := new(big.Int)
	indexBigInt = indexBigInt.Mod(hashBigInt, countBigInt)
	var indexInt64, _ = new(big.Int).SetString(indexBigInt.String(), 10)
	index := int(indexInt64.Int64())
	selectedValidator := goodList[index]
	return selectedValidator
}

func (list *ValidatorList) FindGood(app Application) []id.Identity {
	var goodIdentities []id.Identity
	for _, entry := range list.Signers {
		entryIsBad := IsByzantine(entry.Validator, list.Byzantine)
		if !entryIsBad {
			formatted := hex.EncodeToString(entry.Validator.Address)
			identity := app.Identities.FindTendermint(formatted)
			goodIdentities = append(goodIdentities, identity)
		}
	}
	return goodIdentities
}

func IsByzantine(validator types.Validator, badValidators []types.Evidence) (result bool) {
	for _, entry := range badValidators {
		if bytes.Equal(validator.Address, entry.Validator.Address) {
			return true
		}
	}
	return false
}
