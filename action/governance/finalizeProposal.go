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
	"github.com/Oneledger/protocol/data/keys"
)

type FinalizeProposal struct {
	ProposalID governance.ProposalID
	Proposer   keys.Address
}

func (p FinalizeProposal) Signers() []action.Address {
	return []action.Address{p.Proposer}
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
		Key:   []byte("tx.proposer"),
		Value: p.Proposer.Bytes(),
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
	passedProposal := FinalizeProposal{}
	err := passedProposal.Unmarshal(signedTx.Data)
	if err != nil {
		return false, errors.Wrap(action.ErrWrongTxType, err.Error())
	}

	//Validate basic signature
	err = action.ValidateBasic(signedTx.RawBytes(), passedProposal.Signers(), signedTx.Signatures)
	if err != nil {
		return false, err
	}
	//Validate Fee for funding request
	err = action.ValidateFee(ctx.FeePool.GetOpt(), signedTx.Fee)
	if err != nil {
		return false, err
	}

	//Check if Funder address is valid oneLedger address
	err = passedProposal.Proposer.Err()
	if err != nil {
		return false, errors.Wrap(err, "invalid proposer address")
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
	passedProposal := FinalizeProposal{}
	err := passedProposal.Unmarshal(tx.Data)
	if err != nil {
		return false, action.Response{}
	}
	proposal, err := ctx.ProposalMasterStore.Proposal.Get(passedProposal.ProposalID)
	if err != nil {
		result := action.Response{
			Events: action.GetEvent(passedProposal.Tags(), "finalize_proposal_not_found"),
		}
		return false, result
	}
	if proposal.Status != governance.ProposalStatusCompleted {
		result := action.Response{
			Events: action.GetEvent(passedProposal.Tags(), "finalize_proposal_not_completed"),
		}
		return false, result
	}

	voteResult, err := ctx.ProposalMasterStore.ProposalVote.ResultSoFar(proposal.ProposalID, proposal.PassPercentage)
	if err != nil {
		result := action.Response{
			Events: action.GetEvent(passedProposal.Tags(), "finalize_proposal_vote_result_unavailable"),
		}
		return false, result
	}
	if voteResult == governance.VOTE_RESULT_TBD {
		result := action.Response{
			Events: action.GetEvent(passedProposal.Tags(), "finalize_proposal_vote_result_TBD"),
		}
		return false, result
	}
	if voteResult == governance.VOTE_RESULT_PASSED {
		proposalDistribution := ctx.ProposalMasterStore.Proposal.GetOptionsByType(proposal.Type).PassedFundDistribution
		err := distributeFunds(ctx, proposal, &proposalDistribution)
		if err != nil {
			result := action.Response{
				Events: action.GetEvent(passedProposal.Tags(), "finalize_proposal_passed_distribution_failed"),
			}
			return false, result
		}
	}
	if voteResult == governance.VOTE_RESULT_FAILED {
		proposalDistribution := ctx.ProposalMasterStore.Proposal.GetOptionsByType(proposal.Type).FailedFundDistribution
		err := distributeFunds(ctx, proposal, &proposalDistribution)
		if err != nil {
			result := action.Response{
				Events: action.GetEvent(passedProposal.Tags(), "finalize_proposal_failed_distribution_failed"),
			}
			return false, result
		}
	}
	result := action.Response{
		Events: action.GetEvent(passedProposal.Tags(), "finalize_proposal_success"),
	}
	return true, result
}

func distributeFunds(ctx *action.Context, proposal *governance.Proposal, proposalDistribution *governance.ProposalFundDistribution) error {
	// Required Perimeters for Fund Distribution
	//proposalDistribution := ctx.ProposalMasterStore.Proposal.GetOptionsByType(proposal.Type).PassedFundDistribution
	totalFunding := governance.GetCurrentFunds(proposal.ProposalID, ctx.ProposalMasterStore.ProposalFund).Float()
	c, ok := ctx.Currencies.GetCurrencyByName("OLT")
	if !ok {
		return errors.New("fund_proposal_olt_unavailable")
	}
	//Fund Distribution
	//Validators
	validatorList, err := ctx.Validators.GetValidatorSet()
	if err != nil {
		return err
	}
	validatorEarningOLT := getPercentageCoin(c, totalFunding, proposalDistribution.Validators).Divide(len(validatorList))
	//TODO:Instead of break continue sending to rest of the validators ,and report failed Validator
	for _, v := range validatorList {
		err = ctx.Balances.AddToAddress(v.Address, validatorEarningOLT)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("Unable to send funds to Validator :%s", v.Name))
		}
	}

	//Fee Pool
	err = ctx.FeePool.AddToPool(getPercentageCoin(c, totalFunding, proposalDistribution.FeePool))
	if err != nil {
		return errors.Wrap(err, "Failed in adding to feepool")
	}
	//Reward for Proposer
	err = ctx.Balances.AddToAddress(proposal.Proposer, getPercentageCoin(c, totalFunding, proposalDistribution.ProposerReward))
	if err != nil {
		return errors.Wrap(err, "Failed in rewarding Proposer")
	}
	//Bounty Program
	bountyAddress := action.Address(ctx.ProposalMasterStore.Proposal.GetOptions().BountyProgramAddr)
	err = ctx.Balances.AddToAddress(bountyAddress, getPercentageCoin(c, totalFunding, proposalDistribution.BountyPool))
	if err != nil {
		return errors.Wrap(err, "Failed in adding to Adding to bounty program")
	}
	//Burn
	burnAmount := getPercentageCoin(c, totalFunding, proposalDistribution.Burn)
	fmt.Println(burnAmount)
	return nil

}

func getPercentageCoin(c balance.Currency, totalFunding *big.Float, percentage float64) balance.Coin {
	// TODO : How to deal with accuracy
	amount, _ := totalFunding.Mul(totalFunding, big.NewFloat(percentage/100)).Float64()
	return c.NewCoinFromFloat64(amount)
}
