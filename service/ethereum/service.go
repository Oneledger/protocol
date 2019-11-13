package ethereum

import (
	"crypto/ecdsa"
	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/app/node"
	"github.com/Oneledger/protocol/chains/ethereum"
<<<<<<< HEAD
	"github.com/Oneledger/protocol/config"
=======
>>>>>>> 98ea16d3f77a9a18800e62754b70d9ae27263893
	"github.com/Oneledger/protocol/data/accounts"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/identity"
	"github.com/Oneledger/protocol/log"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

type Service struct {
<<<<<<< HEAD
	config      *config.EthereumChainDriverConfig
=======
	config      *ethereum.Config
>>>>>>> 98ea16d3f77a9a18800e62754b70d9ae27263893
	router      action.Router
	accounts    accounts.Wallet
	logger      *log.Logger
	nodeContext node.Context
	validators  *identity.ValidatorStore
	// trackerStore *ethereum.TrackerStore
}

// Returns a new Service, should be passed as an RPC handler
<<<<<<< HEAD
func NewEthereumService(
//balances *balance.Store,
	config config.EthereumChainDriverConfig,
=======
func NewService(
//balances *balance.Store,
	config ethereum.Config,
>>>>>>> 98ea16d3f77a9a18800e62754b70d9ae27263893
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
	RawTx []byte `json:"rawTx"`
	Address     keys.Address
	Fee         action.Amount `json:"fee"`
	Gas         int64         `json:"gas"`
}

type OLTLockReply struct {
	RawTX []byte `json:"UnsignedOLTLock"`
}

type ETHLockRequest struct {
	PublicKey *ecdsa.PublicKey `json:"public_key"`
	Amount    *big.Int         `json:"amount"`
}

type LockRawTX struct {
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
