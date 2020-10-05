package client

import (
	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/data/keys"
)

//-------------------TX-----------------
type FinalizeRewardsRequest struct {
	Delegator keys.Address   `json:"delegator"`
	Amount    action.Amount `json:"amount"`
	GasPrice  action.Amount  `json:"gasPrice"`
	Gas       int64          `json:"gas"`
}
