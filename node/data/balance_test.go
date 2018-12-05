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

func TestBalance_IsEqual(t *testing.T) {

	a := NewBalance()
	b := NewBalance()

	a.AddAmmount(NewCoin(100, "OLT"))
	a.AddAmmount(NewCoin(42, "VT"))
	result1 := a.IsEqual(b)
	assert.Equal(t, result1, false)

	a.MinusAmmount(NewCoin(100, "OLT"))
	a.MinusAmmount(NewCoin(42, "VT"))
	result2 := a.IsEqual(b)
	assert.Equal(t, result2, true)
}

func TestBalance_IsEnough(t *testing.T) {
	a := NewBalance()

	a.AddAmmount(NewCoin(100, "OLT"))

	c := NewCoin(101, "OLT")
	d := NewCoin(99, "OLT")

	result1 := a.IsEnough(c)
	result2 := a.IsEnough(d)
	assert.Equal(t, result1, false)

	assert.Equal(t, result2, true)
}
