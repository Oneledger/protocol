package keys

import (
	"encoding/hex"
	"fmt"
)

//Address to be used as to reference a key-pair.
type Address []byte

func (a Address) String() string {
	return fmt.Sprint("0x", hex.EncodeToString(a))
}

func (a Address) Bytes() []byte {
	return a
}
