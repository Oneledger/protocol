package identity

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"

	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/delegation"
	"github.com/Oneledger/protocol/data/fees"
	"github.com/Oneledger/protocol/data/governance"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/serialize"
)

type Validator struct {
	Address      keys.Address   `json:"address"`
	StakeAddress keys.Address   `json:"stakeAddress"`
	PubKey       keys.PublicKey `json:"pubKey"`
	ECDSAPubKey  keys.PublicKey `json:"ecdsaPubkey"`
	Power        int64          `json:"power"`
	Name         string         `json:"name"`
	Staking      balance.Amount `json:"staking,string"`
}

func NewValidator(address keys.Address, stakeAddress keys.Address, pubKey keys.PublicKey, ecdsaPubKey keys.PublicKey, amount balance.Amount, name string) *Validator {
	return &Validator{
		Address:      address,
		StakeAddress: stakeAddress,
		PubKey:       pubKey,
		ECDSAPubKey:  ecdsaPubKey,
		Power:        calculatePower(amount),
		Name:         name,
		Staking:      amount,
	}
}

func (v *Validator) Bytes() []byte {
	value, err := serialize.GetSerializer(serialize.PERSISTENT).Serialize(v)
	if err != nil {
		logger.Error("validator not serializable", err)
		return []byte{}
	}
	return value
}

func (v *Validator) FromBytes(msg []byte) (*Validator, error) {
	err := serialize.GetSerializer(serialize.PERSISTENT).Deserialize(msg, v)
	if err != nil {
		logger.Error("failed to deserialize account from bytes", err)
		return nil, err
	}
	return v, nil
}

func (v *Validator) GetBTCScriptAddress(params *chaincfg.Params) (keys.Address, error) {

	h, err := v.ECDSAPubKey.GetHandler()
	if err != nil {
		return nil, err
	}

	apk, err := btcutil.NewAddressPubKey(h.Bytes(), params)
	if err != nil {
		return nil, err
	}

	return apk.ScriptAddress(), nil
}

type Stake struct {
	ValidatorAddress keys.Address
	StakeAddress     keys.Address
	Pubkey           keys.PublicKey
	ECDSAPubKey      keys.PublicKey
	Name             string
	Amount           balance.Amount
}

type Unstake struct {
	Address keys.Address
	Amount  balance.Amount
}

type ValidatorContext struct {
	Balances   *balance.Store
	FeePool    *fees.Store
	Delegators *delegation.DelegationStore
	Govern     *governance.Store
	// TODO: add necessary config
}

func NewValidatorContext(balances *balance.Store, feePool *fees.Store, delegators *delegation.DelegationStore, govern *governance.Store) *ValidatorContext {
	return &ValidatorContext{
		Balances:   balances,
		FeePool:    feePool,
		Delegators: delegators,
		Govern:     govern,
	}
}
