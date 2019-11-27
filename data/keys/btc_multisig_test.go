/*

 */

package keys

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBTCMultiSig(t *testing.T) {

	msg := hexDump("a22b23c32d534")

	addrs := []Address{
		hexDump("aaaaaa"),
		hexDump("bbbbbb"),
		hexDump("cccccc"),
		hexDump("dddddd"),
	}
	m, err := NewBTCMultiSig(msg, 3, addrs)
	assert.NoError(t, err)

	err = m.AddSignature(&BTCSignature{
		3,
		hexDump("dddddd"),
		hexDump("dddddd"),
	})
	assert.NoError(t, err)

	err = m.AddSignature(&BTCSignature{
		1,
		hexDump("aaaaaa"),
		hexDump("aaaaaa"),
	})
	assert.Error(t, err)

	err = m.AddSignature(&BTCSignature{
		0,
		hexDump("aaaaaa"),
		hexDump("aaaaaa"),
	})
	assert.NoError(t, err)

	assert.False(t, m.IsValid())

	err = m.AddSignature(&BTCSignature{
		1,
		hexDump("bbbbbb"),
		hexDump("bbbbbb"),
	})
	assert.NoError(t, err)
	assert.True(t, m.IsValid())

	assert.True(t, m.HasAddressSigned(hexDump("bbbbbb")))

	assert.False(t, m.HasAddressSigned(hexDump("cccccc")))

	assert.True(t, bytes.Equal(m.Signers[0], hexDump("aaaaaa")))
	assert.True(t, bytes.Equal(m.Signers[1], hexDump("bbbbbb")))

	signatures := make([][]byte, len(m.Signatures))
	for i, signed := range m.GetSignatures() {
		signatures[i] = signed.Sign
	}

	assert.Equal(t, "aaaaaa", hex.EncodeToString(signatures[0]))
	assert.Equal(t, "bbbbbb", hex.EncodeToString(signatures[1]))
	assert.Equal(t, "dddddd", hex.EncodeToString(signatures[3]))

	sass := m.GetSignaturesInOrder()

	assert.Equal(t, len(sass), 3)
	assert.Equal(t, "aaaaaa", hex.EncodeToString(sass[0]))
	assert.Equal(t, "bbbbbb", hex.EncodeToString(sass[1]))
	assert.Equal(t, "dddddd", hex.EncodeToString(sass[2]))
}

func hexDump(a string) []byte {
	b, _ := hex.DecodeString(a)
	return b
}
