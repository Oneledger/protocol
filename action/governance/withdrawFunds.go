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

var _ action.Msg = &WithdrawFunds{}

type WithdrawFunds struct {
	ProposalID    governance.ProposalID `json:"proposalId"`
	Funder        keys.Address          `json:"funderAddress"`
	WithdrawValue action.Amount         `json:"withdrawValue"`
	Beneficiary   keys.Address          `json:"beneficiaryAddress"`
}

func (wp WithdrawFunds) Validate(ctx *action.Context, signedTx action.SignedTx) (bool, error) {
	withdrawFunds := WithdrawFunds{}
	err := withdrawFunds.Unmarshal(signedTx.Data)
	if err != nil {
		return false, errors.Wrap(action.ErrWrongTxType, err.Error())
	}

	//validate basic signature
	err = action.ValidateBasic(signedTx.RawBytes(), withdrawFunds.Signers(), signedTx.Signatures)
	if err != nil {
		return false, err
	}
	//validate fee

	err = action.ValidateFee(ctx.FeePool.GetOpt(), signedTx.Fee)

	if err != nil {
		return false, err
	}

	// the currency should be OLT
	currency, ok := ctx.Currencies.GetCurrencyById(0)
	if !ok {
		panic("no default currency available in the network")
	}
	if currency.Name != withdrawFunds.WithdrawValue.Currency {
		return false, errors.Wrap(action.ErrInvalidAmount, withdrawFunds.WithdrawValue.String())
	}

	//Check if fund funder address is valid oneLedger address
	err = withdrawFunds.Funder.Err()
	if err != nil {
		return false, errors.Wrap(governance.ErrInvalidFunderAddr, err.Error())
	}

	//Check if Proposal ID is valid
	if err = withdrawFunds.ProposalID.Err(); err != nil {
		return false, governance.ErrInvalidProposalId
	}

	//Check if withdraw beneficiary address is valid oneLedger address
	err = withdrawFunds.Beneficiary.Err()
	if err != nil {
		return false, errors.Wrap(governance.ErrInvalidBeneficiaryAddr, err.Error())
	}

	return true, nil
}

func (wp WithdrawFunds) ProcessCheck(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	ctx.Logger.Debug("Processing WithdrawFunds Transaction for CheckTx", tx)
	return runWithdraw(ctx, tx)
}

func (wp WithdrawFunds) ProcessDeliver(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	ctx.Logger.Debug("Processing WithdrawFunds Transaction for DeliverTx", tx)
	return runWithdraw(ctx, tx)
}

func (wp WithdrawFunds) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas, gasUsed action.Gas) (bool, action.Response) {
	return action.BasicFeeHandling(ctx, signedTx, start, size, 1)
}

func runWithdraw(ctx *action.Context, signedTx action.RawTx) (bool, action.Response) {
	withdrawProposal := WithdrawFunds{}
	err := withdrawProposal.Unmarshal(signedTx.Data)
	if err != nil {
		return false, action.Response{
			Log: action.ErrWrongTxType.Wrap(err).Marshal(),
		}
	}

	// 1. Check if Proposal already exists, if so, check the withdraw requirement:
	//    a. if the proposal outcome is cancelled or insufficient funds
	//    or
	//    b. the funding goal is not reached and the funding height is reached
	proposal, _, err := ctx.ProposalMasterStore.Proposal.QueryAllStores(withdrawProposal.ProposalID)
	if err != nil {
		ctx.Logger.Error("Proposal does not exist :", withdrawProposal.ProposalID)
		result := action.Response{
			Events: action.GetEvent(withdrawProposal.Tags(), "withdraw_proposal_does_not_exist"),
			Log:    governance.ErrProposalNotExists.Wrap(err).Marshal(),
		}
		return false, result
	}

	fundStore := ctx.ProposalMasterStore.ProposalFund
	currentFundsForProposal := fundStore.GetCurrentFundsForProposal(proposal.ProposalID)
	// if outcome is not cancelled or insufficient funds
	if proposal.Outcome != governance.ProposalOutcomeCancelled && proposal.Outcome != governance.ProposalOutcomeInsufficientFunds {
		// if funding goal is reached or there is still time for funding,
		// this will also cover InsufficientVotes and Completed
		if currentFundsForProposal.BigInt().Cmp(proposal.FundingGoal.BigInt()) >= 0 || ctx.Header.Height <= proposal.FundingDeadline {
			ctx.Logger.Error("Proposal does not meet withdraw requirement", withdrawProposal.ProposalID)
			result := action.Response{
				Events: action.GetEvent(withdrawProposal.Tags(), "withdraw_proposal_does_not_meet_withdraw_requirement"),
				Log:    governance.ErrProposalWithdrawNotEligible.Marshal(),
			}
			return false, result
		}
		// 2. change outcome, status, state
		// when it reaches here means proposal outcome is only to be InProgress
		proposal.Outcome = governance.ProposalOutcomeInsufficientFunds
		proposal.Status = governance.ProposalStatusCompleted

		err = ctx.ProposalMasterStore.Proposal.WithPrefixType(governance.ProposalStateFailed).Set(proposal)
		if err != nil {
			ctx.Logger.Error("Failed to add proposal to FAILED store :", proposal.ProposalID)
			result := action.Response{
				Events: action.GetEvent(withdrawProposal.Tags(), "failed_to_add_proposal_to_failed_store"),
				Log:    governance.ErrAddingProposalToFailedStore.Wrap(err).Marshal(),
			}
			return false, result
		}
		ok, err := ctx.ProposalMasterStore.Proposal.WithPrefixType(governance.ProposalStateActive).Delete(proposal.ProposalID)
		if err != nil || !ok {
			ctx.Logger.Error("Failed to delete proposal from ACTIVE store :", proposal.ProposalID)
			result := action.Response{
				Events: action.GetEvent(withdrawProposal.Tags(), "failed_to_delete_proposal_from_active_store"),
				Log:    governance.ErrDeletingProposalFromActiveStore.Wrap(err).Marshal(),
			}
			return false, result
		}
	}

	// 3. Check if the funder has funded this proposal
	isFundedByFunder := fundStore.IsFundedByFunder(proposal.ProposalID, withdrawProposal.Funder)
	ctx.Logger.Detail("isFundedByFunder: ", isFundedByFunder)
	if !isFundedByFunder {
		ctx.Logger.Error("No such funder funded this proposal :", withdrawProposal.Funder)
		result := action.Response{
			Events: action.GetEvent(withdrawProposal.Tags(), "no_such_funder_funded_this_proposal"),
			Log:    governance.ErrNoSuchFunder.Marshal(),
		}
		return false, result
	}

	// 4. withdraw
	// deduct from proposal fund and check if the funder has sufficient funds to withdraw for that proposal
	withdrawAmount := balance.NewAmountFromBigInt(withdrawProposal.WithdrawValue.Value.BigInt())
	err = ctx.ProposalMasterStore.ProposalFund.DeductFunds(proposal.ProposalID, withdrawProposal.Funder, withdrawAmount)
	if err != nil {
		ctx.Logger.Error("Failed to deduct funds from proposal:", withdrawProposal.ProposalID)
		result := action.Response{
			Events: action.GetEvent(withdrawProposal.Tags(), "withdraw_proposal_deduct_fund_failed"),
			Log:    governance.ErrDeductFunding.Wrap(err).Marshal(),
		}
		return false, result
	}
	// add to beneficiary address
	coin := withdrawProposal.WithdrawValue.ToCoin(ctx.Currencies)
	err = ctx.Balances.AddToAddress(withdrawProposal.Beneficiary.Bytes(), coin)
	if err != nil {
		result := action.Response{
			Events: action.GetEvent(withdrawProposal.Tags(), "withdraw_proposal_addition_failed"),
			Log:    governance.ErrAddFunding.Marshal(),
		}
		return false, result
	}

	result := action.Response{
		Events: action.GetEvent(withdrawProposal.Tags(), "withdraw_proposal_success"),
	}

	return true, result
}

func (wp WithdrawFunds) Signers() []action.Address {
	return []action.Address{wp.Funder}
}

func (wp WithdrawFunds) Type() action.Type {
	return action.PROPOSAL_WITHDRAW_FUNDS
}

func (wp WithdrawFunds) Tags() kv.Pairs {
	tags := make([]kv.Pair, 0)

	tag := kv.Pair{
		Key:   []byte("tx.type"),
		Value: []byte(wp.Type().String()),
	}
	tag2 := kv.Pair{
		Key:   []byte("tx.proposalID"),
		Value: []byte(string(wp.ProposalID)),
	}
	tag3 := kv.Pair{
		Key:   []byte("tx.funder"),
		Value: wp.Funder.Bytes(),
	}
	tag4 := kv.Pair{
		Key:   []byte("tx.withdrawValue"),
		Value: []byte(wp.WithdrawValue.String()),
	}
	tag5 := kv.Pair{
		Key:   []byte("tx.beneficiary"),
		Value: wp.Beneficiary.Bytes(),
	}

	tags = append(tags, tag, tag2, tag3, tag4, tag5)
	return tags
}

func (wp WithdrawFunds) Marshal() ([]byte, error) {
	return json.Marshal(wp)
}

func (wp *WithdrawFunds) Unmarshal(bytes []byte) error {
	return json.Unmarshal(bytes, wp)
}
