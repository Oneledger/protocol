package ethereum

import (
	"math/big"
)

type RedeemRequest struct {
	Amount *big.Int
}

type LockErcRequest struct {
	Receiver string
  	TokenAmount *big.Int
}

type LockRequest struct {
	Amount *big.Int
}
