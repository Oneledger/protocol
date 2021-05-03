package action

import (
	"encoding/binary"

	"github.com/tendermint/tendermint/libs/kv"
)

func UintTag(key string, value uint64) kv.Pair {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, value)
	tag := kv.Pair{
		Key:   []byte(key),
		Value: []byte(b),
	}
	return tag
}
