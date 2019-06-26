package keys

import (
	"bytes"
	"encoding"
	"encoding/hex"
	"strings"

	"github.com/Oneledger/protocol/utils"
	"github.com/pkg/errors"
)

//Address to be used as to reference a key-pair.
type Address []byte

func (a Address) String() string {
	return utils.PrefixHex(hex.EncodeToString(a))
}

func (a Address) Bytes() []byte {
	return a
}

func (a Address) Equal(b Address) bool {
	return bytes.Equal(a, b)
}

func (a Address) Humanize() string {
	return utils.PrefixHex(strings.ToLower(hex.EncodeToString(a)))
}

func (a Address) Err() error {
	switch {
	case len(a) == 0:
		return errors.New("address is empty")
	case len(a) != 20:
		return errors.New("address is the incorrect length: must be 20-bytes (40 characters after 0x prefix)")
	}
	return nil
}

var _ encoding.TextMarshaler = Address{}
var _ encoding.TextUnmarshaler = &Address{}

// MarshalText returns the text form for an Address. It returns a byteslice containing the hex encoding
// of the address including the 0x prefix
func (a Address) MarshalText() ([]byte, error) {
	addrHex := hex.EncodeToString(a)
	return []byte(utils.PrefixHex(addrHex)), nil
}

// UnmarshalText decodes the given text in a byteslice of characters,
// and works regardless of whether the 0x prefix is present or not.
func (a *Address) UnmarshalText(text []byte) error {
	if a == nil {
		*a = Address{}
	}
	addrStr := string(text)

	// Cut off the hex prefix if it exists
	if strings.HasPrefix(addrStr, utils.HexPrefix) {
		addrStr = addrStr[len(utils.HexPrefix):]
	}

	addrRaw, err := hex.DecodeString(addrStr)
	if err != nil {
		return errors.Wrap(err, "address text unmarshal failed, not a hex address")
	}

	*a = addrRaw

	return nil
}
