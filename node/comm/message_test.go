/*
	Copywrite 2017-2018 OneLedger

	Tests to validate the way wire handles various struct conversions.
*/
package comm

import (
	"testing"

	"github.com/Oneledger/protocol/node/log"
	"github.com/stretchr/testify/assert"
	wire "github.com/tendermint/go-wire"
)

// Test just an integer
func TestInt(t *testing.T) {
	log.Info("Testing int")
	variable := 5

	buffer, err := Serialize(variable)

	if err != nil {
		log.Debug("Serialized failed", "err", err)
	} else {
		log.Debug("buffer", "buffer", buffer)
	}

	var integer int

	// result is probably not an integer
	result, err := Deserialize(buffer, integer)

	if err != nil {
		log.Debug("Deserialized failed", "err", err)
	} else {
		log.Debug("result", "result", result)
	}

	assert.Equal(t, variable, result, "These should be equal")
}

type BasicType interface {
	IsBasicType()
}
type Basic struct {
	Pad    string
	Number int
	Name   string
}

func (Basic) IsBasicType() {}

const typeBasic = 0x88

func Register() {

	// Tell wire how to find the underlying types for a given interface
	var _ = wire.RegisterInterface(
		struct{ BasicType }{},
		wire.ConcreteType{Basic{}, typeBasic},
	)
}

// Basic structural conversion
func TestStruct(t *testing.T) {
	log.Info("Testing Struct")

	Register()

	basic := &Basic{Pad: "xxxx", Number: 123456, Name: "A Name"}
	log.Debug("The basic type", "basic", basic)

	buffer, err := Serialize(basic)

	if err != nil {
		assert.FailNow(t, "Serialization failed ", err.Error())
	} else {
		log.Debug("buffer", "buffer", buffer)
	}

	result := new(Basic)
	base, err := Deserialize(buffer, result)
	value := base.(*Basic)

	if err != nil {
		assert.Fail(t, "Deserialization failed", err.Error())
	} else {
		log.Debug("result", "value", value)
	}

	assert.Equal(t, basic, value, "These should be equal")
}

// Basic structural conversion
func TestStruct2(t *testing.T) {
	log.Info("Testing Struct2")

	Register()

	basic := &Basic{Pad: "xxxx", Number: 123456, Name: "A Name"}
	log.Debug("The basic type", "basic", basic)

	buffer, err := Serialize(basic)

	if err != nil {
		log.Debug("Serialized failed", "err", err)
	} else {
		log.Debug("buffer", "buffer", buffer)
	}

	result, err := Deserialize(buffer, struct{ *Basic }{})
	value := result.(struct{ *Basic }).Basic

	if err != nil {
		log.Debug("Deserialized failed", "err", err)
	} else {
		log.Debug("result", "value", value)
	}

	assert.Equal(t, basic, value, "These should be equal")
}

func TestInterface(t *testing.T) {
	log.Info("Testing Interface")

	Register()

	var generic BasicType

	generic = &Basic{Pad: "-------", Number: 654321, Name: "A rose by any other name"}
	log.Debug("generic", "generic", generic)

	buffer, err := Serialize(generic)
	if err != nil {
		assert.FailNow(t, "Serialization failed", err.Error())
	} else {
		log.Debug("buffer", "buffer", buffer)
	}

	result, err := Deserialize(buffer, struct{ BasicType }{})
	if err != nil {
		assert.Fail(t, "Deserialization failed", err.Error())
	} else {
		log.Debug("result", "result", result)
	}

	value := result.(struct{ BasicType }).BasicType
	log.Debug("value", "value", value)
	assert.Equal(t, generic, value, "These should be equal")
}
