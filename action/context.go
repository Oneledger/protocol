package action

import (
	"github.com/Oneledger/protocol/data/accounts"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/log"
)

type Context struct {
	Router     Router
	Accounts   accounts.Wallet
	Balances   *balance.Store
	Currencies *balance.CurrencyList
	Logger     *log.Logger
}

func NewContext(r Router, wallet accounts.Wallet, balances *balance.Store, currencies *balance.CurrencyList, logger *log.Logger) *Context {

	return &Context{
		Router:     r,
		Accounts:   wallet,
		Balances:   balances,
		Currencies: currencies,
		Logger:     logger,
	}
}

// enable sendTx
func (ctx *Context) EnableSend() *Context {

	err := ctx.Router.AddHandler(SEND, sendTx{})
	//todo: ignore the err for now because register the handler path the second doesn't case any problem
	if err != nil {

	}
	return ctx
}
