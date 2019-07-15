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

	if coin.Amount.Int.Cmp(&value.Amount.Int) < 0 {
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

	if coin.Amount.Int.Cmp(&value.Amount.Int) <= 0 {
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
		return coin.Amount.Int.Cmp(big.NewInt(0)) >= 0
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
	if coin.Amount.Int.Cmp(&value.Amount.Int) == 0 {
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
		logger.Error("Mismatching currencies", coin, value)
		return coin, ErrMismatchingCurrency
	}

	base := NewAmount(0)
	result := Coin{
		Currency: coin.Currency,
		Amount:   &Amount{*base.Int.Sub(&coin.Amount.Int, &value.Amount.Int)},
	}
	if result.Amount.Int.Cmp(big.NewInt(0)) == -1 {
		return result, ErrInsufficientBalance
	}
	return result, nil
}

// Plus two coins
func (coin Coin) Plus(value Coin) (Coin, error) {
	if coin.Amount == nil {
		debug.PrintStack()
		logger.Fatal("Invalid Coin", "coin", coin)
	}

	if coin.Currency.Name != value.Currency.Name {
		logger.Error("Mismatching currencies", coin, value)
		return coin, ErrMismatchingCurrency
	}

	base := big.NewInt(0)
	result := Coin{
		Currency: coin.Currency,
		Amount:   &Amount{*base.Add(&coin.Amount.Int, &value.Amount.Int)},
	}
	return result, nil
}

func (coin Coin) Divide(value int) Coin {
	if coin.Amount == nil {
		coin.Amount = NewAmount(0)
	}

	base := big.NewInt(0)
	divisor := big.NewInt(int64(value))
	result := Coin{
		Currency: coin.Currency,
		Amount:   &Amount{*base.Div(&coin.Amount.Int, divisor)},
	}
	return result

}

// Multiply one coin by another
func (coin Coin) MultiplyInt(value int) Coin {
	if coin.Amount == nil {
		coin.Amount = NewAmount(0)
	}

	multiplier := big.NewInt(int64(value))
	base := big.NewInt(0)
	result := Coin{
		Currency: coin.Currency,
		Amount:   &Amount{*base.Mul(&coin.Amount.Int, multiplier)},
	}
	return result
}

// Turn a coin into a readable, floating point string with the currency
func (coin Coin) String() string {
	/*
		if coin.Amount == nil {
			debug.PrintStack()
			logger.Fatal("Invalid Coin", "err", "Amount is nil")
		}
	*/
	float := new(big.Float).SetInt(&coin.Amount.Int)
	value := float.Quo(float, new(big.Float).SetInt(coin.Currency.Base()))

	return fmt.Sprintf("%s %s", value.String(), coin.Currency.Name)
}
