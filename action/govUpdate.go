package action

import (
	"fmt"
	"strconv"

	"github.com/pkg/errors"

	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/governance"
)

type FunctionBehaviour int

const (
	ValidateOnly      FunctionBehaviour = 0
	ValidateAndUpdate FunctionBehaviour = 1
)

type GovernaceUpdateAndValidate struct {
	GovernanceUpdateFunction map[string]func(interface{}, *Context, FunctionBehaviour) (bool, error)
}

func NewGovUpdate() *GovernaceUpdateAndValidate {

	newUpdateAndValidate := GovernaceUpdateAndValidate{
		GovernanceUpdateFunction: make(map[string]func(interface{}, *Context, FunctionBehaviour) (bool, error)),
	}
	newUpdateAndValidate.inititalizize()
	return &newUpdateAndValidate
}

func (g GovernaceUpdateAndValidate) inititalizize() {
	g.GovernanceUpdateFunction["feeOption.minFeeDecimal"] = feeOptionminFeeDecimal
	g.GovernanceUpdateFunction["onsOptions.perBlockFees"] = onsOptionsperBlockFees
	g.GovernanceUpdateFunction["onsOptions.baseDomainPrice"] = onsOptionsbaseDomainPrice
	g.GovernanceUpdateFunction["stakingOptions.minSelfDelegationAmount"] = stakingOptionsminSelfDelegationAmount
	g.GovernanceUpdateFunction["stakingOptions.topValidatorCount"] = stakingOptionstopValidatorCount
	g.GovernanceUpdateFunction["stakingOptions.maturityTime"] = stakingOptionsmaturityTime
	// Initial Funding Updates may also include updates for funding goal
	g.GovernanceUpdateFunction["propOptions.configUpdate.initialFunding"] = propOptionsconfigUpdateinitialFunding
	g.GovernanceUpdateFunction["propOptions.codeChange.initialFunding"] = propOptionscodeChangeinitialFunding
	g.GovernanceUpdateFunction["propOptions.general.initialFunding"] = propOptionsgeneralinitialFunding
	g.GovernanceUpdateFunction["propOptions.configUpdate.fundingGoal"] = propOptionsconfigUpdatefundingGoal
	g.GovernanceUpdateFunction["propOptions.codeChange.fundingGoal"] = propOptionscodeChangefundingGoal
	g.GovernanceUpdateFunction["propOptions.general.fundingGoal"] = propOptionsgeneralfundingGoal
	g.GovernanceUpdateFunction["propOptions.configUpdate.votingDeadline"] = propOptionsconfigUpdatevotingDeadline
	g.GovernanceUpdateFunction["propOptions.codeChange.votingDeadline"] = propOptionscodeChangevotingDeadline
	g.GovernanceUpdateFunction["propOptions.general.votingDeadline"] = propOptionsgeneralvotingDeadline
	g.GovernanceUpdateFunction["propOptions.configUpdate.fundingDeadline"] = propOptionsconfigUpdatefundingDeadline
	g.GovernanceUpdateFunction["propOptions.codeChange.fundingDeadline"] = propOptionscodeChangefundingDeadline
	g.GovernanceUpdateFunction["propOptions.general.fundingDeadline"] = propOptionsgeneralfundingDeadline
	g.GovernanceUpdateFunction["propOptions.configUpdate.passPercentage"] = propOptionsconfigUpdatepassPercentage
	g.GovernanceUpdateFunction["propOptions.codeChange.passPercentage"] = propOptionscodeChangepassPercentage
	g.GovernanceUpdateFunction["propOptions.general.passPercentage"] = propOptionsgeneralpassPercentage
	//g.GovernanceUpdateFunction["propOptions.configUpdate.passedFundDistribution"] = propOptionsconfigUpdatepassedFundDistribution
	//g.GovernanceUpdateFunction["propOptions.codeChange.passedFundDistribution"] = propOptionscodeChangepassedFundDistribution
	//g.GovernanceUpdateFunction["propOptions.general.passedFundDistribution"] = propOptionsgeneralpassedFundDistribution
	//g.GovernanceUpdateFunction["propOptions.configUpdate.failedFundDistribution"] = propOptionsconfigUpdatefailedFundDistribution
	//g.GovernanceUpdateFunction["propOptions.codeChange.failedFundDistribution"] = propOptionscodeChangefailedFundDistribution
	//g.GovernanceUpdateFunction["propOptions.general.failedFundDistribution"] = propOptionsgeneralfailedFundDistribution
	g.GovernanceUpdateFunction["evidenceOptions.minVotesRequired"] = evidenceOptionsminVotesRequired
	g.GovernanceUpdateFunction["evidenceOptions.blockVotesDiff"] = evidenceOptionsblockVotesDiff
	g.GovernanceUpdateFunction["evidenceOptions.penaltyBasePercentage"] = evidenceOptionspenaltyBasePercentage
	//g.GovernanceUpdateFunction["evidenceOptions.penaltyPercentage"] = evidenceOptionspenaltyPercentage

	//MinVotesRequired: 2, // should be atleast 70% or greater of block votes diff
	//	BlockVotesDiff:   4,// min limit - 1000, max limit - 100,000
	//		PenaltyBasePercentage: 30, // min limit - 10, max limit - 40
	//		PenaltyBountyPercentage: 50,// can be between 0 - 100, PenaltyBurnPercentage + PenaltyBountyPercentage is always 100
	//		PenaltyBurnPercentage: 50,// can be between 0 -100, PenaltyBurnPercentage + PenaltyBountyPercentage is always 100
}

func evidenceOptionsminVotesRequired(value interface{}, ctx *Context, validationOnly FunctionBehaviour) (bool, error) {
	Options, err := ctx.GovernanceStore.GetEvidenceOptions()
	if err != nil {
		return false, err
	}
	newValue, err := getNewValueInt64(value)
	if err != nil {
		return false, err
	}
	Options.MinVotesRequired = newValue
	ok, err := ctx.GovernanceStore.ValidateEvidence(Options)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, errors.New("Validation Failed")
	}
	if validationOnly == ValidateOnly {
		return true, nil
	}
	err = ctx.GovernanceStore.WithHeight(ctx.Header.Height).SetEvidenceOptions(*Options)
	if err != nil {
		return false, errors.Wrap(err, "Setup Evidence Options")
	}
	err = ctx.GovernanceStore.WithHeight(ctx.Header.Height).SetLUH(governance.LAST_UPDATE_HEIGHT_EVIDENCE)
	if err != nil {
		return false, errors.Wrap(err, "Unable to set last Update height ")
	}
	ctx.Logger.Debug("Governance options set at height : ", ctx.Header.Height, "| evidenceOptions.minVotesRequired :", newValue)
	return true, nil
}
func evidenceOptionsblockVotesDiff(value interface{}, ctx *Context, validationOnly FunctionBehaviour) (bool, error) {
	Options, err := ctx.GovernanceStore.GetEvidenceOptions()
	if err != nil {
		return false, err
	}
	newValue, err := getNewValueInt64(value)
	if err != nil {
		return false, err
	}
	Options.BlockVotesDiff = newValue
	ok, err := ctx.GovernanceStore.ValidateEvidence(Options)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, errors.New("Validation Failed")
	}
	if validationOnly == ValidateOnly {
		return true, nil
	}
	err = ctx.GovernanceStore.WithHeight(ctx.Header.Height).SetEvidenceOptions(*Options)
	if err != nil {
		return false, errors.Wrap(err, "Setup Evidence Options")
	}
	err = ctx.GovernanceStore.WithHeight(ctx.Header.Height).SetLUH(governance.LAST_UPDATE_HEIGHT_EVIDENCE)
	if err != nil {
		return false, errors.Wrap(err, "Unable to set last Update height ")
	}
	ctx.Logger.Debug("Governance options set at height : ", ctx.Header.Height, "| evidenceOptions.blockVotesDiff :", newValue)
	return true, nil

}
func evidenceOptionspenaltyBasePercentage(value interface{}, ctx *Context, validationOnly FunctionBehaviour) (bool, error) {
	Options, err := ctx.GovernanceStore.GetEvidenceOptions()
	if err != nil {
		return false, err
	}
	newValue, err := getNewValueInt64(value)
	if err != nil {
		return false, err
	}
	Options.PenaltyBasePercentage = newValue
	ok, err := ctx.GovernanceStore.ValidateEvidence(Options)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, errors.New("Validation Failed")
	}
	if validationOnly == ValidateOnly {
		return true, nil
	}
	err = ctx.GovernanceStore.WithHeight(ctx.Header.Height).SetEvidenceOptions(*Options)
	if err != nil {
		return false, errors.Wrap(err, "Setup Evidence Options")
	}
	err = ctx.GovernanceStore.WithHeight(ctx.Header.Height).SetLUH(governance.LAST_UPDATE_HEIGHT_EVIDENCE)
	if err != nil {
		return false, errors.Wrap(err, "Unable to set last Update height ")
	}
	ctx.Logger.Debug("Governance options set at height : ", ctx.Header.Height, "| evidenceOptions.penaltyBasePercentage :", newValue)
	return true, nil
}

//func evidenceOptionspenaltyPercentage(value interface{}, ctx *Context, validationOnly FunctionBehaviour) (bool, error) {
//	Options, err := ctx.GovernanceStore.GetEvidenceOptions()
//	if err != nil {
//		return false, err
//	}
//	newValue, err := getEvidenceOpt(value)
//	if err != nil {
//		fmt.Println(err)
//		return false, err
//	}
//	fmt.Println(newValue)
//	Options.PenaltyBountyPercentage = newValue.PenaltyBountyPercentage
//	Options.PenaltyBurnPercentage = newValue.PenaltyBurnPercentage
//	ok, err := ctx.GovernanceStore.ValidateEvidence(Options)
//	if err != nil {
//		return false, err
//	}
//	if !ok {
//		return false, errors.New("Validation Failed")
//	}
//	if validationOnly == ValidateOnly {
//		return true, nil
//	}
//	err = ctx.GovernanceStore.WithHeight(ctx.Header.Height).SetEvidenceOptions(*Options)
//	if err != nil {
//		return false, errors.Wrap(err, "Setup Evidence Options")
//	}
//	//Staking not part of context
//	err = ctx.GovernanceStore.WithHeight(ctx.Header.Height).SetLUH(governance.LAST_UPDATE_HEIGHT_EVIDENCE)
//	if err != nil {
//		return false, errors.Wrap(err, "Unable to set last Update height ")
//	}
//	ctx.Logger.Debug("Governance options set at height : ", ctx.Header.Height, "| evidenceOptions.penaltyBountyPercentage :", newValue)
//	return true, nil
//}

func stakingOptionsminSelfDelegationAmount(value interface{}, ctx *Context, validationOnly FunctionBehaviour) (bool, error) {
	Options, err := ctx.GovernanceStore.GetStakingOptions()
	if err != nil {
		return false, err
	}
	newValue, err := getNewValueInt64(value)
	if err != nil {
		return false, err
	}
	Options.MinSelfDelegationAmount = *balance.NewAmount(newValue)

	ok, err := ctx.GovernanceStore.ValidateStaking(Options)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, errors.New("Validation Failed")
	}
	if validationOnly == ValidateOnly {
		return true, nil
	}
	err = ctx.GovernanceStore.WithHeight(ctx.Header.Height).SetStakingOptions(*Options)
	if err != nil {
		return false, errors.Wrap(err, "Setup Staking Options")
	}
	//Staking not part of context
	err = ctx.GovernanceStore.WithHeight(ctx.Header.Height).SetLUH(governance.LAST_UPDATE_HEIGHT_STAKING)
	if err != nil {
		return false, errors.Wrap(err, "Unable to set last Update height ")
	}
	ctx.Logger.Debug("Governance options set at height : ", ctx.Header.Height, "| stakingOptions.minSelfDelegationAmount :", newValue)
	return true, nil
}
func stakingOptionstopValidatorCount(value interface{}, ctx *Context, validationOnly FunctionBehaviour) (bool, error) {
	Options, err := ctx.GovernanceStore.GetStakingOptions()
	if err != nil {
		return false, err
	}
	newValue, err := getNewValueInt64(value)
	if err != nil {
		return false, err
	}

	Options.TopValidatorCount = newValue

	ok, err := ctx.GovernanceStore.ValidateStaking(Options)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, errors.New("Validation Failed")
	}
	if validationOnly == ValidateOnly {
		return true, nil
	}
	err = ctx.GovernanceStore.WithHeight(ctx.Header.Height).SetStakingOptions(*Options)
	if err != nil {
		return false, errors.Wrap(err, "Setup Staking Options")
	}
	//Staking not part of context
	err = ctx.GovernanceStore.WithHeight(ctx.Header.Height).SetLUH(governance.LAST_UPDATE_HEIGHT_STAKING)
	if err != nil {
		return false, errors.Wrap(err, "Unable to set last Update height ")
	}
	ctx.Logger.Debug("Governance options set at height : ", ctx.Header.Height, "| takingOptions.topValidatorCount :", newValue)
	return true, nil
}
func stakingOptionsmaturityTime(value interface{}, ctx *Context, validationOnly FunctionBehaviour) (bool, error) {
	Options, err := ctx.GovernanceStore.GetStakingOptions()
	if err != nil {
		return false, err
	}

	newValue, err := getNewValueInt64(value)
	if err != nil {
		return false, err
	}
	fmt.Println("NewValue :", newValue)
	Options.MaturityTime = newValue

	ok, err := ctx.GovernanceStore.ValidateStaking(Options)
	fmt.Println("Validate Staking :", ok, err)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, errors.New("Validation Failed")
	}
	if validationOnly == ValidateOnly {
		return true, nil
	}
	err = ctx.GovernanceStore.WithHeight(ctx.Header.Height).SetStakingOptions(*Options)
	if err != nil {
		return false, errors.Wrap(err, "Setup Staking Options")
	}
	//Staking not part of context
	err = ctx.GovernanceStore.WithHeight(ctx.Header.Height).SetLUH(governance.LAST_UPDATE_HEIGHT_STAKING)
	if err != nil {
		return false, errors.Wrap(err, "Unable to set last Update height ")
	}
	ctx.Logger.Debug("Governance options set at height : ", ctx.Header.Height, "| stakingOptions.maturityTime :", newValue)
	return true, nil
}
func propOptionsconfigUpdateinitialFunding(value interface{}, ctx *Context, validationOnly FunctionBehaviour) (bool, error) {
	Options, err := ctx.GovernanceStore.GetProposalOptions()
	if err != nil {
		return false, err
	}
	newValue, err := getNewValueInt64(value)
	if err != nil {
		return false, err
	}
	Options.ConfigUpdate.InitialFunding = balance.NewAmount(newValue)

	ok, err := ctx.GovernanceStore.ValidateProposal(Options)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, errors.New("Validation Failed")
	}
	if validationOnly == ValidateOnly {
		return true, nil
	}
	err = ctx.GovernanceStore.WithHeight(ctx.Header.Height).SetProposalOptions(*Options)
	if err != nil {
		return false, errors.Wrap(err, "Setup Proposal Options")
	}
	ctx.ProposalMasterStore.Proposal.SetOptions(Options)
	err = ctx.GovernanceStore.WithHeight(ctx.Header.Height).SetLUH(governance.LAST_UPDATE_HEIGHT_PROPOSAL)
	if err != nil {
		return false, errors.Wrap(err, "Unable to set last Update height ")
	}
	ctx.Logger.Debug("Governance options set at height : ", ctx.Header.Height, "| propOptions.configUpdate.initialFunding :", newValue)
	return true, nil
}
func propOptionscodeChangeinitialFunding(value interface{}, ctx *Context, validationOnly FunctionBehaviour) (bool, error) {
	Options, err := ctx.GovernanceStore.GetProposalOptions()
	if err != nil {
		return false, err
	}
	newValue, err := getNewValueInt64(value)
	if err != nil {
		return false, err
	}
	Options.CodeChange.InitialFunding = balance.NewAmount(newValue)

	ok, err := ctx.GovernanceStore.ValidateProposal(Options)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, errors.New("Validation Failed")
	}
	if validationOnly == ValidateOnly {
		return true, nil
	}
	err = ctx.GovernanceStore.WithHeight(ctx.Header.Height).SetProposalOptions(*Options)
	if err != nil {
		return false, errors.Wrap(err, "Setup Proposal Options")
	}
	ctx.ProposalMasterStore.Proposal.SetOptions(Options)
	err = ctx.GovernanceStore.WithHeight(ctx.Header.Height).SetLUH(governance.LAST_UPDATE_HEIGHT_PROPOSAL)
	if err != nil {
		return false, errors.Wrap(err, "Unable to set last Update height ")
	}
	ctx.Logger.Debug("Governance options set at height : ", ctx.Header.Height, "| propOptions.codeChange.initialFunding :", newValue)
	return true, nil
}
func propOptionsgeneralinitialFunding(value interface{}, ctx *Context, validationOnly FunctionBehaviour) (bool, error) {
	Options, err := ctx.GovernanceStore.GetProposalOptions()
	if err != nil {
		return false, err
	}
	newValue, err := getNewValueInt64(value)
	if err != nil {
		return false, err
	}
	Options.General.InitialFunding = balance.NewAmount(newValue)

	ok, err := ctx.GovernanceStore.ValidateProposal(Options)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, errors.New("Validation Failed")
	}
	if validationOnly == ValidateOnly {
		return true, nil
	}
	err = ctx.GovernanceStore.WithHeight(ctx.Header.Height).SetProposalOptions(*Options)
	if err != nil {
		return false, errors.Wrap(err, "Setup Proposal Options")
	}
	ctx.ProposalMasterStore.Proposal.SetOptions(Options)
	err = ctx.GovernanceStore.WithHeight(ctx.Header.Height).SetLUH(governance.LAST_UPDATE_HEIGHT_PROPOSAL)
	if err != nil {
		return false, errors.Wrap(err, "Unable to set last Update height ")
	}
	ctx.Logger.Debug("Governance options set at height : ", ctx.Header.Height, "| propOptions.general.initialFunding :", newValue)
	return true, nil
}
func propOptionsconfigUpdatefundingGoal(value interface{}, ctx *Context, validationOnly FunctionBehaviour) (bool, error) {
	Options, err := ctx.GovernanceStore.GetProposalOptions()
	if err != nil {
		return false, err
	}
	newValue, err := getNewValueInt64(value)
	if err != nil {
		return false, err
	}
	Options.ConfigUpdate.FundingGoal = balance.NewAmount(newValue)

	ok, err := ctx.GovernanceStore.ValidateProposal(Options)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, errors.New("Validation Failed")
	}
	if validationOnly == ValidateOnly {
		return true, nil
	}
	err = ctx.GovernanceStore.WithHeight(ctx.Header.Height).SetProposalOptions(*Options)
	if err != nil {
		return false, errors.Wrap(err, "Setup Proposal Options")
	}
	ctx.ProposalMasterStore.Proposal.SetOptions(Options)
	err = ctx.GovernanceStore.WithHeight(ctx.Header.Height).SetLUH(governance.LAST_UPDATE_HEIGHT_PROPOSAL)
	if err != nil {
		return false, errors.Wrap(err, "Unable to set last Update height ")
	}
	ctx.Logger.Debug("Governance options set at height : ", ctx.Header.Height, "| propOptions.configUpdate.fundingGoal :", newValue)
	return true, nil
}
func propOptionscodeChangefundingGoal(value interface{}, ctx *Context, validationOnly FunctionBehaviour) (bool, error) {
	Options, err := ctx.GovernanceStore.GetProposalOptions()
	if err != nil {
		return false, err
	}
	newValue, err := getNewValueInt64(value)
	if err != nil {
		return false, err
	}
	Options.CodeChange.FundingGoal = balance.NewAmount(newValue)

	ok, err := ctx.GovernanceStore.ValidateProposal(Options)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, errors.New("Validation Failed")
	}
	if validationOnly == ValidateOnly {
		return true, nil
	}
	err = ctx.GovernanceStore.WithHeight(ctx.Header.Height).SetProposalOptions(*Options)
	if err != nil {
		return false, errors.Wrap(err, "Setup Proposal Options")
	}
	ctx.ProposalMasterStore.Proposal.SetOptions(Options)
	err = ctx.GovernanceStore.WithHeight(ctx.Header.Height).SetLUH(governance.LAST_UPDATE_HEIGHT_PROPOSAL)
	if err != nil {
		return false, errors.Wrap(err, "Unable to set last Update height ")
	}
	ctx.Logger.Debug("Governance options set at height : ", ctx.Header.Height, "| propOptions.codeChange.fundingGoal :", newValue)
	return true, nil
}
func propOptionsgeneralfundingGoal(value interface{}, ctx *Context, validationOnly FunctionBehaviour) (bool, error) {
	Options, err := ctx.GovernanceStore.GetProposalOptions()
	if err != nil {
		return false, err
	}
	newValue, err := getNewValueInt64(value)
	if err != nil {
		return false, err
	}
	Options.General.FundingGoal = balance.NewAmount(newValue)

	ok, err := ctx.GovernanceStore.ValidateProposal(Options)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, errors.New("Validation Failed")
	}
	if validationOnly == ValidateOnly {
		return true, nil
	}
	err = ctx.GovernanceStore.WithHeight(ctx.Header.Height).SetProposalOptions(*Options)
	if err != nil {
		return false, errors.Wrap(err, "Setup Proposal Options")
	}
	ctx.ProposalMasterStore.Proposal.SetOptions(Options)
	err = ctx.GovernanceStore.WithHeight(ctx.Header.Height).SetLUH(governance.LAST_UPDATE_HEIGHT_PROPOSAL)
	if err != nil {
		return false, errors.Wrap(err, "Unable to set last Update height ")
	}
	ctx.Logger.Debug("Governance options set at height : ", ctx.Header.Height, "| propOptions.general.fundingGoal :", newValue)
	return true, nil
}
func propOptionsconfigUpdatevotingDeadline(value interface{}, ctx *Context, validationOnly FunctionBehaviour) (bool, error) {
	Options, err := ctx.GovernanceStore.GetProposalOptions()
	if err != nil {
		return false, err
	}
	newValue, err := getNewValueInt64(value)
	if err != nil {
		return false, err
	}
	Options.ConfigUpdate.VotingDeadline = newValue

	ok, err := ctx.GovernanceStore.ValidateProposal(Options)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, errors.New("Validation Failed")
	}
	if validationOnly == ValidateOnly {
		return true, nil
	}
	err = ctx.GovernanceStore.WithHeight(ctx.Header.Height).SetProposalOptions(*Options)
	if err != nil {
		return false, errors.Wrap(err, "Setup Proposal Options")
	}
	ctx.ProposalMasterStore.Proposal.SetOptions(Options)
	err = ctx.GovernanceStore.WithHeight(ctx.Header.Height).SetLUH(governance.LAST_UPDATE_HEIGHT_PROPOSAL)
	if err != nil {
		return false, errors.Wrap(err, "Unable to set last Update height ")
	}
	ctx.Logger.Debug("Governance options set at height : ", ctx.Header.Height, "| propOptions.configUpdate.votingDeadline :", newValue)
	return true, nil
}
func propOptionscodeChangevotingDeadline(value interface{}, ctx *Context, validationOnly FunctionBehaviour) (bool, error) {
	Options, err := ctx.GovernanceStore.GetProposalOptions()
	if err != nil {
		return false, err
	}
	newValue, err := getNewValueInt64(value)
	if err != nil {
		return false, err
	}
	Options.CodeChange.VotingDeadline = newValue

	ok, err := ctx.GovernanceStore.ValidateProposal(Options)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, errors.New("Validation Failed")
	}
	if validationOnly == ValidateOnly {
		return true, nil
	}
	err = ctx.GovernanceStore.WithHeight(ctx.Header.Height).SetProposalOptions(*Options)
	if err != nil {
		return false, errors.Wrap(err, "Setup Proposal Options")
	}
	ctx.ProposalMasterStore.Proposal.SetOptions(Options)
	err = ctx.GovernanceStore.WithHeight(ctx.Header.Height).SetLUH(governance.LAST_UPDATE_HEIGHT_PROPOSAL)
	if err != nil {
		return false, errors.Wrap(err, "Unable to set last Update height ")
	}
	ctx.Logger.Debug("Governance options set at height : ", ctx.Header.Height, "| propOptions.codeChange.votingDeadline :", newValue)
	return true, nil
}
func propOptionsgeneralvotingDeadline(value interface{}, ctx *Context, validationOnly FunctionBehaviour) (bool, error) {
	Options, err := ctx.GovernanceStore.GetProposalOptions()
	if err != nil {
		return false, err
	}
	newValue, err := getNewValueInt64(value)
	if err != nil {
		return false, err
	}
	Options.General.VotingDeadline = newValue

	ok, err := ctx.GovernanceStore.ValidateProposal(Options)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, errors.New("Validation Failed")
	}
	if validationOnly == ValidateOnly {
		return true, nil
	}
	err = ctx.GovernanceStore.WithHeight(ctx.Header.Height).SetProposalOptions(*Options)
	if err != nil {
		return false, errors.Wrap(err, "Setup Proposal Options")
	}
	ctx.ProposalMasterStore.Proposal.SetOptions(Options)
	err = ctx.GovernanceStore.WithHeight(ctx.Header.Height).SetLUH(governance.LAST_UPDATE_HEIGHT_PROPOSAL)
	if err != nil {
		return false, errors.Wrap(err, "Unable to set last Update height ")
	}
	ctx.Logger.Debug("Governance options set at height : ", ctx.Header.Height, "| propOptions.general.votingDeadline :", newValue)
	return true, nil
}
func propOptionsconfigUpdatefundingDeadline(value interface{}, ctx *Context, validationOnly FunctionBehaviour) (bool, error) {
	Options, err := ctx.GovernanceStore.GetProposalOptions()
	if err != nil {
		return false, err
	}
	newValue, err := getNewValueInt64(value)
	if err != nil {
		return false, err
	}
	Options.ConfigUpdate.FundingDeadline = newValue

	ok, err := ctx.GovernanceStore.ValidateProposal(Options)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, errors.New("Validation Failed")
	}
	if validationOnly == ValidateOnly {
		return true, nil
	}
	err = ctx.GovernanceStore.WithHeight(ctx.Header.Height).SetProposalOptions(*Options)
	if err != nil {
		return false, errors.Wrap(err, "Setup Proposal Options")
	}
	ctx.ProposalMasterStore.Proposal.SetOptions(Options)
	err = ctx.GovernanceStore.WithHeight(ctx.Header.Height).SetLUH(governance.LAST_UPDATE_HEIGHT_PROPOSAL)
	if err != nil {
		return false, errors.Wrap(err, "Unable to set last Update height ")
	}
	ctx.Logger.Debug("Governance options set at height : ", ctx.Header.Height, "| propOptions.configUpdate.fundingDeadline :", newValue)
	return true, nil
}
func propOptionscodeChangefundingDeadline(value interface{}, ctx *Context, validationOnly FunctionBehaviour) (bool, error) {
	Options, err := ctx.GovernanceStore.GetProposalOptions()
	if err != nil {
		return false, err
	}
	newValue, err := getNewValueInt64(value)
	if err != nil {
		return false, err
	}
	Options.CodeChange.FundingDeadline = newValue

	ok, err := ctx.GovernanceStore.ValidateProposal(Options)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, errors.New("Validation Failed")
	}
	if validationOnly == ValidateOnly {
		return true, nil
	}
	err = ctx.GovernanceStore.WithHeight(ctx.Header.Height).SetProposalOptions(*Options)
	if err != nil {
		return false, errors.Wrap(err, "Setup Proposal Options")
	}
	ctx.ProposalMasterStore.Proposal.SetOptions(Options)
	err = ctx.GovernanceStore.WithHeight(ctx.Header.Height).SetLUH(governance.LAST_UPDATE_HEIGHT_PROPOSAL)
	if err != nil {
		return false, errors.Wrap(err, "Unable to set last Update height ")
	}
	ctx.Logger.Debug("Governance options set at height : ", ctx.Header.Height, "| propOptions.codeChange.fundingDeadline :", newValue)
	return true, nil
}
func propOptionsgeneralfundingDeadline(value interface{}, ctx *Context, validationOnly FunctionBehaviour) (bool, error) {
	Options, err := ctx.GovernanceStore.GetProposalOptions()
	if err != nil {
		return false, err
	}
	newValue, err := getNewValueInt64(value)
	if err != nil {
		return false, err
	}
	Options.General.FundingDeadline = newValue

	ok, err := ctx.GovernanceStore.ValidateProposal(Options)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, errors.New("Validation Failed")
	}
	if validationOnly == ValidateOnly {
		return true, nil
	}
	err = ctx.GovernanceStore.WithHeight(ctx.Header.Height).SetProposalOptions(*Options)
	if err != nil {
		return false, errors.Wrap(err, "Setup Proposal Options")
	}
	ctx.ProposalMasterStore.Proposal.SetOptions(Options)
	err = ctx.GovernanceStore.WithHeight(ctx.Header.Height).SetLUH(governance.LAST_UPDATE_HEIGHT_PROPOSAL)
	if err != nil {
		return false, errors.Wrap(err, "Unable to set last Update height ")
	}
	ctx.Logger.Debug("Governance options set at height : ", ctx.Header.Height, "| propOptions.general.fundingDeadline :", newValue)
	return true, nil
}
func propOptionsconfigUpdatepassPercentage(value interface{}, ctx *Context, validationOnly FunctionBehaviour) (bool, error) {
	Options, err := ctx.GovernanceStore.GetProposalOptions()
	if err != nil {
		return false, err
	}
	newValue, err := getNewValueInt64(value)
	if err != nil {
		return false, err
	}
	Options.ConfigUpdate.PassPercentage = int(newValue)

	ok, err := ctx.GovernanceStore.ValidateProposal(Options)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, errors.New("Validation Failed")
	}
	if validationOnly == ValidateOnly {
		return true, nil
	}
	err = ctx.GovernanceStore.WithHeight(ctx.Header.Height).SetProposalOptions(*Options)
	if err != nil {
		return false, errors.Wrap(err, "Setup Proposal Options")
	}
	ctx.ProposalMasterStore.Proposal.SetOptions(Options)
	err = ctx.GovernanceStore.WithHeight(ctx.Header.Height).SetLUH(governance.LAST_UPDATE_HEIGHT_PROPOSAL)
	if err != nil {
		return false, errors.Wrap(err, "Unable to set last Update height ")
	}
	ctx.Logger.Debug("Governance options set at height : ", ctx.Header.Height, "| propOptions.configUpdate.passPercentage :", newValue)
	return true, nil
}
func propOptionscodeChangepassPercentage(value interface{}, ctx *Context, validationOnly FunctionBehaviour) (bool, error) {
	Options, err := ctx.GovernanceStore.GetProposalOptions()
	if err != nil {
		return false, err
	}
	newValue, err := getNewValueInt64(value)
	if err != nil {
		return false, err
	}
	Options.CodeChange.PassPercentage = int(newValue)

	ok, err := ctx.GovernanceStore.ValidateProposal(Options)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, errors.New("Validation Failed")
	}
	if validationOnly == ValidateOnly {
		return true, nil
	}
	err = ctx.GovernanceStore.WithHeight(ctx.Header.Height).SetProposalOptions(*Options)
	if err != nil {
		return false, errors.Wrap(err, "Setup Proposal Options")
	}
	ctx.ProposalMasterStore.Proposal.SetOptions(Options)
	err = ctx.GovernanceStore.WithHeight(ctx.Header.Height).SetLUH(governance.LAST_UPDATE_HEIGHT_PROPOSAL)
	if err != nil {
		return false, errors.Wrap(err, "Unable to set last Update height ")
	}
	ctx.Logger.Debug("Governance options set at height : ", ctx.Header.Height, "| propOptions.codeChange.passPercentage :", newValue)
	return true, nil
}
func propOptionsgeneralpassPercentage(value interface{}, ctx *Context, validationOnly FunctionBehaviour) (bool, error) {
	Options, err := ctx.GovernanceStore.GetProposalOptions()
	if err != nil {
		return false, err
	}
	newValue, err := getNewValueInt64(value)
	if err != nil {
		return false, err
	}
	Options.General.PassPercentage = int(newValue)

	ok, err := ctx.GovernanceStore.ValidateProposal(Options)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, errors.New("Validation Failed")
	}
	if validationOnly == ValidateOnly {
		return true, nil
	}
	err = ctx.GovernanceStore.WithHeight(ctx.Header.Height).SetProposalOptions(*Options)
	if err != nil {
		return false, errors.Wrap(err, "Setup Proposal Options")
	}
	ctx.ProposalMasterStore.Proposal.SetOptions(Options)
	err = ctx.GovernanceStore.WithHeight(ctx.Header.Height).SetLUH(governance.LAST_UPDATE_HEIGHT_PROPOSAL)
	if err != nil {
		return false, errors.Wrap(err, "Unable to set last Update height ")
	}
	ctx.Logger.Debug("Governance options set at height : ", ctx.Header.Height, "| propOptions.general.passPercentage :", newValue)
	return true, nil
}

//func propOptionsconfigUpdatepassedFundDistribution(value interface{}, ctx *Context, validationOnly FunctionBehaviour) (bool, error) {
//	Options, err := ctx.GovernanceStore.GetProposalOptions()
//	if err != nil {
//		return false, err
//	}
//	newDistribution, err := getDistribution(value)
//	if err != nil {
//		return false, err
//	}
//	Options.ConfigUpdate.PassedFundDistribution = newDistribution
//	ok, err := ctx.GovernanceStore.ValidateProposal(Options)
//	if err != nil {
//		return false, err
//	}
//	if !ok {
//		return false, errors.New("Validation Failed")
//	}
//	if validationOnly == ValidateOnly {
//		return true, nil
//	}
//	err = ctx.GovernanceStore.WithHeight(ctx.Header.Height).SetProposalOptions(*Options)
//	if err != nil {
//		return false, errors.Wrap(err, "Setup Proposal Options")
//	}
//	ctx.ProposalMasterStore.Proposal.SetOptions(Options)
//	err = ctx.GovernanceStore.WithHeight(ctx.Header.Height).SetLUH(governance.LAST_UPDATE_HEIGHT_PROPOSAL)
//	if err != nil {
//		return false, errors.Wrap(err, "Unable to set last Update height ")
//	}
//	ctx.Logger.Debug("Governance options set at height : ", ctx.Header.Height, "| propOptions.configUpdate.passedFundDistribution :", newDistribution)
//	return true, nil
//}
//func propOptionscodeChangepassedFundDistribution(value interface{}, ctx *Context, validationOnly FunctionBehaviour) (bool, error) {
//	Options, err := ctx.GovernanceStore.GetProposalOptions()
//	if err != nil {
//		return false, err
//	}
//	newDistribution, err := getDistribution(value)
//	if err != nil {
//		return false, err
//	}
//	Options.CodeChange.PassedFundDistribution = newDistribution
//	ok, err := ctx.GovernanceStore.ValidateProposal(Options)
//	if err != nil {
//		return false, err
//	}
//	if !ok {
//		return false, errors.New("Validation Failed")
//	}
//	if validationOnly == ValidateOnly {
//		return true, nil
//	}
//	err = ctx.GovernanceStore.WithHeight(ctx.Header.Height).SetProposalOptions(*Options)
//	if err != nil {
//		return false, errors.Wrap(err, "Setup Proposal Options")
//	}
//	ctx.ProposalMasterStore.Proposal.SetOptions(Options)
//	err = ctx.GovernanceStore.WithHeight(ctx.Header.Height).SetLUH(governance.LAST_UPDATE_HEIGHT_PROPOSAL)
//	if err != nil {
//		return false, errors.Wrap(err, "Unable to set last Update height ")
//	}
//	ctx.Logger.Debug("Governance options set at height : ", ctx.Header.Height, "| propOptions.codeChange.passedFundDistribution :", newDistribution)
//	return true, nil
//}
//func propOptionsgeneralpassedFundDistribution(value interface{}, ctx *Context, validationOnly FunctionBehaviour) (bool, error) {
//	Options, err := ctx.GovernanceStore.GetProposalOptions()
//	if err != nil {
//		return false, err
//	}
//	newDistribution, err := getDistribution(value)
//	if err != nil {
//		return false, err
//	}
//	Options.General.PassedFundDistribution = newDistribution
//	ok, err := ctx.GovernanceStore.ValidateProposal(Options)
//	if err != nil {
//		return false, err
//	}
//	if !ok {
//		return false, errors.New("Validation Failed")
//	}
//	if validationOnly == ValidateOnly {
//		return true, nil
//	}
//	err = ctx.GovernanceStore.WithHeight(ctx.Header.Height).SetProposalOptions(*Options)
//	if err != nil {
//		return false, errors.Wrap(err, "Setup Proposal Options")
//	}
//	ctx.ProposalMasterStore.Proposal.SetOptions(Options)
//	err = ctx.GovernanceStore.WithHeight(ctx.Header.Height).SetLUH(governance.LAST_UPDATE_HEIGHT_PROPOSAL)
//	if err != nil {
//		return false, errors.Wrap(err, "Unable to set last Update height ")
//	}
//	ctx.Logger.Debug("Governance options set at height : ", ctx.Header.Height, "| propOptions.general.passedFundDistribution :", newDistribution)
//	return true, nil
//}
//func propOptionsconfigUpdatefailedFundDistribution(value interface{}, ctx *Context, validationOnly FunctionBehaviour) (bool, error) {
//	Options, err := ctx.GovernanceStore.GetProposalOptions()
//	if err != nil {
//		return false, err
//	}
//	newDistribution, err := getDistribution(value)
//	if err != nil {
//		return false, err
//	}
//	Options.ConfigUpdate.FailedFundDistribution = newDistribution
//	ok, err := ctx.GovernanceStore.ValidateProposal(Options)
//	if err != nil {
//		return false, err
//	}
//	if !ok {
//		return false, errors.New("Validation Failed")
//	}
//	if validationOnly == ValidateOnly {
//		return true, nil
//	}
//	err = ctx.GovernanceStore.WithHeight(ctx.Header.Height).SetProposalOptions(*Options)
//	if err != nil {
//		return false, errors.Wrap(err, "Setup Proposal Options")
//	}
//	ctx.ProposalMasterStore.Proposal.SetOptions(Options)
//	err = ctx.GovernanceStore.WithHeight(ctx.Header.Height).SetLUH(governance.LAST_UPDATE_HEIGHT_PROPOSAL)
//	if err != nil {
//		return false, errors.Wrap(err, "Unable to set last Update height ")
//	}
//	ctx.Logger.Debug("Governance options set at height : ", ctx.Header.Height, "| propOptions.configUpdate.failedFundDistribution :", newDistribution)
//	return true, nil
//}
//func propOptionscodeChangefailedFundDistribution(value interface{}, ctx *Context, validationOnly FunctionBehaviour) (bool, error) {
//	Options, err := ctx.GovernanceStore.GetProposalOptions()
//	if err != nil {
//		return false, err
//	}
//	newDistribution, err := getDistribution(value)
//	if err != nil {
//		return false, err
//	}
//	Options.CodeChange.FailedFundDistribution = newDistribution
//	ok, err := ctx.GovernanceStore.ValidateProposal(Options)
//	if err != nil {
//		return false, err
//	}
//	if !ok {
//		return false, errors.New("Validation Failed")
//	}
//	if validationOnly == ValidateOnly {
//		return true, nil
//	}
//	err = ctx.GovernanceStore.WithHeight(ctx.Header.Height).SetProposalOptions(*Options)
//	if err != nil {
//		return false, errors.Wrap(err, "Setup Proposal Options")
//	}
//	ctx.ProposalMasterStore.Proposal.SetOptions(Options)
//	err = ctx.GovernanceStore.WithHeight(ctx.Header.Height).SetLUH(governance.LAST_UPDATE_HEIGHT_PROPOSAL)
//	if err != nil {
//		return false, errors.Wrap(err, "Unable to set last Update height ")
//	}
//	ctx.Logger.Debug("Governance options set at height : ", ctx.Header.Height, "| propOptions.codeChange.failedFundDistribution :", newDistribution)
//	return true, nil
//}
//func propOptionsgeneralfailedFundDistribution(value interface{}, ctx *Context, validationOnly FunctionBehaviour) (bool, error) {
//	Options, err := ctx.GovernanceStore.GetProposalOptions()
//	if err != nil {
//		return false, err
//	}
//	newDistribution, err := getDistribution(value)
//	if err != nil {
//		return false, err
//	}
//	Options.General.FailedFundDistribution = newDistribution
//	ok, err := ctx.GovernanceStore.ValidateProposal(Options)
//	if err != nil {
//		return false, err
//	}
//	if !ok {
//		return false, errors.New("Validation Failed")
//	}
//	if validationOnly == ValidateOnly {
//		return true, nil
//	}
//	err = ctx.GovernanceStore.WithHeight(ctx.Header.Height).SetProposalOptions(*Options)
//	if err != nil {
//		return false, errors.Wrap(err, "Setup Proposal Options")
//	}
//	ctx.ProposalMasterStore.Proposal.SetOptions(Options)
//	err = ctx.GovernanceStore.WithHeight(ctx.Header.Height).SetLUH(governance.LAST_UPDATE_HEIGHT_PROPOSAL)
//	if err != nil {
//		return false, errors.Wrap(err, "Unable to set last Update height ")
//	}
//	ctx.Logger.Debug("Governance options set at height : ", ctx.Header.Height, "| propOptions.general.failedFundDistribution :", newDistribution)
//	return true, nil
//}

func onsOptionsperBlockFees(value interface{}, ctx *Context, validationOnly FunctionBehaviour) (bool, error) {
	onsOptions, err := ctx.GovernanceStore.GetONSOptions()
	if err != nil {
		return false, err
	}
	newValue, err := getNewValueInt64(value)
	if err != nil {
		return false, err
	}
	onsOptions.PerBlockFees = *balance.NewAmount(newValue)

	ok, err := ctx.GovernanceStore.ValidateONS(onsOptions)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, errors.New("Validation Failed")
	}
	if validationOnly == ValidateOnly {
		return true, nil
	}
	err = ctx.GovernanceStore.WithHeight(ctx.Header.Height).SetONSOptions(*onsOptions)
	if err != nil {
		return false, errors.Wrap(err, "Setup ONS Options")
	}
	ctx.Domains.SetOptions(onsOptions)
	err = ctx.GovernanceStore.WithHeight(ctx.Header.Height).SetLUH(governance.LAST_UPDATE_HEIGHT_ONS)
	if err != nil {
		return false, errors.Wrap(err, "Unable to set last Update height ")
	}
	ctx.Logger.Debug("Governance options set at height : ", ctx.Header.Height, "| onsOptions.perBlockFees :", newValue)
	return true, nil
}

func onsOptionsbaseDomainPrice(value interface{}, ctx *Context, validationOnly FunctionBehaviour) (bool, error) {
	onsOptions, err := ctx.GovernanceStore.GetONSOptions()
	if err != nil {
		return false, err
	}
	newValue, err := getNewValueInt64(value)
	if err != nil {
		return false, err
	}
	onsOptions.BaseDomainPrice = *balance.NewAmount(newValue)

	ok, err := ctx.GovernanceStore.ValidateONS(onsOptions)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, errors.New("Validation Failed")
	}
	if validationOnly == ValidateOnly {
		return true, nil
	}
	err = ctx.GovernanceStore.WithHeight(ctx.Header.Height).SetONSOptions(*onsOptions)
	if err != nil {
		return false, errors.Wrap(err, "Setup ONS Options")
	}
	ctx.Domains.SetOptions(onsOptions)
	err = ctx.GovernanceStore.WithHeight(ctx.Header.Height).SetLUH(governance.LAST_UPDATE_HEIGHT_ONS)
	if err != nil {
		return false, errors.Wrap(err, "Unable to set last Update height ")
	}
	ctx.Logger.Debug("Governance options set at height : ", ctx.Header.Height, "| onsOptions.baseDomainPrice :", newValue)
	return true, nil
}

func feeOptionminFeeDecimal(value interface{}, ctx *Context, validationOnly FunctionBehaviour) (bool, error) {
	feeOptions, err := ctx.GovernanceStore.GetFeeOption()
	if err != nil {
		return false, err
	}
	newValue, err := getNewValueInt64(value)
	if err != nil {
		return false, err
	}
	feeOptions.MinFeeDecimal = newValue

	ok, err := ctx.GovernanceStore.ValidateFee(feeOptions)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, errors.New("Validation Failed")
	}
	if validationOnly == ValidateOnly {
		return true, nil
	}
	err = ctx.GovernanceStore.WithHeight(ctx.Header.Height).SetFeeOption(*feeOptions)
	if err != nil {
		return false, errors.Wrap(err, "Setup Fee Options")
	}
	ctx.FeePool.SetupOpt(feeOptions)
	err = ctx.GovernanceStore.WithHeight(ctx.Header.Height).SetLUH(governance.LAST_UPDATE_HEIGHT_FEE)
	if err != nil {
		return false, errors.Wrap(err, "Unable to set last Update height ")
	}
	ctx.Logger.Debug("Governance options set at height : ", ctx.Header.Height, "feeOption.minFeeDecimal :", newValue)
	return true, nil
}

func getNewValueInt64(value interface{}) (int64, error) {
	newValue, ok := value.(string)
	if !ok {
		return 0, errors.New("Type assertion failed")
	}
	intval, err := strconv.Atoi(newValue)
	if err != nil {
		return 0, err
	}
	return int64(intval), nil
}

//func getDistribution(value interface{}) (newDistribution governance.ProposalFundDistribution, err error) {
//	newValue, ok := value.(map[string]interface{})
//	if !ok {
//		return
//	}
//	err = mapstructure.Decode(newValue, &newDistribution)
//	if err != nil {
//		return
//	}
//	return
//}
//
//func getEvidenceOpt(value interface{}) (newDistribution evidence.Options, err error) {
//	newValue, ok := value.(map[string]interface{})
//	if !ok {
//		return
//	}
//	err = mapstructure.Decode(newValue, &newDistribution)
//	if err != nil {
//		return
//	}
//	return
//}
