package app

import (
	"fmt"

	"github.com/google/uuid"
	abciTypes "github.com/tendermint/tendermint/abci/types"

	"github.com/Oneledger/protocol/action"
	gov_action "github.com/Oneledger/protocol/action/governance"
	"github.com/Oneledger/protocol/data/governance"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/data/transactions"
)

// Functions for block Beginner
func AddInternalTX(proposalMasterStore *governance.ProposalMasterStore, validator keys.Address, height int64, transaction *transactions.TransactionStore) {
	proposals := proposalMasterStore.Proposal
	activeProposals := proposals.WithPrefixType(governance.ProposalStateActive)
	activeProposals.Iterate(func(id governance.ProposalID, proposal *governance.Proposal) bool {
		//If the proposal is in Voting state and voting period expired, trigger internal tx to handle expiry
		if proposal.Status == governance.ProposalStatusVoting && proposal.VotingDeadline < height {
			tx, err := GetExpireTX(proposal.ProposalID, validator)
			if err != nil {
				return true
			}
			fmt.Println(tx.String()) // Add tx to store (KEY EXPIRE)
		}
		return false
	})

	passedProposals := proposals.WithPrefixType(governance.ProposalStatePassed)
	passedProposals.Iterate(func(id governance.ProposalID, proposal *governance.Proposal) bool {
		if proposal.Status == governance.ProposalStatusCompleted && proposal.Outcome == governance.ProposalOutcomeCompletedYes {
			tx, err := GetFinalizeTX(proposal.ProposalID, validator)
			if err != nil {
				return true
			}
			transaction.AddFinalized(string(proposal.ProposalID), &tx)
			transaction.State.Commit()
		}
		return false
	})

	failedProposals := proposals.WithPrefixType(governance.ProposalStateFailed)
	failedProposals.Iterate(func(id governance.ProposalID, proposal *governance.Proposal) bool {
		if proposal.Status == governance.ProposalStatusCompleted && proposal.Outcome == governance.ProposalOutcomeCompletedNo {
			tx, err := GetFinalizeTX(proposal.ProposalID, validator)
			if err != nil {
				return true
			}
			transaction.AddFinalized(string(proposal.ProposalID), &tx)
		}
		return false
	})
}

func GetFinalizeTX(proposalId governance.ProposalID, validatorAddress keys.Address) (abciTypes.RequestDeliverTx, error) {
	finalizeProposal := &gov_action.FinalizeProposal{
		ProposalID:       proposalId,
		ValidatorAddress: validatorAddress,
	}

	txData, err := finalizeProposal.Marshal()
	if err != nil {
		return RequestDeliverTx{}, err
	}

	internalFinalizeTx := abciTypes.RequestDeliverTx{
		Tx:                   txData,
		XXX_NoUnkeyedLiteral: struct{}{},
		XXX_unrecognized:     nil,
		XXX_sizecache:        0,
	}
	return internalFinalizeTx, nil
}

func GetExpireTX(proposalId governance.ProposalID, validatorAddress keys.Address) (abciTypes.RequestDeliverTx, error) {
	expireVote := &gov_action.ExpireVotes{
		ProposalID:       proposalId,
		ValidatorAddress: validatorAddress,
	}

	txData, err := expireVote.Marshal()
	if err != nil {
		return RequestDeliverTx{}, err
	}

	expireVotes := abciTypes.RequestDeliverTx{
		Tx:                   txData,
		XXX_NoUnkeyedLiteral: struct{}{},
		XXX_unrecognized:     nil,
		XXX_sizecache:        0,
	}
	return expireVotes, nil
}

// Functions for block Ender
func FinalizeProposals(header *Header, ctx *context) bool {
	var finalizeProposals []abciTypes.RequestDeliverTx // Get this from the store
	ctx.transaction.IterateFinalized(func(key string, tx *abciTypes.RequestDeliverTx) bool {
		finalizeProposals = append(finalizeProposals, *tx)
		return false
	})
	for _, proposal := range finalizeProposals {
		actionctx := ctx.Action(header, ctx.deliver)
		txData := proposal.Tx
		newFinalize := gov_action.FinalizeProposal{}
		err := newFinalize.Unmarshal(txData)
		if err != nil {
			fmt.Println("Unable to UnMarshal TX")
			return false
		}
		uuidNew, _ := uuid.NewUUID()
		rawTx := action.RawTx{
			Type: action.PROPOSAL_FINALIZE,
			Data: txData,
			Fee:  action.Fee{},
			Memo: uuidNew.String(),
		}
		ok, _ := newFinalize.ProcessDeliver(actionctx, rawTx)
		if !ok {
			return false
		}
		ctx.deliver.Commit()
	}
	//Delete all proposals
	ctx.transaction.IterateFinalized(func(key string, tx *abciTypes.RequestDeliverTx) bool {
		ctx.transaction.DeleteFinalized(key)
		return false
	})
	return true
}

func ExpireProposals(header *Header, ctx *context) error {
	var finalizePropossals []abciTypes.RequestDeliverTx // Get this from the store
	for _, proposal := range finalizePropossals {
		ctx := ctx.Action(header, ctx.deliver)
		txData := proposal.Tx
		newExpire := gov_action.ExpireVotes{}
		err := newExpire.Unmarshal(txData)
		if err != nil {
			return err
		}
		uuidNew, _ := uuid.NewUUID()
		rawTx := action.RawTx{
			Type: action.EXPIRE_VOTES,
			Data: txData,
			Fee:  action.Fee{},
			Memo: uuidNew.String(),
		}
		newExpire.ProcessDeliver(ctx, rawTx)
	}
	return nil
}
