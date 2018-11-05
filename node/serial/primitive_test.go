/*
	Copywrite 2017-2018 OneLedger

	Tests to validate the way wire handles various struct conversions.
*/
package serial

import (
	"testing"

	"github.com/Oneledger/protocol/node/log"
	"github.com/stretchr/testify/assert"
)

// TODO: This should really work...
// Test just an integer
func TestInt(t *testing.T) {
	log.Info("Testing int")
	variable := 5

	buffer, err := Serialize(variable, CLIENT)

	if err != nil {
		log.Fatal("Serialized failed", "err", err)
	}

	var prototype int

	// result is probably not an integer
	result, err := DumpDeserialize(buffer, prototype, CLIENT)

	if err != nil {
		log.Fatal("Deserialized failed", "err", err)
	}

	assert.Equal(t, variable, result, "These should be equal")
}

// TODO: This should really work...
// Test just an integer
func TestString(t *testing.T) {
	log.Info("Testing string")
	variable := "Hello World"

	buffer, err := Serialize(variable, CLIENT)

	if err != nil {
		log.Fatal("Serialized failed", "err", err)
	}

	var prototype string

	// result is probably not an integer
	result, err := DumpDeserialize(buffer, prototype, CLIENT)

	if err != nil {
		log.Fatal("Deserialized failed", "err", err)
	}

	assert.Equal(t, variable, result, "These should be equal")
}
