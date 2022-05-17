package utils

import (
	"crypto/ecdsa"
	"errors"
	"math/big"
	"os"
	"sort"
	"sync"

	"github.com/Oneledger/protocol/log"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"golang.org/x/crypto/sha3"
)

var logger = log.NewLoggerWithPrefix(os.Stdout, "utils")

func ToUncompressedSig(R, S, Vb *big.Int) []byte {
	// encode the signature in uncompressed format
	r, s := R.Bytes(), S.Bytes()
	V := byte(Vb.Uint64() - 27)
	sig := make([]byte, crypto.SignatureLength)
	copy(sig[32-len(r):32], r)
	copy(sig[64-len(s):64], s)
	sig[64] = V
	return sig
}

func RecoverPlain(sighash common.Hash, R, S, Vb *big.Int, homestead bool) (*ecdsa.PublicKey, error) {
	if Vb.BitLen() > 8 {
		return nil, types.ErrInvalidSig
	}
	V := byte(Vb.Uint64() - 27)
	if !crypto.ValidateSignatureValues(V, R, S, homestead) {
		return nil, types.ErrInvalidSig
	}
	sig := ToUncompressedSig(R, S, Vb)
	// recover the public key from the signature
	pub, err := crypto.Ecrecover(sighash[:], sig)
	if err != nil {
		return nil, err
	}
	if len(pub) == 0 || pub[0] != 4 {
		return nil, errors.New("invalid public key")
	}
	return crypto.UnmarshalPubkey(pub)
}

// hasherPool holds LegacyKeccak256 hashers for rlpHash.
var hasherPool = sync.Pool{
	New: func() interface{} { return sha3.NewLegacyKeccak256() },
}

// rlpHash encodes x and hashes the encoded bytes.
func RlpHash(x interface{}) (h common.Hash) {
	sha := hasherPool.Get().(crypto.KeccakState)
	defer hasherPool.Put(sha)
	sha.Reset()
	rlp.Encode(sha, x)
	sha.Read(h[:])
	return h
}

func PrintStringMap(dict map[string]interface{}, msg string, sorted bool) {
	keys := make([]string, 0, len(dict))
	for key := range dict {
		keys = append(keys, key)
	}
	if sorted {
		sort.Strings(keys)
	}
	for _, key := range keys {
		logger.Infof(msg, key, dict[key])
	}
}
