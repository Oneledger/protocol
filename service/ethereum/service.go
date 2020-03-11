package ethereum

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/app/node"
	chain "github.com/Oneledger/protocol/chains/ethereum"
	"github.com/Oneledger/protocol/config"
	"github.com/Oneledger/protocol/data/accounts"
	ethTracker "github.com/Oneledger/protocol/data/ethereum"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/identity"
	"github.com/Oneledger/protocol/log"
)

func Name() string {
	return "eth"
}

type Service struct {
	config       *config.EthereumChainDriverConfig
	router       action.Router
	accounts     accounts.Wallet
	logger       *log.Logger
	nodeContext  node.Context
	validators   *identity.ValidatorStore
	trackerStore *ethTracker.TrackerStore
}

// Returns a new Service, should be passed as an RPC handler
func NewService(
	//balances *balance.Store,
	config *config.EthereumChainDriverConfig,
	router action.Router,
	accounts accounts.Wallet,
	nodeCtx node.Context,
	validators *identity.ValidatorStore,
	trackerStore *ethTracker.TrackerStore,
	logger *log.Logger,
) *Service {
	return &Service{
		//balances:     balances,
		config:       config,
		router:       router,
		nodeContext:  nodeCtx,
		accounts:     accounts,
		validators:   validators,
		trackerStore: trackerStore,
		logger:       logger,
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

type OLTERC20LockRequest struct {
	RawTx   []byte `json:"rawTx"`
	Address keys.Address
	Fee     action.Amount `json:"fee"`
	Gas     int64         `json:"gas"`
}

type OLTReply struct {
	RawTX []byte `json:"rawTx"`
}

type RedeemRequest struct {
	UserOLTaddress action.Address `json:"userOLTAddress"`
	UserETHaddress action.Address `json:"userETHAddress"`
	ETHTxn         []byte         `json:"ethTxn"`
	Fee            action.Amount  `json:"fee"`
	Gas            int64          `json:"gas"`
}

type OLTERC20RedeemRequest struct {
	UserOLTaddress action.Address `json:"userOLTAddress"`
	UserETHaddress action.Address `json:"userETHAddress"`
	ETHTxn         []byte         `json:"ethTxn"`
	Fee            action.Amount  `json:"fee"`
	Gas            int64          `json:"gas"`
}

type ETHLockRequest struct {
	UserAddress common.Address `json:"userETHAddress"`
	Amount      *big.Int       `json:"amount"`
}

type ETHLockRawTX struct {
	UnsignedRawTx []byte `json:"unsignedRawTx"`
}

type SignRequest struct {
	Wei       *big.Int       `json:"wei"`
	Recipient common.Address `json:"recipient"`
}

type SignReply struct {
	TxHash common.Hash `json:"txTash"`
}

type BalanceRequest struct {
	Address chain.Address `json:"address"`
}

type BalanceReply struct {
	Address chain.Address `json:"address"`
	Amount  *big.Int      `json:"amount"`
}

type TrackerStatusRequest struct {
	TrackerName chain.TrackerName `json:"tracker_name"`
}

type TrackerStatusReply struct {
	Status string `json:"status"`
}
