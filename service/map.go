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
	ethTracker "github.com/Oneledger/protocol/data/ethereum"
	"github.com/Oneledger/protocol/data/fees"
	"github.com/Oneledger/protocol/data/ons"
	"github.com/Oneledger/protocol/data/passport"
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
	FeePool      *fees.Store
	AuthTokens   *passport.AuthTokenStore
	Tests        *passport.TestInfoStore
	ValidatorSet *identity.ValidatorStore
	WitnessSet   *identity.WitnessStore
	Trackers     *bitcoin.TrackerStore
	EthTrackers  *ethTracker.TrackerStore
	// configurations
	Cfg        config.Server
	Currencies *balance.CurrencySet

	NodeContext node.Context

	Router   action.Router
	Services client.ExtServiceContext
	Logger   *log.Logger

	TxTypes *[]action.TxTypeDescribe
}

// Map of services, keyed by the name/prefix of the service
type Map map[string]interface{}

func NewMap(ctx *Context) (Map, error) {

	defaultMap := Map{
		broadcast.Name(): broadcast.NewService(ctx.Services, ctx.Router, ctx.Currencies, ctx.FeePool, ctx.Domains, ctx.Logger, ctx.Trackers, ctx.AuthTokens),
		nodesvc.Name():   nodesvc.NewService(ctx.NodeContext, &ctx.Cfg, ctx.Logger),
		owner.Name():     owner.NewService(ctx.Accounts, ctx.Logger),
		query.Name():     query.NewService(ctx.Services, ctx.Balances, ctx.Currencies, ctx.Tests, ctx.AuthTokens, ctx.ValidatorSet, ctx.WitnessSet, ctx.Domains, ctx.FeePool, ctx.Logger, ctx.TxTypes),
		tx.Name():        tx.NewService(ctx.Balances, ctx.Router, ctx.Accounts, ctx.FeePool.GetOpt(), ctx.NodeContext, ctx.Logger),
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
