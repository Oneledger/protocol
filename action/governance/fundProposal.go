package governance

import (
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/kv"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/governance"
	"github.com/Oneledger/protocol/data/keys"
)

var _ action.Msg = &FundProposal{}

type FundProposal struct {
	proposalId    governance.ProposalID
	funderAddress keys.Address
	fundValue     action.Amount
}

func (fp FundProposal) Signers() []action.Address {
	return []action.Address{fp.funderAddress}
}

func (fp FundProposal) Type() action.Type {
	return action.PROPOSAL_FUND
}

func (fp FundProposal) Tags() kv.Pairs {
	tags := make([]kv.Pair, 0)

	tag := kv.Pair{
		Key:   []byte("tx.type"),
		Value: []byte(fp.Type().String()),
	}
	tag2 := kv.Pair{
		Key:   []byte("tx.proposer"),
		Value: fp.funderAddress.Bytes(),
	}
	tag3 := kv.Pair{
		Key:   []byte("tx.proposalID"),
		Value: []byte(string(fp.proposalId)),
	}
	tag4 := kv.Pair{
		Key:   []byte("tx.fundValue"),
		Value: []byte(fp.fundValue.String()),
	}

	tags = append(tags, tag, tag2, tag3, tag4)
	return tags
}

func (fp FundProposal) Marshal() ([]byte, error) {
	return json.Marshal(fp)
}

func (fp *FundProposal) Unmarshal(bytes []byte) error {
	return json.Unmarshal(bytes, fp)
}

type fundProposalTx struct {
}

var _ action.Tx = fundProposalTx{}

func (fundProposalTx) Validate(ctx *action.Context, signedTx action.SignedTx) (bool, error) {
	fundProposal := FundProposal{}
	err := fundProposal.Unmarshal(signedTx.Data)
	if err != nil {
		return false, errors.Wrap(action.ErrWrongTxType, err.Error())
	}

	//Validate basic signature
	err = action.ValidateBasic(signedTx.RawBytes(), fundProposal.Signers(), signedTx.Signatures)
	if err != nil {
		return false, err
	}
	//Validate Fee for funding request
	err = action.ValidateFee(ctx.FeePool.GetOpt(), signedTx.Fee)
	if err != nil {
		return false, err
	}

	// Funding currency should be OLT
	currency, ok := ctx.Currencies.GetCurrencyById(0)
	if !ok {
		panic("no default currency available in the network")
	}
	if currency.Name != fundProposal.fundValue.Currency {
		return false, errors.Wrap(action.ErrInvalidAmount, fundProposal.fundValue.String())
	}

	//Check if Funder address is valid oneLedger address
	err = fundProposal.funderAddress.Err()
	if err != nil {
		return false, errors.Wrap(err, "invalid proposer address")
	}

	return true, nil
}

func (fundProposalTx) ProcessCheck(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	ctx.Logger.Debug("Processing FundProposal Transaction for CheckTx", tx)
	return runFundProposal(ctx, tx)
}

func (fundProposalTx) ProcessDeliver(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	ctx.Logger.Debug("Processing FundProposal Transaction for DeliverTx", tx)
	return runFundProposal(ctx, tx)
}

func (fundProposalTx) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas) (bool, action.Response) {
	return action.BasicFeeHandling(ctx, signedTx, start, size, 1)
}

func runFundProposal(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	fundProposal := FundProposal{}
	err := fundProposal.Unmarshal(tx.Data)
	if err != nil {
		return false, action.Response{}
	}
	// 1 check if proposal exists
	proposal, err := ctx.ProposalMasterStore.Proposal.Get(fundProposal.proposalId)
	if err != nil {
		ctx.Logger.Error("Proposal does not exist :", fundProposal.proposalId)
		result := action.Response{
			Events: action.GetEvent(fundProposal.Tags(), "fund_proposal_does_not_exist"),
		}
		return false, result
	}

	//3. Check if the Funding height is already reached
	//6. If the proposal has already passed Funding height, reject the transaction
	if ctx.Header.Height > proposal.FundingDeadline {
		ctx.Logger.Debug("Funding Height has already been reached")
		result := action.Response{
			Events: action.GetEvent(fundProposal.Tags(), "fund_proposal_height_crossed"),
		}
		return false, result

	}
	//2. Check if the Proposal is in funding stage
	//5. If the proposal is not in FUNDING state, reject the transaction
	if proposal.Status != governance.ProposalStatusFunding {
		ctx.Logger.Debug("Cannot fund proposal , Current proposal state : ", proposal.Status)
		result := action.Response{
			Events: action.GetEvent(fundProposal.Tags(), "fund_proposal_not_funding_state"),
		}
		return false, result
	}
	//3. Check if the Proposal has reached funding goal, when this contribution is added
	//4. Change the state of the proposal to VOTING, if funding goal is met
	fundingAmount := balance.NewAmountFromBigInt(fundProposal.fundValue.Value.BigInt())
	currentFundsforProposal := governance.GetCurrentFunds(proposal.ProposalID, ctx.ProposalMasterStore.ProposalFund)
	newAmount := fundingAmount.Plus(currentFundsforProposal)
	if newAmount.BigInt().Cmp(proposal.FundingGoal.BigInt()) >= 0 {
		proposal.Status = governance.ProposalStatusVoting
		//7. If the proposal moves into Voting state, take a snap shot of Validator Set,
		//at that instinct and add entries for each and every validator at the point into the Proposal_Vote_Store.
		//In the value, we will just update the Voting power. The vote / opinion field remains empty for now
	}

	//8. If the Proposal is still in Funding Stage (Or just moved to Voting) and Funding Goal is not met, add an entry into Proposal Fund Store. No change of State required
	err = ctx.ProposalMasterStore.ProposalFund.AddFunds(proposal.ProposalID, fundProposal.funderAddress, fundingAmount)
	if err != nil {
		ctx.Logger.Error("Faied to add funds to proposal:", fundProposal.proposalId)
		result := action.Response{
			Events: action.GetEvent(fundProposal.Tags(), "fund_proposal_AddFund_failed"),
		}
		return false, result
	}
	result := action.Response{
		Events: action.GetEvent(fundProposal.Tags(), "fund_proposal_success"),
	}
	return true, result
}
