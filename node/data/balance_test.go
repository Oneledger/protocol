package data

import (
	"fmt"
	"testing"

	"github.com/Oneledger/protocol/node/log"
	"github.com/Oneledger/protocol/node/serial"
	"github.com/magiconair/properties/assert"
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

	assert.Equal(t, a.Amounts[3].Amount.Uint64()+10, b.Amounts[3].Amount.Uint64())

}

func TestBalance_MinusAmmount(t *testing.T) {
	a := NewBalanceFromInt(100, "VT")
	b := NewBalanceFromInt(100, "VT")

	b.MinusAmount(NewCoinFromInt(20, "VT"))
	assert.Equal(t, a.Amounts[3].Amount.Uint64()-20, b.Amounts[3].Amount.Uint64())
}
