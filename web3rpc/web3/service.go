package web3

import (
	"github.com/Oneledger/protocol/version"
	"github.com/Oneledger/protocol/web3rpc"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

// ClientVersion returns the client version in the Web3 user agent format.
func (svc *web3rpc.Servic) ClientVersion() string {
	return version.Client.String()
}

// Sha3 returns the keccak-256 hash of the passed-in input.
func (svc *web3rpc.Service) Sha3(input hexutil.Bytes) hexutil.Bytes {
	return crypto.Keccak256(input)
}
