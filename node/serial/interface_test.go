package serial

import (
	"testing"

	"github.com/Oneledger/protocol/node/log"
	"github.com/stretchr/testify/assert"
)

type Wrapper struct {
	Interface ArrayInterface
}

type Array1 [32]byte

type ArrayInterface interface {
	Key() []byte
}

func (a Array1) Key() []byte { return a[:] }

func init() {
	Register(Wrapper{})
	Register(Array1{})
	var p ArrayInterface
	RegisterInterface(&p)
}

func TestSerializeByte(t *testing.T) {
	DumpTypes()

	a := Array1{20, 19, 18, 17, 16, 15, 14, 13, 12}

	buffer, err := Serialize(a, PERSISTENT)
	if err != nil {
		log.Fatal("Serialize failed", "err", err)
	}

	//var b ArrayInterface
	var b Array1

	result, err := Deserialize(buffer, b, PERSISTENT)
	if err != nil {
		log.Fatal("Deserialize failed", "err", err)
	}

	assert.Equal(t, a, result, "These should be equal")
}

func TestSerializeInterface(t *testing.T) {
	DumpTypes()

	a := Array1{20, 19, 18, 17, 16, 15, 14, 13, 12}

	buffer, err := Serialize(a, PERSISTENT)
	if err != nil {
		log.Fatal("Serialize failed", "err", err)
	}

	var b ArrayInterface

	result, err := Deserialize(buffer, b, PERSISTENT)

	if err != nil {
		log.Fatal("Deserialize failed", "err", err)
	}

	assert.Equal(t, a, result, "These should be equal")

}

func TestSerializeWrappedInterface(t *testing.T) {
	DumpTypes()

	a := Array1{20, 19, 18, 17, 16, 15, 14, 13, 12}
	s := Wrapper{a}

	buffer, err := Serialize(s, PERSISTENT)
	if err != nil {
		log.Fatal("Serialize failed", "err", err)
	}

	var b Wrapper

	result, err := Deserialize(buffer, b, PERSISTENT)

	if err != nil {
		log.Fatal("Deserialize failed", "err", err)
	}

	assert.Equal(t, s, result, "These should be equal")

}
