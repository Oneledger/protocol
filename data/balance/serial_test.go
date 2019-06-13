/*
   ____             _              _                      _____           _                  _
  / __ \           | |            | |                    |  __ \         | |                | |
 | |  | |_ __   ___| |     ___  __| | __ _  ___ _ __     | |__) | __ ___ | |_ ___   ___ ___ | |
 | |  | | '_ \ / _ \ |    / _ \/ _` |/ _` |/ _ \ '__|    |  ___/ '__/ _ \| __/ _ \ / __/ _ \| |
 | |__| | | | |  __/ |___|  __/ (_| | (_| |  __/ |       | |   | | | (_) | || (_) | (_| (_) | |
  \____/|_| |_|\___|______\___|\__,_|\__, |\___|_|       |_|   |_|  \___/ \__\___/ \___\___/|_|
                                      __/ |
                                     |___/


Copyright 2017 - 2019 OneLedger
*/

package balance

import (
	//logger "github.com/Oneledger/protocol/log"
	"github.com/stretchr/testify/assert"
	"testing"
)

//
//var log = logger.NewDefaultLogger(os.Stdout)

/*
func TestCoin(t *testing.T) {
	// temp; remove this logger later
	var log = logger.NewDefaultLogger(os.Stdout)
	var coin Coin

	coin = NewCoinFromInt(100000, "BTC")

	// Serialize the go data structure
	buffer, err := pSzlr.Serialize(&coin)

	if err != nil {
		log.Fatal("Serialized failed", "err", err)
	}

	var result = &Coin{}

	// Deserialize back into a go data structure
	err = pSzlr.Deserialize(buffer, result)
	if err != nil {
		log.Fatal("Deserialized failed", "err", err)
	}

	assert.Equal(t, &coin, result, "These should be equal")
}

*/

func TestBalance(t *testing.T) {
	// temp; remove this logger later

	var balance = NewBalance()

	// Serialize the go data structure
	buffer, err := pSzlr.Serialize(balance)

	if err != nil {
		logger.Fatal("Serialized failed", "err", err)
	}

	var result = &Balance{}

	// Deserialize back into a go data structure
	err = pSzlr.Deserialize(buffer, result)
	if err != nil {
		logger.Fatal("Deserialized failed", "err", err)
	}

	assert.Equal(t, balance, result, "These should be equal")
}
