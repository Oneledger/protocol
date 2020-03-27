package ethereum

import (
	"github.com/ethereum/go-ethereum/common"
)

type ChainDriverOption struct {
	ContractABI     string
	ContractAddress common.Address
	TokenList       []ERC20Token
	ERCContractABI  string
	//Not needed ,same as to address in Lock Transaction
	ERCContractAddress common.Address
	TotalSupply        string
	TotalSupplyAddr    string
	BlockConfirmation  int64
}

type ERC20Token struct {
	TokName        string
	TokAddr        common.Address
	TokAbi         string
	TokTotalSupply string
}

type ChainDriver interface {
}
