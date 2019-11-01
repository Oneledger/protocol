package ethereum

import (
	"math/big"
)

type Service struct {
	access *Access
}

// Returns a new Service, should be passed as an RPC handler
func NewService(access *Access) *Service {
	return &Service{access: access}
}

type LockRequest struct {
	// RawTransaction of a Lock call from the user to the smart contract
	// This should be signed and RLP encoded with the ethereum address of the user
	OLTAddress string  `json:"oltAddress"`
	Amount     int64 `json:"amount"`
}
type LockReply struct {
	//Amount *big.Int `json:"amount"`
	VerifyBalance     bool     `json:"ok"`
	Reason string  `json:"reason"`
}

func (svc *Service) CheckLock(req LockRequest, out *LockReply) error {
	verifyBalance, err := svc.access.CheckLock(req.Amount,req.OLTAddress)
	if err != nil {
		return err
	}
	if !verifyBalance {
		*out = LockReply{
			verifyBalance,
			"Balance Calculations do not match",
		}
	}
	*out = LockReply{
			VerifyBalance:verifyBalance,
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
