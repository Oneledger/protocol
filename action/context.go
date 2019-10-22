package action

import (
	"github.com/Oneledger/protocol/data/accounts"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/fees"
	"github.com/Oneledger/protocol/data/ons"
	"github.com/Oneledger/protocol/identity"
	"github.com/Oneledger/protocol/log"
	"github.com/Oneledger/protocol/storage"
	abci "github.com/tendermint/tendermint/abci/types"
)

type Context struct {
	Router     Router
	State      *storage.State
	Header     *abci.Header
	Accounts   accounts.Wallet
	Balances   *balance.Store
	Domains    *ons.DomainStore
	FeePool    *fees.Store
	Currencies *balance.CurrencySet
	FeeOpt     *fees.FeeOption
	Validators *identity.ValidatorStore
	Logger     *log.Logger
}

func NewContext(r Router, header *abci.Header, state *storage.State,
	wallet accounts.Wallet, balances *balance.Store,
	currencies *balance.CurrencySet, feeOpt *fees.FeeOption, feePool *fees.Store,
	validators *identity.ValidatorStore, domains *ons.DomainStore,
	logger *log.Logger) *Context {
	return &Context{
		Router:     r,
		State:      state,
		Header:     header,
		Accounts:   wallet,
		Balances:   balances,
		Domains:    domains,
		FeePool:    feePool,
		Currencies: currencies,
		FeeOpt:     feeOpt,
		Validators: validators,
		Logger:     logger,
	}
}
