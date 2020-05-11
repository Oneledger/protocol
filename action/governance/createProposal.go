package governance

import (
	"github.com/Oneledger/protocol/action"
	"github.com/tendermint/tendermint/libs/kv"
)

var _ action.Msg = &CreateProposal{}

type CreateProposal struct {
}

func (c CreateProposal) Validate(ctx *action.Context, signedTx action.SignedTx) (bool, error) {
	panic("implement me")
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
