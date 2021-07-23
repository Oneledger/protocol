package keys

import (
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

type Algorithm int

const (
	UNKNOWN Algorithm = iota
	ED25519
	SECP256K1
	BTCECSECP

	ETHSECP

	ED25519_PUB_SIZE  int = ed25519.PubKeyEd25519Size
	ED25519_PRIV_SIZE int = 64

	SECP256K1_PUB_SIZE  int = secp256k1.PubKeySecp256k1Size
	SECP256K1_PRIV_SIZE int = 32
)

func (a Algorithm) String() string {
	switch a {
	case ETHSECP:
		return "ethsecp"
	case ED25519:
		return "ed25519"
	case SECP256K1:
		return "secp256k1"
	case BTCECSECP:
		return "btcecsecp"
	}
	return "Unknown algorithm"
}

func (a Algorithm) MarshalText() ([]byte, error) {
	switch a {
	case ETHSECP, ED25519, SECP256K1, BTCECSECP:
		return []byte(a.String()), nil
	case UNKNOWN:
		return []byte{}, nil
	default:
		return nil, errors.New("Unknown Algorithm specified")
	}
}

func (a *Algorithm) UnmarshalText(name []byte) error {
	algo := GetAlgorithmFromTmKeyName(string(name))
	if algo == UNKNOWN {
		return nil
	}
	*a = algo
	return nil
}

func GetAlgorithmFromTmKeyName(name string) Algorithm {
	switch name {
	case "ethsecp":
		return ETHSECP
	case "ed25519":
		return ED25519
	case "secp256k1":
		return SECP256K1
	case "btcecsecp":
		return BTCECSECP
	}

	return UNKNOWN
}
