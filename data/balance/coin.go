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
	Amount   *Amount  `json:"amount,string"`
}

// See if the coin is one of a list of currencies
func (coin Coin) IsCurrency(currencies ...string) bool {
	if coin.Amount == nil {
		return false
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
		coin.Amount = NewAmount(0)
	}

	if coin.Currency.Chain != value.Currency.Chain {
		logger.Fatal("Compare two different coin", coin, value)
	}

	if coin.Amount.BigInt().Cmp(value.Amount.BigInt()) < 0 {
		return true
	}
	return false
}

// LessThanEqual, for coins...
func (coin Coin) LessThanEqualCoin(value Coin) bool {
	if coin.Amount == nil || value.Amount == nil {
		coin.Amount = NewAmount(0)
	}

	if coin.Currency.Chain != value.Currency.Chain {
		logger.Fatal("Compare two different coin", coin, value)
	}

	//logger.Dump("LessThanEqualCoin", value, coin)

	if coin.Amount.BigInt().Cmp(value.Amount.BigInt()) <= 0 {
		return true
	}
	return false
}

// IsValid coin or is it broken
func (coin Coin) IsValid() bool {
	switch {
	case coin.Amount == nil:
		return false
	case coin.Currency.Name == "":
		return false
	default:
		return coin.Amount.BigInt().Cmp(big.NewInt(0)) >= 0
	}
}

// Equals another coin
func (coin Coin) Equals(value Coin) bool {
	if coin.Amount == nil {
		coin.Amount = NewAmount(0)
	}

	if coin.Currency.Chain != value.Currency.Chain {
		return false
	}
	if coin.Amount.BigInt().Cmp(value.Amount.BigInt()) == 0 {
		return true
	}
	return false
}

// Minus two coins
func (coin Coin) Minus(value Coin) (Coin, error) {
	if coin.Amount == nil {
		coin.Amount = NewAmount(0)
	}

	if coin.Currency.Name != value.Currency.Name {
		logger.Fatal("Mismatching currencies", coin, value)
	}

	base := NewAmount(0)
	result := Coin{
		Currency: coin.Currency,
		Amount:   (*Amount)(base.BigInt().Sub(coin.Amount.BigInt(), value.Amount.BigInt())),
	}
	if result.Amount.BigInt().Cmp(big.NewInt(0)) == -1 {
		return result, ErrInsufficientBalance
	}
	return result, nil
}

// Plus two coins
func (coin Coin) Plus(value Coin) Coin {
	if coin.Amount == nil {
		debug.PrintStack()
		logger.Fatal("Invalid Coin", "coin", coin)
	}

	if coin.Currency.Name != value.Currency.Name {
		logger.Fatal("Mismatching currencies", coin, value)
	}

	base := big.NewInt(0)
	result := Coin{
		Currency: coin.Currency,
		Amount:   (*Amount)(base.Add(coin.Amount.BigInt(), value.Amount.BigInt())),
	}
	return result
}

func (coin Coin) Divide(value int) Coin {
	return coin.DivideInt64(int64(value))
}

func (coin Coin) DivideInt64(value int64) Coin {
	if coin.Amount == nil {
		coin.Amount = NewAmount(0)
	}

	base := big.NewInt(0)
	divisor := big.NewInt(value)
	result := Coin{
		Currency: coin.Currency,
		Amount:   (*Amount)(base.Div(coin.Amount.BigInt(), divisor)),
	}
	return result

}

// Multiply one coin by another
func (coin Coin) MultiplyInt(value int) Coin {
	return coin.MultiplyInt64(int64(value))
}

func (coin Coin) MultiplyInt64(value int64) Coin {
	if coin.Amount == nil {
		coin.Amount = NewAmount(0)
		return coin
	}

	multiplier := big.NewInt(value)
	base := big.NewInt(0)
	result := Coin{
		Currency: coin.Currency,
		Amount:   (*Amount)(base.Mul(coin.Amount.BigInt(), multiplier)),
	}
	return result
}

// Turn a coin into a readable, floating point string with the currency
func (coin Coin) String() string {
	return fmt.Sprintf("%s %s", coin.Humanize(), coin.Currency.Name)
}

func (coin Coin) Humanize() string {
	return PrintDecimal(coin.Amount.BigInt(), int(coin.Currency.Decimal))
}

func PrintDecimal(i *big.Int, decimal int) string {
	str := i.String()
	l := len(str)

	if l <= decimal {
		padZero := ""
		for i := 0; i < decimal-l; i++ {
			padZero += "0"
		}
		return removeTrailingZero("0." + padZero + str)
	}

	return removeTrailingZero(str[:l-decimal] + "." + str[l-decimal:])
}

func removeTrailingZero(str string) string {
	l := len(str)
	if l == 0 {
		return ""
	}

	if str[l-1] == '.' {
		return str[:l-1]
	}

	if str[l-1] == '0' {
		return removeTrailingZero(str[:l-1])
	}

	return str
}
