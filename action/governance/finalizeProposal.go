package governance

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/kv"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/governance"
)

type FinalizeProposal struct {
	ProposalID       governance.ProposalID
	ValidatorAddress action.Address
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

var _ action.Tx = &finalizeProposalTx{}

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
	//Validate Fee for funding request
	err = action.ValidateFee(ctx.FeePool.GetOpt(), signedTx.Fee)
	if err != nil {
		return false, err
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
	return action.BasicFeeHandling(ctx, signedTx, start, size, 1)
}

func runFinalizeProposal(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	finalizedProposal := FinalizeProposal{}
	err := finalizedProposal.Unmarshal(tx.Data)
	if err != nil {
		return false, action.Response{}
	}
	proposal, err := ctx.ProposalMasterStore.Proposal.WithPrefixType(governance.ProposalStatePassed).Get(finalizedProposal.ProposalID)
	if err != nil {
		result := action.Response{
			Events: action.GetEvent(finalizedProposal.Tags(), action.ErrProposalNotFound.Msg),
			Log:    action.ErrProposalNotFound.Error(),
		}
		return false, result
	}
	if proposal.Status != governance.ProposalStatusCompleted {
		result := action.Response{
			Events: action.GetEvent(finalizedProposal.Tags(), action.ErrStatusNotCompleted.Msg),
			Log:    action.ErrStatusNotCompleted.Error(),
		}
		return false, result
	}

	voteResult, err := ctx.ProposalMasterStore.ProposalVote.ResultSoFar(proposal.ProposalID, proposal.PassPercentage)
	if err != nil {
		result := action.Response{
			Events: action.GetEvent(finalizedProposal.Tags(), action.ErrUnabletoQueryVoteResult.Msg),
			Log:    action.ErrUnabletoQueryVoteResult.Error(),
		}
		return false, result
	}
	if voteResult == governance.VOTE_RESULT_TBD {
		result := action.Response{
			Events: action.GetEvent(finalizedProposal.Tags(), action.ErrVotingTBD.Msg),
			Log:    action.ErrVotingTBD.Error(),
		}
		return false, result
	}
	if voteResult == governance.VOTE_RESULT_PASSED {
		proposalDistribution := ctx.ProposalMasterStore.Proposal.GetOptionsByType(proposal.Type).PassedFundDistribution
		err := distributeFunds(ctx, proposal, &proposalDistribution)
		if err != nil {
			result := action.Response{
				Events: action.GetEvent(finalizedProposal.Tags(), action.ErrFinalizeDistributtionFailed.Msg),
				Log:    errors.Wrap(action.ErrFinalizeDistributtionFailed, err.Error()).Error(),
			}
			return false, result
		}
		if proposal.Type == governance.ProposalTypeConfigUpdate {
			err := executeConfigUpdate(ctx, proposal)
			if err != nil {
				result := action.Response{
					Events: action.GetEvent(finalizedProposal.Tags(), action.ErrFinalizeConfigUpdateFailed.Msg),
					Log:    errors.Wrap(action.ErrFinalizeConfigUpdateFailed, err.Error()).Error(),
				}
				return false, result
			}
		}
		proposal.Status = governance.ProposalStatusFinalized
		err = ctx.ProposalMasterStore.Proposal.WithPrefixType(governance.ProposalStatePassed).Set(proposal)
		if err != nil {
			result := action.Response{
				Events: action.GetEvent(finalizedProposal.Tags(), action.ErrStatusUnableToSetFinalized.Msg),
				Log:    action.ErrStatusUnableToSetFinalized.Error(),
			}
			return false, result

		}
	}
	if voteResult == governance.VOTE_RESULT_FAILED {
		proposalDistribution := ctx.ProposalMasterStore.Proposal.GetOptionsByType(proposal.Type).FailedFundDistribution
		err := distributeFunds(ctx, proposal, &proposalDistribution)
		if err != nil {
			result := action.Response{
				Events: action.GetEvent(finalizedProposal.Tags(), action.ErrFinalizeDistributtionFailed.Msg),
				Log:    errors.Wrap(action.ErrFinalizeDistributtionFailed, err.Error()).Error(),
			}
			return false, result
		}
		proposal.Status = governance.ProposalStatusFinalized
		err = ctx.ProposalMasterStore.Proposal.WithPrefixType(governance.ProposalStateFailed).Set(proposal)
		if err != nil {
			result := action.Response{
				Events: action.GetEvent(finalizedProposal.Tags(), action.ErrStatusUnableToSetFinalized.Msg),
				Log:    action.ErrStatusUnableToSetFinalized.Error(),
			}
			return false, result

		}
	}
	result := action.Response{
		Events: action.GetEvent(finalizedProposal.Tags(), "finalize_proposal_success"),
	}
	return true, result
}

func distributeFunds(ctx *action.Context, proposal *governance.Proposal, proposalDistribution *governance.ProposalFundDistribution) error {
	// Required Perimeters for Fund Distribution

	totalFunding := governance.GetCurrentFunds(proposal.ProposalID, ctx.ProposalMasterStore.ProposalFund).Float()
	c, ok := ctx.Currencies.GetCurrencyByName("OLT")
	if !ok {
		return errors.New("fund_proposal_olt_unavailable")
	}
	fundTracker, _ := totalFunding.Float64()
	//Starting Fund Distribution
	//Validators
	validatorList, err := ctx.Validators.GetValidatorSet()
	if err != nil {
		return err
	}
	validatorEarningOLT := getPercentageCoin(c, totalFunding, &fundTracker, proposalDistribution.Validators).Divide(len(validatorList))

	for _, v := range validatorList {
		err = ctx.Balances.AddToAddress(v.Address, validatorEarningOLT)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("Unable to send funds to Validator :%s", v.Name))
		}
	}

	//Fee Pool
	err = ctx.FeePool.AddToPool(getPercentageCoin(c, totalFunding, &fundTracker, proposalDistribution.FeePool))
	if err != nil {
		return errors.Wrap(err, "Failed in adding to feepool")
	}

	//Reward for Proposer
	err = ctx.Balances.AddToAddress(proposal.Proposer, getPercentageCoin(c, totalFunding, &fundTracker, proposalDistribution.ProposerReward))
	if err != nil {
		return errors.Wrap(err, "Failed in rewarding Proposer")
	}

	//Bounty Program

	bountyAddress := action.Address(ctx.ProposalMasterStore.Proposal.GetOptions().BountyProgramAddr)
	err = ctx.Balances.AddToAddress(bountyAddress, getPercentageCoin(c, totalFunding, &fundTracker, proposalDistribution.BountyPool))
	if err != nil {
		return errors.Wrap(err, "Failed in adding to Adding to bounty program")
	}

	//ExecutionCost
	executionAddress := action.Address(ctx.ProposalMasterStore.Proposal.GetOptionsByType(proposal.Type).ProposalExecutionCost)
	err = ctx.Balances.AddToAddress(executionAddress, getPercentageCoin(c, totalFunding, &fundTracker, proposalDistribution.ExecutionCost))
	if err != nil {
		return errors.Wrap(err, "Failed in adding to Adding to Execution Address")
	}

	//Burn
	getPercentageCoin(c, totalFunding, &fundTracker, proposalDistribution.Burn)
	if fundTracker != 0 {
		return errors.New(fmt.Sprintf("Extra Funding Amount Left %s", fundTracker))
	}
	err = governance.DeleteAllFunds(proposal.ProposalID, ctx.ProposalMasterStore.ProposalFund)
	if err != nil {
		return errors.Wrap(err, "Unable to Burn all funds")
	}
	return nil
}
func executeConfigUpdate(ctx *action.Context, proposal *governance.Proposal) error {
	return nil
}
func getPercentageCoin(c balance.Currency, totalFunding *big.Float, fundTracker *float64, percentage float64) balance.Coin {
	// TODO : How to deal with accuracy
	amount, _ := big.NewFloat(1.0).Mul(totalFunding, big.NewFloat(percentage/100)).Float64()
	*fundTracker = *fundTracker - amount
	return c.NewCoinFromFloat64(amount)
}
