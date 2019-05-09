package action

import (
	"fmt"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/keys"
	"math/big"
)

type Address = keys.Address

type Balance = balance.Balance


type Amount struct {
	Value []byte
	Currency string
}

func (a Amount) ToCoin(ctx *Context) (balance.Coin, error) {
	currency, ok := ctx.Currencies[a.Currency]
	if !ok {
		return balance.Coin{}, ErrInvalidAmount
	}
	return currency.NewCoinFromBytes(a.Value), nil
}

func (a Amount) IsValid(ctx *Context) (balance.Coin, error) {
	coin, err := a.ToCoin(ctx)
	if err != nil {
		return balance.Coin{}, err
	}
	return coin, nil
}

func (a Amount) String() string {
	n := big.NewInt(0).SetBytes(a.Value)
	return fmt.Sprintf(n.String(), a.Currency)
}