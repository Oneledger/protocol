package serialize

import (
	"github.com/stretchr/testify/assert"
	"github.com/vmihailenco/msgpack"
	"testing"
)


// registered concrete type should be deserialized through interfaces
func TestMsgpackStrategy_SerializeInterface(t *testing.T) {
	msgpack.RegisterExt(14, new(testStuff))


	f := &testStuff{"asdf", 1, 123123, []byte("4983h4tsdof"), 56.77867342}
	fb, err := ms.Serialize(f)

	assert.Nil(t, err)
	eb, err := msgpack.Marshal(f)
	assert.Equal(t, eb, fb)

	var fn interface{}
	err = ms.Deserialize(fb, &fn)

	fnn, ok := fn.(*testStuff)
	assert.True(t, ok)

	assert.Nil(t, err)
	assert.Equal(t, fnn, f)
}

// unregistered concrete type should not be deserialized through interfaces
func TestMsgpackStrategy_SerializeInterface2(t *testing.T) {

	f := &testStuffAdData{"asdf", 1, 123123, []byte("4983h4tsdof"), "56.77867342"}
	fb, err := ms.Serialize(f)
	assert.Nil(t, err)

	eb, err := msgpack.Marshal(f)
	assert.Nil(t, err)
	assert.Equal(t, eb, fb)

	var fn interface{}
	err = ms.Deserialize(fb, &fn)
	assert.Nil(t, err)

	fnn, ok := fn.(*testStuffAdData)
	assert.False(t, ok)

	assert.NotEqual(t, fnn, f)
}
