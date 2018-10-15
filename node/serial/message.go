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

	//log.Dump("Serialized format", string(buffer))

	return buffer, err
}

// Given a serialized slice, put it back into the original struct format
func Deserialize(input []byte, output interface{}, medium Format) (msg interface{}, err error) {

	//log.Dump("Deserialize the string", string(input))

	//wrapper := &(map[string]interface{}{})
	wrapper := &SerialWrapper{}

	switch medium {

	case PERSISTENT:
		err = json.Unmarshal(input, wrapper)

	case NETWORK:
		err = json.Unmarshal(input, wrapper)

	case CLIENT:
		err = json.Unmarshal(input, wrapper)

	case JSON:
		err = json.Unmarshal(input, output)

		// Exit before trying to contract
		if err == nil {
			//log.Dump("JSON Deserialized to", output)
			return output, err
		}
	}

	if err != nil {
		return nil, err
	}

	result := Contract(wrapper)
	//log.Dump("Deserialized to", result)

	return result, err
}

type SerialWrapper struct {
	Type   string
	Fields map[string]interface{}
	Size   int
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

	if len(smap) != 3 {
		return false
	}

	if _, ok := smap["Type"]; !ok {
		return false
	}

	if _, ok := smap["Fields"]; !ok {
		return false
	}

	if _, ok := smap["Size"]; !ok {
		return false
	}
	return true
}
