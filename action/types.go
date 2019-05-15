package action

import (
	"encoding/json"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/keys"
	"strconv"
)

type Address = keys.Address

type Balance = balance.Balance

type Amount struct {
	Currency string `json:"currency"`
	Value    string `json:"value"`
}

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

func (a Amount) String() string {
	result, _ := json.Marshal(a)
	return string(result)
}

func (a Amount) ToCoin(ctx *Context) balance.Coin {
	currency, ok := ctx.Currencies.GetCurrencyByName(a.Currency)
	if !ok {
		return balance.Coin{}
	}
	f, err := strconv.ParseFloat(a.Value, 64)
	if err != nil {
		return balance.Coin{}
	}
	return currency.NewCoinFromFloat64(f)
}
