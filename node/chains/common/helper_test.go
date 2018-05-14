package common


import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseInt(t *testing.T) {
	i, err := ParseInt("0x41")
	assert.Nil(t, err)
	assert.Equal(t, 65, i)

	i, err = ParseInt("41")
	assert.Nil(t, err)
	assert.Equal(t, 65, i)

	i, err = ParseInt("0xabc")
	assert.Nil(t, err)
	assert.Equal(t, 2748, i)

	i, err = ParseInt("1*29")
	assert.NotNil(t, err)
	assert.Equal(t, 0, i)
}

func TestParseBigInt(t *testing.T) {
	i, err := ParseBigInt("0xde0b6b3a7640000")
	assert.Nil(t, err)
	assert.Equal(t, int64(1000000000000000000), i.Int64())

	i, err = ParseBigInt("$%1")
	assert.NotNil(t, err)
}

func TestIntToHex(t *testing.T) {
	assert.Equal(t, "0xabc", IntToHex(2748))
	assert.Equal(t, "0x41", IntToHex(65))
}

func TestBigToHex(t *testing.T) {
	i1, _ := big.NewInt(0).SetString("1000000000000000000", 10)
	assert.Equal(t, "0xde0b6b3a7640000", BigToHex(*i1))

	i2, _ := big.NewInt(0).SetString("100000000000000000000", 10)
	assert.Equal(t, "0x56bc75e2d63100000", BigToHex(*i2))
}