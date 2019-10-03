package balance

import (
	"github.com/Oneledger/protocol/log"
)

type Context struct {
	balances   *Store
	currencies *CurrencySet
	logger     *log.Logger
}

func (ctx *Context) Store() *Store {
	return ctx.balances
}

func (ctx *Context) Currencies() *CurrencySet {
	return ctx.currencies
}

func NewContext(logger *log.Logger, balances *Store, currencies *CurrencySet) *Context {
	return &Context{
		logger:     logger,
		balances:   balances,
		currencies: currencies,
	}
}
