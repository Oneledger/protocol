/*
	Copyright 2017 - 2018 OneLedger
*/
package data

import "github.com/Oneledger/protocol/node/log"

// Coin is the basic amount, specified in integers, at the smallest increment (i.e. a satoshi, not a bitcoin)
type Coin struct {
	Currency string `json:"currency"`
	Amount   int64  `json:"amount"`
}

type Coins []Coin

func (coin Coin) Minus(value Coin) Coin {

	if coin.Currency != value.Currency {
		log.Error("Mismatching Currencies", "coin", coin, "value", value)
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
		log.Error("Mismatching Currencies", "coin", coin, "value", value)
		return coin
	}

	result := Coin{
		Currency: coin.Currency,
		Amount:   coin.Amount + value.Amount,
	}
	return result
}
