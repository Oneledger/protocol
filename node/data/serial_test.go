package data

import (
	"testing"

	"github.com/Oneledger/protocol/node/log"
	"github.com/stretchr/testify/assert"
)

func TestCoin(t *testing.T) {
	var coin Coin

	coin = NewCoinFromInt(100000, "BTC")

	// Serialize the go data structure
	buffer, err := pSzlr.Serialize(coin)

	if err != nil {
		log.Fatal("Serialized failed", "err", err)
	}

	var result = &Coin{}

	// Deserialize back into a go data structure
	err = pSzlr.Deserialize(buffer, result)
	if err != nil {
		log.Fatal("Deserialized failed", "err", err)
	}

	assert.Equal(t, coin, result, "These should be equal")
}

func TestBalance(t *testing.T) {
	var balance Balance

	// Serialize the go data structure
	buffer, err := pSzlr.Serialize(balance)

	if err != nil {
		log.Fatal("Serialized failed", "err", err)
	}

	var result = &Balance{}

	// Deserialize back into a go data structure
	err = pSzlr.Deserialize(buffer, result)
	if err != nil {
		log.Fatal("Deserialized failed", "err", err)
	}

	assert.Equal(t, balance, result, "These should be equal")
}
