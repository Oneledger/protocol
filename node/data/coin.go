/*
	Copyright 2017 - 2018 OneLedger

	Encapsulate the coins, allow int64 for interfacing and big.Int as base type
*/
package data

import (
	"fmt"
	"math/big"

	"github.com/Oneledger/protocol/node/log"
)

// Coin is the basic amount, specified in integers, at the smallest increment (i.e. a satoshi, not a bitcoin)
type Coin struct {
	Currency string   `json:"currency"`
	Amount   *big.Int `json:"amount"` // TODO: Switch to math/big
	//Amount   int64  `json:"amount"` // TODO: Switch to math/big
}

type Coins []Coin

var OLTBase *big.Float = big.NewFloat(1000000000000000000)

var Currencies map[string]bool = map[string]bool{
	"OLT": true,
	"BLT": true,
	"ETH": true,
}

func NewCoin(amount int64, currency string) Coin {
	value := big.NewInt(amount)
	coin := Coin{
		Currency: currency,
		Amount:   value,
	}
	if !coin.IsValid() {
		log.Warn("Create Invalid Coin", "coin", coin)
	}
	return coin
}

func (coin Coin) LessThanEqual(value int64) bool {
	if coin.Amount.Cmp(big.NewInt(value)) <= 0 {
		return true
	}
	return false
}

func (coin Coin) LessThan(value int64) bool {
	if coin.Amount.Cmp(big.NewInt(value)) < 0 {
		return true
	}
	return false
}

func (coin Coin) IsValid() bool {
	if coin.Currency == "" {
		return false
	}

	// TODO: Combine this with convert.GetCurrency...
	if coin.Currency == "OLT" {
		return true
	}
	if coin.Currency == "BTC" {
		return true
	}
	if coin.Currency == "ETH" {
		return true
	}

	return false
}

func (coin Coin) Equals(value Coin) bool {
	if coin.Amount.Cmp(value.Amount) == 0 {
		return true
	}
	return false
}

func (coin Coin) EqualsInt64(value int64) bool {
	if coin.Amount.Cmp(big.NewInt(int64(value))) == 0 {
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

	base := big.NewInt(0)
	result := Coin{
		Currency: coin.Currency,
		Amount:   base.Sub(coin.Amount, value.Amount),
	}
	return result
}

func (coin Coin) Plus(value Coin) Coin {
	if coin.Currency != value.Currency {
		//log.Error("Mismatching Currencies", "coin", coin, "value", value)
		log.Fatal("Mismatching Currencies", "coin", coin, "value", value)
		return coin
	}

	base := big.NewInt(0)
	result := Coin{
		Currency: coin.Currency,
		Amount:   base.Add(coin.Amount, value.Amount),
	}
	return result
}

func (coin Coin) AsString() string {
	value := new(big.Float).SetInt(coin.Amount)
	result := value.Quo(value, OLTBase)
	text := fmt.Sprintf("%f.3", result)
	return text
}
