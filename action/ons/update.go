package ons

import (
	"github.com/Oneledger/protocol/action"
)

var _ action.Msg = DomainUpdate{}

type DomainUpdate struct {
}

func (du DomainUpdate) Signers() []action.Address {
	panic("implement me")
}

func (du DomainUpdate) Type() action.Type {
	panic("implement me")
}

var _ action.Tx = domainUpdateTx{}

type domainUpdateTx struct {
}

func (domainUpdateTx) Validate(ctx *action.Context, msg action.Msg, fee action.Fee, memo string, signatures []action.Signature) (bool, error) {
	panic("implement me")
}

func (domainUpdateTx) ProcessCheck(ctx *action.Context, msg action.Msg, fee action.Fee) (bool, action.Response) {
	panic("implement me")
}

func (domainUpdateTx) ProcessDeliver(ctx *action.Context, msg action.Msg, fee action.Fee) (bool, action.Response) {
	panic("implement me")
}

func (domainUpdateTx) ProcessFee(ctx *action.Context, fee action.Fee) (bool, action.Response) {
	panic("implement me")
}
