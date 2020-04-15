package identity

import (
	"encoding/hex"
	"testing"

	"github.com/Oneledger/protocol/data/keys"
	"github.com/stretchr/testify/assert"
)

var (
	nodeAddr, _ = hex.DecodeString("F2143ADE3D941025468792311A0AB38D5085E15A")
	witness     = &Witness{
		Address: nodeAddr,
		PubKey: keys.PublicKey{
			KeyType: keys.ED25519,
			Data:    nil,
		},
		Name: "test_node",
	}
)

func TestETHWitness_Bytes(t *testing.T) {
	assert.NotEqual(t, []byte{}, witness.Bytes())
}

func TestETHWitness_FromBytes(t *testing.T) {
	value := witness.Bytes()
	witness_after := &Witness{}
	witness_after, err := witness_after.FromBytes(value)
	if assert.NoError(t, err) {
		assert.Equal(t, witness, witness_after)
	}
}
