/*
	Copyright 2017-2018 OneLedger
*/

package data

import (
	"math/big"

	"github.com/Oneledger/protocol/node/serial"
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
	amounts := make([]Coin, 0)

	/*
		for _, value := range Currencies {
			amounts[value.Id] = NewCoinFromInt(0, value.Name)
		}
	*/
	result := Balance{
		Amounts: amounts,
	}
	return result
}

func NewBalanceFromString(amount string, currency string) Balance {
	coin := NewCoinFromString(amount, currency)
	balance := NewBalance()
	balance.AddCoin(coin)
	return balance
}

func NewBalanceFromInt(amount int64, currency string) Balance {
	coin := NewCoinFromInt(amount, currency)
	balance := NewBalance()
	balance.AddCoin(coin)
	return balance
}

func NewBalanceFromCoin(coin Coin) Balance {
	balance := NewBalance()
	balance.AddCoin(coin)
	return balance
}

func (b *Balance) FindCoin(currency Currency) *Coin {
	for _, next := range b.Amounts {
		if next.Currency == currency {
			return &next
		}
	}
	return nil
}

// Add a new or existing coin
func (b *Balance) AddCoin(coin Coin) {
	result := b.FindCoin(coin.Currency)
	if result == nil {
		b.Amounts = append(b.Amounts, coin)
		return
	}
	result.Plus(coin)
	return
}

func (b *Balance) MinusCoin(coin Coin) {
	result := b.FindCoin(coin.Currency)
	if result == nil {
		// TODO: This results in a negative coin
		base := NewCoinFromInt(0, coin.Currency.Name)
		b.Amounts = append(b.Amounts, base.Minus(coin))
		return
	}
	result.Minus(coin)
	return
}

/*
func (b *Balance) FromCoin(coin Coin) {
	b.Amounts[coin.Currency.Id] = coin
}
*/

func (b *Balance) GetCoin(currency Currency) Coin {
	result := b.FindCoin(currency)
	if result == nil {
		// Missing coins are actually zero value coins.
		return NewCoinFromInt(0, currency.Name)
	}
	return b.Amounts[currency.Id]
}

// TODO: GetCoinByName?
func (b *Balance) GetAmountByName(name string) Coin {
	currency := NewCurrency(name)
	result := b.FindCoin(currency)
	if result == nil {
		// Missing coins are actually zero value coins.
		return NewCoinFromInt(0, name)
	}
	return b.Amounts[Currencies[name].Id]
}

func (b *Balance) SetAmount(coin Coin) {
	b.AddCoin(coin)
	return
}

func (b *Balance) AddAmount(coin Coin) {
	b.AddCoin(coin)
	return
}

func (b *Balance) MinusAmount(coin Coin) {
	b.MinusCoin(coin)
	return
}

//String used in fmt and Dump
func (balance Balance) String() string {
	buffer := ""
	for _, coin := range balance.Amounts {
		// TODO: Why?
		if coin.Amount.Cmp(big.NewInt(0)) == 1 || coin.Currency.Id == 0 {
			if buffer != "" {
				buffer += ", "
			}
			buffer += coin.String()
		}
	}

	/*
		if len(buffer) > 0 {
			buffer = buffer[:len(buffer)]
		}
	*/
	return buffer
}
