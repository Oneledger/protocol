package service

import (
	"github.com/pkg/errors"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/app/node"
	"github.com/Oneledger/protocol/client"
	"github.com/Oneledger/protocol/config"
	"github.com/Oneledger/protocol/data/accounts"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/bitcoin"
	"github.com/Oneledger/protocol/data/fees"
	"github.com/Oneledger/protocol/data/ons"
	"github.com/Oneledger/protocol/identity"
	"github.com/Oneledger/protocol/log"
	"github.com/Oneledger/protocol/service/broadcast"
	"github.com/Oneledger/protocol/service/btc"
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
	Trackers     *bitcoin.TrackerStore

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

func NewMap(ctx *Context) (Map, error) {

	defaultMap := Map{
		broadcast.Name(): broadcast.NewService(ctx.Services, ctx.Router, ctx.Currencies, ctx.FeeOpt, ctx.Logger, ctx.Trackers),
		nodesvc.Name():   nodesvc.NewService(ctx.NodeContext, &ctx.Cfg, ctx.Logger),
		owner.Name():     owner.NewService(ctx.Accounts, ctx.Logger),
		query.Name():     query.NewService(ctx.Services, ctx.Balances, ctx.Currencies, ctx.ValidatorSet, ctx.Domains, ctx.Logger),
		tx.Name():        tx.NewService(ctx.Balances, ctx.Router, ctx.Accounts, ctx.FeeOpt, ctx.NodeContext, ctx.Logger),
		btc.Name(): btc.NewService(ctx.Balances, ctx.Accounts, ctx.NodeContext, ctx.ValidatorSet, ctx.Trackers, ctx.Logger,
			ctx.Cfg.ChainDriver.BlockCypherToken, ctx.Cfg.ChainDriver.BitcoinChainType),
	}

	serviceMap := Map{}
	for _, serviceName := range ctx.Cfg.Node.Services {
		if _, ok := defaultMap[serviceName]; ok {
			serviceMap[serviceName] = defaultMap[serviceName]
		} else {
			return serviceMap, errors.Wrap(errors.New("Service doesn't exist "), serviceName)
		}
	}

	return serviceMap, nil
}
