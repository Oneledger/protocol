package key

import (
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

type Algorithm int

const (
	ED25519 Algorithm = iota
	SECP256K1

	ED25519_PUB_SIZE   int = ed25519.PubKeyEd25519Size
	SECP256K1_PUB_SIZE int = secp256k1.PubKeySecp256k1Size
)
