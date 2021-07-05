package web3

import (
	"github.com/Oneledger/protocol/version"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

// PublicWeb3API is the web3_ prefixed set of APIs in the Web3 JSON-RPC spec.
type PublicWeb3API struct{}

// NewService creates an instance of the Web3 API.
func NewService() *PublicWeb3API {
	return &PublicWeb3API{}
}

// ClientVersion returns the client version in the Web3 user agent format.
func (svc PublicWeb3API) ClientVersion() string {
	return version.Client.String()
}

// Sha3 returns the keccak-256 hash of the passed-in input.
func (svc PublicWeb3API) Sha3(input hexutil.Bytes) hexutil.Bytes {
	return crypto.Keccak256(input)
}
