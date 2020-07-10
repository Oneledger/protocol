package governance

import (
	"bytes"
	"reflect"

	"github.com/Oneledger/protocol/chains/bitcoin"
	ethchain "github.com/Oneledger/protocol/chains/ethereum"
	"github.com/Oneledger/protocol/data/delegation"
	"github.com/Oneledger/protocol/data/fees"
	"github.com/Oneledger/protocol/data/ons"
	"github.com/Oneledger/protocol/data/rewards"
)

func (st *Store) ValidateGov(govstate GovernanceState) bool {
	return st.ValidateONS(&govstate.ONSOptions) ||
		st.ValidateFee(&govstate.FeeOption) ||
		st.ValidateBTC(&govstate.BTCCDOption) ||
		st.ValidateETH(&govstate.ETHCDOption) ||
		st.ValidateStaking(&govstate.StakingOptions) ||
		st.ValidateProposal(&govstate.PropOptions) ||
		st.ValidateRewards(&govstate.RewardOptions)
}

func (st *Store) ValidateONS(opt *ons.Options) bool {
	oldOptions, err := st.GetONSOptions()
	if err != nil {
		return false
	}
	if oldOptions.Currency != opt.Currency {
		return false
	}
	if !reflect.DeepEqual(oldOptions.FirstLevelDomains, opt.FirstLevelDomains) {
		return false
	}
	if opt.BaseDomainPrice.BigInt().Int64() < 0 {
		return false
	}

	return true
}

func (st *Store) ValidateFee(opt *fees.FeeOption) bool {
	oldOptions, err := st.GetFeeOption()
	if err != nil {
		return false
	}
	if !reflect.DeepEqual(opt.FeeCurrency, oldOptions.FeeCurrency) {
		return false
	}
	if opt.MinFeeDecimal < 0 || opt.MinFeeDecimal > 18 {
		return false
	}
	return true
}

func (st *Store) ValidateStaking(opt *delegation.Options) bool {
	return true
}

func (st *Store) ValidateETH(opt *ethchain.ChainDriverOption) bool {
	oldOptions, err := st.GetETHChainDriverOption()
	if err != nil {
		return false
	}
	if oldOptions.ContractABI != opt.ContractABI {
		return false
	}
	if oldOptions.ERCContractABI != opt.ContractABI {
		return false
	}
	if oldOptions.TotalSupply != opt.TotalSupply {
		return false
	}
	if oldOptions.TotalSupplyAddr != opt.TotalSupplyAddr {
		return false
	}
	if !bytes.Equal(opt.ContractAddress.Bytes(), oldOptions.ContractAddress.Bytes()) {
		return false
	}
	if !bytes.Equal(opt.ERCContractAddress.Bytes(), oldOptions.ERCContractAddress.Bytes()) {
		return false
	}
	if opt.BlockConfirmation < 0 || opt.BlockConfirmation > 50 {
		return false
	}
	return true
}

func (st *Store) ValidateBTC(opt *bitcoin.ChainDriverOption) bool {
	oldOptions, err := st.GetBTCChainDriverOption()
	if err != nil {
		return false
	}
	return reflect.DeepEqual(oldOptions, opt)
}
func (st *Store) ValidateProposal(opt *ProposalOptionSet) bool {
	return true
}

func (st *Store) ValidateRewards(opt *rewards.Options) bool {
	return true
}
