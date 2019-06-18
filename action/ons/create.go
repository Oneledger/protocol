package ons

import "github.com/Oneledger/protocol/action"

type DomainCreate struct {
}

func (dc DomainCreate) Signers() []action.Address {
	panic("implement me")
}

func (dc DomainCreate) Type() action.Type {
	panic("implement me")
}

var _ action.Tx = domainCreateTx{}

type domainCreateTx struct {
}

func (domainCreateTx) Validate(ctx *action.Context, msg action.Msg, fee action.Fee, memo string, signatures []action.Signature) (bool, error) {
	panic("implement me")
}

func (domainCreateTx) ProcessCheck(ctx *action.Context, msg action.Msg, fee action.Fee) (bool, action.Response) {
	panic("implement me")
}

func (domainCreateTx) ProcessDeliver(ctx *action.Context, msg action.Msg, fee action.Fee) (bool, action.Response) {
	panic("implement me")
}

func (domainCreateTx) ProcessFee(ctx *action.Context, fee action.Fee) (bool, action.Response) {
	panic("implement me")
}
