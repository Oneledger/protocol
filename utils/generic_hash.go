/*
   ____             _              _                      _____           _                  _
  / __ \           | |            | |                    |  __ \         | |                | |
 | |  | |_ __   ___| |     ___  __| | __ _  ___ _ __     | |__) | __ ___ | |_ ___   ___ ___ | |
 | |  | | '_ \ / _ \ |    / _ \/ _` |/ _` |/ _ \ '__|    |  ___/ '__/ _ \| __/ _ \ / __/ _ \| |
 | |__| | | | |  __/ |___|  __/ (_| | (_| |  __/ |       | |   | | | (_) | || (_) | (_| (_) | |
  \____/|_| |_|\___|______\___|\__,_|\__, |\___|_|       |_|   |_|  \___/ \__\___/ \___\___/|_|
                                      __/ |
                                     |___/


Copyright 2017 - 2019 OneLedger
*/

package utils

import (
	"crypto/sha256"
	"hash/fnv"
	"math/big"

	"golang.org/x/crypto/ripemd160"
)

// Hash returns ripemd160 hash of the given input
func Hash(result []byte) []byte {
	hasher := ripemd160.New()
	hasher.Write(result)
	return hasher.Sum(nil)
}

func SHA2(data []byte) []byte {
	h := sha256.Sum256(data)
	return h[:]
}

func GetTransactionHash(tx []byte) []byte {
	return SHA2(tx)
}

// hashToBigInt used to convert mostly chain id which is a string
func HashToBigInt(s string) *big.Int {
	h := fnv.New64a()
	h.Write([]byte(s))
	return new(big.Int).SetUint64(h.Sum64())
}
