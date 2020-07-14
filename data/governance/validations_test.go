package governance

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	db "github.com/tendermint/tm-db"

	"github.com/Oneledger/protocol/chains/bitcoin"
	"github.com/Oneledger/protocol/chains/ethereum"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/chain"
	"github.com/Oneledger/protocol/data/delegation"
	"github.com/Oneledger/protocol/data/fees"
	"github.com/Oneledger/protocol/data/ons"
	"github.com/Oneledger/protocol/data/rewards"
	"github.com/Oneledger/protocol/storage"
)

var (
	proposalInitialFunding, _   = balance.NewAmountFromString("1000000000", 10)
	proposalFundingGoal, _      = balance.NewAmountFromString("10000000000", 10)
	proposalFundingDeadline     = int64(36401)
	proposalVotingDeadline      = int64(36401)
	proposalPassPercentage      = 51
	bountyProgramAddr           = "oneledgerBountyProgram"
	executionCostAddrConfig     = "executionCostConfig"
	executionCostAddrCodeChange = "executionCostCodeChange"
	executionCostAddrGeneral    = "executionCostGeneral"
	passedProposalDistribution  = ProposalFundDistribution{
		Validators:     18.00,
		FeePool:        18.00,
		Burn:           18.00,
		ExecutionCost:  18.00,
		BountyPool:     10.00,
		ProposerReward: 18.00,
	}
	failedProposalDistribution = ProposalFundDistribution{
		Validators:     10.00,
		FeePool:        10.00,
		Burn:           10.00,
		ExecutionCost:  20.00,
		BountyPool:     50.00,
		ProposerReward: 00.00,
	}

	blockRewardsCalcInterval    = int64(1000)
	annualBlockRewardsSupply, _ = balance.NewAmountFromString("50000000000000000000000000", 10)
	yearsOfBlockRewardsSupply   = int64(5)
	vStore                      *Store
)

func init() {
	fmt.Println("####### TESTING VALIDATIONS #######")

	//Create Test DB
	newDB := db.NewDB("test", db.MemDBBackend, "")
	cs := storage.NewState(storage.NewChainState("chainstate", newDB))

	//Create Governance store
	vStore = NewStore("g", cs)
}

func TestStore_ValidateProposal(t *testing.T) {
	//Tests for Initial Funding
	updates := generateGov()
	updates.PropOptions.ConfigUpdate.InitialFunding = balance.NewAmountFromInt(1000000)
	ok, err := vStore.ValidateProposal(&updates.PropOptions)
	assert.Error(t, err, "Initial Funding too high")
	assert.False(t, ok)
	updates.PropOptions.ConfigUpdate.InitialFunding = balance.NewAmountFromInt(0)
	ok, err = vStore.ValidateProposal(&updates.PropOptions)
	assert.Error(t, err, "Initial Funding too low")
	assert.False(t, ok)
	updates.PropOptions.ConfigUpdate.InitialFunding = balance.NewAmountFromInt(5000)
	ok, err = vStore.ValidateProposal(&updates.PropOptions)
	assert.NoError(t, err, "Should Pass within Range ConfigUpdate")
	assert.True(t, ok)
	updates.PropOptions.CodeChange.InitialFunding = balance.NewAmountFromInt(1000000000)
	ok, err = vStore.ValidateProposal(&updates.PropOptions)
	assert.NoError(t, err, "Should Pass max is IntMax CodeChange")
	assert.True(t, ok)
	//Tests for Funding Goal
	updates.PropOptions.ConfigUpdate.InitialFunding = balance.NewAmountFromInt(5000)
	updates.PropOptions.ConfigUpdate.FundingGoal = balance.NewAmountFromInt(5000)
	ok, err = vStore.ValidateProposal(&updates.PropOptions)
	assert.Error(t, err, "Funding goal is not 3 times initial funding ")
	assert.False(t, ok)
	updates.PropOptions.ConfigUpdate.InitialFunding = balance.NewAmountFromInt(5000)
	updates.PropOptions.ConfigUpdate.FundingGoal = balance.NewAmountFromInt(15000)
	ok, err = vStore.ValidateProposal(&updates.PropOptions)
	assert.NoError(t, err, "Should Pass Funding Goal is 3 times initial funding")
	assert.True(t, ok)
	//Tests for Deadline
	updates.PropOptions.ConfigUpdate.FundingDeadline = int64(156001)
	ok, err = vStore.ValidateProposal(&updates.PropOptions)
	assert.Error(t, err, "Funding Deadline is too high")
	assert.False(t, ok)

	updates.PropOptions.ConfigUpdate.VotingDeadline = int64(15)
	ok, err = vStore.ValidateProposal(&updates.PropOptions)
	assert.Error(t, err, "Voting Deadline is too low")
	assert.False(t, ok)

	updates.PropOptions.ConfigUpdate.VotingDeadline = int64(36401)
	updates.PropOptions.ConfigUpdate.FundingDeadline = int64(36401)
	ok, err = vStore.ValidateProposal(&updates.PropOptions)
	assert.NoError(t, err, "Should pass Voting and Funding Deadline")
	assert.True(t, ok)
	//Tests for pass percentage
	updates.PropOptions.ConfigUpdate.PassPercentage = 50
	ok, err = vStore.ValidateProposal(&updates.PropOptions)
	assert.Error(t, err, "Pass percentage too low")
	assert.False(t, ok)
	updates.PropOptions.ConfigUpdate.PassPercentage = 70
	ok, err = vStore.ValidateProposal(&updates.PropOptions)
	assert.Error(t, err, "Pass percentage too high")
	assert.False(t, ok)
	updates.PropOptions.ConfigUpdate.PassPercentage = 55
	ok, err = vStore.ValidateProposal(&updates.PropOptions)
	assert.NoError(t, err, "Should pass /Pass percentage")
	assert.True(t, ok)

	// Test distribution
	updates.PropOptions.ConfigUpdate.PassedFundDistribution.FeePool = 0
	ok, err = vStore.ValidateProposal(&updates.PropOptions)
	assert.Error(t, err, "Distribution is not 100%")
	assert.False(t, ok)

	updates.PropOptions.ConfigUpdate.PassedFundDistribution.FeePool = 18.000000001
	updates.PropOptions.ConfigUpdate.PassedFundDistribution.Validators = 18.000000001
	updates.PropOptions.ConfigUpdate.PassedFundDistribution.BountyPool = 10.000000001
	updates.PropOptions.ConfigUpdate.PassedFundDistribution.Burn = 18.0000000001
	updates.PropOptions.ConfigUpdate.PassedFundDistribution.ProposerReward = 18.000000001
	updates.PropOptions.ConfigUpdate.PassedFundDistribution.ExecutionCost = 17.9999999959
	ok, err = vStore.ValidateProposal(&updates.PropOptions)

	assert.Error(t, err, "Failing , Because of decimal conversion")
	assert.False(t, ok)
	updates.PropOptions.ConfigUpdate.PassedFundDistribution.FeePool = 18.00000000
	updates.PropOptions.ConfigUpdate.PassedFundDistribution.Validators = 18.00000000
	updates.PropOptions.ConfigUpdate.PassedFundDistribution.BountyPool = 10.00000000
	updates.PropOptions.ConfigUpdate.PassedFundDistribution.Burn = 18.000000000
	updates.PropOptions.ConfigUpdate.PassedFundDistribution.ProposerReward = 18.00000000
	updates.PropOptions.ConfigUpdate.PassedFundDistribution.ExecutionCost = 18.00000000
	ok, err = vStore.ValidateProposal(&updates.PropOptions)
	assert.NoError(t, err, "Should pass , decimals rounded off ")
	assert.True(t, ok)
}
func TestStore_ValidateStaking(t *testing.T) {
	//Tests for Initial Funding
	updates := generateGov()
	updates.StakingOptions.MaturityTime = 2340000
	ok, err := vStore.ValidateStaking(&updates.StakingOptions)
	assert.Error(t, err, "Maturity time too high")
	assert.False(t, ok)
	updates.StakingOptions.MaturityTime = 2340
	ok, err = vStore.ValidateStaking(&updates.StakingOptions)
	assert.Error(t, err, "Maturity time too low")
	assert.False(t, ok)
	updates.StakingOptions.MaturityTime = 109200
	ok, err = vStore.ValidateStaking(&updates.StakingOptions)
	assert.NoError(t, err, "Should Pass")
	assert.True(t, ok)
	updates.StakingOptions.MinDelegationAmount = *balance.NewAmountFromInt(100000)
	ok, err = vStore.ValidateStaking(&updates.StakingOptions)
	assert.Error(t, err, "Delegation amount too less")
	assert.False(t, ok)
	updates.StakingOptions.MinDelegationAmount = *balance.NewAmountFromInt(1000000)
	ok, err = vStore.ValidateStaking(&updates.StakingOptions)
	assert.NoError(t, err, "Should Pass")
	assert.True(t, ok)
	updates.StakingOptions.TopValidatorCount = 2
	ok, err = vStore.ValidateStaking(&updates.StakingOptions)
	assert.Error(t, err, "Validator count is too low")
	assert.False(t, ok)
	updates.StakingOptions.TopValidatorCount = 10
	ok, err = vStore.ValidateStaking(&updates.StakingOptions)
	assert.NoError(t, err, "Should Pass")
	assert.True(t, ok)
	updates = generateGov()
	ok, err = vStore.ValidateStaking(&updates.StakingOptions)
	assert.NoError(t, err, "Should Pass")
	assert.True(t, ok)
}

func TestStore_ValidateONS(t *testing.T) {
	updates := generateGov()
	updates.ONSOptions.Currency = "Test"
	ok, err := vStore.ValidateONS(&updates.ONSOptions)
	assert.Error(t, err, "Currency Cannot be changed")
	assert.False(t, ok)
	updates = generateGov()
	updates.ONSOptions.FirstLevelDomains = []string{"Test"}
	ok, err = vStore.ValidateONS(&updates.ONSOptions)
	assert.Error(t, err, "First level Domains Cannot be changed")
	assert.False(t, ok)
	updates = generateGov()
	updates.ONSOptions.BaseDomainPrice = *balance.NewAmountFromInt(int64(-1))
	ok, err = vStore.ValidateONS(&updates.ONSOptions)
	assert.Error(t, err, "BaseDomainPrice must be greater than 0")
	assert.False(t, ok)
	updates = generateGov()
	ok, err = vStore.ValidateONS(&updates.ONSOptions)
	assert.NoError(t, err, "Cannot be changed")
	assert.True(t, ok)
}

func TestStore_ValidateFee(t *testing.T) {
	updates := generateGov()
	updates.ONSOptions.Currency = "Test"
	ok, err := vStore.ValidateFee(&updates.FeeOption)
	assert.Error(t, err, "Currency Cannot be changed")
	assert.False(t, ok)
}
func generateGov() *GovernanceState {
	perblock, _ := big.NewInt(0).SetString("100000000000000", 10)
	baseDomainPrice, _ := big.NewInt(0).SetString("1000000000000000000000", 10)
	onsOp := ons.Options{
		Currency:          "OLT",
		PerBlockFees:      *balance.NewAmountFromBigInt(perblock),
		FirstLevelDomains: []string{"ol"},
		BaseDomainPrice:   *balance.NewAmountFromBigInt(baseDomainPrice),
	}
	btccdo := bitcoin.ChainDriverOption{
		"testnet3",
		"1000000000",
		"oneledgerSupplyAddress",
		int64(6),
	}

	propOpt := ProposalOptionSet{
		ConfigUpdate: ProposalOption{
			InitialFunding:         proposalInitialFunding,
			FundingGoal:            proposalFundingGoal,
			FundingDeadline:        proposalFundingDeadline,
			VotingDeadline:         proposalVotingDeadline,
			PassPercentage:         proposalPassPercentage,
			PassedFundDistribution: passedProposalDistribution,
			FailedFundDistribution: failedProposalDistribution,
			ProposalExecutionCost:  executionCostAddrConfig,
		},
		CodeChange: ProposalOption{
			InitialFunding:         proposalInitialFunding,
			FundingGoal:            proposalFundingGoal,
			FundingDeadline:        proposalFundingDeadline,
			VotingDeadline:         proposalVotingDeadline,
			PassPercentage:         proposalPassPercentage,
			PassedFundDistribution: passedProposalDistribution,
			FailedFundDistribution: failedProposalDistribution,
			ProposalExecutionCost:  executionCostAddrCodeChange,
		},
		General: ProposalOption{
			InitialFunding:         proposalInitialFunding,
			FundingGoal:            proposalFundingGoal,
			FundingDeadline:        proposalFundingDeadline,
			VotingDeadline:         proposalVotingDeadline,
			PassPercentage:         proposalPassPercentage,
			PassedFundDistribution: passedProposalDistribution,
			FailedFundDistribution: failedProposalDistribution,
			ProposalExecutionCost:  executionCostAddrGeneral,
		},
		BountyProgramAddr: bountyProgramAddr,
	}

	rewzOpt := rewards.Options{
		RewardInterval:    150,
		RewardPoolAddress: "rewardpool",
		RewardCurrency:    "OLT",
		CalculateInterval: blockRewardsCalcInterval,
		AnnualSupply:      *annualBlockRewardsSupply,
		YearsOfSupply:     yearsOfBlockRewardsSupply,
	}
	stakingOption := delegation.Options{
		MinSelfDelegationAmount: *balance.NewAmount(1000000),
		MinDelegationAmount:     *balance.NewAmount(10000000),
		TopValidatorCount:       8,
		MaturityTime:            109200,
	}
	olt := balance.Currency{Id: 0, Name: "OLT", Chain: chain.ONELEDGER, Decimal: 18, Unit: "nue"}
	feeOpt := fees.FeeOption{
		FeeCurrency:   olt,
		MinFeeDecimal: 9,
	}
	ethOpt := ethereum.ChainDriverOption{
		ContractABI:        "",
		ContractAddress:    common.Address{},
		TokenList:          nil,
		ERCContractABI:     "",
		ERCContractAddress: common.Address{},
		TotalSupply:        "",
		TotalSupplyAddr:    "",
		BlockConfirmation:  12,
	}
	vStore.WithHeight(0).SetFeeOption(feeOpt)
	vStore.WithHeight(0).SetStakingOptions(stakingOption)
	vStore.WithHeight(0).SetONSOptions(onsOp)
	vStore.WithHeight(0).SetETHChainDriverOption(ethOpt)
	vStore.WithHeight(0).SetBTCChainDriverOption(btccdo)
	vStore.WithHeight(0).SetRewardOptions(rewzOpt)
	vStore.WithHeight(0).SetProposalOptions(propOpt)
	vStore.WithHeight(0).SetLUH()
	return &GovernanceState{
		FeeOption:      feeOpt,
		ETHCDOption:    ethOpt,
		BTCCDOption:    btccdo,
		ONSOptions:     onsOp,
		PropOptions:    propOpt,
		StakingOptions: stakingOption,
		RewardOptions:  rewzOpt,
	}
}
