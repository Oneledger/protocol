package web3

import (
	"sync"

	"github.com/Oneledger/protocol/log"
	"github.com/Oneledger/protocol/version"
	web3types "github.com/Oneledger/protocol/web3/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

type Service struct {
	ctx    web3types.Web3Context
	logger *log.Logger

	mu sync.Mutex
}

func NewService(ctx web3types.Web3Context) *Service {
	return &Service{ctx: ctx, logger: ctx.GetLogger()}
}

// ClientVersion returns the client version in the Web3 user agent format.
func (svc *Service) ClientVersion() string {
	return version.Client.String()
}

// Sha3 returns the keccak-256 hash of the passed-in input.
func (svc *Service) Sha3(input hexutil.Bytes) hexutil.Bytes {
	return crypto.Keccak256(input)
}
