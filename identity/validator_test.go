package identity

import (
	"encoding/hex"
	"testing"

	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/stretchr/testify/assert"
)

var hexString, _ = hex.DecodeString("89507C7ABC6D1E9124FE94101A0AB38D5085E15A")

var validator = &Validator{
	Address:      hexString,
	StakeAddress: hexString,
	PubKey: keys.PublicKey{
		KeyType: keys.ED25519,
		Data:    nil,
	},
	Power:   500,
	Name:    "test_node",
	Staking: *balance.NewAmount(100),
}

func TestValidator_Bytes(t *testing.T) {
	assert.NotEqual(t, []byte{}, validator.Bytes())
}

func TestValidator_FromBytes(t *testing.T) {
	validator, err := validator.FromBytes(validator.Bytes())
	if assert.NoError(t, err) {
		assert.Equal(t, validator, validator)
	}
}
