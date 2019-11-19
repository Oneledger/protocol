package ethereum

import (
	"crypto/ecdsa"
	"math/big"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/app/node"
	"github.com/Oneledger/protocol/chains/ethereum"
	"github.com/Oneledger/protocol/config"
	"github.com/Oneledger/protocol/data/accounts"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/identity"
	"github.com/Oneledger/protocol/log"
	"github.com/ethereum/go-ethereum/common"
)

type Service struct {
	config      *config.EthereumChainDriverConfig
	router      action.Router
	accounts    accounts.Wallet
	logger      *log.Logger
	nodeContext node.Context
	validators  *identity.ValidatorStore
	// trackerStore *ethereum.TrackerStore
}

// Returns a new Service, should be passed as an RPC handler
func NewEthereumService(
	//balances *balance.Store,
	config config.EthereumChainDriverConfig,
	router action.Router,
	accounts accounts.Wallet,
	nodeCtx node.Context,
	validators *identity.ValidatorStore,
	//trackerStore *bitcoin.TrackerStore,
	logger *log.Logger,
) *Service {
	return &Service{
		//balances:     balances,
		config:      &config,
		router:      router,
		nodeContext: nodeCtx,
		accounts:    accounts,
		validators:  validators,
		//	trackerStore: trackerStore,
		logger: logger,
	}
}

type OLTLockRequest struct {
	// RawTransaction of a Lock call from the user to the smart contract
	// This should be signed and RLP encoded with the ethereum address of the user
	//OLTAddress common.Address  `json:"oltAddress"`
	RawTx   []byte `json:"rawTx"`
	Address keys.Address
	Fee     action.Amount `json:"fee"`
	Gas     int64         `json:"gas"`
}

type OLTLockReply struct {
	RawTX []byte `json:"UnsignedOLTLock"`
}

type ETHLockRequest struct {
	PublicKey *ecdsa.PublicKey `json:"public_key"`
	Amount    *big.Int         `json:"amount"`
}

type ETHLockRawTX struct {
	UnsignedRawTx []byte `json:"unsigned_raw_tx"`
}

type SignRequest struct {
	wei       *big.Int       `json:"wei"`
	recepient common.Address `json:"recepient"`
}

type SignReply struct {
	txHash common.Hash `json:"tx_hash"`
}

type BalanceRequest struct {
	Address ethereum.Address `json:"address"`
}

type BalanceReply struct {
	Address ethereum.Address `json:"address"`
	Amount  *big.Int         `json:"amount"`
}
