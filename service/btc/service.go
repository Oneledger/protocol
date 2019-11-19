/*

 */

package btc

import (
	"github.com/Oneledger/protocol/app/node"
	"github.com/Oneledger/protocol/data/accounts"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/bitcoin"
	"github.com/Oneledger/protocol/identity"
	"github.com/Oneledger/protocol/log"
)

func Name() string {
	return "btc"
}

type Service struct {
	balances *balance.Store
	//	router      action.Router
	accounts    accounts.Wallet
	logger      *log.Logger
	nodeContext node.Context

	validators   *identity.ValidatorStore
	trackerStore *bitcoin.TrackerStore

	blockCypherToken string
	btcChainType     string
}

func NewService(
	balances *balance.Store,
	//	router action.Router,
	accounts accounts.Wallet,
	nodeCtx node.Context,
	validators *identity.ValidatorStore,
	trackerStore *bitcoin.TrackerStore,
	logger *log.Logger,
	blockCypherToken, btcChainType string,
) *Service {

	chain := "test3"
	switch btcChainType {
	case "testnet3":
		chain = "test3"
	case "testnet":
		chain = "test"
	case "mainnet":
		chain = "main"
	}

	return &Service{
		balances: balances,
		//		router:       router,
		nodeContext:  nodeCtx,
		accounts:     accounts,
		validators:   validators,
		trackerStore: trackerStore,
		logger:       logger,

		blockCypherToken: blockCypherToken,
		btcChainType:     chain,
	}
}
