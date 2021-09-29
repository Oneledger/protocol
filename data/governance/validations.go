package governance

import (
	"math"
	"math/big"
	"reflect"

	"github.com/pkg/errors"

	"github.com/Oneledger/protocol/chains/bitcoin"
	ethchain "github.com/Oneledger/protocol/chains/ethereum"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/delegation"
	"github.com/Oneledger/protocol/data/evidence"
	"github.com/Oneledger/protocol/data/fees"
	"github.com/Oneledger/protocol/data/network_delegation"
	"github.com/Oneledger/protocol/data/ons"
	"github.com/Oneledger/protocol/data/rewards"
)

var (
	infiniteInt        = int64(math.MaxInt64)
	infiniteMaxBalance = balance.NewAmountFromInt(infiniteInt)
	//Proposal
	//Initial Funding
	initialFundingConfigMin  = balance.NewAmountFromInt(1)
	initialFundingConfigMax  = infiniteMaxBalance
	initialFundingCodeMin    = balance.NewAmountFromInt(1)
	initialFundingCodeMax    = infiniteMaxBalance
	initialFundingGeneralMin = balance.NewAmountFromInt(10000)
	initialFundingGeneralMax = infiniteMaxBalance
	//Funding for execution
	initialToFundingMultiplier = int64(3)
	//Deadlines in number of blocks
	//Config
	minDeadlineFundingConfig = getDeadline(10000)
	maxDeadlineFundingConfig = getDeadline(infiniteInt)
	minDeadlineVotingConfig  = getDeadline(10000)
	maxDeadlineVotingConfig  = getDeadline(infiniteInt)
	//Code Change
	minDeadlineFundingCode = getDeadline(10000)
	maxDeadlineFundingCode = getDeadline(infiniteInt)
	minDeadlineVotingCode  = getDeadline(150000)
	maxDeadlineVotingCode  = getDeadline(450000)
	//General
	minDeadlineFundingGeneral = getDeadline(75000)
	maxDeadlineFundingGeneral = getDeadline(150000)
	minDeadlineVotingeGeneral = getDeadline(75000)
	maxDeadlineVotingGeneral  = getDeadline(150000)
	minPassPercentage         = int64(51)
	maxPassPercentage         = int64(80)
	decimalsAllowed           = 2
	//ONS
	minPerBlockFee     = balance.NewAmountFromInt(1)
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
	minSelfDelegationAmount = balance.NewAmountFromInt(500_000)
	maxSelfDelegationAmount = balance.NewAmountFromInt(10_000_000)
	minValidatorCount       = int64(8)
	maxValidatorCount       = int64(64)
	minMaturityTime         = int64(109200)
	maxMaturityTime         = int64(468000)
	//Evidence
	MinVotesRequiredPercentage = int64(70)
	minBlockVotesDiff          = int64(1000)
	maxBlockVotesDiff          = int64(100000)
	minPenaltyBasePercentage   = int64(10)
	maxPenaltyBasePercentage   = int64(40)
	minValidatorVotePercentage = int64(50)
	maxValidatorVotePercentage = int64(100)
	// can be between 0 -100, PenaltyBurnPercentage + PenaltyBountyPercentage is always 100
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
	ok, err = st.ValidateEvidence(&govstate.EvidenceOptions)
	if err != nil || !ok {
		return false, err
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

func (st *Store) ValidateNetwkDeleg(opt *network_delegation.Options) (bool, error) {
	oldOptions, err := st.GetNetworkDelegOptions()
	if err != nil {
		return false, err
	}
	return reflect.DeepEqual(oldOptions, opt), nil
}

func (st *Store) ValidateETH(opt *ethchain.ChainDriverOption) (bool, error) {
	oldOptions, err := st.GetETHChainDriverOption()
	if err != nil {
		return false, err
	}
	return reflect.DeepEqual(oldOptions, opt), nil
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
		return false, errors.New("voting deadline for Config update is not within range")
	}
	if !verifyRangeInt64(code.FundingDeadline, minDeadlineFundingCode, maxDeadlineFundingCode) {
		return false, errors.New("funding deadline for code update is not within range")
	}
	if !verifyRangeInt64(code.VotingDeadline, minDeadlineVotingCode, maxDeadlineVotingCode) {
		return false, errors.New("voting deadline for code update is not within range")
	}
	if !verifyRangeInt64(general.FundingDeadline, minDeadlineFundingGeneral, maxDeadlineFundingGeneral) {
		return false, errors.New("funding deadline for general update is not within range")
	}
	if !verifyRangeInt64(general.VotingDeadline, minDeadlineVotingeGeneral, maxDeadlineVotingGeneral) {
		return false, errors.New("voting deadline for general update is not within range")
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

	ok = reflect.DeepEqual(config.PassedFundDistribution, oldOptions.ConfigUpdate.PassedFundDistribution)
	if !ok {
		return ok, errors.New("Proposal Distribution Cannot be changed ")
	}
	ok = reflect.DeepEqual(config.FailedFundDistribution, oldOptions.ConfigUpdate.FailedFundDistribution)
	if !ok {
		return ok, errors.New("Proposal Distribution Cannot be changed ")
	}
	ok = reflect.DeepEqual(code.PassedFundDistribution, oldOptions.CodeChange.PassedFundDistribution)
	if !ok {
		return ok, errors.New("Proposal Distribution Cannot be changed ")
	}
	ok = reflect.DeepEqual(code.FailedFundDistribution, oldOptions.CodeChange.FailedFundDistribution)
	if !ok {
		return ok, errors.New("Proposal Distribution Cannot be changed ")
	}
	ok = reflect.DeepEqual(general.PassedFundDistribution, oldOptions.General.PassedFundDistribution)
	if !ok {
		return ok, errors.New("Proposal Distribution Cannot be changed ")
	}
	ok = reflect.DeepEqual(general.FailedFundDistribution, oldOptions.General.FailedFundDistribution)
	if !ok {
		return ok, errors.New("Proposal Distribution Cannot be changed ")
	}

	if code.ProposalExecutionCost != oldOptions.CodeChange.ProposalExecutionCost {
		return false, errors.New("OneLedger Address cannot be changed | Code change proposal execution cost address")
	}
	if config.ProposalExecutionCost != oldOptions.ConfigUpdate.ProposalExecutionCost {
		return false, errors.New("OneLedger Address cannot be changed | Config update proposal execution cost address")
	}
	if general.ProposalExecutionCost != oldOptions.General.ProposalExecutionCost {
		return false, errors.New("OneLedger Address cannot be changed  | General proposal execution cost address")
	}
	if opt.BountyProgramAddr != oldOptions.BountyProgramAddr {
		return false, errors.New("OneLedger Address cannot be changed | Bounty Program address")
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

func (st *Store) ValidateEvidence(opt *evidence.Options) (bool, error) {
	oldOptions, err := st.GetEvidenceOptions()
	if err != nil {
		return false, err
	}
	ok := verifyRangeInt64(opt.BlockVotesDiff, minBlockVotesDiff, maxBlockVotesDiff)
	if !ok {
		return false, errors.New("Block Votes Diff is not in range")
	}
	if opt.MinVotesRequired < (MinVotesRequiredPercentage * (opt.BlockVotesDiff / 100)) {
		return false, errors.New("Min Required Votes cannot be less that 70% of BlockVotesDiff")
	}
	if opt.MinVotesRequired > opt.BlockVotesDiff {
		return false, errors.New("Min Required Votes cannot be more than BlockVotesDiff")
	}
	ok = verifyRangeInt64(opt.PenaltyBasePercentage/(opt.PenaltyBaseDecimals/100), minPenaltyBasePercentage, maxPenaltyBasePercentage)
	if !ok {
		return false, errors.New("PenaltyBasePercentage not in range")
	}
	ok = verifyRangeInt64(opt.ValidatorVotePercentage/(opt.ValidatorVoteDecimals/100), minValidatorVotePercentage, maxValidatorVotePercentage)
	if !ok {
		return false, errors.New("Validator Vote Percentage not in range")
	}

	if opt.PenaltyBountyPercentage != oldOptions.PenaltyBountyPercentage {
		return false, errors.New("Bounty percentage cannot be changed")
	}
	if opt.PenaltyBountyDecimals != oldOptions.PenaltyBountyDecimals {
		return false, errors.New("Bounty percentage cannot be changed")
	}
	if opt.PenaltyBurnPercentage != oldOptions.PenaltyBurnPercentage {
		return false, errors.New("Burn percentage cannot be changed")
	}
	if opt.PenaltyBurnDecimals != oldOptions.PenaltyBurnDecimals {
		return false, errors.New("Burn percentage cannot be changed")
	}
	if opt.ValidatorReleaseTime != oldOptions.ValidatorReleaseTime {
		return false, errors.New("Validator release time cannot be changed")
	}
	if opt.AllegationPercentage != oldOptions.AllegationPercentage {
		return false, errors.New("AllegationPercentage cannot be changed")
	}
	if opt.AllegationDecimals != oldOptions.AllegationDecimals {
		return false, errors.New("AllegationDecimals cannot be changed")
	}
	return true, nil
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
