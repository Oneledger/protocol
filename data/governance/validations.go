package governance

import (
	"github.com/Oneledger/protocol/chains/bitcoin"
	ethchain "github.com/Oneledger/protocol/chains/ethereum"
	"github.com/Oneledger/protocol/data/delegation"
	"github.com/Oneledger/protocol/data/fees"
	"github.com/Oneledger/protocol/data/ons"
	"github.com/Oneledger/protocol/data/rewards"
)

func ValidateGov(govstate GovernanceState) bool {
	return ValidateONS(&govstate.ONSOptions) ||
		ValidateFee(&govstate.FeeOption) ||
		ValidateBTC(&govstate.BTCCDOption) ||
		ValidateETH(&govstate.ETHCDOption) ||
		ValidateStaking(&govstate.StakingOptions) ||
		ValidateProposal(&govstate.PropOptions) ||
		ValidateRewards(&govstate.RewardOptions)
}

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
func ValidateProposal(opt *ProposalOptionSet) bool {
	return true
}

func ValidateRewards(opt *rewards.Options) bool {
	return true
}
