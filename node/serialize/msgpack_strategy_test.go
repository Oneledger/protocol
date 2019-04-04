package serialize

import (
	"github.com/stretchr/testify/assert"
	"github.com/vmihailenco/msgpack"
	"testing"
)


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

	temp := testStuffAdData{"asdf", 1, 123123, []byte("4983h4tsdof"), 56.77867342}
	eb, err := msgpack.Marshal(temp)
	assert.Nil(t, err)
	assert.Equal(t, eb, fb)

	var fn = &testStuffAd{}
	err = ms.Deserialize(fb, fn)

	assert.Nil(t, err)
	assert.Equal(t, fn, f)

}
