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
	Amounts map[string]Coin
}

func init() {
	serial.Register(Balance{})
}

func NewBalance() *Balance {
	amounts := make(map[string]Coin)
	coin := NewCoin(0, "OLT")
	amounts[string(coin.Currency.Key())] = coin
	result := &Balance{
		Amounts: amounts,
	}
	return result
}

func NewBalanceFromString(amount int64, currency string) *Balance {
	coin := NewCoin(amount, currency)
	balance := NewBalance()
	balance.AddAmount(coin)
	return balance
}

func NewBalanceFromCoin(coin Coin) *Balance {
	balance := NewBalance()
	balance.AddAmount(coin)
	return balance
}

func (b *Balance) FromCoin(coin Coin) {
	b.Amounts[string(coin.Currency.Key())] = coin
}

func (b *Balance) GetAmountByCurrency(currency Currency) Coin {
	v, ok := b.Amounts[currency.Key()]
	if !ok {
		return NewCoin(0, currency.Name)
	}
	return v
}

func (b *Balance) GetAmountByName(name string) Coin {
	v, ok := b.Amounts[Currencies[name].Key()]
	if !ok {
		return NewCoin(0, name)
	}
	return v
}

func (b *Balance) SetAmount(coin Coin) {
	b.Amounts[coin.Currency.Key()] = coin
	return
}

func (b *Balance) AddAmount(coin Coin) {
	key := coin.Currency.Key()
	v, ok := b.Amounts[key]
	if ok {
		coin = v.Plus(coin)
	}
	b.Amounts[key] = coin
	return
}

func (b *Balance) MinusAmount(coin Coin) {
	key := coin.Currency.Key()
	v, ok := b.Amounts[key]
	if ok {
		coin = v.Minus(coin)
	}
	b.Amounts[key] = coin
	return
}

func (b Balance) IsEnoughBalance(balance Balance) bool {
	for i, coin := range balance.Amounts {
		v, ok := b.Amounts[i]
		if !ok {
			v = NewCoin(0, coin.Currency.Name)
		}

		if v.Minus(coin).LessThan(0) {
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
