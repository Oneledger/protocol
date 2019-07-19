package action

import (
	"github.com/Oneledger/protocol/data/accounts"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/ons"
	"github.com/Oneledger/protocol/identity"
	"github.com/Oneledger/protocol/log"
	abci "github.com/tendermint/tendermint/abci/types"
)

type Context struct {
	Router     Router
	Header     *abci.Header
	Accounts   accounts.Wallet
	Balances   *balance.Store
	Domains    *ons.DomainStore
	Currencies *balance.CurrencyList
	Validators *identity.ValidatorStore
	Logger     *log.Logger
}

func NewContext(r Router, header *abci.Header, wallet accounts.Wallet, balances *balance.Store, currencies *balance.CurrencyList,
	validators *identity.ValidatorStore, domains *ons.DomainStore, logger *log.Logger) *Context {
	return &Context{
		Router:     r,
		Header:     header,
		Accounts:   wallet,
		Balances:   balances,
		Domains:    domains,
		Currencies: currencies,
		Validators: validators,
		Logger:     logger,
	}
}
