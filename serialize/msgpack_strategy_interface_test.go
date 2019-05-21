package serialize

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vmihailenco/msgpack"
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

type Executer interface {
	Execute()
}

type printTask struct {
	Payload string
}

func (p *printTask) Execute() {
	fmt.Println(p.Payload)
}

// unregistered concrete type should not be deserialized through interfaces
func TestMsgpackStrategy_SerializeInterface3(t *testing.T) {
	RegisterInterface(new(Executer))
	RegisterConcrete(new(printTask), "print_task")

	f := &printTask{"string to print"}
	fb, err := ms.Serialize(f)

	assert.Nil(t, err)
	eb, err := msgpack.Marshal(f)
	assert.Equal(t, eb, fb)

	var fn Executer
	err = ms.Deserialize(fb, &fn)

	fnn, ok := fn.(*printTask)
	assert.True(t, ok)

	var i interface{}
	err = ms.Deserialize(fb, i)
	assert.Equal(t, err, ErrIncorrectWrapper)

	err = ms.Deserialize(fb, &i)
	assert.Nil(t, err)
	assert.Equal(t, fnn, i)
}
