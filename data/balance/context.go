package balance

import (
	"errors"

	"github.com/Oneledger/protocol/log"
)

var (
	ErrDuplicateCurrency = errors.New("provided currency has already been registered")
)

type Context struct {
	balances        *Store
	currencies      map[string]Currency

	logger *log.Logger
}

func (ctx *Context) Store() *Store {
	return ctx.balances
}

func (ctx *Context) Currencies() map[string]Currency {
	return ctx.currencies
}

func NewContext(logger *log.Logger, balances *Store, currencies map[string]Currency) *Context {
	return &Context{
		logger:          logger,
		balances:        balances,
		currencies:      currencies,
	}
}

// Register registers a new type of currency
func (ctx *Context) RegisterCurrency(currency Currency) error {
	_, ok := ctx.currencies[currency.Name]
	if ok { // If the currency is already registered, return a duplicate error
		return ErrDuplicateCurrency
	}
	ctx.currencies[currency.Name] = currency
	return nil
}

