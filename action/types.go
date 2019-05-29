package action

import (
	"encoding/json"
	"strconv"

	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/keys"
)

// Address an action package over Address in data/keys package
type Address = keys.Address

// Balance an action package over Balance in data/balance
type Balance = balance.Balance

// Amount is an easily serializable representation of coin. Nodes can create coin from the Amount object
// received over the network
type Amount struct {
	Currency string `json:"currency"`
	Value    string `json:"value"`
}

// New Amount creates a new amount account object
func NewAmount(currency, value string) *Amount {
	return &Amount{currency, value}
}

// IsValid checks the validity of the currency and the amount string in the account object, which may be received
// over a network.
func (a Amount) IsValid(ctx *Context) bool {
	currency, ok := ctx.Currencies.GetCurrencyByName(a.Currency)
	if !ok {
		return false
	}
	f, err := strconv.ParseFloat(a.Value, 64)
	if err != nil {
		return false
	}
	coin := currency.NewCoinFromFloat64(f)
	return coin.IsValid()
}

// String returns a string representation of the Amount object.
func (a Amount) String() string {
	result, _ := json.Marshal(a)
	return string(result)
}

// ToCoin converts an easier to transport Amount object to a Coin object in Oneledger protocol.
// It takes the action context to determine the currency from which to create the coin.
func (a Amount) ToCoin(ctx *Context) balance.Coin {

	// get currency of Amount a
	currency, ok := ctx.Currencies.GetCurrencyByName(a.Currency)
	if !ok {
		return balance.Coin{}
	}

	// parse float string
	f, err := strconv.ParseFloat(a.Value, 64)
	if err != nil {
		return balance.Coin{}
	}
	return currency.NewCoinFromFloat64(f)
}
