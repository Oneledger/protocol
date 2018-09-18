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

	switch medium {

	case PERSISTENT:
		buffer, err = json.Marshal(Wrap(input))

	case NETWORK:
		buffer, err = json.Marshal(Wrap(input))

	case CLIENT:
		buffer, err = json.Marshal(Wrap(input))
	}

	return buffer, err
}

// Given something in wire format, stick it back into the original golang types
// If output is a struct, make sure it is a pointer to a struct
func Deserialize(input []byte, output interface{}, medium Format) (msg interface{}, err error) {

	//var raw json.RawMessage
	wrapper := &SerialWrapper{
		Value: output,
	}

	switch medium {

	case PERSISTENT:
		err = json.Unmarshal(input, wrapper)

	case NETWORK:
		err = json.Unmarshal(input, wrapper)

	case CLIENT:
		err = json.Unmarshal(input, wrapper)
	}

	log.Debug("Unmarshal", "err", err, "wrapper", wrapper)

	return wrapper.Value, err
}

type SerialWrapper struct {
	Type      string
	IsPointer bool
	Value     interface{}
}

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

var structures map[string]reflect.Type

func Register(base interface{}) {
	if structures == nil {
		structures = make(map[string]reflect.Type)
	}
	structures[reflect.TypeOf(base).String()] = reflect.TypeOf(base)
}

func NewStruct(name string) interface{} {
	struct_type := structures[name]
	if struct_type == nil {
		log.Debug("Missing structure type", "name", name, "structures", structures)
		return nil
	}
	return reflect.Zero(struct_type)
}

func (w SerialWrapper) MarshalJSON() (value []byte, err error) {
	log.Debug("MarshalJSON for SerialWrapper")
	buffer, err := json.Marshal(w.Value)

	wrapper := PartialWrapper{
		Type:      w.Type,
		IsPointer: w.IsPointer,
		Value:     string(buffer),
	}

	return json.Marshal(wrapper)
}

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

// Wrap an interface and return
func Wrap(input interface{}) interface{} {
	log.Debug("Type is", "type", reflect.TypeOf(input))

	return SerialWrapper{
		Type:      reflect.TypeOf(input).String(),
		IsPointer: false,
		Value:     input,
	}
}
