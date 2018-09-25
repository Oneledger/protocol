/*
	Copyright 2017-2018 OneLedger

	We want encapsulated conversions for the three different pathways:

		- transactions sent from Tendermint (is this a mix between wire and JSON?)
		- data stored in LevelDB
		- queries coming in from http

*/
package comm

import (
	"encoding/json"
	"reflect"
	"strings"

	"github.com/Oneledger/protocol/node/log"
)

type Format int

const (
	PERSISTENT Format = iota
	NETWORK
	CLIENT
)

type Message = []byte

// Given any type of input (except Maps), convert it into wire format
func Serialize(input interface{}, medium Format) (buffer []byte, err error) {
	copy := Extend(input)
	log.Debug("Extended To Serialize", "input", input, "copy", copy)

	switch medium {

	case PERSISTENT:
		buffer, err = json.Marshal(copy)

	case NETWORK:
		buffer, err = json.Marshal(copy)

	case CLIENT:
		buffer, err = json.Marshal(copy)
	}

	log.Debug("Serialized", "buffer", buffer, "err", err)
	return buffer, err
}

// Given something in wire format, stick it back into the original golang types
// If output is a struct, make sure it is a pointer to a struct
func Deserialize(input []byte, output interface{}, medium Format) (msg interface{}, err error) {

	//var raw json.RawMessage
	wrapper := &SerialWrapper{}

	switch medium {

	case PERSISTENT:
		err = json.Unmarshal(input, wrapper)

	case NETWORK:
		err = json.Unmarshal(input, wrapper)

	case CLIENT:
		err = json.Unmarshal(input, wrapper)
	}

	if err != nil {
		return nil, err
	}

	result := Contract(wrapper)
	log.Debug("Contracted for Deserialize", "result", result, "wrapper", wrapper)

	return result, err
}

type SerialWrapper struct {
	Type   string
	Fields map[string]interface{}
}

var prototype = SerialWrapper{}

func IsSerialWrapper(input interface{}) bool {
	if reflect.TypeOf(input) == reflect.TypeOf(prototype) {
		return true
	}
	return false
}

func IsSerialWrapperMap(input interface{}) bool {
	if reflect.TypeOf(input).Kind() != reflect.Map {
		return false
	}

	smap := input.(map[string]interface{})

	if smap == nil {
		return false
	}

	if len(smap) != 2 {
		return false
	}

	if _, ok := smap["Type"]; !ok {
		return false
	}

	if _, ok := smap["Fields"]; !ok {
		return false
	}

	return true
}

/*
type PartialWrapper struct {
	Type      string
	IsPointer bool
	Value     string
}

type SerialUnwrapper struct {
	Type      string
	IsPointer bool
	Value     interface{}
}
*/

var structures map[string]reflect.Type

// Register a structure by its name
func Register(base interface{}) {

	// Allocate on the first call
	if structures == nil {
		structures = make(map[string]reflect.Type)
	}

	structures[reflect.TypeOf(base).String()] = reflect.TypeOf(base)
}

// Dynamically create a structure from its name
func NewStruct(name string) interface{} {
	name = strings.TrimPrefix(name, "*")

	struct_type := structures[name]
	if struct_type == nil {
		log.Warn("Missing structure type", "name", name, "structures", structures)
		return nil
	}
	base := reflect.New(struct_type)
	log.Debug("Create A Dynamic Struct", "name", name, "base", base.Interface())

	return base.Interface()
}

/*
func (w SerialWrapper) MarshalJSON() (value []byte, err error) {
	log.Debug("MarshalJSON for SerialWrapper")
	//buffer, err := json.Marshal(w.Value)

		wrapper := PartialWrapper{
			Type:      w.Type,
			IsPointer: w.IsPointer,
			Value:     string(buffer),
		}

	return json.Marshal(value)
}
*/

/*
func (w *SerialWrapper) UnmarshalJSON(value []byte) (err error) {
	log.Debug("UnmarshalJSON for SerialWrapper")

	var raw json.RawMessage
	wrapper := &SerialUnwrapper{Value: &raw}

	// Unmarshal the top-level wrapper
	status := json.Unmarshal(value, wrapper)
	if status != nil {
		return status
	}
	proto := NewStruct(wrapper.Type)

	// Now, properly unmarshal the interface
	status = json.Unmarshal(raw, &proto)

	w.Type = wrapper.Type
	w.IsPointer = wrapper.IsPointer
	w.Value = proto

	return status
}
*/

// Wrap an interface and return
/*
func Wrap(input interface{}) interface{} {
	log.Debug("Type is", "type", reflect.TypeOf(input))

	return SerialWrapper{
		Type:      reflect.TypeOf(input).String(),
		IsPointer: false,
		Value:     input,
	}
}
*/
