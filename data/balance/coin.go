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
	"math/big"
	"runtime/debug"
)

/*
	Coin starts here
*/
// Coin is the basic amount, specified in integers, at the smallest increment (i.e. a satoshi, not a bitcoin)
type Coin struct {
	Currency Currency `json:"currency"`
	Amount   *big.Int `json:"amount"`
}

func NewCoinFromUnits(amount int64, currency string) Coin {
	value := units2bint(amount, GetBase(currency))

	coin := Coin{
		Currency: currencies[currency],
		Amount:   value,
	}
	if !coin.IsValid() {
		logger.Warn("Create Invalid Coin", coin)
	}

	return coin
}

// Create a coin from integer (not fractional)
func NewCoinFromInt(amount int64, currency string) Coin {

	value := int2bint(amount, GetBase(currency))

	coin := Coin{
		Currency: currencies[currency],
		Amount:   value,
	}
	if !coin.IsValid() {
		logger.Warn("Create Invalid Coin", coin)
	}

	return coin
}

// Create a coin from floating point
func NewCoinFromFloat(amount float64, currency string) Coin {
	value := float2bint(amount, GetBase(currency))

	coin := Coin{
		Currency: currencies[currency],
		Amount:   value,
	}
	if !coin.IsValid() {
		logger.Warn("Create Invalid Coin", amount, coin)
	}

	return coin
}

// Create a coin from string
func NewCoinFromString(amount string, currency string) Coin {
	value := parseString(amount, GetBase(currency))

	coin := Coin{
		Currency: currencies[currency],
		Amount:   value,
	}
	if !coin.IsValid() {
		logger.Warn("Create Invalid Coin", coin)
	}

	return coin
}

func (coin Coin) Float64() float64 {
	return bint2float(coin.Amount, GetBase(coin.Currency.Name))
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

// LessThanEqual, just for OLTs...
func (coin Coin) LessThanEqual(value float64) bool {
	if coin.Amount == nil {
		debug.PrintStack()
		logger.Fatal("Invalid Coin", coin)
	}

	compare := float2bint(value, GetBase("OLT"))
	//logger.Dump("LessThanEqual", value, compare)

	if coin.Amount.Cmp(compare) <= 0 {
		return true
	}
	return false
}

// LessThan, just for OLTs...
func (coin Coin) LessThan(value float64) bool {
	if coin.Amount == nil {
		debug.PrintStack()
		logger.Fatal("Invalid Coin", coin)
	}

	compare := float2bint(value, GetBase("OLT"))
	//logger.Dump("LessThanEqual", value, compare)

	if coin.Amount.Cmp(compare) < 0 {
		return true
	}
	return false
}

// LessThan, for coins...
func (coin Coin) LessThanCoin(value Coin) bool {
	if coin.Amount == nil || value.Amount == nil {
		debug.PrintStack()
		logger.Fatal("Invalid Coin", coin)
	}

	if coin.Currency.Id != value.Currency.Id {
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

	if coin.Currency.Id != value.Currency.Id {
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

	if _, ok := currencies[coin.Currency.Name]; ok {
		return true
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

	if coin.Currency.Id != value.Currency.Id {
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

	currency := coin.Currency.Name
	extra := GetExtra(currency)
	float := new(big.Float).SetInt(coin.Amount)
	value := float.Quo(float, extra.Units)
	text := value.Text(extra.Format, extra.Decimal) + " " + currency

	return text
}
