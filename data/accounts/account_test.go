package accounts

import (
	"testing"

	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/serialize"
	"github.com/stretchr/testify/assert"
)

func TestAccount(t *testing.T) {
	acc, err := NewAccount(0, "accountName", nil, nil)
	assert.Equal(t, acc, Account{})
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "empty keys")

	pubkey, privkey, err := keys.NewKeyPairFromTendermint()
	assert.NoError(t, err)

	acc, err = NewAccount(0, "accountName", &privkey, &pubkey)
	assert.NoError(t, err)
	assert.NotEqual(t, acc, Account{})

	address := acc.Address()
	assert.Len(t, address, 20)

	assert.NotNil(t, acc.Bytes())

	msg := []byte("981h9th983hf32894h09aish4089h2930ihjd2893h4r9283h8wiejd0923jd0923jed")
	signed, err := acc.Sign(msg)
	assert.NoError(t, err)
	assert.NotNil(t, signed)

	str := acc.String()
	assert.NotEqual(t, len(str), 0)

	b := acc.Bytes()
	accNew := &Account{}
	accNew = accNew.FromBytes(b)

	assert.NotNil(t, accNew)
	assert.Equal(t, accNew, &acc)
}

func TestAccountSerialize(t *testing.T) {

	pubkey, privkey, err := keys.NewKeyPairFromTendermint()
	assert.NoError(t, err)

	acc, err := NewAccount(0, "accountName", &privkey, &pubkey)
	assert.NoError(t, err)
	assert.NotEqual(t, acc, Account{})

	dat, err := serialize.GetSerializer(serialize.NETWORK).Serialize(&acc)
	assert.NoError(t, err)
	assert.NotEqual(t, len(dat), 0)

	accNew := &Account{}
	err = serialize.GetSerializer(serialize.NETWORK).Deserialize(dat, accNew)
	assert.NoError(t, err)
	assert.Equal(t, accNew, &acc)

}
