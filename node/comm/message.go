/*
	Copyright 2017-2018 OneLedger

	TODO: We want configurable, switchable conversions for the different pathways
		- transactions sent from Tendermint (is this a mix between wire and JSON?)
		- data stored in LevelDB
		- queries coming in from http
*/
package comm

import (
	"encoding/json"

	"github.com/Oneledger/protocol/node/log"
)

type Message = []byte

// Given any type of input (except Maps), convert it into wire format
func Serialize(input interface{}) (msg []byte, err error) {
	buffer, err := json.Marshal(input)
	if err != nil {
		log.Error("Marshal", "err", err)
	}

	return buffer, err

	// TODO: Do we really need wire here?
	/*
		var count int

		buffer := new(bytes.Buffer)

		wire.WriteBinary(input, buffer, &count, &err)

		return buffer.Bytes(), err
	*/
}

// Given something in wire format, stick it back into the original golang types
// If output is a struct, make sure it is a pointer to a struct
func Deserialize(input []byte, output interface{}) (msg interface{}, err error) {

	err = json.Unmarshal(input, output)
	if err != nil {
		log.Error("Unmarshal", "err", err)
	}

	return output, err

	// TODO: Do we really need wire here?
	/*
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
	*/
}
