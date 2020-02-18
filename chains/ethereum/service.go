package ethereum

import (
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type Service struct {
	access *ETHChainDriver
}

// Returns a new Service, should be passed as an RPC handler
func NewService(access *ETHChainDriver) *Service {
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
	PublicKey *ecdsa.PublicKey `json:"publicKey"`
	Amount    *big.Int         `json:"amount"`
}

type OfflineLockRawTX struct {
	UnsignedRawTx []byte `json:"unsignedRawTx"`
}
type LockReply struct {
	Amount *big.Int `json:"amount"`
	Ok     bool     `json:"ok"`
	//VerifyBalance     bool     `json:"ok"`
	//Reason string  `json:"reason"`
}

type SignRequest struct {
	wei       *big.Int       `json:"wei"`
	recipient common.Address `json:"recipient"`
}

type SignReply struct {
	txHash common.Hash `json:"txHash"`
}

type BalanceRequest struct {
	Address Address `json:"address"`
}

type BalanceReply struct {
	Address Address  `json:"address"`
	Amount  *big.Int `json:"amount"`
}
