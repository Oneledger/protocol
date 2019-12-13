package ons

import (
	"github.com/Oneledger/protocol/data/balance"
)
type OnsOptions struct {
	PerBlockFees int64
	FirstLevelDomain string
	SubDomainPrice balance.Coin
}
