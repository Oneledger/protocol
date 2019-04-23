/*

   ____             _          _
  / __ \           | |        | |
 | |  | |_ __   ___| | ___  __| | __ _  ___ _ __
 | |  | | '_ \ / _ \ |/ _ \/ _` |/ _` |/ _ \ '__|
 | |__| | | | |  __/ |  __/ (_| | (_| |  __/ |
  \____/|_| |_|\___|_|\___|\__,_|\__, |\___|_|
                                  __/ |
                                 |___/

	Copyright 2017 - 2019 OneLedger

*/

package balance

import (
	"github.com/Oneledger/protocol/data"
)

// Wrap the amount with owner information
type Balance struct {
	Amounts   map[int]Coin
	coinOrder []int // this field helps to maintain order during serialization ;
	// so that all the nodes have the same hash of account balances
}

/*
	balance Generators start here
 */
func NewBalance() *Balance {
	amounts := make(map[int]Coin, 0)
	result := &Balance{
		Amounts:   amounts,
		coinOrder: []int{},
	}
	return result
}

func NewBalanceFromString(amount string, currency string) *Balance {
	coin := NewCoinFromString(amount, currency)
	b := NewBalance()
	b.AddCoin(coin)
	return b
}

func NewBalanceFromInt(amount int64, currency string) *Balance {
	coin := NewCoinFromInt(amount, currency)
	b := NewBalance()
	b.AddCoin(coin)
	return b
}

func GetBalanceFromDb(db data.GetSetter, accountKey data.DatastoreKey) (*Balance, error) {
	data, err := db.Get(accountKey)
	if err != nil {
		return nil, err
	}

	var b = &Balance{}
	err = pSzlr.Deserialize(data, b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

/*
	methods for balance start here
 */
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
		b.coinOrder = append(b.coinOrder, coin.Currency.Id)
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
		b.coinOrder = append(b.coinOrder, coin.Currency.Id)
		return
	}
	b.Amounts[coin.Currency.Id] = result.Minus(coin)
	return
}

func (b *Balance) GetCoin(currency Currency) Coin {
	result := b.FindCoin(currency)
	if result == nil {
		// NOTE: Missing coins are actually zero value coins.
		return NewCoinFromInt(0, currency.Name)
	}
	return b.Amounts[currency.Id]
}

// GetCoinByName
func (b *Balance) GetCoinByName(name string) Coin {
	currency := NewCurrency(name)
	result := b.FindCoin(currency)
	if result == nil {
		// NOTE: Missing coins are actually zero value coins.
		return NewCoinFromInt(0, name)
	}
	return b.Amounts[currencies[name].Id]
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

// String method used in fmt and Dump
func (b Balance) String() string {
	buffer := ""
	for _, coin := range b.Amounts {
		if buffer != "" {
			buffer += ", "
		}
		buffer += coin.String()
	}
	return buffer
}
