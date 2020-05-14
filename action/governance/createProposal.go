package governance

import (
	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/data/governance"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/kv"
)

var _ action.Msg = &CreateProposal{}

type CreateProposal struct {
	proposalType   governance.ProposalType
	description    string
	proposer       keys.Address
	initialFunding int64
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

	options := ctx.Proposals.GetOptionsByType(createProposal.proposalType)

	//Check if initial funding is greater than minimum amount based on type.
	if createProposal.initialFunding < options.InitialFunding {
		return false, errors.New("initial funding does not meet requirements")
	}

	//Check if Proposal Type is valid
	switch createProposal.proposalType {
	case governance.ProposalTypeGeneral:
	case governance.ProposalTypeCodeChange:
	case governance.ProposalTypeConfigUpdate:
	default:
		return false, errors.New("invalid proposal type")
	}

	//Check if proposer address is valid oneLedger address
	err = createProposal.proposer.Err()
	if err != nil {
		return false, errors.Wrap(err, "invalid proposer address")
	}

	if len(createProposal.description) == 0 {
		return false, errors.New("invalid description of proposal")
	}

	return true, nil
}

func (c CreateProposal) ProcessCheck(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	ctx.Logger.Debug("Processing CreateProposal Transaction for CheckTx", tx)
	return runTx(ctx, tx)
}

func (c CreateProposal) ProcessDeliver(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	ctx.Logger.Debug("Processing CreateProposal Transaction for DeliverTx", tx)
	return runTx(ctx, tx)
}

func (c CreateProposal) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas) (bool, action.Response) {
	return action.BasicFeeHandling(ctx, signedTx, start, size, 1)
}

func runTx(ctx *action.Context, signedTx action.RawTx) (bool, action.Response) {
	createProposal := CreateProposal{}
	err := createProposal.Unmarshal(signedTx.Data)
	if err != nil {
		return false, action.Response{}
	}

	options := ctx.Proposals.GetOptionsByType(createProposal.proposalType)

	//Check if initial funding is greater than minimum amount based on type.
	if createProposal.initialFunding < options.InitialFunding {
		return false, action.Response{}
	}

	//Check if Proposal Type is valid
	switch createProposal.proposalType {
	case governance.ProposalTypeGeneral:
	case governance.ProposalTypeCodeChange:
	case governance.ProposalTypeConfigUpdate:
	default:
		return false, action.Response{}
	}

	//Check if proposer address is valid oneLedger address
	err = createProposal.proposer.Err()
	if err != nil {
		return false, action.Response{}
	}

	if len(createProposal.description) == 0 {
		return false, action.Response{}
	}

	return true, action.Response{}
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
