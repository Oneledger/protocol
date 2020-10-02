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
	PendingAmount      []SinglePendingAmount `json:"pendingAmount"`
	TotalPendingAmount balance.Amount        `json:"totalPendingAmount"`
	MaturedAmount      balance.Amount        `json:"maturedAmount"`
	Height             int64                 `json:"height"`
}

//-------------TX
type NetUndelegateRequest struct {
	Delegator keys.Address  `json:"delegator"`
	Amount    action.Amount `json:"amount"`
	GasPrice  action.Amount `json:"gasPrice"`
	Gas       int64         `json:"gas"`
}
