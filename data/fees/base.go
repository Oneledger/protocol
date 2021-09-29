package fees

import (
	"math/big"
	"sync"

	"github.com/Oneledger/protocol/data/balance"
)

const (
	POOL_KEY        = "00000000000000000000"
	FEE_LOCK_BLOCKS = int64(3)
)

type FeeOption struct {
	FeeCurrency   balance.Currency `json:"feeCurrency"`
	MinFeeDecimal int64            `json:"minFeeDecimal"`

	minimalFee *balance.Coin
	rmu        sync.RWMutex
}

func (fo *FeeOption) MinFee() balance.Coin {
	// FIXME: temporary fix
	// ugly fix for concurrency to prevent non expected issues
	if fo.minimalFee == nil {
		fo.rmu.RLock()
		defer fo.rmu.RUnlock()
	}
	if fo.minimalFee == nil {
		amount := balance.Amount(*big.NewInt(0).Exp(big.NewInt(10), big.NewInt(fo.FeeCurrency.Decimal-fo.MinFeeDecimal), nil))
		coin := fo.FeeCurrency.NewCoinFromAmount(amount)
		fo.minimalFee = &coin
	}
	return *fo.minimalFee
}
