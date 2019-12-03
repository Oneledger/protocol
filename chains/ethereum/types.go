package ethereum

import (
	"math/big"
)

type RedeemRequest struct {
	Amount *big.Int
}

type LockRequest struct {
	Amount *big.Int
}
