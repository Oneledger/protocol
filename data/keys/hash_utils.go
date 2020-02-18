package keys

import (
	"crypto/sha256"
	"crypto/sha512"
	"hash"
)

const (
	SHA224 = "SHA224"
	SHA256 = "SHA256"
	SHA384 = "SHA384"
	SHA512 = "SHA512"

	TAGLEN  = 6
	SIG_LEN = 64
)

func getSigPrefix(sig []byte) []byte {
	if len(sig) >= TAGLEN {
		return sig[0:TAGLEN]
	}
	return sig
}

func getHash(hash string) hash.Hash {
	switch hash {
	case SHA224:
		return sha256.New224()
	case SHA256:
		return sha256.New()
	case SHA384:
		return sha512.New384()
	case SHA512:
		return sha512.New()
	default:
		return nil
	}
}

func PreHashRequired(sig []byte) (bool, hash.Hash) {
	prefix := getSigPrefix(sig)

	if len(prefix) == TAGLEN {
		prefixStr := string(prefix)
		hash := getHash(prefixStr)

		if hash != nil {
			if len(sig) > SIG_LEN {
				return true, hash
			}
		}
	}

	return false, nil
}
