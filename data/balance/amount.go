package balance

import (
	"encoding/json"
	"math/big"
	"strconv"
	"strings"

	"github.com/Oneledger/protocol/utils"
	"github.com/pkg/errors"
)

// Amount represents an amount of a currency
type Amount struct {
	Int big.Int `json:"bint,string"`
}

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
		return nil, errors.New("failed to set amount from string:" + x + " with the given base:" + strconv.Itoa(base))
	}
	return &Amount{*out}, nil
}

func NewAmountFromInt(x int64) *Amount {
	return &Amount{*big.NewInt(x)}
}

func (a Amount) MarshalJSON() ([]byte, error) {
	v := a.Int.String()
	return json.Marshal(v)
}

func (a *Amount) UnmarshalJSON(b []byte) error {
	v := ""
	err := json.Unmarshal(b, &v)
	if err != nil {
		return err
	}
	i, ok := big.NewInt(0).SetString(v, 0)
	if !ok {
		return errors.New("failed to unmarshal amount" + v)
	}
	a.Int = *i
	return nil
}

func (a Amount) MarshalText() ([]byte, error) {
	return a.Int.MarshalText()
}

func (a *Amount) UnmarshalText(b []byte) error {
	return a.Int.UnmarshalText(b)
}
