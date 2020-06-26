package governance

import (
	"encoding/json"
	"fmt"
	"math/big"

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

type finalizeProposalTx struct {
}

var _ action.Tx = finalizeProposalTx{}

func (finalizeProposalTx) Validate(ctx *action.Context, signedTx action.SignedTx) (bool, error) {
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

func (finalizeProposalTx) ProcessCheck(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	return runFinalizeProposal(ctx, tx)
}

func (finalizeProposalTx) ProcessDeliver(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	return runFinalizeProposal(ctx, tx)
}

func (finalizeProposalTx) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas) (bool, action.Response) {
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
		return helpers.LogAndReturnFalse(ctx.Logger, governance.ErrProposalNotExists, finalizedProposal.Tags(), err)
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

	//Handle Result Passed
	if voteStatus.Result == governance.VOTE_RESULT_PASSED {
		proposalDistribution := ctx.ProposalMasterStore.Proposal.GetOptionsByType(proposal.Type).PassedFundDistribution
		distributeErr := distributeFunds(ctx, proposal, &proposalDistribution)
		if distributeErr != nil {
			err = setToFinalizeFailed(ctx, proposal)
			if err != nil {
				return helpers.LogAndReturnFalse(ctx.Logger, governance.ErrStatusUnableToSetFinalizeFailed, finalizedProposal.Tags(), err)
			}
			ctx.Logger.Error("Distribuition of Funds failed , Set Proposal to Finalize Failed")
			return helpers.LogAndReturnTrue(ctx.Logger, finalizedProposal.Tags(), governance.ErrFinalizeDistributtionFailed.Wrap(distributeErr).Marshal())
		}

		if proposal.Type == governance.ProposalTypeConfigUpdate {
			err := executeConfigUpdate(ctx, proposal)
			if err != nil {
				return helpers.LogAndReturnFalse(ctx.Logger, governance.ErrFinalizeConfigUpdateFailed, finalizedProposal.Tags(), err)
			}
		}
		err = setToFinalize(ctx, proposal)
		if err != nil {
			return helpers.LogAndReturnFalse(ctx.Logger, governance.ErrStatusUnableToSetFinalized, finalizedProposal.Tags(), err)
		}
	}
	//Handle Result Failed
	if voteStatus.Result == governance.VOTE_RESULT_FAILED {
		proposalDistribution := ctx.ProposalMasterStore.Proposal.GetOptionsByType(proposal.Type).FailedFundDistribution
		distributeErr := distributeFunds(ctx, proposal, &proposalDistribution)
		if distributeErr != nil {
			err = setToFinalizeFailed(ctx, proposal)
			if err != nil {
				return helpers.LogAndReturnFalse(ctx.Logger, governance.ErrStatusUnableToSetFinalizeFailed, finalizedProposal.Tags(), err)
			}
			ctx.Logger.Error("Distribution of Funds failed , Set Proposal to Finalize Failed")
			return helpers.LogAndReturnTrue(ctx.Logger, finalizedProposal.Tags(), governance.ErrFinalizeDistributtionFailed.Wrap(distributeErr).Marshal())
		}
		err = setToFinalize(ctx, proposal)
		if err != nil {
			return helpers.LogAndReturnFalse(ctx.Logger, governance.ErrStatusUnableToSetFinalized, finalizedProposal.Tags(), err)
		}
	}
	fmt.Println("Finalized Validator :", finalizedProposal.ValidatorAddress.String(), "Proposal : ", proposal.ProposalID)
	return helpers.LogAndReturnTrue(ctx.Logger, finalizedProposal.Tags(), "finalize_proposal_success")
}

//Function to distribute funds
func distributeFunds(ctx *action.Context, proposal *governance.Proposal, proposalDistribution *governance.ProposalFundDistribution) error {
	// Required Perimeters for Fund Distribution
	fundStore := ctx.ProposalMasterStore.ProposalFund
	totalFunding := fundStore.GetCurrentFundsForProposal(proposal.ProposalID).Float()
	c, ok := ctx.Currencies.GetCurrencyByName("OLT")
	if !ok {
		return action.ErrInvalidCurrency
	}
	fundTracker, _ := totalFunding.Float64()
	//Starting Fund Distribution
	//Validators
	validatorList, err := ctx.Validators.GetValidatorSet()
	if err != nil {
		return action.ErrGettingValidatorList
	}
	ctx.Logger.Detailf("Transferring to Validators ")
	validatorEarningOLT := getPercentageCoin(c, totalFunding, &fundTracker, proposalDistribution.Validators).Divide(len(validatorList))

	for _, v := range validatorList {
		ctx.Logger.Detailf("Validator : \"%v\" : \"%v", v.Address.String(), validatorEarningOLT)
		err = ctx.Balances.AddToAddress(v.Address, validatorEarningOLT)
		if err != nil {
			return err
		}
	}

	//Fee Pool
	ctx.Logger.Detailf("Transferring to Fee Pool")
	err = ctx.FeePool.AddToPool(getPercentageCoin(c, totalFunding, &fundTracker, proposalDistribution.FeePool))
	if err != nil {
		return err
	}
	//Reward for Proposer
	ctx.Logger.Detailf("Transferring to Proposer :\"%v\"", proposal.Proposer.String())
	err = ctx.Balances.AddToAddress(proposal.Proposer, getPercentageCoin(c, totalFunding, &fundTracker, proposalDistribution.ProposerReward))
	if err != nil {
		return err
	}

	//Bounty Program

	bountyAddress := action.Address(ctx.ProposalMasterStore.Proposal.GetOptions().BountyProgramAddr)
	ctx.Logger.Detailf("Transferring to Bounty Program :\"%v", bountyAddress.String())
	err = ctx.Balances.AddToAddress(bountyAddress, getPercentageCoin(c, totalFunding, &fundTracker, proposalDistribution.BountyPool))
	if err != nil {
		return err
	}

	//ExecutionCost
	executionAddress := action.Address(ctx.ProposalMasterStore.Proposal.GetOptionsByType(proposal.Type).ProposalExecutionCost)
	ctx.Logger.Detailf("Transferring to Execution Cost :\"%v", executionAddress.String())
	err = ctx.Balances.AddToAddress(executionAddress, getPercentageCoin(c, totalFunding, &fundTracker, proposalDistribution.ExecutionCost))
	if err != nil {
		return err
	}

	//Burn
	ctx.Logger.Detailf("Transferring to Burn Address ")
	getPercentageCoin(c, totalFunding, &fundTracker, proposalDistribution.Burn)
	if fundTracker != 0 {
		return errors.New(fmt.Sprintf("Extra Funding Amount Left %s", fundTracker))
	}
	err = fundStore.DeleteAllFunds(proposal.ProposalID)
	if err != nil {
		return err
	}
	return nil
}

//Function to execute Config update to governanace
func executeConfigUpdate(ctx *action.Context, proposal *governance.Proposal) error {
	return nil
}

//Helper function to get percentage
func getPercentageCoin(c balance.Currency, totalFunding *big.Float, fundTracker *float64, percentage float64) balance.Coin {
	// TODO : How to deal with accuracy
	amount, _ := big.NewFloat(1.0).Mul(totalFunding, big.NewFloat(percentage/100)).Float64()
	//fmt.Println("-------> Transferred ", c.NewCoinFromFloat64(amount))
	*fundTracker = *fundTracker - amount
	return c.NewCoinFromFloat64(amount)
}

//Helper to set Proposal to Finalized
func setToFinalize(ctx *action.Context, proposal *governance.Proposal) error {
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
