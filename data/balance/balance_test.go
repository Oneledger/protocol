/*
   ____             _          _
  / __ \           | |        | |
 | |  | |_ __   ___| | ___  __| | __ _  ___ _ __
 | |  | | '_ \ / _ \ |/ _ \/ _` |/ _` |/ _ \ '__|
 | |__| | | | |  __/ |  __/ (_| | (_| |  __/ |
  \____/|_| |_|\___|_|\___|\__,_|\__, |\___|_|
                                  __/ |
                                 |___/

Copyright 2017 - 2019 OneLedger
*/

package balance

import (
	"fmt"
	"github.com/Oneledger/protocol/data/chain"
	"testing"

	"github.com/Oneledger/protocol/serialize"
	"github.com/stretchr/testify/assert"
)

var currencies = make(map[string]Currency)

func init() {
	pSzlr = serialize.GetSerializer(serialize.PERSISTENT)
	currencies["OLT"] = Currency{"OLT", chain.Type(0), 18}
	currencies["VT"] = Currency{"VT", chain.Type(1), 0}
}

func xTestNewBalance(t *testing.T) {
	olt := currencies["OLT"]
	vt := currencies["VT"]
	a := NewBalance().AddCoin(olt.NewCoinFromFloat64(2623.232))
	b := vt.NewCoinFromInt(100)
	fmt.Println(a)
	fmt.Println(b)

	a.AddCoin(b)
	fmt.Println(a)

}

func TestSerialize(t *testing.T) {
	olt := currencies["OLT"]
	a := NewBalance().AddCoin(olt.NewCoinFromFloat64(2623.232))
	buffer, err := pSzlr.Serialize(a)
	assert.Nil(t, err)

	var result = &Balance{}

	err = pSzlr.Deserialize(buffer, result)
	assert.Nil(t, err)
	assert.Equal(t, a, result, "These should be equal")
}

func TestBalance_AddCoin(t *testing.T) {

	olt := currencies["OLT"]
	a := NewBalance().AddCoin(olt.NewCoinFromFloat64(2623.232))
	b := olt.NewCoinFromInt(100)
	expect := a.GetCoin(olt).Amount.Uint64() + 100

	a.AddCoin(b)

	assert.Equal(t, expect, a.GetCoin(olt).Amount.Uint64())

}

func TestBalance_MinusCoin(t *testing.T) {
	olt := currencies["OLT"]
	a := NewBalance().AddCoin(olt.NewCoinFromFloat64(2623.232))
	b := olt.NewCoinFromInt(100)
	expect := a.GetCoin(olt).Amount.Uint64() - 100
	a.MinusCoin(b)

	assert.Equal(t, expect, a.GetCoin(olt).Amount.Uint64())
}

func TestBalance_IsEnoughBalance(t *testing.T) {
	olt := currencies["OLT"]
	vt := currencies["VT"]

	a := NewBalance()
	b := NewBalance()

	a.AddCoin(olt.NewCoinFromInt(100))
	b.AddCoin(vt.NewCoinFromInt(42))
	//fmt.Println("init", a, b)
	result1 := a.IsEnoughBalance(*b)
	assert.Equal(t, false, result1, a.String(), b.String())
	//fmt.Println("first", a, b)
	a.AddCoin(vt.NewCoinFromFloat64(43))
	result2 := a.IsEnoughBalance(*b)
	assert.Equal(t, true, result2, a.String(), b.String())
	//fmt.Println("second", a, b)
	b.AddCoin(olt.NewCoinFromInt(101))
	result3 := a.IsEnoughBalance(*b)
	assert.Equal(t, false, result3, a.String(), b.String())
	//fmt.Println("third", a, b)
}
