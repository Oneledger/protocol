package ethereum

import (
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

type Service struct {
	access *EthereumChainDriver
}

// Returns a new Service, should be passed as an RPC handler
func NewService(access *EthereumChainDriver) *Service {
	return &Service{access: access}
}

type OnlineLockRequest struct {
	// RawTransaction of a Lock call from the user to the smart contract
	// This should be signed and RLP encoded with the ethereum address of the user
	//OLTAddress common.Address  `json:"oltAddress"`
	RawTx []byte `json:"rawTx"`
	//Amount     int64 `json:"amount"`
}
type OfflineLockRequest struct {
	PublicKey *ecdsa.PublicKey `json:"public_key"`
	Amount    *big.Int         `json:"amount"`
}

type OfflineLockRawTX struct {
	UnsignedRawTx []byte `json:"unsigned_raw_tx"`
}
type LockReply struct {
	Amount *big.Int `json:"amount"`
	Ok     bool     `json:"ok"`
	//VerifyBalance     bool     `json:"ok"`
	//Reason string  `json:"reason"`
}

type SignRequest struct {
	wei       *big.Int       `json:"wei"`
	recepient common.Address `json:"recepient"`
}

type SignReply struct {
	txHash common.Hash `json:"tx_hash"`
}

type BalanceRequest struct {
	Address Address `json:"address"`
}

type BalanceReply struct {
	Address Address  `json:"address"`
	Amount  *big.Int `json:"amount"`
}
