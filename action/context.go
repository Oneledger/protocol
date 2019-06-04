package action

import (
	"github.com/Oneledger/protocol/data/accounts"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/identity"
	"github.com/Oneledger/protocol/log"
)

type Context struct {
	Router     Router
	Accounts   accounts.Wallet
	Balances   *balance.Store
	Currencies *balance.CurrencyList
	Validators *identity.ValidatorStore
	Logger     *log.Logger
}

func NewContext(r Router, wallet accounts.Wallet, balances *balance.Store, currencies *balance.CurrencyList, validators *identity.ValidatorStore, logger *log.Logger) *Context {

	return &Context{
		Router:     r,
		Accounts:   wallet,
		Balances:   balances,
		Currencies: currencies,
		Validators: validators,
		Logger:     logger,
	}
}
