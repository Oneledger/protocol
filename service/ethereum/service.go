package ethereum

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/app/node"
	chain "github.com/Oneledger/protocol/chains/ethereum"
	"github.com/Oneledger/protocol/config"
	"github.com/Oneledger/protocol/data/accounts"
	tracker "github.com/Oneledger/protocol/data/ethereum"
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
	trackerStore *tracker.TrackerStore
}

// Returns a new Service, should be passed as an RPC handler
func NewService(
	//balances *balance.Store,
	config *config.EthereumChainDriverConfig,
	router action.Router,
	accounts accounts.Wallet,
	nodeCtx node.Context,
	validators *identity.ValidatorStore,
	//trackerStore *bitcoin.TrackerStore,
	logger *log.Logger,
) *Service {
	return &Service{
		//balances:     balances,
		config:      config,
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

type RedeemRequest struct {
	UserOLTaddress action.Address `json:"user_olt_address"`
	UserETHaddress action.Address `json:"user_eth_address"`
	ETHTxn         []byte         `json:"eth_txn"`
	Fee            action.Amount  `json:"fee"`
	Gas            int64          `json:"gas"`
}

type RedeemReply struct {
	OK bool `json:"ok"`
}
type ETHLockRequest struct {
	UserAddress common.Address `json:"user_eth_address"`
	Amount      *big.Int       `json:"amount"`
}

type ETHLockRawTX struct {
	UnsignedRawTx []byte `json:"unsigned_raw_tx"`
}

type SignRequest struct {
	Wei       *big.Int       `json:"wei"`
	Recepient common.Address `json:"recepient"`
}

type SignReply struct {
	TxHash common.Hash `json:"tx_hash"`
}

type BalanceRequest struct {
	Address chain.Address `json:"address"`
}

type BalanceReply struct {
	Address chain.Address `json:"address"`
	Amount  *big.Int      `json:"amount"`
}
