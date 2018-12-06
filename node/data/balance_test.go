package data

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewBalance(t *testing.T) {

	a := NewBalanceFromString(100, "VT")
	b := NewBalanceFromString(120, "OLT")
	fmt.Println(a)
	fmt.Println(b)

	a.AddAmount(NewCoin(10, "OLT"))
	fmt.Println(a)

}

func TestBalance_AddAmmount(t *testing.T) {

	a := NewBalanceFromString(100, "VT")
	b := NewBalanceFromString(100, "VT")
	b.AddAmount(NewCoin(10, "VT"))

	assert.Equal(t, a.GetAmountByName("VT").Amount.Uint64()+10, b.GetAmountByName("VT").Amount.Uint64())

}

func TestBalance_MinusAmmount(t *testing.T) {
	a := NewBalanceFromString(100, "VT")
	b := NewBalanceFromString(100, "VT")

	b.MinusAmount(NewCoin(20, "VT"))
	assert.Equal(t, a.GetAmountByName("VT").Amount.Uint64()-20, b.GetAmountByName("VT").Amount.Uint64())
}

func TestBalance_IsEnoughBalance(t *testing.T) {

	a := NewBalance()
	b := NewBalance()

	a.AddAmount(NewCoin(100, "OLT"))
	b.AddAmount(NewCoin(42, "VT"))
	fmt.Println("init", a, b)
	result1 := a.IsEnoughBalance(b)
	assert.Equal(t, false, result1, a.String(), b.String())
	fmt.Println("first", a, b)
	a.AddAmount(NewCoin(43, "VT"))
	result2 := a.IsEnoughBalance(b)
	assert.Equal(t, true, result2, a.String(), b.String())
	fmt.Println("second", a, b)
	b.AddAmount(NewCoin(101, "OLT"))
	result3 := a.IsEnoughBalance(b)
	assert.Equal(t, false, result3, a.String(), b.String())
	fmt.Println("third", a, b)
}
