package common_rpc_query

import (
	"github.com/Oneledger/protocol/client"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/fees"
	"github.com/Oneledger/protocol/data/governance"
	"github.com/Oneledger/protocol/data/ons"
	"github.com/Oneledger/protocol/external_apps/common/common_data"
	"github.com/Oneledger/protocol/identity"
	"github.com/Oneledger/protocol/log"
)

type Service struct {
	ext            client.ExtServiceContext
	balances       *balance.Store
	currencies     *balance.CurrencySet
	validators     *identity.ValidatorStore
	ons            *ons.DomainStore
	feePool        *fees.Store
	governance     *governance.Store
	logger         *log.Logger
	extStores  	   common_data.Router
}

func Name() string {
	return "ext_query"
}

func NewService(ctx client.ExtServiceContext, balances *balance.Store, currencies *balance.CurrencySet, validators *identity.ValidatorStore,
	domains *ons.DomainStore, govern *governance.Store, feePool *fees.Store, logger *log.Logger, extStores common_data.Router) *Service {
	return &Service{
		ext:            ctx,
		currencies:     currencies,
		balances:       balances,
		validators:     validators,
		ons:            domains,
		feePool:        feePool,
		logger:         logger,
		governance:     govern,
		extStores:      extStores,
	}
}