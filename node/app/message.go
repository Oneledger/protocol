/*
	Copyright 2017-2018 OneLedger

	TODO: We want configurable, switchable conversions for the different pathways
		- transactions sent from Tendermint (is this a mix between wire and JSON?)
		- data stored in LevelDB
		- queries coming in from http
*/
package app

import (
	"bytes"
	"encoding/json"
	"reflect"

	wire "github.com/tendermint/go-wire"
)

// Given any type of input (except Maps), convert it into wire format
func Serialize(input interface{}) (msg []byte, err error) {
	var count int

	buffer := new(bytes.Buffer)

	wire.WriteBinary(input, buffer, &count, &err)

	return buffer.Bytes(), err
}

// Given something in wire format, stick it back into the original golang types
func Deserialize(input []byte, output interface{}) (msg interface{}, err error) {
	var count int

	buffer := bytes.NewBuffer(input)

	valueOf := reflect.ValueOf(output)
	if valueOf.Kind() == reflect.Ptr {
		msg = wire.ReadBinaryPtr(output, buffer, len(input), &count, &err)
		//msg = wire.ReadBinaryPtr(output, buffer, 0, &count, &err)
	} else {
		msg = wire.ReadBinary(output, buffer, len(input), &count, &err)
		//msg = wire.ReadBinary(output, buffer, 0, &count, &err)
	}
	return msg, err
}

// Go's version of JSON
func ConvertToJSON(input interface{}) (msg []byte, err error) {
	bytes, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

// Go's version of JSON
func ConvertFromJSON(input []byte, output interface{}) (err error) {
	err = json.Unmarshal(input, output)
	return err
}

// Convert into wire's version of JSON (which is still non-standard?)
func ConvertToWireJSON(input interface{}) (msg []byte, err error) {
	var count int

	buffer := new(bytes.Buffer)

	wire.WriteJSON(input, buffer, &count, &err)

	return buffer.Bytes(), err
}

// Convert from wire's JSON format back into the original golang type
func ConvertFromWireJSON(input []byte, output interface{}) (err error) {

	valueOf := reflect.ValueOf(output)

	if valueOf.Kind() == reflect.Ptr {
		wire.ReadJSONPtr(output, input, &err)
	} else {
		wire.ReadJSON(output, input, &err)
	}
	return err
}
