package service

import (
	"github.com/Oneledger/protocol/data"
	"github.com/pkg/errors"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/app/node"
	"github.com/Oneledger/protocol/client"
	"github.com/Oneledger/protocol/config"
	"github.com/Oneledger/protocol/data/accounts"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/bitcoin"
	"github.com/Oneledger/protocol/data/delegation"
	ethTracker "github.com/Oneledger/protocol/data/ethereum"
	"github.com/Oneledger/protocol/data/fees"
	"github.com/Oneledger/protocol/data/governance"
	"github.com/Oneledger/protocol/data/ons"
	"github.com/Oneledger/protocol/identity"
	"github.com/Oneledger/protocol/log"
	"github.com/Oneledger/protocol/service/broadcast"
	"github.com/Oneledger/protocol/service/btc"
	"github.com/Oneledger/protocol/service/ethereum"
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
	Govern       *governance.Store
	Delegators   *delegation.DelegationStore
	FeePool      *fees.Store
	ValidatorSet *identity.ValidatorStore
	WitnessSet   *identity.WitnessStore
	Trackers     *bitcoin.TrackerStore
	EthTrackers  *ethTracker.TrackerStore
	// configurations
	Cfg            config.Server
	Currencies     *balance.CurrencySet
	ProposalMaster *governance.ProposalMasterStore
	ExtStores      data.Router

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
		broadcast.Name(): broadcast.NewService(ctx.Services, ctx.Router, ctx.Currencies, ctx.FeePool, ctx.Domains, ctx.Govern, ctx.Delegators, ctx.ValidatorSet, ctx.Logger, ctx.Trackers, ctx.ProposalMaster, ctx.ExtStores),
		nodesvc.Name():   nodesvc.NewService(ctx.NodeContext, &ctx.Cfg, ctx.Logger),
		owner.Name():     owner.NewService(ctx.Accounts, ctx.Logger),
		query.Name(): query.NewService(ctx.Services, ctx.Balances, ctx.Currencies, ctx.ValidatorSet, ctx.WitnessSet, ctx.Domains, ctx.Delegators, ctx.Govern,
			ctx.FeePool, ctx.ProposalMaster, ctx.Logger, ctx.TxTypes),

		tx.Name():       tx.NewService(ctx.Balances, ctx.Router, ctx.Accounts, ctx.ValidatorSet, ctx.Govern, ctx.Delegators, ctx.FeePool.GetOpt(), ctx.NodeContext, ctx.Logger),
		btc.Name():      btc.NewService(ctx.Balances, ctx.Accounts, ctx.NodeContext, ctx.ValidatorSet, ctx.Trackers, ctx.Logger),
		ethereum.Name(): ethereum.NewService(ctx.Cfg.EthChainDriver, ctx.Router, ctx.Accounts, ctx.NodeContext, ctx.ValidatorSet, ctx.EthTrackers, ctx.Logger),
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
