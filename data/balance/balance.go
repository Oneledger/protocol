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
	"github.com/Oneledger/protocol/storage"
)

// Wrap the amount with owner information
type Balance struct {
	Amounts map[string]Coin `json:"amounts"`
}

/*
	balance Generators start here
*/
func NewBalance() *Balance {
	amounts := make(map[string]Coin, 0)
	result := &Balance{
		Amounts: amounts,
	}
	return result
}

/*
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
*/

// GetBalanceFromDb takes a datastore with GetSetter interface and initializes a new Balance
// from the data.
func GetBalanceFromDb(db storage.Store, accountKey storage.StoreKey) (*Balance, error) {
	dat, err := db.Get(accountKey)
	if err != nil {
		return nil, err
	}

	var b = &Balance{}
	err = pSzlr.Deserialize(dat, b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

/*
	methods for balance start here
*/
func (b *Balance) FindCoin(currency Currency) *Coin {
	if coin, ok := b.Amounts[currency.StringKey()]; ok {
		return &coin
	}
	return nil
}

// Add a new or existing coin
func (b *Balance) AddCoin(coin Coin) {
	result := b.FindCoin(coin.Currency)
	if result == nil {
		b.Amounts[coin.Currency.StringKey()] = coin
		return
	}
	b.Amounts[coin.Currency.StringKey()] = result.Plus(coin)
	return
}

func (b *Balance) MinusCoin(coin Coin) {
	result := b.FindCoin(coin.Currency)
	if result == nil {
		// TODO: This results in a negative coin, which is what was asked for...
		base := coin.Currency.NewCoinFromInt(0)
		b.Amounts[coin.Currency.StringKey()] = base.Minus(coin)
		return
	}
	b.Amounts[coin.Currency.StringKey()] = result.Minus(coin)
	return
}

func (b *Balance) GetCoin(currency Currency) Coin {
	result := b.FindCoin(currency)
	if result == nil {
		// NOTE: Missing coins are actually zero value coins.
		return currency.NewCoinFromInt(0)
	}
	return b.Amounts[currency.StringKey()]
}

func (b *Balance) setAmount(coin Coin) {
	b.Amounts[coin.Currency.StringKey()] = coin
	return
}

func (b Balance) IsEnoughBalance(balance Balance) bool {
	for i, coin := range balance.Amounts {
		v, ok := b.Amounts[i]
		if !ok {
			v = coin.Currency.NewCoinFromInt(0)
		}

		if v.Minus(coin).LessThanCoin(coin.Currency.NewCoinFromInt(0)) {
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
