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
	"github.com/Oneledger/protocol/status_codes"
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
	//Validate proposal ID
	if len(finalizedProposal.ProposalID) <= 0 {
		return false, errors.New("invalid proposal id")
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
			Events: action.GetEvent(finalizedProposal.Tags(), "finalize_proposal_success"),
		}
		return true, result
	}
	//Get Proposal
	proposal, err := ctx.ProposalMasterStore.Proposal.WithPrefixType(governance.ProposalStatePassed).Get(finalizedProposal.ProposalID)
	if err != nil {
		return logAndReturnFalse(ctx.Logger, governance.ErrProposalNotFound, finalizedProposal.Tags())
	}
	//Check Status is Completed
	fmt.Println("Validator :", finalizedProposal.ValidatorAddress.String(), "Proposal : ", proposal.ProposalID)

	if proposal.Status != governance.ProposalStatusCompleted {
		return logAndReturnFalse(ctx.Logger, governance.ErrStatusNotCompleted, finalizedProposal.Tags())
	}

	//Get Vote Result
	voteResult, err := ctx.ProposalMasterStore.ProposalVote.ResultSoFar(proposal.ProposalID, proposal.PassPercentage)
	if err != nil {
		return logAndReturnFalse(ctx.Logger, governance.ErrUnabletoQueryVoteResult, finalizedProposal.Tags())
	}

	//Handle Result TBD
	if voteResult == governance.VOTE_RESULT_TBD {
		return logAndReturnFalse(ctx.Logger, governance.ErrVotingTBD, finalizedProposal.Tags())
	}

	//Handle Result Passed
	if voteResult == governance.VOTE_RESULT_PASSED {
		proposalDistribution := ctx.ProposalMasterStore.Proposal.GetOptionsByType(proposal.Type).PassedFundDistribution
		protocolErr := distributeFunds(ctx, proposal, &proposalDistribution)
		if protocolErr != governance.NoError {
			return logAndReturnFalse(ctx.Logger, *governance.ErrFinalizeDistributtionFailed.Wrap(protocolErr), finalizedProposal.Tags())
		}
		if proposal.Type == governance.ProposalTypeConfigUpdate {
			err := executeConfigUpdate(ctx, proposal)
			if err != nil {
				return logAndReturnFalse(ctx.Logger, governance.ErrFinalizeConfigUpdateFailed, finalizedProposal.Tags())
			}
		}
		protocolErr = setToFinalize(ctx, proposal)
		if protocolErr != governance.NoError {
			return logAndReturnFalse(ctx.Logger, protocolErr, finalizedProposal.Tags())
		}
	}

	//Handle Result Failed
	if voteResult == governance.VOTE_RESULT_FAILED {
		proposalDistribution := ctx.ProposalMasterStore.Proposal.GetOptionsByType(proposal.Type).FailedFundDistribution
		protocolErr := distributeFunds(ctx, proposal, &proposalDistribution)
		if protocolErr != governance.NoError {
			return logAndReturnFalse(ctx.Logger, *governance.ErrFinalizeDistributtionFailed.Wrap(protocolErr), finalizedProposal.Tags())
		}
		protocolErr = setToFinalize(ctx, proposal)
		if protocolErr != governance.NoError {
			return logAndReturnFalse(ctx.Logger, protocolErr, finalizedProposal.Tags())
		}
	}
	fmt.Println("Finalized Validator :", finalizedProposal.ValidatorAddress.String(), "Proposal : ", proposal.ProposalID)
	result := action.Response{
		Events: action.GetEvent(finalizedProposal.Tags(), "finalize_proposal_success"),
	}
	return true, result
}

//Function to distribute funds
func distributeFunds(ctx *action.Context, proposal *governance.Proposal, proposalDistribution *governance.ProposalFundDistribution) status_codes.ProtocolError {
	// Required Perimeters for Fund Distribution

	totalFunding := governance.GetCurrentFunds(proposal.ProposalID, ctx.ProposalMasterStore.ProposalFund).Float()
	c, ok := ctx.Currencies.GetCurrencyByName("OLT")
	if !ok {
		return action.ErrInvalidCurrency
	}
	fundTracker, _ := totalFunding.Float64()
	//Starting Fund Distribution
	//Validators
	validatorList, err := ctx.Validators.GetValidatorSet()
	if err != nil {
		return action.ErrValidatorsUnableGetList
	}
	validatorEarningOLT := getPercentageCoin(c, totalFunding, &fundTracker, proposalDistribution.Validators).Divide(len(validatorList))

	for _, v := range validatorList {
		err = ctx.Balances.AddToAddress(v.Address, validatorEarningOLT)
		if err != nil {
			return *balance.ErrBalanceErrorAddFailed.Wrap(errors.New("Distribute to validators"))
		}
	}

	//Fee Pool
	err = ctx.FeePool.AddToPool(getPercentageCoin(c, totalFunding, &fundTracker, proposalDistribution.FeePool))
	if err != nil {
		return *balance.ErrBalanceErrorAddFailed.Wrap(errors.New("Distribute to Fee Pool"))
	}

	//Reward for Proposer
	err = ctx.Balances.AddToAddress(proposal.Proposer, getPercentageCoin(c, totalFunding, &fundTracker, proposalDistribution.ProposerReward))
	if err != nil {
		return *balance.ErrBalanceErrorAddFailed.Wrap(errors.New("Distribute to Proposer"))
	}

	//Bounty Program

	bountyAddress := action.Address(ctx.ProposalMasterStore.Proposal.GetOptions().BountyProgramAddr)
	err = ctx.Balances.AddToAddress(bountyAddress, getPercentageCoin(c, totalFunding, &fundTracker, proposalDistribution.BountyPool))
	if err != nil {
		return *balance.ErrBalanceErrorAddFailed.Wrap(errors.New("Distribute to Bounty Program"))
	}

	//ExecutionCost
	executionAddress := action.Address(ctx.ProposalMasterStore.Proposal.GetOptionsByType(proposal.Type).ProposalExecutionCost)
	err = ctx.Balances.AddToAddress(executionAddress, getPercentageCoin(c, totalFunding, &fundTracker, proposalDistribution.ExecutionCost))
	if err != nil {
		return *balance.ErrBalanceErrorAddFailed.Wrap(errors.New("Distribute to Execution Cost"))
	}

	//Burn
	getPercentageCoin(c, totalFunding, &fundTracker, proposalDistribution.Burn)
	if fundTracker != 0 {
		return *governance.ErrGovFundBalanceMismatch.Wrap(errors.New(fmt.Sprintf("Extra Funding Amount Left %s", fundTracker)))
	}
	err = governance.DeleteAllFunds(proposal.ProposalID, ctx.ProposalMasterStore.ProposalFund)
	if err != nil {
		return governance.ErrGovFundUnableToDelete
	}
	return governance.NoError
}

//Function to execute Config update to governanace
func executeConfigUpdate(ctx *action.Context, proposal *governance.Proposal) error {
	return nil
}

//Helper function to get percentage
func getPercentageCoin(c balance.Currency, totalFunding *big.Float, fundTracker *float64, percentage float64) balance.Coin {
	// TODO : How to deal with accuracy
	amount, _ := big.NewFloat(1.0).Mul(totalFunding, big.NewFloat(percentage/100)).Float64()
	*fundTracker = *fundTracker - amount
	return c.NewCoinFromFloat64(amount)
}

//Helper to set Proposal to Finalized
func setToFinalize(ctx *action.Context, proposal *governance.Proposal) status_codes.ProtocolError {
	err := ctx.ProposalMasterStore.Proposal.WithPrefixType(governance.ProposalStateFinalized).Set(proposal)
	if err != nil {
		return governance.ErrStatusUnableToSetFinalized
	}
	ok, err := ctx.ProposalMasterStore.Proposal.WithPrefixType(governance.ProposalStatePassed).Delete(proposal.ProposalID)
	if !ok {
		ok, err = ctx.ProposalMasterStore.Proposal.WithPrefixType(governance.ProposalStateFinalized).Delete(proposal.ProposalID)
		if !ok {
			panic(errors.Wrap(err, "error deleting proposal from finalize prefix(Potentially two copies of same proposal)"))
		}
		return governance.ErrStatusUnableToSetFinalized
	}
	return governance.NoError
}
