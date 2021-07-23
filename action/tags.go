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

func BoolTag(key string, value bool) kv.Pair {
	t := make([]bool, 1)
	t[0] = value
	b := make([]byte, (len(t)+7)/8)
	for i, x := range t {
		if x {
			b[i/8] |= 0x80 >> uint(i%8)
		}
	}
	tag := kv.Pair{
		Key:   []byte(key),
		Value: []byte(b),
	}
	return tag
}
