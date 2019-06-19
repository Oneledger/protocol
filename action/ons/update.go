package ons

import (
	"bytes"
	"encoding/hex"
	"fmt"

	"github.com/Oneledger/protocol/action"
	"github.com/tendermint/tendermint/libs/common"
)

var _ action.Msg = DomainUpdate{}

type DomainUpdate struct {
	Owner   action.Address
	Account action.Address
	Name    string
	Active  bool
}

func (du DomainUpdate) Signers() []action.Address {
	return []action.Address{du.Owner}
}

func (du DomainUpdate) Type() action.Type {
	return action.DOMAIN_UPDATE
}

func (du DomainUpdate) Tags() common.KVPairs {
	tags := make([]common.KVPair, 0)

	tag := common.KVPair{
		Key:   []byte("tx.type"),
		Value: []byte(du.Type().String()),
	}
	tag2 := common.KVPair{
		Key:   []byte("tx.from"),
		Value: du.Owner.Bytes(),
	}

	tags = append(tags, tag, tag2)
	return tags
}

var _ action.Tx = domainUpdateTx{}

type domainUpdateTx struct {
}

func (domainUpdateTx) Validate(ctx *action.Context, msg action.Msg, fee action.Fee, memo string, signatures []action.Signature) (bool, error) {
	ok, err := action.ValidateBasic(msg, fee, memo, signatures)
	if err != nil {
		return ok, err
	}

	update, ok := msg.(*DomainUpdate)
	if !ok {
		return false, action.ErrWrongTxType
	}

	if update.Owner == nil || len(update.Name) <= 0 {
		return false, action.ErrMissingData
	}

	if update.Active == false && update.Account == nil {
		return false, action.ErrMissingData
	}

	return true, nil
}

func (domainUpdateTx) ProcessCheck(ctx *action.Context, msg action.Msg, fee action.Fee) (bool, action.Response) {
	update, ok := msg.(*DomainUpdate)
	if !ok {
		return false, action.Response{Log: action.ErrWrongTxType.Error()}
	}

	if !ctx.Domains.Exists(update.Name) {
		return false, action.Response{Log: fmt.Sprintf("domain doesn't exist: %s", update.Name)}
	}

	d, err := ctx.Domains.Get(update.Name)
	if err != nil {
		return false, action.Response{Log: fmt.Sprintf("failed to get domain: %s", update.Name)}
	}

	if !bytes.Equal(d.OwnerAddress, update.Owner) {
		return false, action.Response{Log: fmt.Sprintf("domain is not owned by: %s", hex.EncodeToString(update.Owner))}
	}

	return true, action.Response{Tags: update.Tags()}
}

func (domainUpdateTx) ProcessDeliver(ctx *action.Context, msg action.Msg, fee action.Fee) (bool, action.Response) {
	update, ok := msg.(*DomainUpdate)
	if !ok {
		return false, action.Response{Log: action.ErrWrongTxType.Error()}
	}

	if !ctx.Domains.Exists(update.Name) {
		return false, action.Response{Log: fmt.Sprintf("domain doesn't exist: %s", update.Name)}
	}

	d, err := ctx.Domains.Get(update.Name)
	if err != nil {
		return false, action.Response{Log: fmt.Sprintf("failed to get domain: %s", update.Name)}
	}

	if !d.IsChangeable(ctx.Header.Height) {
		return false, action.Response{Log: fmt.Sprintf("domain is not changable: %s, last change: %d", update.Name, d.LastUpdateHeight)}
	}

	if !bytes.Equal(d.OwnerAddress, update.Owner) {
		return false, action.Response{Log: fmt.Sprintf("domain is not owned by: %s", hex.EncodeToString(update.Owner))}
	}

	d.SetAccountAddress(update.Account)
	if update.Active {
		d.Activate()
	} else {
		d.Deactivate()
	}
	d.SetLastUpdatedHeight(ctx.Header.Height)
	return true, action.Response{Tags: update.Tags()}
}

func (domainUpdateTx) ProcessFee(ctx *action.Context, fee action.Fee) (bool, action.Response) {
	panic("implement me")
}
