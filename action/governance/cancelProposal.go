package governance

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/kv"

	"github.com/Oneledger/protocol/action"
	gov "github.com/Oneledger/protocol/data/governance"
	"github.com/Oneledger/protocol/data/keys"
)

var _ action.Msg = &CancelProposal{}

type CancelProposal struct {
	ProposalId gov.ProposalID `json:"proposalId"`
	Proposer   keys.Address   `json:"proposerAddress"`
	Reason     string         `json:"cancelReason"`
}

func (cp CancelProposal) Signers() []action.Address {
	return []action.Address{cp.Proposer.Bytes()}
}

func (cp CancelProposal) Type() action.Type {
	return action.PROPOSAL_CANCEL
}

func (cp CancelProposal) Tags() kv.Pairs {
	tags := make([]kv.Pair, 0)

	tag := kv.Pair{
		Key:   []byte("tx.type"),
		Value: []byte(cp.Type().String()),
	}
	tag2 := kv.Pair{
		Key:   []byte("tx.proposalID"),
		Value: []byte(cp.ProposalId),
	}
	tag3 := kv.Pair{
		Key:   []byte("tx.proposer"),
		Value: cp.Proposer.Bytes(),
	}
	tag4 := kv.Pair{
		Key:   []byte("tx.reason"),
		Value: []byte(cp.Reason),
	}

	tags = append(tags, tag, tag2, tag3, tag4)
	return tags
}

func (cp *CancelProposal) Marshal() ([]byte, error) {
	return json.Marshal(cp)
}

func (cp *CancelProposal) Unmarshal(data []byte) error {
	return json.Unmarshal(data, cp)
}

var _ action.Tx = cancelProposalTx{}

type cancelProposalTx struct {
}

func (c cancelProposalTx) Validate(ctx *action.Context, tx action.SignedTx) (bool, error) {
	ctx.Logger.Debug("Validate CancelProposalTx transaction for CheckTx", tx)

	cc := &CancelProposal{}
	err := cc.Unmarshal(tx.Data)
	if err != nil {
		return false, errors.Wrap(action.ErrWrongTxType, err.Error())
	}

	// validate basic signature
	err = action.ValidateBasic(tx.RawBytes(), cc.Signers(), tx.Signatures)
	if err != nil {
		return false, err
	}

	err = action.ValidateFee(ctx.FeePool.GetOpt(), tx.Fee)
	if err != nil {
		return false, err
	}

	// validate params
	if err = cc.ProposalId.Err(); err != nil {
		return false, gov.ErrInvalidProposalId
	}
	if err = cc.Proposer.Err(); err != nil {
		return false, action.ErrInvalidAddress
	}

	return true, nil
}

func (c cancelProposalTx) ProcessCheck(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	ctx.Logger.Debug("ProcessCheck CancelProposalTx transaction for CheckTx", tx)
	return runCancel(ctx, tx)
}

func (c cancelProposalTx) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas, gasUsed action.Gas) (bool, action.Response) {
	return action.BasicFeeHandling(ctx, signedTx, start, size, 1)
}

func (c cancelProposalTx) ProcessDeliver(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	ctx.Logger.Debug("ProcessDeliver CancelProposalTx transaction for DeliverTx", tx)
	return runCancel(ctx, tx)
}

func runCancel(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	cc := &CancelProposal{}
	err := cc.Unmarshal(tx.Data)
	if err != nil {
		return false, action.Response{
			Log: action.ErrWrongTxType.Wrap(err).Marshal(),
		}
	}

	// Get Proposal from proposal ACTIVE store
	pms := ctx.ProposalMasterStore
	proposal, err := pms.Proposal.WithPrefixType(gov.ProposalStateActive).Get(cc.ProposalId)
	if err != nil {
		return false, action.Response{
			Log: gov.ErrProposalNotExists.Wrap(err).Marshal(),
		}
	}

	// Check if proposal is in FUNDING status
	if proposal.Status != gov.ProposalStatusFunding {
		return false, action.Response{
			Log: gov.ErrStatusNotFunding.Marshal(),
		}
	}

	// Check if proposal funding height is passed
	if ctx.Header.Height > proposal.FundingDeadline {
		return false, action.Response{
			Log: gov.ErrFundingDeadlineCrossed.Marshal(),
		}
	}

	// Check if proposer matches
	if !proposal.Proposer.Equal(cc.Proposer) {
		return false, action.Response{
			Log: gov.ErrUnmatchedProposer.Marshal(),
		}
	}

	// Update fields and add it to FAILED store
	proposal.Status = gov.ProposalStatusCompleted
	proposal.Outcome = gov.ProposalOutcomeCancelled
	proposal.Description += " - Canceled"
	if cc.Reason != "" {
		proposal.Description += fmt.Sprintf(" for reason: %v", cc.Reason)
	}
	err = pms.Proposal.WithPrefixType(gov.ProposalStateFailed).Set(proposal)
	if err != nil {
		return false, action.Response{
			Log: gov.ErrAddingProposalToFailedStore.Wrap(err).Marshal(),
		}
	}

	// Delete proposal in ACTIVE store
	ok, err := pms.Proposal.WithPrefixType(gov.ProposalStateActive).Delete(cc.ProposalId)
	if err != nil || !ok {
		return false, action.Response{Log: gov.ErrDeletingProposalFromActiveStore.Marshal()}
	}

	return true, action.Response{Events: action.GetEvent(cc.Tags(), "cancel_proposal_success")}
}
