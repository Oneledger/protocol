package governance

import (
	"bytes"
	"math"
	"math/big"
	"reflect"

	"github.com/pkg/errors"

	"github.com/Oneledger/protocol/chains/bitcoin"
	ethchain "github.com/Oneledger/protocol/chains/ethereum"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/delegation"
	"github.com/Oneledger/protocol/data/fees"
	"github.com/Oneledger/protocol/data/ons"
	"github.com/Oneledger/protocol/data/rewards"
)

var (
	blocksPerDay               = int64(5200)
	infiniteInt                = int64(math.MaxInt64)
	infiniteMaxBalance         = balance.NewAmountFromInt(infiniteInt)
	initialFundingConfigMin    = balance.NewAmountFromInt(1)
	initialFundingConfigMax    = infiniteMaxBalance
	initialFundingCodeMin      = balance.NewAmountFromInt(1)
	initialFundingCodeMax      = infiniteMaxBalance
	initialFundingGeneralMin   = balance.NewAmountFromInt(10000)
	initialFundingGeneralMax   = infiniteMaxBalance
	initialToFundingMultiplier = int64(3)
	minDeadlineFundingConfig   = getDeadline(10000)
	maxDeadlineFundingConfig   = getDeadline(infiniteInt)
	minDeadlineVotingConfig    = getDeadline(10000)
	maxDeadlineVotingConfig    = getDeadline(infiniteInt)
	minDeadlineFundingCode     = getDeadline(10000)
	maxDeadlineFundingCode     = getDeadline(infiniteInt)
	minDeadlineVotingCode      = getDeadline(150000)
	maxDeadlineVotingCode      = getDeadline(450000)
	minDeadlineFundingGeneral  = getDeadline(75000)
	maxDeadlineFundingGeneral  = getDeadline(150000)
	minDeadlineVotingeGeneral  = getDeadline(75000)
	maxDeadlineVotingGeneral   = getDeadline(150000)
	minPassPercentage          = int64(51)
	maxPassPercentage          = int64(67)
	decimalsAllowed            = 2
	//ONS
	minPerBlockFee     = balance.NewAmountFromInt(0)
	maxPerBlockFee     = infiniteMaxBalance
	minBaseDomainPrice = balance.NewAmountFromInt(0)
	maxBaseDomainPrice = infiniteMaxBalance
	//FEE
	minFeeDecimal = int64(0)
	maxFeeDecimal = int64(18)
	//ETH
	minBlockConfirmation = int64(0)
	maxBlockConfirmation = int64(50)
	//Staking
	minSelfDelegationAmount = balance.NewAmountFromInt(3000000)
	maxSelfDelegationAmount = balance.NewAmountFromInt(10000000)
	minValidatorCount       = int64(8)
	maxValidatorCount       = int64(64)
	minMaturityTime         = int64(109200)
	maxMaturityTime         = int64(234000)
)

func (st *Store) ValidateGov(govstate GovernanceState) (bool, error) {
	ok, err := st.ValidateONS(&govstate.ONSOptions)
	if err != nil || !ok {
		return false, err
	}
	ok, err = st.ValidateFee(&govstate.FeeOption)
	if err != nil || !ok {
		return false, err
	}
	ok, err = st.ValidateBTC(&govstate.BTCCDOption)
	if err != nil || !ok {
		return false, errors.New("BTC options cannot be changed")
	}
	ok, err = st.ValidateETH(&govstate.ETHCDOption)
	if err != nil || !ok {
		return false, err
	}
	ok, err = st.ValidateStaking(&govstate.StakingOptions)
	if err != nil || !ok {
		return false, err
	}
	ok, err = st.ValidateProposal(&govstate.PropOptions)
	if err != nil || !ok {
		return false, err
	}
	ok, err = st.ValidateRewards(&govstate.RewardOptions)
	if err != nil || !ok {
		return false, errors.New("reward options cannot be changed ")
	}
	return true, nil
}

func (st *Store) ValidateONS(opt *ons.Options) (bool, error) {
	oldOptions, err := st.GetONSOptions()
	if err != nil {
		return false, err
	}
	if oldOptions.Currency != opt.Currency {
		return false, errors.New("currency cannot be changed")
	}
	if !reflect.DeepEqual(oldOptions.FirstLevelDomains, opt.FirstLevelDomains) {
		return false, errors.New("first level domains cannot be changed")
	}
	ok, err := opt.PerBlockFees.CheckInRange(*minPerBlockFee, *maxPerBlockFee)
	if err != nil || !ok {
		return false, errors.Wrap(err, "Per Block Fee")
	}
	ok, err = opt.BaseDomainPrice.CheckInRange(*minBaseDomainPrice, *maxBaseDomainPrice)
	if err != nil || !ok {
		return false, errors.Wrap(err, "Base Domain Price")
	}

	return true, nil
}

func (st *Store) ValidateFee(opt *fees.FeeOption) (bool, error) {
	oldOptions, err := st.GetFeeOption()
	if err != nil {
		return false, err
	}
	if !reflect.DeepEqual(opt.FeeCurrency, oldOptions.FeeCurrency) {
		return false, errors.New("fee currency cannot be changed")
	}
	if !verifyRangeInt64(opt.MinFeeDecimal, minFeeDecimal, maxFeeDecimal) {
		return false, errors.New("fee Decimal should be between 0 and 18")
	}
	return true, nil
}

func (st *Store) ValidateStaking(opt *delegation.Options) (bool, error) {
	ok, err := opt.MinSelfDelegationAmount.CheckInRange(*minSelfDelegationAmount, *maxSelfDelegationAmount)
	if err != nil || !ok {
		return false, errors.Wrap(err, "Min Delegation Amount")
	}
	if !verifyRangeInt64(opt.TopValidatorCount, minValidatorCount, maxValidatorCount) {
		return false, errors.New("validator count not within range")
	}
	if !verifyRangeInt64(opt.MaturityTime, minMaturityTime, maxMaturityTime) {
		return false, errors.New("maturity time not with range")
	}

	return true, nil
}

func (st *Store) ValidateETH(opt *ethchain.ChainDriverOption) (bool, error) {
	oldOptions, err := st.GetETHChainDriverOption()
	if err != nil {
		return false, err
	}
	if oldOptions.ContractABI != opt.ContractABI {
		return false, errors.New("contract abi cannot be changed")
	}
	if oldOptions.ERCContractABI != opt.ContractABI {
		return false, errors.New("erc contract abi cannot be changed")
	}
	if oldOptions.TotalSupply != opt.TotalSupply {
		return false, errors.New("total supply cannot be changed")
	}
	if oldOptions.TotalSupplyAddr != opt.TotalSupplyAddr {
		return false, errors.New("total supply address cannot be changed")
	}
	if !bytes.Equal(opt.ContractAddress.Bytes(), oldOptions.ContractAddress.Bytes()) {
		return false, errors.New("smart contract address cannot be changed")
	}
	if !bytes.Equal(opt.ERCContractAddress.Bytes(), oldOptions.ERCContractAddress.Bytes()) {
		return false, errors.New("ERC20 smart contract address cannot be changed")
	}
	if !verifyRangeInt64(opt.BlockConfirmation, minBlockConfirmation, maxBlockConfirmation) {
		return false, errors.New("BlockConfirmations should be between 0 and 50")
	}
	return true, nil
}

func (st *Store) ValidateBTC(opt *bitcoin.ChainDriverOption) (bool, error) {
	oldOptions, err := st.GetBTCChainDriverOption()
	if err != nil {
		return false, errors.New("unable to get BTC options")
	}
	return reflect.DeepEqual(oldOptions, opt), nil
}
func (st *Store) ValidateProposal(opt *ProposalOptionSet) (bool, error) {
	config := opt.ConfigUpdate
	code := opt.CodeChange
	general := opt.General
	oldOptions, err := st.GetProposalOptions()
	// Verify Initial Funding
	ok, err := config.InitialFunding.CheckInRange(*initialFundingConfigMin, *initialFundingConfigMax)
	if err != nil || !ok {
		return false, err
	}
	ok, err = code.InitialFunding.CheckInRange(*initialFundingCodeMin, *initialFundingCodeMax)
	if err != nil || !ok {
		return false, err
	}
	ok, err = general.InitialFunding.CheckInRange(*initialFundingGeneralMin, *initialFundingGeneralMax)
	if err != nil || !ok {
		return false, err
	}
	// Verify Min Funding
	ok, err = verifyMinFunding(config.InitialFunding, opt.ConfigUpdate.FundingGoal)
	if err != nil || !ok {
		return false, err
	}
	ok, err = verifyMinFunding(code.InitialFunding, opt.CodeChange.FundingGoal)
	if err != nil || !ok {
		return false, err
	}
	ok, err = verifyMinFunding(general.InitialFunding, opt.General.FundingGoal)
	if err != nil || !ok {
		return false, err
	}
	// Verify Funding and Voting Deadlines
	if !verifyRangeInt64(config.FundingDeadline, minDeadlineFundingConfig, maxDeadlineFundingConfig) {
		return false, errors.New("funding deadline for Config update is not within range")
	}
	if !verifyRangeInt64(config.VotingDeadline, minDeadlineVotingConfig, maxDeadlineVotingConfig) {
		return false, errors.New("funding deadline for Config update is not within range")
	}
	if !verifyRangeInt64(code.FundingDeadline, minDeadlineFundingCode, maxDeadlineFundingCode) {
		return false, errors.New("funding deadline for code update is not within range")
	}
	if !verifyRangeInt64(code.VotingDeadline, minDeadlineVotingCode, maxDeadlineVotingCode) {
		return false, errors.New("funding deadline for code update is not within range")
	}
	if !verifyRangeInt64(general.FundingDeadline, minDeadlineFundingGeneral, maxDeadlineFundingGeneral) {
		return false, errors.New("funding deadline for general update is not within range")
	}
	if !verifyRangeInt64(general.VotingDeadline, minDeadlineVotingeGeneral, maxDeadlineVotingGeneral) {
		return false, errors.New("funding deadline for general update is not within range")
	}
	if !verifyRangeInt64(int64(config.PassPercentage), minPassPercentage, maxPassPercentage) {
		return false, errors.New("pass percentage for config update is not within range")
	}
	if !verifyRangeInt64(int64(code.PassPercentage), minPassPercentage, maxPassPercentage) {
		return false, errors.New("pass percentage for code update is not within range")
	}
	if !verifyRangeInt64(int64(general.PassPercentage), minPassPercentage, maxPassPercentage) {
		return false, errors.New("pass percentage for general update is not within range")
	}
	ok, err = verifyDistribution(config.PassedFundDistribution)
	if err != nil || !ok {
		return false, err
	}
	ok, err = verifyDistribution(config.FailedFundDistribution)
	if err != nil || !ok {
		return false, err
	}
	ok, err = verifyDistribution(code.PassedFundDistribution)
	if err != nil || !ok {
		return false, err
	}
	ok, err = verifyDistribution(code.FailedFundDistribution)
	if err != nil || !ok {
		return false, err
	}
	ok, err = verifyDistribution(general.PassedFundDistribution)
	if err != nil || !ok {
		return false, err
	}
	ok, err = verifyDistribution(general.FailedFundDistribution)
	if err != nil || !ok {
		return false, err
	}
	if code.ProposalExecutionCost != oldOptions.CodeChange.ProposalExecutionCost {
		return false, errors.New("OnLedger Address cannot me changed code update execution cost")
	}
	if config.ProposalExecutionCost != oldOptions.ConfigUpdate.ProposalExecutionCost {
		return false, errors.New("OnLedger Address cannot me changed config update execution cost")
	}
	if general.ProposalExecutionCost != oldOptions.General.ProposalExecutionCost {
		return false, errors.New("OnLedger Address cannot me changed general update execution cost")
	}
	if opt.BountyProgramAddr != oldOptions.BountyProgramAddr {
		return false, errors.New("OnLedger Address cannot me changed Bounty Program address")
	}
	return true, nil
}

func (st *Store) ValidateRewards(opt *rewards.Options) (bool, error) {
	oldOptions, err := st.GetRewardOptions()
	if err != nil {
		return false, err
	}
	return reflect.DeepEqual(oldOptions, opt), nil
}

func verifyMinFunding(intialFunding *balance.Amount, fundingGoal *balance.Amount) (bool, error) {
	minFundingGoalConfig := big.NewInt(0)
	minFundingGoalConfig = minFundingGoalConfig.Mul(intialFunding.BigInt(), big.NewInt(initialToFundingMultiplier))
	return fundingGoal.CheckInRange(*balance.NewAmountFromBigInt(minFundingGoalConfig), *infiniteMaxBalance)
}

func verifyRangeInt64(value int64, min int64, max int64) bool {
	return value >= min && value <= max
}

func verifyDistribution(distribution ProposalFundDistribution) (bool, error) {
	if !(checkDecimalPlaces(distribution.ExecutionCost) ||
		checkDecimalPlaces(distribution.Burn) || checkDecimalPlaces(distribution.BountyPool) ||
		checkDecimalPlaces(distribution.ProposerReward) || checkDecimalPlaces(distribution.Validators) || checkDecimalPlaces(distribution.Validators)) {
		return false, errors.New("only two decimal places allowed")
	}

	if distribution.Burn+distribution.BountyPool+distribution.ExecutionCost+distribution.ProposerReward+distribution.FeePool+distribution.Validators != 100.00 {
		return false, errors.New("sum of amounts does not equal 100")
	}
	return true, nil
}

func checkDecimalPlaces(value float64) bool {
	valuef := value * math.Pow(10.0, float64(decimalsAllowed))
	extra := valuef - float64(int(valuef))
	return extra == 0
}

func getDeadline(numberofblocks int64) int64 {
	return numberofblocks
}
