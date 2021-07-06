package service

import (
	"github.com/Oneledger/protocol/data/network_delegation"
	"github.com/Oneledger/protocol/external_apps/common"
	"github.com/pkg/errors"

	"github.com/Oneledger/protocol/data"
	"github.com/Oneledger/protocol/data/governance"
	netwkDeleg "github.com/Oneledger/protocol/data/network_delegation"
	"github.com/Oneledger/protocol/data/rewards"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/app/node"
	"github.com/Oneledger/protocol/client"
	"github.com/Oneledger/protocol/config"
	"github.com/Oneledger/protocol/data/accounts"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/bitcoin"
	"github.com/Oneledger/protocol/data/delegation"
	ethTracker "github.com/Oneledger/protocol/data/ethereum"
	"github.com/Oneledger/protocol/data/evidence"
	"github.com/Oneledger/protocol/data/evm"
	"github.com/Oneledger/protocol/data/fees"
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
	Accounts        accounts.Wallet
	Balances        *balance.Store
	Domains         *ons.DomainStore
	Delegators      *delegation.DelegationStore
	NetwkDelegators *netwkDeleg.MasterStore
	EvidenceStore   *evidence.EvidenceStore
	FeePool         *fees.Store
	ValidatorSet    *identity.ValidatorStore
	WitnessSet      *identity.WitnessStore
	Trackers        *bitcoin.TrackerStore
	EthTrackers     *ethTracker.TrackerStore
	// configurations
	Cfg                   config.Server
	Currencies            *balance.CurrencySet
	ProposalMaster        *governance.ProposalMasterStore
	RewardMaster          *rewards.RewardMasterStore
	Govern                *governance.Store
	ExtStores             data.Router
	ExtServiceMap         common.ExtServiceMap
	GovUpdate             *action.GovernaceUpdateAndValidate
	NetwkDelegatorsMaster *network_delegation.MasterStore
	NodeContext           node.Context

	Router   action.Router
	Services client.ExtServiceContext
	Logger   *log.Logger

	TxTypes *[]action.TxTypeDescribe

	// evm
	Contracts     *evm.ContractStore
	AccountKeeper balance.AccountKeeper
	StateDB       *action.CommitStateDB
}

// Map of services, keyed by the name/prefix of the service
type Map map[string]interface{}

func NewMap(ctx *Context) (Map, error) {

	defaultMap := Map{
		broadcast.Name(): broadcast.NewService(ctx.Services, ctx.Router, ctx.Currencies, ctx.FeePool, ctx.Domains, ctx.Govern, ctx.Delegators, ctx.EvidenceStore, ctx.NetwkDelegators,
			ctx.ValidatorSet, ctx.Logger, ctx.Trackers, ctx.ProposalMaster, ctx.RewardMaster, ctx.ExtStores, ctx.GovUpdate, ctx.StateDB),
		nodesvc.Name(): nodesvc.NewService(ctx.NodeContext, &ctx.Cfg, ctx.Logger),
		owner.Name():   owner.NewService(ctx.Accounts, ctx.Logger),
		query.Name(): query.NewService(ctx.Services, ctx.Balances, ctx.Currencies, ctx.ValidatorSet, ctx.WitnessSet, ctx.Domains, ctx.Delegators, ctx.NetwkDelegators, ctx.EvidenceStore,
			ctx.Govern, ctx.FeePool, ctx.ProposalMaster, ctx.RewardMaster, ctx.Logger, ctx.TxTypes, ctx.Contracts, ctx.AccountKeeper),
		tx.Name():       tx.NewService(ctx.Balances, ctx.Router, ctx.Accounts, ctx.ValidatorSet, ctx.Govern, ctx.Delegators, ctx.EvidenceStore, ctx.FeePool.GetOpt(), ctx.NodeContext, ctx.Logger),
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
	for name, service := range ctx.ExtServiceMap {
		if _, ok := defaultMap[name]; ok {
			return serviceMap, errors.Wrap(errors.New("Error adding external service, conflict service exist: "), name)
		} else {
			serviceMap[name] = service
		}
	}

	return serviceMap, nil
}
