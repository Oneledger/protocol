/*
	Copywrite 2017-2018 OneLedger

	Tests to validate the way wire handles various struct conversions.
*/
package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
	wire "github.com/tendermint/go-wire"
)

// Test just an integer
func TestInt(t *testing.T) {
	Log.Info("Testing int")
	variable := 5

	buffer, err := Serialize(variable)

	if err != nil {
		Log.Debug("Serialized failed", "err", err)
	} else {
		Log.Debug("buffer", "buffer", buffer)
	}

	var integer int

	// result is probably not an integer
	result, err := Deserialize(buffer, integer)

	if err != nil {
		Log.Debug("Deserialized failed", "err", err)
	} else {
		Log.Debug("result", "result", result)
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
	Log.Info("Testing Struct")

	//Register()

	basic := &Basic{Pad: "xxxx", Number: 123456, Name: "A Name"}
	Log.Debug("The basic type", "basic", basic)

	buffer, err := Serialize(basic)

	if err != nil {
		Log.Debug("Serialized failed", "err", err)
		assert.FailNow(t, "Serialization failed ")
	} else {
		Log.Debug("buffer", "buffer", buffer)
	}

	result := new(Basic)
	base, err := Deserialize(buffer, result)
	value := base.(*Basic)

	if err != nil {
		assert.FailNow(t, "Deserialization failed", err.Error())
	} else {
		Log.Debug("result", "value", value)
	}

	assert.Equal(t, basic, value, "These should be equal")
}

// Basic structural conversion
func TestStruct2(t *testing.T) {
	Log.Info("Testing Struct2")

	//Register()

	basic := &Basic{Pad: "xxxx", Number: 123456, Name: "A Name"}
	Log.Debug("The basic type", "basic", basic)

	buffer, err := Serialize(basic)

	if err != nil {
		Log.Debug("Serialized failed", "err", err)
	} else {
		Log.Debug("buffer", "buffer", buffer)
	}

	result, err := Deserialize(buffer, struct{ *Basic }{})
	value := result.(struct{ *Basic }).Basic

	if err != nil {
		Log.Debug("Deserialized failed", "err", err)
	} else {
		Log.Debug("result", "value", value)
	}

	assert.Equal(t, basic, value, "These should be equal")
}

func TestInterface(t *testing.T) {
	Log.Info("Testing Interface")

	Register()

	var generic BasicType

	generic = &Basic{Pad: "-------", Number: 654321, Name: "A rose by any other name"}

	buffer, err := Serialize(generic)

	if err != nil {
		Log.Debug("Serialized failed", "err", err)
	} else {
		Log.Debug("buffer", "buffer", buffer)
	}

	result, err := Deserialize(buffer, generic)
	//value := result.(BasicType).(*Basic)
	value := result.(*Basic)

	if err != nil {
		Log.Debug("Deserialized failed", "err", err)
	} else {
		Log.Debug("result", "result", result)
		Log.Debug("result", "value", value)
	}

	assert.Equal(t, generic, value, "These should be equal")
}
