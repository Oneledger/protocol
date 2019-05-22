package balance

import (
	"github.com/Oneledger/protocol/log"
)

type Context struct {
	balances   *Store
	currencies *CurrencyList
	logger     *log.Logger
}

func (ctx *Context) Store() *Store {
	return ctx.balances
}

func (ctx *Context) Currencies() *CurrencyList {
	return ctx.currencies
}

func NewContext(logger *log.Logger, balances *Store, currencies *CurrencyList) *Context {
	return &Context{
		logger:     logger,
		balances:   balances,
		currencies: currencies,
	}
}
