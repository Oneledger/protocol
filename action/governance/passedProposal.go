package governance

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/kv"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/data/keys"
)

type PassedProposal struct {
	Proposer keys.Address
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

	tags = append(tags, tag, tag2)
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
	panic("implement me")
}

func (passedProposalTx) ProcessDeliver(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	panic("implement me")
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
	proposalOptions := ctx.ProposalMasterStore.Proposal.GetOptions()
	fmt.Println(proposalOptions)

}
