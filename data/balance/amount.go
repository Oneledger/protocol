package balance

import (
	"math/big"
	"strings"

	"github.com/Oneledger/protocol/utils"
	"github.com/pkg/errors"
)

// Amount represents an amount of a currency
type Amount = big.Int

func NewAmount(x int64) *Amount {
	return NewAmountFromInt(x)
}

// NewAmountFromString parses the amount as a string with the given base. For example, if base is 10, then it expects
// the given string to be base-10 notation. If the base is 16, then it expects the string in hexadecimal notation. If
// the base is set to 16, it ignores any "0x" prefix if present.
func NewAmountFromString(x string, base int) (*Amount, error) {
	// Hex prefix handling
	if base == 16 && strings.HasPrefix(x, utils.HexPrefix) {
		x = x[len(utils.HexPrefix):]
	}
	out, ok := big.NewInt(0).SetString(x, base)
	if !ok {
		return nil, errors.New("failed to set amount from string with the given base")
	}
	return out, nil
}

func NewAmountFromInt(x int64) *Amount {
	return big.NewInt(x)
}
