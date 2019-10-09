package service

import (
	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/app/node"
	"github.com/Oneledger/protocol/client"
	"github.com/Oneledger/protocol/config"
	"github.com/Oneledger/protocol/data/accounts"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/fees"
	"github.com/Oneledger/protocol/data/ons"
	"github.com/Oneledger/protocol/identity"
	"github.com/Oneledger/protocol/log"
	"github.com/Oneledger/protocol/service/broadcast"
	nodesvc "github.com/Oneledger/protocol/service/node"
	"github.com/Oneledger/protocol/service/owner"
	"github.com/Oneledger/protocol/service/query"
	"github.com/Oneledger/protocol/service/tx"
)

// Context is the master context for creating new contexts
type Context struct {
	//stores
	Accounts     accounts.Wallet
	Balances     *balance.Store
	Domains      *ons.DomainStore
	ValidatorSet *identity.ValidatorStore

	// configurations
	Cfg         config.Server
	Currencies  *balance.CurrencySet
	FeeOpt      *fees.FeeOption
	NodeContext node.Context

	Router   action.Router
	Services client.ExtServiceContext
	Logger   *log.Logger
}

// Map of services, keyed by the name/prefix of the service
type Map map[string]interface{}

func NewMap(ctx *Context) Map {
	return Map{
		broadcast.Name(): broadcast.NewService(ctx.Services, ctx.Router, ctx.Currencies, ctx.FeeOpt, ctx.Logger),
		nodesvc.Name():   nodesvc.NewService(ctx.NodeContext, &ctx.Cfg, ctx.Logger),
		owner.Name():     owner.NewService(ctx.Accounts, ctx.Logger),
		query.Name():     query.NewService(ctx.Services, ctx.Balances, ctx.Currencies, ctx.ValidatorSet, ctx.Domains, ctx.Logger),
		tx.Name():        tx.NewService(ctx.Balances, ctx.Router, ctx.Accounts, ctx.FeeOpt, ctx.NodeContext, ctx.Logger),
	}
}
