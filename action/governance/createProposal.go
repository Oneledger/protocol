package governance

import (
	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/kv"
)

var _ action.Msg = &CreateProposal{}

type CreateProposal struct {
	proposalType    int //TODO: Change types to governance types
	description     string
	proposer        keys.Address
	initialFunding  int
	fundingDeadline int
	fundingGoal     int
	votingDeadline  int
}

func (c CreateProposal) Validate(ctx *action.Context, signedTx action.SignedTx) (bool, error) {
	createProposal := CreateProposal{}
	err := createProposal.Unmarshal(signedTx.Data)
	if err != nil {
		return false, errors.Wrap(action.ErrWrongTxType, err.Error())
	}

	//validate basic signature
	err = action.ValidateBasic(signedTx.RawBytes(), createProposal.Signers(), signedTx.Signatures)
	if err != nil {
		return false, err
	}

	err = action.ValidateFee(ctx.FeePool.GetOpt(), signedTx.Fee)
	if err != nil {
		return false, err
	}

	//validate transaction specific field
	//Check if initial funding is greater than minimum amount based on type.

	return true, nil
}

func (c CreateProposal) ProcessCheck(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	panic("implement me")
}

func (c CreateProposal) ProcessDeliver(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	panic("implement me")
}

func (c CreateProposal) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas) (bool, action.Response) {
	panic("implement me")
}

func (c CreateProposal) Signers() []action.Address {
	panic("implement me")
}

func (c CreateProposal) Type() action.Type {
	panic("implement me")
}

func (c CreateProposal) Tags() kv.Pairs {
	panic("implement me")
}

func (c CreateProposal) Marshal() ([]byte, error) {
	panic("implement me")
}

func (c CreateProposal) Unmarshal(bytes []byte) error {
	panic("implement me")
}
