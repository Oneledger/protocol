package ethereum

import (
	"github.com/ethereum/go-ethereum/common"
)

type ChainDriverOption struct {
	ContractABI      string
	TestTokenABI     string
	ERCContractABI   string
	ContractAddress  common.Address
	TestTokenAddress common.Address // Not needed ,same as to address in Lock Transaction
 	ERCContractAddress  common.Address
}
