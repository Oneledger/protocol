package client

import (
	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/keys"
)

//-------------query

type GetUndelegatedRequest struct {
	Delegator keys.Address `json:"delegator"`
}

type SinglePendingAmount struct {
	Amount       balance.Amount `json:"amount"`
	MatureHeight int64          `json:"matureHeight"`
}

type GetUndelegatedReply struct {
	PendingAmounts []SinglePendingAmount `json:"pendingAmount"`
	MaturedAmount  balance.Amount        `json:"maturedAmount"`
	TotalAmount    balance.Amount        `json:"totalAmount"`
	Height         int64                 `json:"height"`
}

type GetTotalNetwkDelgReply struct {
	Amount balance.Amount `json:"amount"`
}

//-------------TX
type NetUndelegateRequest struct {
	Delegator keys.Address  `json:"delegator"`
	Amount    action.Amount `json:"amount"`
	GasPrice  action.Amount `json:"gasPrice"`
	Gas       int64         `json:"gas"`
}
