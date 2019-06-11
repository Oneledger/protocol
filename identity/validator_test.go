package identity

import (
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

var v = &Validator{
	Address: []byte("89507C7ABC6D1E9124FE94101A0AB38D5085E15A"),
	StakeAddress: []byte("89507C7ABC6D1E9124FE94101A0AB38D5085E15A"),
	PubKey: keys.PublicKey{0, []byte("89507C7ABC6D1E9124FE94101A0AB38D5085E15A")},
	Power: 500,
	Name: "test node",
	Staking: balance.Coin{balance.Currency{"VT", 1, 18}, big.NewInt(100.0)},
}

func TestValidator_Bytes(t *testing.T) {
	t.Run("run Bytes test case", func(t *testing.T) {
		assert.NotEqual(t, []byte{}, v.Bytes())
	})

}

func TestValidator_FromBytes(t *testing.T) {
	t.Run("run FromBytes test case", func(t *testing.T) {
		assert.NotEqual(t, &Validator{}, v.FromBytes(v.Bytes()))
	})
}

