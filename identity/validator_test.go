package identity

import (
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/stretchr/testify/assert"
)

var hexstring, _ = hex.DecodeString("89507C7ABC6D1E9124FE94101A0AB38D5085E15A")

var v = &Validator{
	Address:      hexstring,
	StakeAddress: hexstring,
	PubKey: keys.PublicKey{
		KeyType: keys.ED25519,
		Data:    nil,
	},
	Power:   500,
	Name:    "test node",
	Staking: balance.Coin{balance.Currency{"VT", 1, 18}, big.NewInt(100.0)},
}

func TestValidator_Bytes(t *testing.T) {
	assert.NotEqual(t, []byte{}, v.Bytes())
}

func TestValidator_FromBytes(t *testing.T) {
	validator, err := v.FromBytes(v.Bytes())
	if assert.NoError(t, err) {
		assert.Equal(t, v, validator)
	}
}

func TestNewValidatorContext(t *testing.T) {
	balance := &balance.Store{}
	vc := &ValidatorContext{
		Balances: balance,
	}
	ValidatorContext := NewValidatorContext(balance)
	assert.Equal(t, ValidatorContext, vc)
}

