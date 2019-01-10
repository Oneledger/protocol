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
	Amounts map[int]Coin
}

func init() {
	serial.Register(Balance{})
}

func NewBalance() *Balance {
	amounts := make(map[int]Coin, 0)
	result := &Balance{
		Amounts: amounts,
	}
	return result
}

func NewBalanceFromString(amount string, currency string) *Balance {
	coin := NewCoinFromString(amount, currency)
	balance := NewBalance()
	balance.AddCoin(coin)
	return balance
}

func NewBalanceFromInt(amount int64, currency string) *Balance {
	coin := NewCoinFromInt(amount, currency)
	balance := NewBalance()
	balance.AddCoin(coin)
	return balance
}

func NewBalanceFromCoin(coin Coin) *Balance {
	balance := NewBalance()
	balance.AddCoin(coin)
	return balance
}

func (b *Balance) FindCoin(currency Currency) *Coin {
	if coin, ok := b.Amounts[currency.Id]; ok {
		return &coin
	}
	return nil
}

// Add a new or existing coin
func (b *Balance) AddCoin(coin Coin) {
	result := b.FindCoin(coin.Currency)
	if result == nil {
		b.Amounts[coin.Currency.Id] = coin
		return
	}
	b.Amounts[coin.Currency.Id] = result.Plus(coin)
	return
}

func (b *Balance) MinusCoin(coin Coin) {
	result := b.FindCoin(coin.Currency)
	if result == nil {
		// TODO: This results in a negative coin, which is what was asked for...
		base := NewCoinFromInt(0, coin.Currency.Name)
		b.Amounts[coin.Currency.Id] = base.Minus(coin)
		return
	}
	b.Amounts[coin.Currency.Id] = result.Minus(coin)
	return
}

/*
func (b *Balance) FromCoin(coin Coin) {
	b.Amounts[string(coin.Currency.Key())] = coin
}
*/

func (b *Balance) GetCoin(currency Currency) Coin {
	result := b.FindCoin(currency)
	if result == nil {
		// NOTE: Missing coins are actually zero value coins.
		return NewCoinFromInt(0, currency.Name)
	}
	return b.Amounts[currency.Id]
}

// TODO: GetCoinByName?
func (b *Balance) GetAmountByName(name string) Coin {
	currency := NewCurrency(name)
	result := b.FindCoin(currency)
	if result == nil {
		// NOTE: Missing coins are actually zero value coins.
		return NewCoinFromInt(0, name)
	}
	return b.Amounts[Currencies[name].Id]
}

func (b *Balance) SetAmount(coin Coin) {
	b.Amounts[coin.Currency.Id] = coin
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

func (b Balance) IsEnoughBalance(balance Balance) bool {
	for i, coin := range balance.Amounts {
		v, ok := b.Amounts[i]
		if !ok {
			v = NewCoinFromInt(0, coin.Currency.Name)
		}

		if v.Minus(coin).LessThan(0) {
			return false
		}
	}
	return true
}

//String used in fmt and Dump
func (balance Balance) String() string {
	buffer := ""
	for _, coin := range balance.Amounts {
		if coin.Amount.Cmp(big.NewInt(0)) != 0 || coin.Currency.Id == 0 {
			if buffer != "" {
				buffer += ", "
			}
			buffer += coin.String()
		}
	}
	return buffer
}
