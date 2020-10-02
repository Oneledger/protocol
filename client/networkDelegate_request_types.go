package client

import (
	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/data/keys"
)

type NetworkDelegateRequest struct {
	UserAddress       keys.Address  `json:"userAddress"`
	DelegationAddress keys.Address  `json:"delegationAddress"`
	Amount            action.Amount `json:"amount"`
	GasPrice          action.Amount `json:"gasPrice"`
	Gas               int64         `json:"gas"`
}
