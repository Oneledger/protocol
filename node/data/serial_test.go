package data

import (
	"testing"

	"github.com/Oneledger/protocol/node/log"
	"github.com/Oneledger/protocol/node/serial"
	"github.com/stretchr/testify/assert"
)

func TestCoin(t *testing.T) {
	var coin Coin

	coin = NewCoin(100000, "BTC")

	// Serialize the go data structure
	buffer, err := serial.Serialize(coin, serial.PERSISTENT)

	if err != nil {
		log.Fatal("Serialized failed", "err", err)
	}

	var opp2 Coin

	// Deserialize back into a go data structure
	result, err := serial.Deserialize(buffer, opp2, serial.PERSISTENT)

	if err != nil {
		log.Fatal("Deserialized failed", "err", err)
	}

	assert.Equal(t, coin, result, "These should be equal")
}

func TestBalance(t *testing.T) {
	var balance Balance

	// Serialize the go data structure
	buffer, err := serial.Serialize(balance, serial.PERSISTENT)

	if err != nil {
		log.Fatal("Serialized failed", "err", err)
	}

	var opp2 Balance

	// Deserialize back into a go data structure
	result, err := serial.Deserialize(buffer, opp2, serial.PERSISTENT)

	if err != nil {
		log.Fatal("Deserialized failed", "err", err)
	}

	assert.Equal(t, balance, result, "These should be equal")
}
