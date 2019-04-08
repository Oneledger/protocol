package serialize

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
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
