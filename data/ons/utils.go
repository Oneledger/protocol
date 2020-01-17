/*

 */

package ons

import (
	"math/big"

	"github.com/Oneledger/protocol/data/balance"
)

func CalculateDomainExpiry(amount balance.Amount, perBlockFees int64) uint64 {

	x := amount.BigInt()
	y := big.NewInt(perBlockFees)

	dividend := big.NewInt(0).Div(x, y)
	return dividend.Uint64()
}
