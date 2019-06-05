package identity

import (
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/serialize"
)

type Validator struct {
	Address keys.Address
	PubKey  keys.PublicKey
	Power   int64
	Name    string
	Staking balance.Coin
}

func (v Validator) Bytes() []byte {
	value, err := serialize.GetSerializer(serialize.PERSISTENT).Serialize(v)
	if err != nil {
		logger.Error("validator not serializable", err)
		return nil
	}
	return value
}

func (v *Validator) FromBytes(msg []byte) *Validator {
	err := serialize.GetSerializer(serialize.PERSISTENT).Deserialize(msg, v)
	if err != nil {
		logger.Error("failed to deserialize account from bytes", err)
		return &Validator{}
	}
	return v
}

type Stake struct {
	Address keys.Address
	Pubkey  keys.PublicKey
	Name    string
	Amount  balance.Coin
}

type Unstake struct {
	Address keys.Address
	Amount  balance.Coin
}

type ValidatorContext struct {
	Balances *balance.Store
	//todo: add necessary config
}

func NewValidatorContext(balances *balance.Store) *ValidatorContext {
	return &ValidatorContext{
		Balances: balances,
	}
}

func transferVT(ctx ValidatorContext, validator Validator) bool {
	logger.Debug("Processing Transfer of VT to Payment Account")

	//todo: implement transfer vt from account balance to some where else

	return true
}
