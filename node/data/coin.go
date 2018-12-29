/*
	Copyright 2017 - 2018 OneLedger

	Encapsulate the coins, allow int64 for interfacing and big.Int as base type
*/
package data

import (
	"encoding/hex"
	"math/big"
	"runtime/debug"
	"strconv"

	"golang.org/x/crypto/ripemd160"

	"github.com/Oneledger/protocol/node/log"
	"github.com/Oneledger/protocol/node/serial"
)

// Coin is the basic amount, specified in integers, at the smallest increment (i.e. a satoshi, not a bitcoin)
type Coin struct {
	Currency Currency `json:"currency"`
	Amount   *big.Int `json:"amount"`
}

func init() {
	serial.Register(Coin{})
	serial.Register(Currency{})
}

type Coins []Coin

// TODO: These need to be driven from a domain database, also they are many-to-one with chains
var Currencies map[string]Currency = map[string]Currency{
	"OLT": Currency{"OLT", ONELEDGER, 0},
	"BTC": Currency{"BTC", BITCOIN, 1},
	"ETH": Currency{"ETH", ETHEREUM, 2},
	"VT":  Currency{"VT", ONELEDGER, 3},
}

type Extra struct {
	Units   *big.Float
	Decimal int
	Format  uint8
}

// TODO: Separated from Currency to avoid serializing big floats and giving out this info
var CurrenciesExtra map[string]Extra = map[string]Extra{
	"OLT": Extra{big.NewFloat(1000000000000000000), 6, 'f'},
	"BTC": Extra{big.NewFloat(1), 0, 'f'}, // TODO: This needs to be set correctly
	"ETH": Extra{big.NewFloat(1), 0, 'f'}, // TODO: This needs to be set correctly
	"VT":  Extra{big.NewFloat(1), 0, 'f'},
}

var defaultBase *big.Float = big.NewFloat(1)

func GetBase(currency string) *big.Float {
	return GetExtra(currency).Units
}

func GetExtra(currency string) Extra {
	if value, ok := CurrenciesExtra[currency]; ok {
		return value
	}
	return CurrenciesExtra["OLT"]
}

type Currency struct {
	Name  string    `json:"name"`
	Chain ChainType `json:"chain"`
	Id    int       `json:"id"`
}

// Key sets a encodable key for the currency entry, we may end up using currencyCodes instead.
func (c Currency) Key() string {
	hasher := ripemd160.New()

	buffer, err := serial.Serialize(c, serial.JSON)
	if err != nil {
		log.Fatal("hash serialize failed", "err", err)
	}
	_, err = hasher.Write(buffer)
	if err != nil {
		log.Fatal("hasher failed", "err", err)
	}
	buffer = hasher.Sum(nil)

	return hex.EncodeToString(buffer)
}

// Look up the currency
func NewCurrency(currency string) Currency {
	return Currencies[currency]
}

func NewCoinFromUnits(amount int64, currency string) Coin {
	value := units2bint(amount, GetBase(currency))
	coin := Coin{
		Currency: Currencies[currency],
		Amount:   value,
	}
	if !coin.IsValid() {
		log.Warn("Create Invalid Coin", "coin", coin)
	}
	return coin
}

// Create a coin from integer (not fractional)
func NewCoinFromInt(amount int64, currency string) Coin {

	value := int2bint(amount, GetBase(currency))
	coin := Coin{
		Currency: Currencies[currency],
		Amount:   value,
	}
	if !coin.IsValid() {
		log.Warn("Create Invalid Coin", "coin", coin)
	}
	return coin
}

// Create a coin from floating point
func NewCoinFromFloat(amount float64, currency string) Coin {
	value := float2bint(amount, GetBase(currency))
	coin := Coin{
		Currency: Currencies[currency],
		Amount:   value,
	}
	if !coin.IsValid() {
		log.Warn("Create Invalid Coin", "amount", amount, "coin", coin)
	}
	return coin
}

// Create a coin from string
func NewCoinFromString(amount string, currency string) Coin {
	value := parseString(amount, GetBase(currency))
	coin := Coin{
		Currency: Currencies[currency],
		Amount:   value,
	}
	if !coin.IsValid() {
		log.Warn("Create Invalid Coin", "coin", coin)
	}
	return coin
}

// Handle an incoming string
func parseString(amount string, base *big.Float) *big.Int {
	log.Dump("Parsing Amount", amount)
	if amount == "" {
		log.Error("Empty Amount String", "amount", amount)
		return nil
	}

	/*
		value := new(big.Float)

		_, err := fmt.Sscan(amount, value)
		if err != nil {
			log.Error("Invalid Float String", "err", err, "amount", amount)
			return nil
		}
		result := bfloat2bint(value, base)
	*/

	value, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		log.Error("Invalid Float String", "err", err, "amount", amount)
		return nil
	}

	result := float2bint(value, base)

	return result
}

func units2bint(amount int64, base *big.Float) *big.Int {
	value := new(big.Int).SetInt64(amount)
	return value
}

// Handle an incoming int (often used for comaparisons)
func int2bint(amount int64, base *big.Float) *big.Int {
	value := new(big.Float).SetInt64(amount)

	interim := value.Mul(value, base)
	result, _ := interim.Int(nil)

	//log.Dump("int2bint", amount, result)
	return result
}

// Handle an incoming float
func float2bint(amount float64, base *big.Float) *big.Int {
	value := big.NewFloat(amount)

	interim := value.Mul(value, base)
	result, _ := interim.Int(nil)

	return result
}

// Handle an big float to big int conversion
func bfloat2bint(value *big.Float, base *big.Float) *big.Int {
	interim := value.Mul(value, base)
	result, _ := interim.Int(nil)

	return result
}

// Handle a big int to outgoing float
func bint2float(amount *big.Int, base *big.Float) float64 {
	value := new(big.Float).SetInt(amount)

	interim := value.Quo(value, base)
	result, _ := interim.Float64()

	return result
}

func (coin Coin) Float64() float64 {
	return bint2float(coin.Amount, GetBase(coin.Currency.Name))
}

// See if the coin is one of a list of currencies
func (coin Coin) IsCurrency(currencies ...string) bool {
	if coin.Amount == nil {
		debug.PrintStack()
		log.Fatal("Invalid Coin", "coin", coin)
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
		log.Fatal("Invalid Coin", "coin", coin)
	}

	compare := float2bint(value, GetBase("OLT"))
	//log.Dump("LessThanEqual", value, compare)

	if coin.Amount.Cmp(compare) <= 0 {
		return true
	}
	return false
}

// LessThan, just for OLTs...
func (coin Coin) LessThan(value float64) bool {
	if coin.Amount == nil {
		debug.PrintStack()
		log.Fatal("Invalid Coin", "coin", coin)
	}

	compare := float2bint(value, GetBase("OLT"))
	//log.Dump("LessThanEqual", value, compare)

	if coin.Amount.Cmp(compare) < 0 {
		return true
	}
	return false
}

// LessThan, for coins...
func (coin Coin) LessThanCoin(value Coin) bool {
	if coin.Amount == nil || value.Amount == nil {
		debug.PrintStack()
		log.Fatal("Invalid Coin", "coin", coin)
	}

	if coin.Currency.Id != value.Currency.Id {
		log.Fatal("Compare two different coin", "coin", coin, "value", value)
	}

	//log.Dump("LessThanCoin", value, coin)

	if coin.Amount.Cmp(value.Amount) < 0 {
		return true
	}
	return false
}

// LessThanEqual, for coins...
func (coin Coin) LessThanEqualCoin(value Coin) bool {
	if coin.Amount == nil || value.Amount == nil {
		debug.PrintStack()
		log.Fatal("Invalid Coin", "coin", coin)
	}

	if coin.Currency.Id != value.Currency.Id {
		log.Fatal("Compare two different coin", "coin", coin, "value", value)
	}

	//log.Dump("LessThanEqualCoin", value, coin)

	if coin.Amount.Cmp(value.Amount) <= 0 {
		return true
	}
	return false
}

/*
func (coin Coin) EqualsInt64(value int64) bool {
	if coin.Amount == nil {
		log.Fatal("Invalid Coin", "coin", coin)
	}

	if coin.Amount.Cmp(big.NewInt(int64(value))) == 0 {
		return true
	}
	return false
}
*/

// IsValid coin or is it broken
func (coin Coin) IsValid() bool {
	if coin.Amount == nil {
		debug.PrintStack()
		log.Fatal("Invalid Coin", "coin", coin)
	}

	if coin.Currency.Name == "" {
		return false
	}

	if _, ok := Currencies[coin.Currency.Name]; ok {
		return true
	}

	// TODO: Combine this with convert.GetCurrency...
	return false
}

// Equals another coin
func (coin Coin) Equals(value Coin) bool {
	if coin.Amount == nil {
		debug.PrintStack()
		log.Fatal("Invalid Coin", "coin", coin)
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
		log.Fatal("Invalid Coin", "coin", coin)
	}

	if coin.Currency.Name != value.Currency.Name {
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

// Plus two coins
func (coin Coin) Plus(value Coin) Coin {
	if coin.Amount == nil {
		debug.PrintStack()
		log.Fatal("Invalid Coin", "coin", coin)
	}

	if coin.Currency.Name != value.Currency.Name {
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

// Quotient of one coin by another (divide without remainder, modulus, etc)
func (coin Coin) Quotient(value Coin) Coin {
	if coin.Amount == nil {
		debug.PrintStack()
		log.Fatal("Invalid Coin", "coin", coin)
	}

	if coin.Currency.Name != value.Currency.Name {
		//log.Error("Mismatching Currencies", "coin", coin, "value", value)
		log.Fatal("Mismatching Currencies", "coin", coin, "value", value)
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
		log.Fatal("Invalid Coin", "coin", coin)
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
		log.Fatal("Invalid Coin", "coin", coin)
	}

	if coin.Currency.Name != value.Currency.Name {
		//log.Error("Mismatching Currencies", "coin", coin, "value", value)
		log.Fatal("Mismatching Currencies", "coin", coin, "value", value)
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
		log.Fatal("Invalid Coin", "coin", coin)
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
		log.Fatal("Invalid Coin", "err", "Amount is nil")
	}

	currency := coin.Currency.Name
	extra := GetExtra(currency)
	float := new(big.Float).SetInt(coin.Amount)
	value := float.Quo(float, extra.Units)
	text := value.Text(extra.Format, extra.Decimal) + " " + currency

	return text
}
