package keys

import (
	"bytes"
	"encoding"
	"encoding/hex"
	"strings"

	"github.com/pkg/errors"
)

const hexPrefix = "0x"

//Address to be used as to reference a key-pair.
type Address []byte

func (a Address) String() string {
	return hexPrefix + hex.EncodeToString(a)
}

func (a Address) Bytes() []byte {
	return a
}

func (a Address) Equal(b Address) bool {
	return bytes.Equal(a, b)
}

func (a Address) Humanize() string {
	return hexPrefix + strings.ToUpper(hex.EncodeToString(a))
}

var _ encoding.TextMarshaler = Address{}
var _ encoding.TextUnmarshaler = &Address{}

// MarshalText returns the text form for an Address. It returns a byteslice containing the hex encoding
// of the address excluding the prefix
func (a Address) MarshalText() ([]byte, error) {
	addrHex := hex.EncodeToString(a)
	return []byte(hexPrefix + addrHex), nil
}

// UnmarshalText decodes the given text in a byteslice of UTF8 characters. This implementation ignores hexPrefixprefixes
// and works regardless of whether the prefix is present or not
func (a *Address) UnmarshalText(text []byte) error {
	if a == nil {
		*a = Address{}
	}
	addrStr := string(text)
	// Cut off the hex prefix if it exists

	if strings.HasPrefix(addrStr, hexPrefix) {
		addrStr = addrStr[len(hexPrefix):]
	}

	addrRaw, err := hex.DecodeString(addrStr)
	if err != nil {
		return errors.Wrap(err, "address text unmarshal failed, not a hex address")
	}

	*a = addrRaw

	return nil
}
