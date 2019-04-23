package serialize

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type testStuff2 struct {
	F string
	A int
	B int64
	C []byte
}

var as *aminoStrategy

func init() {
	as = NewAminoStrategy(aminoCodec)
}

func TestAminoStrategy_Serialize(t *testing.T) {

	RegisterConcrete(new(testStuff2), "testStuff2")
	f := &testStuff2{"asdf", 1, 123123, []byte("4983h4tsdof")}
	fb, err := as.Serialize(f)

	assert.Nil(t, err)
	eb, err := aminoCodec.MarshalBinaryLengthPrefixed(f)
	assert.Equal(t, eb, fb)

	var fn = &testStuff2{}
	err = as.Deserialize(fb, fn)

	assert.Nil(t, err)
	assert.Equal(t, fn, f)
}
