/*
	Copyright 2017-2018 OneLedger

	Convert bytes and interfaces into JSON
*/
package convert

import (
	"encoding/json"
)

// Go's version of JSON
func ToJSON(input interface{}) (msg []byte, err error) {
	bytes, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

// Go's version of JSON
func FromJSON(input []byte, output interface{}) (err error) {
	err = json.Unmarshal(input, output)
	return err
}

// Convert into wire's version of JSON (which is still non-standard?)
//func ToWireJSON(input interface{}) (msg []byte, status error) {
//	var count int
//
//	buffer := new(bytes.Buffer)
//
//	wire.WriteJSON(input, buffer, &count, &status)
//
//	return buffer.Bytes(), status
//}

// Convert from wire's JSON format back into the original golang type
//func FromWireJSON(input []byte, output interface{}) (status error) {
//
//	valueOf := reflect.ValueOf(output)
//
//	if valueOf.Kind() == reflect.Ptr {
//		wire.ReadJSONPtr(output, input, &status)
//	} else {
//		wire.ReadJSON(output, input, &status)
//	}
//	return status
//}
