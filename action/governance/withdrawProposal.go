package governance

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/governance"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/kv"
)

var _ action.Msg = &WithdrawProposal{}

type WithdrawProposal struct {
	ProposalID     governance.ProposalID
	Description    string
	Contributor    keys.Address
	WithDrawAmount action.Amount
	Beneficiary    keys.Address

}

func (c WithdrawProposal) Validate(ctx *action.Context, signedTx action.SignedTx) (bool, error) {
	withdrawProposal := WithdrawProposal{}
	err := withdrawProposal.Unmarshal(signedTx.Data)
	if err != nil {
		return false, errors.Wrap(action.ErrWrongTxType, err.Error())
	}

	//validate basic signature
	err = action.ValidateBasic(signedTx.RawBytes(), withdrawProposal.Signers(), signedTx.Signatures)
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
	if currency.Name != withdrawProposal.WithDrawAmount.Currency {
		return false, errors.Wrap(action.ErrInvalidAmount, withdrawProposal.WithDrawAmount.String())
	}

	//Check if fund contributor address is valid oneLedger address
	err = withdrawProposal.Contributor.Err()
	if err != nil {
		return false, errors.Wrap(err, "invalid withdraw contributor address")
	}

	//Check if withdraw beneficiary address is valid oneLedger address
	err = withdrawProposal.Beneficiary.Err()
	if err != nil {
		return false, errors.Wrap(err, "invalid withdraw beneficiary address")
	}

	return true, nil
}

func (c WithdrawProposal) ProcessCheck(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	ctx.Logger.Debug("Processing CreateProposal Transaction for CheckTx", tx)
	return runTx(ctx, tx)
}

func (c WithdrawProposal) ProcessDeliver(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	ctx.Logger.Debug("Processing CreateProposal Transaction for DeliverTx", tx)
	return runTx(ctx, tx)
}

func (c WithdrawProposal) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas) (bool, action.Response) {
	return action.BasicFeeHandling(ctx, signedTx, start, size, 1)
}

func runWithdraw(ctx *action.Context, signedTx action.RawTx) (bool, action.Response) {
	withDrawProposal := WithdrawProposal{}
	err := withDrawProposal.Unmarshal(signedTx.Data)
	if err != nil {
		return false, action.Response{}
	}

	// 1. Check if Proposal already exists,
	// and a. the funding height is reached
	// b. the funding goal is not reached
	if !ctx.ProposalMasterStore.Proposal.Exists(withDrawProposal.ProposalID) {
		result := action.Response{
			Events: action.GetEvent(withDrawProposal.Tags(), "withdraw_proposal_failed"),
		}
		return false, result
	} else {
		proposal, state, err := ctx.ProposalMasterStore.Proposal.QueryAllStores(withDrawProposal.ProposalID)

		if err != nil {
			svc.logger.Error("error getting proposal", err)
			return codes.ErrGetProposal
		}
	}

	// 2. Check if the this contributor has funded this specific proposal, if so, get the amount of funds

	proposalFund := ctx.ProposalMasterStore.ProposalFund.GetFundersForProposalID(id, func(proposalID governance.ProposalID, fundingAddr keys.Address, amt *balance.Amount) governance.ProposalFund {
		return governance.ProposalFund{
			id:            proposalID,
			address:       fundingAddr,
			fundingAmount: amt,
		}
	})



	if !bytes.Equal(proposal., update.Owner) {
		return false, action.Response{Log: fmt.Sprintf("domain is not owned by: %s", hex.EncodeToString(update.Owner))}
	}


	//TODO Check the proposal, and make sure that the proposal is in CANCELLED state. Else reject the transaction

	//TODO Check if the contributor has sufficient funds to withdraw for that proposal. Check for the corresponding value in Proposal Fund Store




	return true, action.Response{}
}

func (c WithdrawProposal) Signers() []action.Address {
	return []action.Address{c.Contributor}
}

func (c WithdrawProposal) Type() action.Type {
	panic("implement me")
}

func (c WithdrawProposal) Tags() kv.Pairs {
	panic("implement me")
}

func (c WithdrawProposal) Marshal() ([]byte, error) {
	panic("implement me")
}

func (c WithdrawProposal) Unmarshal(bytes []byte) error {
	panic("implement me")
}
