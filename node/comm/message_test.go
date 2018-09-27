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

func XTestStack(t *testing.T) {
	var stack Stack = *NewStack()

	log.Debug("Current", "stack", stack)

	stack.Push("Bottom")
	stack.Push("Middle")
	stack.Push("Top")

	stack.Print()
	log.Debug("Top", "is", stack.Pop())
	log.Debug("Middle", "is", stack.Pop())
	log.Debug("Bottom", "is", stack.Pop())
}

type OpaqueRoot interface {
	AFunction(value int) bool
}

type OpaqueNode interface {
	AnotherFunction(value int) bool
}

type RootStruct1 struct {
	RootName string
	Node     OpaqueNode
}

func (p RootStruct1) AFunction(value int) bool {
	return false
}

type RootStruct2 struct {
	RootName string
	Count    int
	Node1    OpaqueNode
	Node2    OpaqueNode
}

func (p RootStruct2) AFunction(value int) bool {
	return false
}

type NodeStruct1 struct {
	NodeName string
	Size     int
}

func (c NodeStruct1) AnotherFunction(value int) bool {
	return false
}

type NodeStruct2 struct {
	NodeName string
	Size     int
	Size2    int
}

func (c NodeStruct2) AnotherFunction(value int) bool {
	return false
}

func init() {
	// Let the deserilization code know how to create these structures
	Register(RootStruct1{})
	Register(RootStruct2{})
	Register(NodeStruct1{})
	Register(NodeStruct2{})
}

var node1 = NodeStruct1{
	NodeName: "The Node 1",
	Size:     1000,
}

var node2 = NodeStruct2{
	NodeName: "The Child 2",
	Size:     3000,
	Size2:    -1,
}

var opc1 = OpaqueNode(&node1) // Pointer
var opc2 = OpaqueNode(node2)  // By Value

var root = RootStruct2{
	RootName: "The Root 2",
	Count:    2000,
	Node1:    opc1,
	Node2:    opc2,
}

var opp = OpaqueRoot(&root)

func XTestPrint(t *testing.T) {
	Print(opp)
}

func XTestClone(t *testing.T) {
	opp2 := Clone(opp)
	assert.Equal(t, opp, opp2, "These should be equal")
}

func XTestExtend(t *testing.T) {
	opp2 := Extend(opp)
	log.Debug("Extended is", "opp2", opp2, "opp", opp)

	assert.NotEqual(t, opp, opp2, "These should not be equal")
}

func XTestSerialize(t *testing.T) {

	buffer, err := Serialize(opp, PERSISTENT)
	if err != nil {
		log.Debug("Serialized failed", "err", err)
	} else {
		log.Debug("Serialized Worked, return is", "buffer", buffer)
	}
}

func TestPolymorphism(t *testing.T) {

	buffer, err := Serialize(opp, PERSISTENT)

	if err != nil {
		log.Debug("Serialized failed", "err", err)
	} else {
		log.Debug("Serialized Worked, return is", "buffer", buffer)
	}

	log.Debug("########### Serialized Worked, now Deserialize")

	var opp2 OpaqueRoot

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
