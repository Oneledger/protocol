package ethereum

import (
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

type Service struct {
	access *Access
}

// Returns a new Service, should be passed as an RPC handler
func NewService(access *Access) *Service {
	return &Service{access: access}
}

type OnlineLockRequest struct {
	// RawTransaction of a Lock call from the user to the smart contract
	// This should be signed and RLP encoded with the ethereum address of the user
	//OLTAddress common.Address  `json:"oltAddress"`
	RawTx      []byte   `json:"rawTx"`
	//Amount     int64 `json:"amount"`
}
type OfflineLockRequest struct {
	PublicKey *ecdsa.PublicKey `json:"public_key"`
	Amount *big.Int `json:"amount"`
}

type OfflineLockRawTX struct {
	UnsignedRawTx []byte `json:"unsigned_raw_tx"`
}
type LockReply struct {
	Amount *big.Int `json:"amount"`
	Ok bool `json:"ok"`
	//VerifyBalance     bool     `json:"ok"`
	//Reason string  `json:"reason"`
}

type SignRequest struct {
	wei *big.Int `json:"wei"`
	recepient common.Address `json:"recepient"`
}

type SignReply struct{
	txHash common.Hash `json:"tx_hash"`
}


func (svc *Service) OnlineLock(req OnlineLockRequest, out *LockReply) error {
	amount, err := svc.access.LockFromSignedTx(req.RawTx)
	if err != nil {
		return err
	}
	*out = LockReply{
		Amount: amount,
		Ok:     true,
	}
	return nil
}


func (svc *Service) OfflineLock (req OfflineLockRequest,out *OfflineLockRawTX) error {
	rawTx,err := svc.access.GetRawLockTX(req.PublicKey,req.Amount)
	if err != nil {
		return err
	}
	*out = OfflineLockRawTX{UnsignedRawTx:rawTx}
	return nil
}


func (svc *Service) Sign(req SignRequest,out *SignReply) error {
	tx,err := svc.access.Sign(req.wei,req.recepient)
	if err!= nil {
		return err
	}
	*out = SignReply{
		txHash:tx.Hash(),
	}
	return nil
}



type BalanceRequest struct {
	Address Address `json:"address"`
}

type BalanceReply struct {
	Address Address  `json:"address"`
	Amount  *big.Int `json:"amount"`
}

// Balance returns the balance of the requested address
//func (svc *Service) Balance(req BalanceRequest, out *BalanceReply) error {
//	amount, err := svc.access.Contract.getTotalEthBalance(nil, req.Address)
//	if err != nil {
//		return err
//	}
//	*out = BalanceReply{
//		Address: req.Address,
//		Amount:  amount,
//	}
//	return nil
//}

//func (svc *Service) ETHBalance(req BalanceRequest, out *BalanceReply) error {
//
//}
