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
	"testing"

	"github.com/Oneledger/protocol/serialize"
	"github.com/stretchr/testify/assert"
)

func TestCoin(t *testing.T) {

	curr := Currency{1, "OLT", 0, 18, "onz"}

	coin := curr.NewCoinFromInt(2)

	str, err := serialize.GetSerializer(serialize.NETWORK).Serialize(&coin)
	assert.NoError(t, err)
	assert.NotEqual(t, len(str), 0)

	coinNew := Coin{}
	err = serialize.GetSerializer(serialize.NETWORK).Deserialize(str, &coinNew)
	assert.NoError(t, err, "serialization error")

	notOk := coin.LessThanCoin(coinNew)
	assert.False(t, notOk)

	ok := coin.Equals(coinNew)
	assert.True(t, ok)

	ok = coin.IsValid()
	assert.True(t, ok)

	doubleCoin := coin.Plus(coinNew)

	ok = coin.LessThanCoin(doubleCoin)
	assert.True(t, ok)

	coinNew = coinNew.Divide(2)
	ok = coinNew.LessThanCoin(coin)
	assert.True(t, ok)

	coinNew = coinNew.MultiplyInt(4)
	ok = coinNew.Equals(doubleCoin)
	assert.True(t, ok)

	ok = coin.IsCurrency("OLT")
	assert.True(t, ok)

	notOk = coin.IsCurrency("Bitcoin")
	assert.False(t, notOk)

	ok = coin.IsCurrency("Ethereum", "OLT", "Bitcoin")
	assert.True(t, ok)
}

func TestCoin_LessThanEqualCoin(t *testing.T) {

	curr := Currency{1, "OLT", 0, 18, "onz"}

	coin := curr.NewCoinFromInt(1)
	coin2 := curr.NewCoinFromFloat64(0.9)

	ok := coin2.LessThanEqualCoin(coin)
	assert.True(t, ok)

	ok = coin2.LessThanEqualCoin(coin2)
	assert.True(t, ok)
}
