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

type OpaqueParent interface {
	AFunction(value int) bool
}

type OpaqueChild interface {
	AnotherFunction(value int) bool
}

type ParentStruct1 struct {
	ParentName string
	Child      OpaqueChild
}

func (p ParentStruct1) AFunction(value int) bool {
	return false
}

type ParentStruct2 struct {
	ParentName string
	Count      int
	Child1     OpaqueChild
	Child2     OpaqueChild
}

func (p ParentStruct2) AFunction(value int) bool {
	return false
}

type ChildStruct1 struct {
	ChildName string
	Size      int
}

func (c ChildStruct1) AnotherFunction(value int) bool {
	return false
}

type ChildStruct2 struct {
	ChildName string
	Size      int
	Size2     int
}

func (c ChildStruct2) AnotherFunction(value int) bool {
	return false
}

func init() {
	// Let the deserilization code know how to create these structures
	Register(ParentStruct1{})
	Register(ParentStruct2{})
	Register(ChildStruct1{})
	Register(ChildStruct2{})
}

var child1 = ChildStruct1{
	ChildName: "The Child",
	Size:      1000,
}

var child2 = ChildStruct2{
	ChildName: "The Child",
	Size:      3000,
}

var opc1 = OpaqueChild(&child1)
var opc2 = OpaqueChild(child2)

var parent = ParentStruct2{
	ParentName: "The Parent",
	Count:      2000,
	Child1:     opc1,
	Child2:     opc2,
}

var opp = OpaqueParent(&parent)

func TestIteration(t *testing.T) {
	log.Debug("CloneIt")
	opp2 := Clone(opp)

	assert.Equal(t, opp, opp2, "These should be equal")

	opp3 := Extend(opp)
	log.Debug("Extend", "result", opp3)

	assert.Equal(t, opp, opp3, "These should be equal")

	/*
		log.Debug("ExtendIt")
		Iterate(opp, Action{ProcessField: ExtendIt})

		log.Debug("ContractIt")
		Iterate(opp, Action{ProcessField: ContractIt})
	*/
}

func XTestPolymorphism(t *testing.T) {

	buffer, err := Serialize(opp, PERSISTENT)

	if err != nil {
		log.Debug("Serialized failed", "err", err)
	} else {
		log.Debug("buffer", "buffer", buffer)
	}

	log.Debug("Serialized Worked, now Deserialize")

	var opp2 OpaqueParent

	result, err := Deserialize(buffer, opp2, PERSISTENT)

	if err != nil {
		log.Debug("Deserialized failed", "err", err)
	} else {
		log.Debug("result", "result", result)
	}

	assert.Equal(t, opp, result, "These should be equal")
}

// Test just an integer
func XTestInt(t *testing.T) {
	log.Info("Testing int")
	variable := 5

	buffer, err := Serialize(variable, CLIENT)

	if err != nil {
		log.Debug("Serialized failed", "err", err)
	} else {
		log.Debug("buffer", "buffer", buffer)
	}

	var integer int

	// result is probably not an integer
	result, err := Deserialize(buffer, integer, CLIENT)

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

func WireRegister() {

	// Tell wire how to find the underlying types for a given interface
	var _ = wire.RegisterInterface(
		struct{ BasicType }{},
		wire.ConcreteType{Basic{}, typeBasic},
	)
}

// Basic structural conversion
func XTestStruct(t *testing.T) {
	log.Info("Testing Struct")

	WireRegister()

	basic := &Basic{Pad: "xxxx", Number: 123456, Name: "A Name"}
	log.Debug("The basic type", "basic", basic)

	buffer, err := Serialize(basic, CLIENT)

	if err != nil {
		assert.FailNow(t, "Serialization failed ", err.Error())
	} else {
		log.Debug("buffer", "buffer", buffer)
	}

	result := new(Basic)
	base, err := Deserialize(buffer, result, CLIENT)
	value := base.(*Basic)

	if err != nil {
		assert.Fail(t, "Deserialization failed", err.Error())
	} else {
		log.Debug("result", "value", value)
	}

	assert.Equal(t, basic, value, "These should be equal")
}

// Basic structural conversion
func XTestStruct2(t *testing.T) {
	log.Info("Testing Struct2")

	WireRegister()

	basic := &Basic{Pad: "xxxx", Number: 123456, Name: "A Name"}
	log.Debug("The basic type", "basic", basic)

	buffer, err := Serialize(basic, CLIENT)

	if err != nil {
		log.Debug("Serialized failed", "err", err)
	} else {
		log.Debug("buffer", "buffer", buffer)
	}

	result, err := Deserialize(buffer, struct{ *Basic }{}, CLIENT)
	value := result.(struct{ *Basic }).Basic

	if err != nil {
		log.Debug("Deserialized failed", "err", err)
	} else {
		log.Debug("result", "value", value)
	}

	assert.Equal(t, basic, value, "These should be equal")
}

func XTestInterface(t *testing.T) {
	log.Info("Testing Interface")

	WireRegister()

	var generic BasicType

	generic = &Basic{Pad: "-------", Number: 654321, Name: "A rose by any other name"}
	log.Debug("generic", "generic", generic)

	buffer, err := Serialize(generic, CLIENT)
	if err != nil {
		assert.FailNow(t, "Serialization failed", err.Error())
	} else {
		log.Debug("buffer", "buffer", buffer)
	}

	result, err := Deserialize(buffer, struct{ BasicType }{}, CLIENT)
	if err != nil {
		assert.Fail(t, "Deserialization failed", err.Error())
	} else {
		log.Debug("result", "result", result)
	}

	value := result.(struct{ BasicType }).BasicType
	log.Debug("value", "value", value)
	assert.Equal(t, generic, value, "These should be equal")
}
