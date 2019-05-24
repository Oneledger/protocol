package serialize

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vmihailenco/msgpack"
)

type testStuff struct {
	F string
	A int
	B int64
	C []byte
	H float64
}

var js = jsonStrategy{}

func TestJsonStrategy_Serialize(t *testing.T) {
	f := &testStuff{"asdf", 1, 123123, []byte("4983h4tsdof"), 56.77867342}
	fb, err := js.Serialize(f)

	assert.Nil(t, err)
	eb, err := json.Marshal(f)
	assert.Equal(t, eb, fb)

	var fn = &testStuff{}
	err = js.Deserialize(fb, fn)

	assert.Nil(t, err)
	assert.Equal(t, fn, f)
}

func TestJsonStrategy_SerializeForAdapter(t *testing.T) {

	f := &testStuffAd{"asdf", 1, 123123, []byte("4983h4tsdof"), 56.77867342}
	fb, err := js.Serialize(f)
	assert.Nil(t, err)

	temp := testStuffAdData{"asdf", 1, 123123, []byte("4983h4tsdof"), "56.77867342"}
	eb, err := json.Marshal(temp)
	assert.Nil(t, err)
	assert.Equal(t, eb, fb)

	var fn = &testStuffAd{}
	err = js.Deserialize(fb, fn)

	assert.Nil(t, err)
	assert.Equal(t, fn, f)

}

var ms = msgpackStrategy{}

func TestMsgpackStrategy_Serialize(t *testing.T) {
	f := &testStuff{"asdf", 1, 123123, []byte("4983h4tsdof"), 56.77867342}
	fb, err := ms.Serialize(f)

	assert.Nil(t, err)
	eb, err := msgpack.Marshal(f)
	assert.Equal(t, eb, fb)

	var fn = &testStuff{}
	err = ms.Deserialize(fb, fn)

	assert.Nil(t, err)
	assert.Equal(t, fn, f)
}

func TestMsgpackStrategy_SerializeForAdapters(t *testing.T) {

	f := &testStuffAd{"asdf", 1, 123123, []byte("4983h4tsdof"), 56.77867342}
	fb, err := ms.Serialize(f)
	assert.Nil(t, err)

	temp := testStuffAdData{"asdf", 1, 123123, []byte("4983h4tsdof"), "56.77867342"}
	eb, err := msgpack.Marshal(temp)
	assert.Nil(t, err)
	assert.Equal(t, eb, fb)

	var fn = &testStuffAd{}
	err = ms.Deserialize(fb, fn)

	assert.Nil(t, err)
	assert.Equal(t, fn, f)

}

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

func TestGobStrategy_Serialize(t *testing.T) {
	gs := gobStrategy{}

	f := &testStuff2{"asdf", 1, 123123, []byte("4983h4tsdof")}
	fb, err := gs.Serialize(f)

	assert.NoError(t, err)

	var fn = &testStuff2{}
	err = gs.Deserialize(fb, fn)

	assert.Nil(t, err)
	assert.Equal(t, fn, f)

}
