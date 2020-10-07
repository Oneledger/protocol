package client

import (
	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/data/keys"
)

type NetworkDelegateRequest struct {
	DelegationAddress keys.Address  `json:"delegationAddress"`
	Amount            action.Amount `json:"amount"`
	GasPrice          action.Amount `json:"gasPrice"`
	Gas               int64         `json:"gas"`
}

type ListDelegationRequest struct {
	DelegationAddress keys.Address `json:"delegationAddress"`
}

type ListDelegationReply struct {
	DelegationStats DelegationStats `json:"delegationStats"`
	Height          int64           `json:"height"`
}

type DelegationStats struct {
	Active  string `json:"active"`
	Pending string `json:"pending"`
	Matured string `json:"matured"`
}
