package common_rpc_tx

import (
	"github.com/Oneledger/protocol/app/node"
	"github.com/Oneledger/protocol/data/accounts"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/fees"
	"github.com/Oneledger/protocol/data/governance"
	"github.com/Oneledger/protocol/external_apps/common/common_action"
	"github.com/Oneledger/protocol/identity"
	"github.com/Oneledger/protocol/log"
)

func Name() string {
	return "ext_tx"
}

type Service struct {
	balances    *balance.Store
	router      common_action.Router
	accounts    accounts.Wallet
	validators  *identity.ValidatorStore
	govern      *governance.Store
	feeOpt      *fees.FeeOption
	logger      *log.Logger
	nodeContext node.Context
}

func NewService(
	balances *balance.Store,
	router common_action.Router,
	accounts accounts.Wallet,
	validators *identity.ValidatorStore,
	govern *governance.Store,
	feeOpt *fees.FeeOption,
	nodeCtx node.Context,
	logger *log.Logger,
) *Service {
	return &Service{
		balances:    balances,
		router:      router,
		nodeContext: nodeCtx,
		accounts:    accounts,
		validators:  validators,
		govern:      govern,
		feeOpt:      feeOpt,
		logger:      logger,
	}
}