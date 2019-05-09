package action

import (
	"github.com/Oneledger/protocol/data/accounts"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/log"

)

type Context struct {
	Router   Router
	Accounts accounts.Wallet
	Balances *balance.Store
	Currencies map[string]balance.Currency
	Logger *log.Logger
}

func NewContext(r Router, wallet accounts.Wallet, balances *balance.Store, currencies map[string]balance.Currency,logger *log.Logger) *Context{

	return &Context{
		Router: r,
		Accounts: wallet,
		Balances: balances,
		Currencies: currencies,
		Logger: logger,
	}
}

// enable sendTx
func (ctx *Context) EnableSend() *Context{

	err := ctx.Router.AddHandler(SEND, sendTx{})
	if err != nil {
		ctx.Logger.Warn("error enable send", err)
	}
	return ctx
}

