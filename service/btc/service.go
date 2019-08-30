/*

 */

package btc

import (
	"log"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/app/node"
	"github.com/Oneledger/protocol/data/accounts"
	"github.com/Oneledger/protocol/data/balance"
)

func Name() string {
	return "btc"
}

type Service struct {
	balances    *balance.Store
	router      action.Router
	accounts    accounts.Wallet
	logger      *log.Logger
	nodeContext node.Context
}

func NewService(
	balances *balance.Store,
	router action.Router,
	accounts accounts.Wallet,
	nodeCtx node.Context,
	logger *log.Logger,
) *Service {
	return &Service{
		balances:    balances,
		router:      router,
		nodeContext: nodeCtx,
		accounts:    accounts,
		logger:      logger,
	}
}
