package keys

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddress_MarshalUnmarshalRoundTrip(t *testing.T) {
	var in Address
	in = []byte{6, 9, 6, 9}
	inBz, err := in.MarshalText()
	assert.Nil(t, err)

	out := new(Address)
	err = out.UnmarshalText(inBz)
	assert.Nil(t, err)

	assert.Equal(t, in, *out)
}
