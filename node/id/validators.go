package id

import (
	"bytes"
	"encoding/hex"
	"math/big"
	"strings"

	"github.com/Oneledger/protocol/node/data"
	"github.com/Oneledger/protocol/node/log"
	"github.com/Oneledger/protocol/node/serial"
	"github.com/tendermint/tendermint/abci/types"
)

type Validators struct {
	Signers            []types.SigningValidator
	Byzantines         []types.Evidence
	Approved           []Identity
	ApprovedValidators []types.Validator
	SelectedValidator  Identity
	NewValidators      []types.Validator
	ToBeRemoved        []types.Validator
}

type ValidatorInfo struct {
	Address string
	PubKey  string
}

func init() {
	serial.Register(Validators{})
	serial.Register(ValidatorInfo{})
}

func NewValidatorList() *Validators {
	return &Validators{}
}

func NewValidatorInfo(address []byte, pubkey types.PubKey) *ValidatorInfo {
	return &ValidatorInfo{
		Address: hex.EncodeToString(address),
		PubKey:  hex.EncodeToString(pubkey.Data),
	}
}

func (list *Validators) Set(app interface{}, validators []types.SigningValidator, badValidators []types.Evidence, hash []byte) {
	if validators == nil {
		return
	}
	list.Signers = validators
	list.Byzantines = badValidators
	list.ApprovedValidators = make([]types.Validator, 0)
	list.Approved = list.FindApproved(app)
	list.NewValidators = make([]types.Validator, 0)
	list.ToBeRemoved = make([]types.Validator, 0)
	if hash != nil {
		list.SelectedValidator = list.FindSelectedValidator(app, hash)
	}

}

func (list *Validators) FindSelectedValidator(app interface{}, hash []byte) Identity {
	log.Debug("FindSelectedValidator", "hash", hash, "approved", len(list.Approved))

	if len(list.Approved) < 1 {
		return Identity{}
	}

	countBigInt := big.NewInt(int64(len(list.Approved)))

	hashBigInt := new(big.Int).SetBytes(hash)

	indexBigInt := new(big.Int)
	indexBigInt = indexBigInt.Mod(hashBigInt, countBigInt)

	var indexInt64, _ = new(big.Int).SetString(indexBigInt.String(), 10)
	index := int(indexInt64.Int64())

	log.Dump("Calcs", countBigInt, indexBigInt, index, len(list.Approved))

	selectedValidator := list.Approved[index]

	log.Dump("Approved", list.Approved)

	log.Debug("Selected", "validator", selectedValidator)
	return selectedValidator
}

func (list *Validators) FindApproved(app interface{}) []Identity {
	var approvedIdentities []Identity
	for _, entry := range list.Signers {
		entryIsBad := IsByzantine(entry.Validator, list.Byzantines)
		if !entryIsBad {
			identities := GetIdentities(app)

			formatted := hex.EncodeToString(entry.Validator.Address)
			identity := identities.FindTendermint(formatted)
			if identity.Name == "" {
				log.Debug("Unable to Find Tendermint identity", "formatted", formatted,
					"raw", entry.Validator.Address)
				continue
			}

			approvedIdentities = append(approvedIdentities, identity)

			//add the approved validators to the NewValidators list to be used in the EndBlock call
			tmp := entry.Validator
			validator := GetTendermintValidator(identity.TendermintAddress, identity.TendermintPubKey, entry.Validator.Power)
			if validator != nil {
				tmp = *validator
			}
			list.ApprovedValidators = append(list.ApprovedValidators, tmp)
		}
	}
	return approvedIdentities
}

func IsByzantine(validator types.Validator, badValidators []types.Evidence) (result bool) {
	for _, entry := range badValidators {
		if bytes.Equal(validator.Address, entry.Validator.Address) {
			return true
		}
	}
	return false
}

func (list Validators) IsValidAccountKey(key AccountKey, index int) bool {
	if index >= len(list.Approved) || index < 0 {
		return false
	}

	id := list.Approved[index]
	if bytes.Equal(id.AccountKey.Bytes(), key.Bytes()) {
		return true
	}

	return false
}

func GetTendermintValidator(address string, pubkey string, power int64) *types.Validator {
	buffer, err := hex.DecodeString(pubkey)
	if err != nil {
		log.Debug("Failed to decode the pubkey", "pubkey", pubkey, "err", err)
		return nil
	}

	if len([]byte(buffer)) != ED25519.Size() {
		log.Debug("Wrong PubKey string length", "buffer", buffer, "size", ED25519.Size())
		return nil
	}

	key, err := ImportBytesKey(buffer, ED25519)
	if err != nil {
		log.Debug("Failed to convert the pubkey", "buffer", pubkey, "err", err)
		return nil
	}
	tpubkey := types.PubKey{
		Type: strings.ToLower(ED25519.String()),
		Data: key.Bytes(),
	}

	addr, err := hex.DecodeString(address)
	if err != nil {
		log.Debug("Failed to decode address", "addres", address)
		return nil
	}
	return &types.Validator{
		Address: addr,
		PubKey:  tpubkey,
		Power:   power,
	}
}

func HasValidatorToken(app interface{}, validator types.Validator) bool {
	identities := GetIdentities(app)
	balances := GetBalances(app)

	formatted := hex.EncodeToString(validator.Address)
	identity := identities.FindTendermint(formatted)

	validatorBalance := balances.Get(identity.AccountKey)
	coin := validatorBalance.FindCoin(data.NewCurrency("VT"))
	if coin.LessThanEqual(0) {
		return false
	}

	return true
}
