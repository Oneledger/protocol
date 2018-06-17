/*
	Copyright 2017 - 2018 OneLedger
*/
package data

import "github.com/Oneledger/protocol/node/log"

// Coin is the basic amount, specified in integers, at the smallest increment (i.e. a satoshi, not a bitcoin)
type Coin struct {
	Currency string `json:"currency"`
	Amount   int64  `json:"amount"` // TODO: Switch to math/big
}

type Coins []Coin

func NewCoin(amount int64, currency string) Coin {
	return Coin{Currency: currency, Amount: amount}
}

func (coin Coin) LessThanEqual(value int) bool {
	if coin.Amount <= int64(value) {
		return true
	}
	return false
}

func (coin Coin) LessThan(value int) bool {
	if coin.Amount < int64(value) {
		return true
	}
	return false
}

func (coin Coin) IsValid() bool {
	if coin.Currency == "" {
		return false
	}
	// TODO: Double check this
	if coin.LessThan(0) {
		return false
	}
	return true
}

func (coin Coin) Equals(value Coin) bool {
	if coin.Amount == value.Amount {
		return true
	}
	return false
}

func (coin Coin) EqualsInt64(value int64) bool {
	if coin.Amount == value {
		return true
	}
	return false
}

func (coin Coin) Minus(value Coin) Coin {

	if coin.Currency != value.Currency {
		//log.Error("Mismatching Currencies", "coin", coin, "value", value)
		log.Fatal("Mismatching Currencies", "coin", coin, "value", value)
		return coin
	}

	result := Coin{
		Currency: coin.Currency,
		Amount:   coin.Amount - value.Amount,
	}
	return result
}

func (coin Coin) Plus(value Coin) Coin {
	if coin.Currency != value.Currency {
		//log.Error("Mismatching Currencies", "coin", coin, "value", value)
		log.Fatal("Mismatching Currencies", "coin", coin, "value", value)
		return coin
	}

	result := Coin{
		Currency: coin.Currency,
		Amount:   coin.Amount + value.Amount,
	}
	return result
}
