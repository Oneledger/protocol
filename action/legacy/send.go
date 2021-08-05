package legacy

import "github.com/Oneledger/protocol/action"

type LegacySend struct{}

var _ action.Tx = legacySendTx{}

type legacySendTx struct{}

func (ltx legacySendTx) Validate(ctx *action.Context, tx action.SignedTx) (bool, error) {
	return true, nil
}

func (legacySendTx) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas, gasUsed action.Gas) (bool, action.Response) {
	return false, action.Response{}
}

func (ltx legacySendTx) ProcessCheck(ctx *action.Context, tx action.RawTx) (ok bool, result action.Response) {
	return
}

func (ltx legacySendTx) ProcessDeliver(ctx *action.Context, tx action.RawTx) (ok bool, result action.Response) {
	return
}
