/*
	Copyright 2017-2018 OneLedger
*/
package app

import (
	"bytes"
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

// Convert into wire's version of JSON (which is still non-standard?)
func ConvertToJSON(input interface{}) (msg []byte, err error) {
	var count int

	buffer := new(bytes.Buffer)

	wire.WriteJSON(input, buffer, &count, &err)

	return buffer.Bytes(), err
}

// Convert from wire's JSON format back into the original golang type
func ConvertFromJSON(input []byte, output interface{}) (err error) {

	valueOf := reflect.ValueOf(output)

	if valueOf.Kind() == reflect.Ptr {
		wire.ReadJSONPtr(output, input, &err)
	} else {
		wire.ReadJSON(output, input, &err)
	}
	return err
}
