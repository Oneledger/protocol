package ethereum

import (
	"github.com/ethereum/go-ethereum/common"
)

type ChainDriverOption struct {
	ContractABI     string
	ContractAddress common.Address
}
