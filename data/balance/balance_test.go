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
	"os"
	"testing"

	"github.com/Oneledger/protocol/serialize"
	"github.com/stretchr/testify/assert"
	logger "github.com/Oneledger/protocol/log"
)

func init() {
	pSzlr = serialize.GetSerializer(serialize.PERSISTENT)
}

func xTestNewBalance(t *testing.T) {

	a := NewBalanceFromInt(100, "VT")
	b := NewBalanceFromInt(241351, "OLT")
	fmt.Println(a)
	fmt.Println(b)

	a.AddAmount(NewCoinFromInt(10, "OLT"))
	fmt.Println(a)

}

func TestSerialize(t *testing.T) {

	// temp; remove this logger later
	var log = logger.NewDefaultLogger(os.Stdout)

	a := NewBalanceFromString("10", "OLT")
	buffer, err := pSzlr.Serialize(a)
	if err != nil {
		logger.Fatal("Serialization Failed", "err", err)
	}
	var result = &Balance{}

	err = pSzlr.Deserialize(buffer, result)
	if err != nil {
		logger.Fatal("Deserialized failed", "err", err)
	}
	assert.Equal(t, result, a, "These should be equal")
}

func TestBalance_AddAmmount(t *testing.T) {

	a := NewBalanceFromInt(100, "VT")
	b := NewBalanceFromInt(100, "VT")
	b.AddAmount(NewCoinFromInt(10, "VT"))

	assert.Equal(t, a.GetCoinByName("VT").Amount.Uint64()+10, b.GetCoinByName("VT").Amount.Uint64())

}

func TestBalance_MinusAmmount(t *testing.T) {
	a := NewBalanceFromInt(100, "VT")
	b := NewBalanceFromInt(100, "VT")

	b.MinusAmount(NewCoinFromInt(20, "VT"))
	//assert.Equal(t, a.Amounts[3].Amount.Uint64()-20, b.Amounts[3].Amount.Uint64())
	assert.Equal(t, a.GetCoinByName("VT").Amount.Uint64()-20, b.GetCoinByName("VT").Amount.Uint64())
}

func TestBalance_IsEnoughBalance(t *testing.T) {

	a := NewBalance()
	b := NewBalance()

	a.AddAmount(NewCoinFromInt(100, "OLT"))
	b.AddAmount(NewCoinFromInt(42, "VT"))
	//fmt.Println("init", a, b)
	result1 := a.IsEnoughBalance(*b)
	assert.Equal(t, false, result1, a.String(), b.String())
	//fmt.Println("first", a, b)
	a.AddAmount(NewCoinFromInt(43, "VT"))
	result2 := a.IsEnoughBalance(*b)
	assert.Equal(t, true, result2, a.String(), b.String())
	//fmt.Println("second", a, b)
	b.AddAmount(NewCoinFromInt(101, "OLT"))
	result3 := a.IsEnoughBalance(*b)
	assert.Equal(t, false, result3, a.String(), b.String())
	//fmt.Println("third", a, b)
}
