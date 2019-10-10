package ons

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"

	"github.com/Oneledger/protocol/action"
	"github.com/tendermint/tendermint/libs/common"
)

var _ Ons = &DomainUpdate{}

type DomainUpdate struct {
	Owner   action.Address `json:"owner"`
	Account action.Address `json:"account"`
	Name    string         `json:"name"`
	Active  bool           `json:"active"`
}

func (du DomainUpdate) Marshal() ([]byte, error) {
	return json.Marshal(du)
}

func (du *DomainUpdate) Unmarshal(data []byte) error {
	return json.Unmarshal(data, du)
}

func (du DomainUpdate) OnsName() string {
	return du.Name
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

func (domainUpdateTx) Validate(ctx *action.Context, tx action.SignedTx) (bool, error) {
	update := &DomainUpdate{}
	err := update.Unmarshal(tx.Data)
	if err != nil {
		return false, errors.Wrap(action.ErrWrongTxType, err.Error())
	}

	err = action.ValidateBasic(tx.RawBytes(), update.Signers(), tx.Signatures)
	if err != nil {
		return false, err
	}

	err = action.ValidateFee(ctx.FeeOpt, tx.Fee)
	if err != nil {
		return false, err
	}

	if update.Owner == nil || len(update.Name) <= 0 {
		return false, action.ErrMissingData
	}

	if update.Active == false && update.Account == nil {
		return false, action.ErrMissingData
	}

	return true, nil
}

func (domainUpdateTx) ProcessCheck(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	update := &DomainUpdate{}
	err := update.Unmarshal(tx.Data)
	if err != nil {
		return false, action.Response{Log: err.Error()}
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

func (domainUpdateTx) ProcessDeliver(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	update := &DomainUpdate{}
	err := update.Unmarshal(tx.Data)
	if err != nil {
		return false, action.Response{Log: err.Error()}
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

	err = ctx.Domains.Set(d)
	if err != nil {
		return false, action.Response{Log: err.Error()}
	}
	return true, action.Response{Tags: update.Tags()}
}

func (domainUpdateTx) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas) (bool, action.Response) {
	panic("implement me")
}
