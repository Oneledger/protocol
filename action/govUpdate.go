package action

import (
	"github.com/pkg/errors"

	"github.com/Oneledger/protocol/data/governance"
)

type FunctionBehaviour int

const (
	ValidateOnly      FunctionBehaviour = 0
	ValidateAndUpdate FunctionBehaviour = 1
)

type GovernaceUpdateAndValidate struct {
	GovernanceUpdateFunction map[string]func(*governance.GovernanceState, *Context, FunctionBehaviour) (bool, error)
}

func NewGovUpdate() *GovernaceUpdateAndValidate {

	newUpdateAndValidate := GovernaceUpdateAndValidate{
		GovernanceUpdateFunction: make(map[string]func(*governance.GovernanceState, *Context, FunctionBehaviour) (bool, error)),
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
	g.GovernanceUpdateFunction["propOptions.configUpdate.passedFundDistribution"] = propOptionsconfigUpdatepassedFundDistribution
	g.GovernanceUpdateFunction["propOptions.codeChange.passedFundDistribution"] = propOptionscodeChangepassedFundDistribution
	g.GovernanceUpdateFunction["propOptions.general.passedFundDistribution"] = propOptionsgeneralpassedFundDistribution
	g.GovernanceUpdateFunction["propOptions.configUpdate.failedFundDistribution"] = propOptionsconfigUpdatefailedFundDistribution
	g.GovernanceUpdateFunction["propOptions.codeChange.failedFundDistribution"] = propOptionscodeChangefailedFundDistribution
	g.GovernanceUpdateFunction["propOptions.general.failedFundDistribution"] = propOptionsgeneralfailedFundDistribution

}
func stakingOptionsminSelfDelegationAmount(newstate *governance.GovernanceState, ctx *Context, validationOnly FunctionBehaviour) (bool, error) {
	return true, nil
}
func stakingOptionstopValidatorCount(newstate *governance.GovernanceState, ctx *Context, validationOnly FunctionBehaviour) (bool, error) {
	return true, nil
}
func stakingOptionsmaturityTime(newstate *governance.GovernanceState, ctx *Context, validationOnly FunctionBehaviour) (bool, error) {
	return true, nil
}
func propOptionsconfigUpdateinitialFunding(newstate *governance.GovernanceState, ctx *Context, validationOnly FunctionBehaviour) (bool, error) {
	return true, nil
}
func propOptionscodeChangeinitialFunding(newstate *governance.GovernanceState, ctx *Context, validationOnly FunctionBehaviour) (bool, error) {
	return true, nil
}
func propOptionsgeneralinitialFunding(newstate *governance.GovernanceState, ctx *Context, validationOnly FunctionBehaviour) (bool, error) {
	return true, nil
}
func propOptionsconfigUpdatefundingGoal(newstate *governance.GovernanceState, ctx *Context, validationOnly FunctionBehaviour) (bool, error) {
	return true, nil
}
func propOptionscodeChangefundingGoal(newstate *governance.GovernanceState, ctx *Context, validationOnly FunctionBehaviour) (bool, error) {
	return true, nil
}
func propOptionsgeneralfundingGoal(newstate *governance.GovernanceState, ctx *Context, validationOnly FunctionBehaviour) (bool, error) {
	return true, nil
}
func propOptionsconfigUpdatevotingDeadline(newstate *governance.GovernanceState, ctx *Context, validationOnly FunctionBehaviour) (bool, error) {
	return true, nil
}
func propOptionscodeChangevotingDeadline(newstate *governance.GovernanceState, ctx *Context, validationOnly FunctionBehaviour) (bool, error) {
	return true, nil
}
func propOptionsgeneralvotingDeadline(newstate *governance.GovernanceState, ctx *Context, validationOnly FunctionBehaviour) (bool, error) {
	return true, nil
}
func propOptionsconfigUpdatefundingDeadline(newstate *governance.GovernanceState, ctx *Context, validationOnly FunctionBehaviour) (bool, error) {
	return true, nil
}
func propOptionscodeChangefundingDeadline(newstate *governance.GovernanceState, ctx *Context, validationOnly FunctionBehaviour) (bool, error) {
	return true, nil
}
func propOptionsgeneralfundingDeadline(newstate *governance.GovernanceState, ctx *Context, validationOnly FunctionBehaviour) (bool, error) {
	return true, nil
}
func propOptionsconfigUpdatepassPercentage(newstate *governance.GovernanceState, ctx *Context, validationOnly FunctionBehaviour) (bool, error) {
	return true, nil
}
func propOptionscodeChangepassPercentage(newstate *governance.GovernanceState, ctx *Context, validationOnly FunctionBehaviour) (bool, error) {
	return true, nil
}
func propOptionsgeneralpassPercentage(newstate *governance.GovernanceState, ctx *Context, validationOnly FunctionBehaviour) (bool, error) {
	return true, nil
}
func propOptionsconfigUpdatepassedFundDistribution(newstate *governance.GovernanceState, ctx *Context, validationOnly FunctionBehaviour) (bool, error) {
	return true, nil
}
func propOptionscodeChangepassedFundDistribution(newstate *governance.GovernanceState, ctx *Context, validationOnly FunctionBehaviour) (bool, error) {
	return true, nil
}
func propOptionsgeneralpassedFundDistribution(newstate *governance.GovernanceState, ctx *Context, validationOnly FunctionBehaviour) (bool, error) {
	return true, nil
}
func propOptionsconfigUpdatefailedFundDistribution(newstate *governance.GovernanceState, ctx *Context, validationOnly FunctionBehaviour) (bool, error) {
	return true, nil
}
func propOptionscodeChangefailedFundDistribution(newstate *governance.GovernanceState, ctx *Context, validationOnly FunctionBehaviour) (bool, error) {
	return true, nil
}
func propOptionsgeneralfailedFundDistribution(newstate *governance.GovernanceState, ctx *Context, validationOnly FunctionBehaviour) (bool, error) {
	return true, nil
}

func onsOptionsperBlockFees(newstate *governance.GovernanceState, ctx *Context, validationOnly FunctionBehaviour) (bool, error) {
	onsOptions, err := ctx.GovernanceStore.GetONSOptions()
	if err != nil {
		return false, err
	}
	onsOptions.PerBlockFees = newstate.ONSOptions.PerBlockFees

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
	ctx.Logger.Debug("Governance options set at height : ", ctx.Header.Height, "| onsOptions.perBlockFees :", newstate.ONSOptions.PerBlockFees)
	return true, nil
}

func onsOptionsbaseDomainPrice(newstate *governance.GovernanceState, ctx *Context, validationOnly FunctionBehaviour) (bool, error) {
	onsOptions, err := ctx.GovernanceStore.GetONSOptions()
	if err != nil {
		return false, err
	}
	onsOptions.BaseDomainPrice = newstate.ONSOptions.BaseDomainPrice

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
	ctx.Logger.Debug("Governance options set at height : ", ctx.Header.Height, "| onsOptions.baseDomainPrice :", newstate.ONSOptions.BaseDomainPrice)
	return true, nil
}

func feeOptionminFeeDecimal(newstate *governance.GovernanceState, ctx *Context, validationOnly FunctionBehaviour) (bool, error) {
	feeOptions, err := ctx.GovernanceStore.GetFeeOption()
	if err != nil {
		return false, err
	}
	feeOptions.MinFeeDecimal = newstate.FeeOption.MinFeeDecimal

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
	ctx.Logger.Debug("Governance options set at height : ", ctx.Header.Height, "feeOption.minFeeDecimal :", newstate.FeeOption.MinFeeDecimal)
	return true, nil
}
