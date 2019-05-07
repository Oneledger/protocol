package balance

import (
	"errors"
	"math/big"

	"github.com/Oneledger/protocol/log"
	"github.com/Oneledger/protocol/storage"
)

var (
	ErrDuplicateCurrency = errors.New("provided currency has already been registered")
)

type Context struct {
	balances        *storage.ChainState
	currencies      map[string]Currency
	currenciesExtra map[string]Extra

	logger *log.Logger
}

func (ctx *Context) Store() *storage.ChainState {
	return ctx.balances
}

func (ctx *Context) Currencies() map[string]Currency {
	return ctx.currencies
}

func NewContext(logger *log.Logger, balances *storage.ChainState, currencies map[string]Currency, currenciesExtra map[string]Extra) *Context {
	return &Context{
		logger:          logger,
		balances:        balances,
		currencies:      currencies,
		currenciesExtra: currenciesExtra,
	}
}

// Register registers a new type of currency
func (ctx *Context) RegisterCurrency(currency Currency) error {
	_, ok := ctx.currencies[currency.Name]
	if ok { // If the currency is already registered, return a duplicate error
		return ErrDuplicateCurrency
	}
	ctx.currencies[currency.Name] = currency
	// TODO: redesign how balance.Extra works
	ctx.currenciesExtra[currency.Name] = Extra{
		Units:   big.NewFloat(1000000000000000000),
		Decimal: 6,
		Format:  'f',
	}
	return nil
}
