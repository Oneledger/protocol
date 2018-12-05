/*
	Copyright 2017-2018 OneLedger
*/

package data

import (
	"fmt"
	"github.com/Oneledger/protocol/node/serial"
	"math/big"
)

// Wrap the amount with owner information
type Balance struct {
	// Address id.Address
	Amounts []Coin
}

func init() {
	serial.Register(Balance{})
}

func NewBalance() Balance {
	amounts := make([]Coin, len(Currencies))
	for _, v := range Currencies {
		amounts[v.Id] = NewCoin(0, v.Name)
	}
	result := Balance{
		Amounts: amounts,
	}
	return result
}

func NewBalanceFromString(amount int64, currency string) Balance {
	coin := NewCoin(amount, currency)
	balance := NewBalance()
	balance.AddAmmount(coin)
	return balance
}

func NewBalanceFromCoin(coin Coin) Balance {
	balance := NewBalance()
	balance.AddAmmount(coin)
	return balance
}

func (b *Balance) FromCoin(coin Coin) {
	b.Amounts[coin.Currency.Id] = coin
}

func (b *Balance) GetAmountByCurrency(currency Currency) Coin {
	return b.Amounts[currency.Id]
}

func (b *Balance) GetAmountByName(name string) Coin {
	return b.Amounts[Currencies[name].Id]
}

func (b *Balance) SetAmmount(coin Coin) {
	b.Amounts[coin.Currency.Id] = coin
	return
}

func (b *Balance) AddAmmount(coin Coin) {
	b.Amounts[coin.Currency.Id] = b.Amounts[coin.Currency.Id].Plus(coin)
	return
}

func (b *Balance) MinusAmmount(coin Coin) {
	b.Amounts[coin.Currency.Id] = b.Amounts[coin.Currency.Id].Minus(coin)
	return
}

func (b Balance) IsEqual(balance Balance) bool {
	for i, v := range b.Amounts {
		if !v.Equals(balance.Amounts[i]) {
			return false
		}
	}
	return true
}

func (b Balance) IsEnough(coins ...Coin) bool {
	for _, coin := range coins {
		b.MinusAmmount(coin)
	}

	for _, coin := range b.Amounts {
		if !coin.LessThan(0) {
			return false
		}
	}
	return true
}

//String used in fmt and Dump
func (b Balance) String() string {
	buffer := ""
	for _, v := range b.Amounts {
		if v.Amount.Cmp(big.NewInt(0)) == 1 || v.Currency.Id == 0 {
			buffer += fmt.Sprintf("%s %s; ", v.Amount.String(), v.Currency.Name)
		}
	}

	if len(buffer) > 0 {
		buffer = buffer[:len(buffer)-1]
	}
	return buffer
}
