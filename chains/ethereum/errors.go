package ethereum

import (
	"errors"
)

var (
	ErrTxFailed      = errors.New("ethereum tx status failed")
	ErrRedeemExpired = errors.New("ethereum redeem request expired")
)
