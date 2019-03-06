package id

import (
	"bytes"
	"encoding/hex"
	"github.com/Oneledger/protocol/node/data"
	"github.com/Oneledger/protocol/node/log"
	"github.com/Oneledger/protocol/node/serial"
	"github.com/tendermint/tendermint/abci/types"
	"math/big"
	"strings"
)

type Validators struct {
	Signers            []types.VoteInfo
	Byzantines         []types.Evidence
	Approved           []Identity
	ApprovedValidators []Validator
	SelectedValidator  Identity
	NewValidators      []ApplyValidator
	ToBeRemoved        []ApplyValidator
}

type Validator struct {
	AccountKey AccountKey
	Address    []byte
	PubKey     types.PubKey
	Power      int64
}

type ApplyValidator struct {
	Validator Validator
	Stake     data.Coin
}

func init() {
	serial.Register(Validators{})
	serial.Register(Validator{})
}

func NewValidatorList() *Validators {
	return &Validators{}
}

func (validator Validator) String() {

}

func (list *Validators) Set(app interface{}, validators []types.VoteInfo, badValidators []types.Evidence, hash []byte) {
	if validators == nil {
		return
	}
	list.Signers = validators
	list.Byzantines = badValidators
	list.ApprovedValidators = make([]Validator, 0)
	list.Approved = list.FindApproved(app)
	list.NewValidators = make([]ApplyValidator, 0)
	list.ToBeRemoved = make([]ApplyValidator, 0)
	if hash != nil {
		list.SelectedValidator = list.FindSelectedValidator(app, hash)
	}

}

func (list *Validators) FindSelectedValidator(app interface{}, hash []byte) Identity {
	if len(list.Approved) < 1 {
		return Identity{}
	}

	countBigInt := big.NewInt(int64(len(list.Approved)))

	hashBigInt := new(big.Int).SetBytes(hash)

	indexBigInt := new(big.Int)
	indexBigInt = indexBigInt.Mod(hashBigInt, countBigInt)

	var indexInt64, _ = new(big.Int).SetString(indexBigInt.String(), 10)
	index := int(indexInt64.Int64())

	selectedValidator := list.Approved[index]

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
			tmp := Validator{
				Address: entry.GetValidator().Address,
				Power:   entry.GetValidator().Power,
			}
			validator := GetValidator(identity.TendermintAddress, identity.TendermintPubKey, entry.Validator.Power)
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

func GetValidator(address string, pubkey string, power int64) *Validator {
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
	return &Validator{
		Address: addr,
		PubKey:  tpubkey,
		Power:   power,
	}
}

func HasValidatorToken(app interface{}, validator Validator) bool {
	identities := GetIdentities(app)
	balances := GetBalances(app)

	formatted := hex.EncodeToString(validator.Address)
	identity := identities.FindTendermint(formatted)

	validatorBalance := balances.Get(identity.AccountKey, false)
	coin := validatorBalance.FindCoin(data.NewCurrency("VT"))
	if coin.LessThanEqual(0) {
		return false
	}

	return true
}
