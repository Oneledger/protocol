package governance

import (
	"github.com/Oneledger/protocol/chains/bitcoin"
	ethchain "github.com/Oneledger/protocol/chains/ethereum"
	"github.com/Oneledger/protocol/data/delegation"
	"github.com/Oneledger/protocol/data/fees"
	"github.com/Oneledger/protocol/data/ons"
	"github.com/Oneledger/protocol/data/rewards"
)

func ValidateONS(opt *ons.Options) bool {
	return true
}

func ValidateFee(opt *fees.FeeOption) bool {
	return true
}

func ValidateStaking(opt *delegation.Options) bool {
	return true
}

func ValidateETH(opt *ethchain.ChainDriverOption) bool {
	return true
}

func ValidateBTC(opt *bitcoin.ChainDriverOption) bool {
	return true
}
func ValidateProposal(opt *ProposalOption) bool {
	return true
}

func RewardOption(opt *rewards.Options) bool {
	return true
}
