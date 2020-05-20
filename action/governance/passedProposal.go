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

type PassedProposal struct {
	ProposalID governance.ProposalID
	Proposer   keys.Address
}

func (p PassedProposal) Signers() []action.Address {
	return []action.Address{p.Proposer}
}

func (p PassedProposal) Type() action.Type {
	return action.PROPOSAL_PASSED
}

func (p PassedProposal) Tags() kv.Pairs {
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

func (p PassedProposal) Marshal() ([]byte, error) {
	return json.Marshal(p)
}

func (p *PassedProposal) Unmarshal(bytes []byte) error {
	return json.Unmarshal(bytes, p)
}

type passedProposalTx struct {
}

var _ action.Tx = &passedProposalTx{}

func (passedProposalTx) Validate(ctx *action.Context, signedTx action.SignedTx) (bool, error) {
	passedProposal := PassedProposal{}
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

func (passedProposalTx) ProcessCheck(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	return runPassedProposal(ctx, tx)
}

func (passedProposalTx) ProcessDeliver(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	return runPassedProposal(ctx, tx)
}

func (passedProposalTx) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas) (bool, action.Response) {
	return action.BasicFeeHandling(ctx, signedTx, start, size, 1)
}

func runPassedProposal(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	passedProposal := PassedProposal{}
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

	ctx.ProposalMasterStore.ProposalVote.IsPassed(proposal.ProposalID, proposal.PassPercentage)
	distributeFunds(ctx, proposal)
	result := action.Response{
		Events: action.GetEvent(passedProposal.Tags(), "finalize_proposal_success"),
	}
	return true, result
}

func distributeFunds(ctx *action.Context, proposal *governance.Proposal) error {
	// Required Perimeters for Fund Distribution
	proposalDistribution := ctx.ProposalMasterStore.Proposal.GetOptionsByType(proposal.Type).PassedFundDistribution
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
	ctx.FeePool.AddToPool(getPercentageCoin(c, totalFunding, proposalDistribution.FeePool))

	//Reward for Proposer
	ctx.Balances.AddToAddress(proposal.Proposer, getPercentageCoin(c, totalFunding, proposalDistribution.ProposerReward))

	//Bounty Program
	bountyAddress := action.Address(ctx.ProposalMasterStore.Proposal.GetOptions().BountyProgramAddr)
	ctx.Balances.AddToAddress(bountyAddress, getPercentageCoin(c, totalFunding, proposalDistribution.BountyPool))

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
