package governance

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/kv"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/action/helpers"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/governance"
)

type FinalizeProposal struct {
	ProposalID       governance.ProposalID `json:"proposalId"`
	ValidatorAddress action.Address        `json:"validatorAddress"`
}

func (p FinalizeProposal) Signers() []action.Address {
	return []action.Address{p.ValidatorAddress}
}

func (p FinalizeProposal) Type() action.Type {
	return action.PROPOSAL_FINALIZE
}

func (p FinalizeProposal) Tags() kv.Pairs {
	tags := make([]kv.Pair, 0)

	tag := kv.Pair{
		Key:   []byte("tx.type"),
		Value: []byte(p.Type().String()),
	}
	tag2 := kv.Pair{
		Key:   []byte("tx.Validator"),
		Value: p.ValidatorAddress.Bytes(),
	}
	tag3 := kv.Pair{
		Key:   []byte("tx.proposalID"),
		Value: []byte(string(p.ProposalID)),
	}

	tags = append(tags, tag, tag2, tag3)
	return tags
}

func (p FinalizeProposal) Marshal() ([]byte, error) {
	return json.Marshal(p)
}

func (p *FinalizeProposal) Unmarshal(bytes []byte) error {
	return json.Unmarshal(bytes, p)
}

func (FinalizeProposal) Validate(ctx *action.Context, signedTx action.SignedTx) (bool, error) {
	finalizedProposal := FinalizeProposal{}
	err := finalizedProposal.Unmarshal(signedTx.Data)
	if err != nil {
		return false, errors.Wrap(action.ErrWrongTxType, err.Error())
	}

	//Validate basic signature
	err = action.ValidateBasic(signedTx.RawBytes(), finalizedProposal.Signers(), signedTx.Signatures)
	if err != nil {
		return false, err
	}
	//Check if Proposal ID is valid
	if err = finalizedProposal.ProposalID.Err(); err != nil {
		return false, governance.ErrInvalidProposalId
	}
	return true, nil
}

func (FinalizeProposal) ProcessCheck(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	return runFinalizeProposal(ctx, tx)
}

func (FinalizeProposal) ProcessDeliver(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	return runFinalizeProposal(ctx, tx)
}

func (FinalizeProposal) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas, gasUsed action.Gas) (bool, action.Response) {
	ctx.State.ConsumeVerifySigGas(1)
	ctx.State.ConsumeStorageGas(size)

	// check the used gas for the tx
	final := ctx.Balances.State.ConsumedGas()
	used := int64(final - start)
	ctx.Logger.Detail("Gas Used : ", used)

	return true, action.Response{}
}

func runFinalizeProposal(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	finalizedProposal := FinalizeProposal{}
	err := finalizedProposal.Unmarshal(tx.Data)
	if err != nil {
		return false, action.Response{}
	}
	//Check if already finalized
	_, err = ctx.ProposalMasterStore.Proposal.WithPrefixType(governance.ProposalStateFinalized).Get(finalizedProposal.ProposalID)
	if err == nil {
		result := action.Response{
			Events: action.GetEvent(finalizedProposal.Tags(), "finalize_proposal_success_already"),
		}
		return true, result
	}

	_, err = ctx.ProposalMasterStore.Proposal.WithPrefixType(governance.ProposalStateFinalizeFailed).Get(finalizedProposal.ProposalID)
	if err == nil {
		result := action.Response{
			Events: action.GetEvent(finalizedProposal.Tags(), "finalize_proposal_failed_already"),
		}
		return true, result
	}

	//Get Proposal
	proposal, err := ctx.ProposalMasterStore.Proposal.WithPrefixType(governance.ProposalStatePassed).Get(finalizedProposal.ProposalID)
	if err != nil {
		proposal, err = ctx.ProposalMasterStore.Proposal.WithPrefixType(governance.ProposalStateFailed).Get(finalizedProposal.ProposalID)
		if err != nil {
			return helpers.LogAndReturnFalse(ctx.Logger, governance.ErrProposalNotExists, finalizedProposal.Tags(), err)
		}
	}
	//Check Status is Completed

	if proposal.Status != governance.ProposalStatusCompleted {
		return helpers.LogAndReturnFalse(ctx.Logger, governance.ErrStatusNotCompleted, finalizedProposal.Tags(), err)
	}

	//Get Vote Result
	voteStatus, err := ctx.ProposalMasterStore.ProposalVote.ResultSoFar(proposal.ProposalID, proposal.PassPercentage)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, governance.ErrUnabletoQueryVoteResult, finalizedProposal.Tags(), err)
	}

	//Handle Result TBD
	if voteStatus.Result == governance.VOTE_RESULT_TBD {
		return helpers.LogAndReturnFalse(ctx.Logger, governance.ErrVotingTBD, finalizedProposal.Tags(), err)
	}
	options, err := ctx.GovernanceStore.GetProposalOptionsByType(proposal.Type)
	if err != nil {
		helpers.LogAndReturnFalse(ctx.Logger, governance.ErrGetProposalOptions, finalizedProposal.Tags(), err)
	}
	//Handle Result Passed
	if voteStatus.Result == governance.VOTE_RESULT_PASSED {
		if proposal.Type == governance.ProposalTypeConfigUpdate {
			updates := proposal.GovernanceStateUpdate
			splitstring := strings.Split(updates, ":")
			if len(splitstring) != 2 {
				return helpers.LogAndReturnFalse(ctx.Logger, governance.ErrInvalidOptions, finalizedProposal.Tags(), errors.New("Invalid options string"))
			}
			updatekey := splitstring[0]
			updateValue := splitstring[1]
			updatefunc, ok := ctx.GovUpdate.GovernanceUpdateFunction[updatekey]
			if !ok {
				return helpers.LogAndReturnFalse(ctx.Logger, governance.ErrFinalizeConfigUpdateFailed, finalizedProposal.Tags(), err)
			}
			ok, err = updatefunc(updateValue, ctx, action.ValidateAndUpdate)
			if err != nil {
				ctx.Logger.Debug("Governance auto update failed ", err)
				err = setToFinalizeFailed(ctx, proposal)
				if err != nil {
					return helpers.LogAndReturnFalse(ctx.Logger, governance.ErrStatusUnableToSetFinalizeFailed, finalizedProposal.Tags(), err)
				}

				return helpers.LogAndReturnTrue(ctx.Logger, finalizedProposal.Tags(), fmt.Sprintf("ConfigUpdate_Validation_Failed | %s", proposal.ProposalID))
			}

		}
		proposalDistribution := options.PassedFundDistribution
		distributeErr := distributeFunds(ctx, proposal, &proposalDistribution)
		if distributeErr != nil {
			err = setToFinalizeFailed(ctx, proposal)
			if err != nil {
				return helpers.LogAndReturnFalse(ctx.Logger, governance.ErrStatusUnableToSetFinalizeFailed, finalizedProposal.Tags(), err)
			}
			ctx.Logger.Error("Distribution of Funds failed , Set Proposal to Finalize Failed")
			return helpers.LogAndReturnTrue(ctx.Logger, finalizedProposal.Tags(), governance.ErrFinalizeDistributtionFailed.Wrap(distributeErr).Marshal())
		}

		err = setToFinalizeFromPassed(ctx, proposal)
		if err != nil {
			return helpers.LogAndReturnFalse(ctx.Logger, governance.ErrStatusUnableToSetFinalized, finalizedProposal.Tags(), err)
		}
	}
	//Handle Result Failed
	if voteStatus.Result == governance.VOTE_RESULT_FAILED {
		proposalDistribution := options.FailedFundDistribution
		distributeErr := distributeFunds(ctx, proposal, &proposalDistribution)
		if distributeErr != nil {
			err = setToFinalizeFailed(ctx, proposal)
			if err != nil {
				return helpers.LogAndReturnFalse(ctx.Logger, governance.ErrStatusUnableToSetFinalizeFailed, finalizedProposal.Tags(), err)
			}
			ctx.Logger.Error("Distribution of Funds failed , Set Proposal to Finalize Failed")
			return helpers.LogAndReturnTrue(ctx.Logger, finalizedProposal.Tags(), governance.ErrFinalizeDistributtionFailed.Wrap(distributeErr).Marshal())
		}
		err = setToFinalizeFromFailed(ctx, proposal)
		if err != nil {
			return helpers.LogAndReturnFalse(ctx.Logger, governance.ErrStatusUnableToSetFinalized, finalizedProposal.Tags(), err)
		}
	}
	ctx.Logger.Debug("Finalized  :", finalizedProposal.ValidatorAddress.String(), "Proposal : ", proposal.ProposalID)
	return helpers.LogAndReturnTrue(ctx.Logger, finalizedProposal.Tags(), "finalize_proposal_success")
}

//Function to distribute funds
func distributeFunds(ctx *action.Context, proposal *governance.Proposal, proposalDistribution *governance.ProposalFundDistribution) error {
	// Required Perimeters for Fund Distribution
	fundStore := ctx.ProposalMasterStore.ProposalFund
	totalFunds := fundStore.GetCurrentFundsForProposal(proposal.ProposalID)
	c, ok := ctx.Currencies.GetCurrencyByName("OLT")
	if !ok {
		return action.ErrInvalidCurrency
	}
	totalFundsCoin := c.NewCoinFromAmount(*totalFunds)
	fundTracker := c.NewCoinFromAmount(*totalFunds)
	ctx.Logger.Detailf("totalFundsCoin: ", totalFundsCoin)
	//Starting Fund Distribution
	//Validators
	validatorList, err := ctx.Validators.GetValidatorSet()
	if err != nil {
		return action.ErrGettingValidatorList
	}
	ctx.Logger.Detailf("Transferring to Validators ")
	validatorEarningOLT := getPercentageCoin(&totalFundsCoin, &fundTracker, proposalDistribution.Validators).Divide(len(validatorList))
	for _, v := range validatorList {
		ctx.Logger.Detailf("Validator : \"%v\" : \"%v", v.Address.String(), validatorEarningOLT)
		err = ctx.Balances.AddToAddress(v.Address, validatorEarningOLT)
		if err != nil {
			return err
		}
	}
	ctx.Logger.Detailf("fundTracker: ", fundTracker)

	//Reward for Proposer
	ctx.Logger.Detailf("Transferring to Proposer :\"%v\"", proposal.Proposer.String())
	err = ctx.Balances.AddToAddress(proposal.Proposer, getPercentageCoin(&totalFundsCoin, &fundTracker, proposalDistribution.ProposerReward))
	if err != nil {
		return err
	}
	ctx.Logger.Detailf("fundTracker: ", fundTracker)

	//Bounty Program
	popt, err := ctx.GovernanceStore.GetProposalOptions()
	if err != nil {
		return err
	}
	bountyAddress := action.Address(popt.BountyProgramAddr)
	ctx.Logger.Detailf("Transferring to Bounty Program :\"%v", bountyAddress.String())
	err = ctx.Balances.AddToAddress(bountyAddress, getPercentageCoin(&totalFundsCoin, &fundTracker, proposalDistribution.BountyPool))
	if err != nil {
		return err
	}
	ctx.Logger.Detailf("fundTracker: ", fundTracker)

	//ExecutionCost
	executionAddress := action.Address(ctx.ProposalMasterStore.Proposal.GetOptionsByType(proposal.Type).ProposalExecutionCost)
	ctx.Logger.Detailf("Transferring to Execution Cost :\"%v", executionAddress.String())
	err = ctx.Balances.AddToAddress(executionAddress, getPercentageCoin(&totalFundsCoin, &fundTracker, proposalDistribution.ExecutionCost))
	if err != nil {
		return err
	}
	ctx.Logger.Detailf("fundTracker: ", fundTracker)

	//Subtract Burning Amount From Fund Tracker
	ctx.Logger.Detailf("Subtract Burning Amount From Fund Tracker")
	getPercentageCoin(&totalFundsCoin, &fundTracker, proposalDistribution.Burn)
	ctx.Logger.Detailf("fundTracker: ", fundTracker)

	//Add What's Left in Fund Tracker to Fee Pool(Including Amount Left Due to Inaccuracy)
	ctx.Logger.Detailf("Transferring to Fee Pool")
	err = ctx.FeePool.AddToPool(fundTracker)
	if err != nil {
		return err
	}

	// Burn
	ctx.Logger.Detailf("Burning Funds  ")
	err = fundStore.DeleteAllFunds(proposal.ProposalID)
	if err != nil {
		return err
	}
	return nil
}

//Helper function to get percentage
func getPercentageCoin(totalFunds *balance.Coin, fundTracker *balance.Coin, percentage float64) balance.Coin {
	percentageInt64 := int64(percentage * 10000)
	amount := totalFunds.MultiplyInt64(percentageInt64).DivideInt64(1000000)
	*fundTracker, _ = fundTracker.Minus(amount)
	return amount
}

//Helper to set Proposal to Finalized
func setToFinalizeFromPassed(ctx *action.Context, proposal *governance.Proposal) error {
	err := ctx.ProposalMasterStore.Proposal.WithPrefixType(governance.ProposalStateFinalized).Set(proposal)
	if err != nil {
		return err
	}
	ok, err := ctx.ProposalMasterStore.Proposal.WithPrefixType(governance.ProposalStatePassed).Delete(proposal.ProposalID)
	if !ok {
		ok, err = ctx.ProposalMasterStore.Proposal.WithPrefixType(governance.ProposalStateFinalized).Delete(proposal.ProposalID)
		if !ok {
			return errors.Wrap(err, "error deleting proposal from finalize prefix")
		}
		return err
	}
	return nil
}

func setToFinalizeFromFailed(ctx *action.Context, proposal *governance.Proposal) error {
	err := ctx.ProposalMasterStore.Proposal.WithPrefixType(governance.ProposalStateFinalized).Set(proposal)
	if err != nil {
		return err
	}
	ok, err := ctx.ProposalMasterStore.Proposal.WithPrefixType(governance.ProposalStateFailed).Delete(proposal.ProposalID)
	if !ok {
		ok, err = ctx.ProposalMasterStore.Proposal.WithPrefixType(governance.ProposalStateFinalized).Delete(proposal.ProposalID)
		if !ok {
			return errors.Wrap(err, "error deleting proposal from finalize prefix")
		}
		return err
	}
	return nil
}

func setToFinalizeFailed(ctx *action.Context, proposal *governance.Proposal) error {
	err := ctx.ProposalMasterStore.Proposal.WithPrefixType(governance.ProposalStateFinalizeFailed).Set(proposal)
	if err != nil {
		return err
	}
	ok, err := ctx.ProposalMasterStore.Proposal.WithPrefixType(governance.ProposalStatePassed).Delete(proposal.ProposalID)
	if !ok {
		ok, err = ctx.ProposalMasterStore.Proposal.WithPrefixType(governance.ProposalStateFinalizeFailed).Delete(proposal.ProposalID)
		if !ok {
			return errors.Wrap(err, "error deleting proposal from finalize prefix")
		}
		return err
	}
	return nil
}
