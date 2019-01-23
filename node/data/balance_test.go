package data

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Oneledger/protocol/node/log"
	"github.com/Oneledger/protocol/node/serial"
)

func xTestNewBalance(t *testing.T) {

	a := NewBalanceFromInt(100, "VT")
	b := NewBalanceFromInt(241351, "OLT")
	fmt.Println(a)
	fmt.Println(b)

	a.AddAmount(NewCoinFromInt(10, "OLT"))
	fmt.Println(a)

}

func TestSerialize(t *testing.T) {
	a := NewBalanceFromString("10", "OLT")
	buffer, err := serial.Serialize(a, serial.PERSISTENT)
	if err != nil {
		log.Fatal("Serialization Failed", "err", err)
	}
	var proto Balance

	result, err := serial.DumpDeserialize(buffer, proto, serial.PERSISTENT)
	if err != nil {
		log.Fatal("Deserialized failed", "err", err)
	}
	assert.Equal(t, result, a, "These should be equal")
}

func TestBalance_AddAmmount(t *testing.T) {

	a := NewBalanceFromInt(100, "VT")
	b := NewBalanceFromInt(100, "VT")
	b.AddAmount(NewCoinFromInt(10, "VT"))

	assert.Equal(t, a.GetAmountByName("VT").Amount.Uint64()+10, b.GetAmountByName("VT").Amount.Uint64())

}

func TestBalance_MinusAmmount(t *testing.T) {
	a := NewBalanceFromInt(100, "VT")
	b := NewBalanceFromInt(100, "VT")

	b.MinusAmount(NewCoinFromInt(20, "VT"))
	//assert.Equal(t, a.Amounts[3].Amount.Uint64()-20, b.Amounts[3].Amount.Uint64())
	assert.Equal(t, a.GetAmountByName("VT").Amount.Uint64()-20, b.GetAmountByName("VT").Amount.Uint64())
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
