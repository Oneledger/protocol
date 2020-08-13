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
	g.GovernanceUpdateFunction["feeOption.minFeeDecimal"] = updateMinDecimal
	g.GovernanceUpdateFunction["onsOptions.perBlockFees"] = updatePerBlockFee
	g.GovernanceUpdateFunction["onsOptions.baseDomainPrice"] = updateBaseDomainPrice
}

func updatePerBlockFee(newstate *governance.GovernanceState, ctx *Context, validationOnly FunctionBehaviour) (bool, error) {
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

func updateBaseDomainPrice(newstate *governance.GovernanceState, ctx *Context, validationOnly FunctionBehaviour) (bool, error) {
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

func updateMinDecimal(newstate *governance.GovernanceState, ctx *Context, validationOnly FunctionBehaviour) (bool, error) {
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
