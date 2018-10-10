/*
	Copyright 2017-2018 OneLedger

	We want encapsulated conversions for the three different pathways:

		- transactions sent from Tendermint (is this a mix between wire and JSON?)
		- data stored in LevelDB
		- queries coming in from http

*/
package serial

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
	JSON
)

type Message = []byte

// Given any type of input (except Maps), convert it into wire format
func Serialize(input interface{}, medium Format) (buffer []byte, err error) {

	var copy interface{}

	if medium == JSON {
		copy = input

	} else if IsPrimitive(input) {
		// Manually wrap the primitive
		typeof := reflect.TypeOf(input).Name()
		dict := make(map[string]interface{})
		dict[""] = input
		copy = SerialWrapper{Type: typeof, Fields: dict}

	} else {
		// Expand all structs with wrappers
		copy = Extend(input)
	}

	switch medium {

	case PERSISTENT:
		buffer, err = json.Marshal(copy)

	case NETWORK:
		buffer, err = json.Marshal(copy)

	case CLIENT:
		buffer, err = json.Marshal(copy)

	case JSON:
		buffer, err = json.Marshal(copy)
	}

	log.Dump("buffer", string(buffer), "err", err)
	return buffer, err
}

// Given a serialized slice, put it back into the original struct format
func Deserialize(input []byte, output interface{}, medium Format) (msg interface{}, err error) {

	log.Dump("Deserialize the string", string(input))

	wrapper := &(map[string]interface{}{})

	switch medium {

	case PERSISTENT:
		err = json.Unmarshal(input, wrapper)

	case NETWORK:
		err = json.Unmarshal(input, wrapper)

	case CLIENT:
		err = json.Unmarshal(input, wrapper)

	case JSON:
		err = json.Unmarshal(input, output)
		if err == nil {
			return output, err
		}
	}

	if err != nil {
		return nil, err
	}

	result := Contract(wrapper)

	log.Dump("final result", result)
	log.Dump("original wrapper", wrapper)

	return result, err
}

type SerialWrapper struct {
	Type   string
	Fields map[string]interface{}
}

var prototype = SerialWrapper{}

// Test to see if this a SerialWrapper struct
func IsSerialWrapper(input interface{}) bool {
	if reflect.TypeOf(input) == reflect.TypeOf(prototype) {
		return true
	}
	return false
}

// Can identify a map created from a SerialWrapper, explicitly depends on the SerialWrapper type
func IsSerialWrapperMap(input interface{}) bool {
	if input == nil {
		return false
	}

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
		log.Dump("structures", structures)
		log.Fatal("Missing structure type", "name", name)
		return nil
	}

	base := reflect.New(struct_type)
	return base.Interface()
}
