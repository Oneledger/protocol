package data

import (
	"fmt"
	"github.com/magiconair/properties/assert"
	"testing"
)

func xTestNewBalance(t *testing.T) {

	a := NewBalanceFromString(100, "VT")
	b := NewBalanceFromString(241351, "OLT")
	fmt.Println(a)
	fmt.Println(b)

	a.AddAmmount(NewCoin(10, "OLT"))
	fmt.Println(a)

}

func TestBalance_AddAmmount(t *testing.T) {

	a := NewBalanceFromString(100, "VT")
	b := NewBalanceFromString(100, "VT")
	b.AddAmmount(NewCoin(10, "VT"))

	assert.Equal(t, a.Amounts[3].Amount.Uint64()+10, b.Amounts[3].Amount.Uint64())

}

func TestBalance_MinusAmmount(t *testing.T) {
	a := NewBalanceFromString(100, "VT")
	b := NewBalanceFromString(100, "VT")

	b.MinusAmmount(NewCoin(20, "VT"))
	assert.Equal(t, a.Amounts[3].Amount.Uint64()-20, b.Amounts[3].Amount.Uint64())
}
