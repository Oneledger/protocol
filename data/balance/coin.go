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
	"fmt"
	"math/big"
	"runtime/debug"
)

/*
	Coin starts here
*/
// Coin is the basic amount, specified in integers, at the smallest increment (i.e. a satoshi, not a bitcoin)
type Coin struct {
	Currency Currency `json:"currency"`
	Amount   *big.Int `json:"amount,string"`
}

// See if the coin is one of a list of currencies
func (coin Coin) IsCurrency(currencies ...string) bool {
	if coin.Amount == nil {
		debug.PrintStack()
		logger.Fatal("Invalid Coin", coin)
	}

	found := false
	for _, currency := range currencies {
		if coin.Currency.Name == currency {
			found = true
			break
		}
	}
	return found
}

// LessThan, for coins...
func (coin Coin) LessThanCoin(value Coin) bool {
	if coin.Amount == nil || value.Amount == nil {
		debug.PrintStack()
		logger.Fatal("Invalid Coin", coin)
	}

	if coin.Currency.Chain != value.Currency.Chain {
		logger.Fatal("Compare two different coin", coin, value)
	}

	//logger.Dump("LessThanCoin", value, coin)

	if coin.Amount.Cmp(value.Amount) < 0 {
		return true
	}
	return false
}

// LessThanEqual, for coins...
func (coin Coin) LessThanEqualCoin(value Coin) bool {
	if coin.Amount == nil || value.Amount == nil {
		debug.PrintStack()
		logger.Fatal("Invalid Coin", coin)
	}

	if coin.Currency.Chain != value.Currency.Chain {
		logger.Fatal("Compare two different coin", coin, value)
	}

	//logger.Dump("LessThanEqualCoin", value, coin)

	if coin.Amount.Cmp(value.Amount) <= 0 {
		return true
	}
	return false
}

// IsValid coin or is it broken
func (coin Coin) IsValid() bool {
	if coin.Amount == nil {
		debug.PrintStack()
		logger.Fatal("Invalid Coin", coin)
	}

	if coin.Currency.Name == "" {
		return false
	}
	// there should not exist minus coin
	if coin.Amount.Cmp(big.NewInt(0)) < 0 {
		return false
	}


	// TODO: Combine this with convert.GetCurrency...
	return false
}

// Equals another coin
func (coin Coin) Equals(value Coin) bool {
	if coin.Amount == nil {
		debug.PrintStack()
		logger.Fatal("Invalid Coin", coin)
	}

	if coin.Currency.Chain != value.Currency.Chain {
		return false
	}
	if coin.Amount.Cmp(value.Amount) == 0 {
		return true
	}
	return false
}

// Minus two coins
func (coin Coin) Minus(value Coin) Coin {
	if coin.Amount == nil {
		debug.PrintStack()
		logger.Fatal("Invalid Coin", coin)
	}

	if coin.Currency.Name != value.Currency.Name {
		//logger.Error("Mismatching currencies", "coin", coin, "value", value)
		logger.Fatal("Mismatching currencies", coin, value)
		return coin
	}

	base := big.NewInt(0)
	result := Coin{
		Currency: coin.Currency,
		Amount:   base.Sub(coin.Amount, value.Amount),
	}
	return result
}

// Plus two coins
func (coin Coin) Plus(value Coin) Coin {
	if coin.Amount == nil {
		debug.PrintStack()
		logger.Fatal("Invalid Coin", "coin", coin)
	}

	if coin.Currency.Name != value.Currency.Name {
		//logger.Error("Mismatching currencies", "coin", coin, "value", value)
		logger.Fatal("Mismatching currencies", coin, value)
		return coin
	}

	base := big.NewInt(0)
	result := Coin{
		Currency: coin.Currency,
		Amount:   base.Add(coin.Amount, value.Amount),
	}
	return result
}

// Quotient of one coin by another (divide without remainder, modulus, etc)
func (coin Coin) Quotient(value Coin) Coin {
	if coin.Amount == nil {
		debug.PrintStack()
		logger.Fatal("Invalid Coin", "coin", coin)
	}

	if coin.Currency.Name != value.Currency.Name {
		//logger.Error("Mismatching currencies", "coin", coin, "value", value)
		logger.Fatal("Mismatching currencies", coin, value)
		return coin
	}

	base := big.NewInt(0)
	result := Coin{
		Currency: coin.Currency,
		Amount:   base.Quo(coin.Amount, value.Amount),
	}
	return result
}

func (coin Coin) Divide(value int) Coin {
	if coin.Amount == nil {
		debug.PrintStack()
		logger.Fatal("Invalid Coin", coin)
	}

	base := big.NewInt(0)
	divisor := big.NewInt(int64(value))
	result := Coin{
		Currency: coin.Currency,
		Amount:   base.Div(coin.Amount, divisor),
	}
	return result

}

// Multiply one coin by another
func (coin Coin) Multiply(value Coin) Coin {
	if coin.Amount == nil {
		debug.PrintStack()
		logger.Fatal("Invalid Coin", coin)
	}

	if coin.Currency.Name != value.Currency.Name {
		//logger.Error("Mismatching currencies", "coin", coin, "value", value)
		logger.Fatal("Mismatching currencies", coin, value)
		return coin
	}

	base := big.NewInt(0)
	result := Coin{
		Currency: coin.Currency,
		Amount:   base.Mul(coin.Amount, value.Amount),
	}
	return result
}

// Multiply one coin by another
func (coin Coin) MultiplyInt(value int) Coin {
	if coin.Amount == nil {
		debug.PrintStack()
		logger.Fatal("Invalid Coin", coin)
	}

	multiplier := big.NewInt(int64(value))
	base := big.NewInt(0)
	result := Coin{
		Currency: coin.Currency,
		Amount:   base.Mul(coin.Amount, multiplier),
	}
	return result
}

// Turn a coin into a readable, floating point string with the currency
func (coin Coin) String() string {
	if coin.Amount == nil {
		debug.PrintStack()
		logger.Fatal("Invalid Coin", "err", "Amount is nil")
	}

	return fmt.Sprintf("%=v", coin)
}
