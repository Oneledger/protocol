package eth

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
	OLTAddress Address  `json:"oltAddress"`
	RawTx      []byte   `json:"rawTx"`
	Amount     *big.Int `json:"amount"`
}
type LockReply struct {
	Amount *big.Int `json:"amount"`
	Ok     bool     `json:"ok"`
	Reason *string  `json:"reason"`
}

func (svc *Service) Lock(req LockRequest, out *LockReply) error {
	amount, err := svc.access.LockFromSignedTx(req.OLTAddress, req.RawTx)
	if err != nil {
		return err
	}
	*out = LockReply{
		Amount: amount,
		Ok:     true,
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
func (svc *Service) Balance(req BalanceRequest, out *BalanceReply) error {
	amount, err := svc.access.Contract.Balances(nil, req.Address)
	if err != nil {
		return err
	}
	*out = BalanceReply{
		Address: req.Address,
		Amount:  amount,
	}
	return nil
}

func (svc *Service) ETHBalance(req BalanceRequest, out *BalanceReply) error {

}
