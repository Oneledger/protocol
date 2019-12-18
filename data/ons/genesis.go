package ons

import (
	"github.com/Oneledger/protocol/data/balance"
)


type Options struct {
	PerBlockFees     int64
	FirstLevelDomain []string
	BaseDomainPrice  balance.Coin

}
