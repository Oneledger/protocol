package service

import (
	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/app/node"
	"github.com/Oneledger/protocol/client"
	"github.com/Oneledger/protocol/config"
	"github.com/Oneledger/protocol/data/accounts"
	"github.com/Oneledger/protocol/data/balance"
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
	Balances     *balance.Store
	Accounts     accounts.Wallet
	Currencies   *balance.CurrencyList
	Cfg          config.Server
	NodeContext  node.Context
	ValidatorSet *identity.ValidatorStore
	Services     client.ExtServiceContext
	Router       action.Router
	Logger       *log.Logger
}

// Map of services, keyed by the name/prefix of the service
type Map map[string]interface{}

func NewMap(ctx *Context) Map {
	return Map{
		broadcast.Name(): broadcast.NewService(ctx.Services, ctx.Logger),
		nodesvc.Name():   nodesvc.NewService(ctx.NodeContext, &ctx.Cfg, ctx.Logger),
		owner.Name():     owner.NewService(ctx.Accounts, ctx.Logger),
		query.Name():     query.NewService(ctx.Balances, ctx.Currencies, ctx.Logger),
		tx.Name():        tx.NewService(ctx.Balances, ctx.Router, ctx.Accounts, ctx.NodeContext, ctx.Logger),
	}
}
