package governance

import (
	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/data/governance"
	"github.com/tendermint/tendermint/libs/kv"
)

var _ action.Msg = &ExpireVotes{}

type ExpireVotes struct {
	ProposalID governance.ProposalID
}

func (e ExpireVotes) Validate(ctx *action.Context, signedTx action.SignedTx) (bool, error) {
	panic("implement me")
}

func (e ExpireVotes) ProcessCheck(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	return runExpireVotes(ctx, tx)
}

func (e ExpireVotes) ProcessDeliver(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	return runExpireVotes(ctx, tx)
}

func (e ExpireVotes) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas) (bool, action.Response) {
	return action.BasicFeeHandling(ctx, signedTx, start, size, 1)
}

func runExpireVotes(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	panic("implement me")
}

func (e ExpireVotes) Signers() []action.Address {
	panic("implement me")
}

func (e ExpireVotes) Type() action.Type {
	panic("implement me")
}

func (e ExpireVotes) Tags() kv.Pairs {
	panic("implement me")
}

func (e ExpireVotes) Marshal() ([]byte, error) {
	panic("implement me")
}

func (e ExpireVotes) Unmarshal(bytes []byte) error {
	panic("implement me")
}
