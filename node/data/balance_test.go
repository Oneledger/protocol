package data

import (
	"fmt"
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestNewBalance(t *testing.T) {

	a := NewBalanceFromString(100, "VT")
	b := NewBalanceFromString(241351, "OLT")
	fmt.Println(a)
	fmt.Println(b)

	a.AddAmmount(NewCoin(10, "OLT"))
	fmt.Println(a)

	amounts := make([]Coin, len(Currencies))
	for _, v := range Currencies {
		if v.Name == "OLT" {
			amounts[v.Id] = NewCoin(10, "OLT")
		} else {
			amounts[v.Id] = NewCoin(0, v.Name)
		}
	}

	for _, v := range amounts {
		fmt.Println(v)
	}
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
