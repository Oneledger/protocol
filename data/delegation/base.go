package delegation

import "github.com/Oneledger/protocol/data/balance"

type Options struct {
	MinSelfDelegationAmount balance.Amount `json:"minSelfDelegationAmount"`
	MinDelegationAmount     balance.Amount `json:"minDelegationAmount"`
	TopValidatorCount       int64          `json:"topValidatorCount"`
	MaturityTime            int64          `json:"maturityTime"`
}
