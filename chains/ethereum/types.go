package ethereum

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type RedeemRequest struct {
	Amount *big.Int
}

type LockErcRequest struct {
	Receiver    common.Address
	TokenAmount *big.Int
}

type LockRequest struct {
	Amount *big.Int
}

type RedeemErcRequest struct {
	Amount *big.Int
	TokenAddress common.Address
}