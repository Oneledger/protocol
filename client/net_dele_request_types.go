package client

import (
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/keys"
)

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
