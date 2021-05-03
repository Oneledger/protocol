package ons

import (
	"bytes"
	"encoding/json"

	"github.com/tendermint/tendermint/libs/kv"

	"github.com/pkg/errors"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/data/ons"
)

/*
		DeleteSub

This transaction deletes the specific sub domain if sub domain name is passed. It deletes all the subdomains if a parent
domain name is passed.

*/
type DeleteSub struct {
	Name  ons.Name       `json:"name"`
	Owner action.Address `json:"owner"`
}

var _ action.Msg = &DeleteSub{}

func (d DeleteSub) Signers() []action.Address {
	return []action.Address{d.Owner}
}

func (d DeleteSub) Type() action.Type {
	return action.DOMAIN_DELETE_SUB
}

func (d DeleteSub) Tags() kv.Pairs {
	tags := make([]kv.Pair, 0)

	tag := kv.Pair{
		Key:   []byte("tx.type"),
		Value: []byte(d.Type().String()),
	}
	tag2 := kv.Pair{
		Key:   []byte("tx.owner"),
		Value: d.Owner.Bytes(),
	}

	tags = append(tags, tag, tag2)
	return tags
}

func (d DeleteSub) Marshal() ([]byte, error) {
	return json.Marshal(d)
}

func (d *DeleteSub) Unmarshal(data []byte) error {
	return json.Unmarshal(data, d)
}

var _ action.Tx = &deleteSubTx{}

type deleteSubTx struct {
}

func (d deleteSubTx) Validate(ctx *action.Context, signedTx action.SignedTx) (bool, error) {
	del := &DeleteSub{}
	err := del.Unmarshal(signedTx.Data)
	if err != nil {
		return false, errors.Wrap(action.ErrWrongTxType, err.Error())
	}

	err = action.ValidateBasic(signedTx.RawBytes(), del.Signers(), signedTx.Signatures)
	if err != nil {
		return false, err
	}

	err = action.ValidateFee(ctx.FeePool.GetOpt(), signedTx.Fee)
	if err != nil {
		return false, err
	}

	if del.Owner == nil || len(del.Name) <= 0 {
		return false, action.ErrMissingData
	}

	// checking if name is valid
	if !del.Name.IsValid() {
		return false, ErrInvalidDomain
	}

	return true, nil
}

func (d deleteSubTx) ProcessCheck(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	return runDeleteSub(ctx, tx)
}

func (d deleteSubTx) ProcessDeliver(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	return runDeleteSub(ctx, tx)
}

func (d deleteSubTx) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas, gasUsed action.Gas) (bool, action.Response) {
	return action.BasicFeeHandling(ctx, signedTx, start, size, 1)
}

func runDeleteSub(ctx *action.Context, tx action.RawTx) (bool, action.Response) {

	del := &DeleteSub{}
	err := del.Unmarshal(tx.Data)
	if err != nil {
		return false, action.Response{Log: err.Error()}
	}

	// if the delete name is a parent domain, then it is the parent name
	parentName := del.Name

	isSub := del.Name.IsSub()
	if isSub {
		// if the delete name is a sub domain, the parentName

		parentName, err = del.Name.GetParentName()
		if err != nil {
			return false, action.Response{Log: err.Error()}
		}
	}

	parent, err := ctx.Domains.Get(parentName)
	//Check if Parent exists in Domain Store
	if err != nil {
		return false, action.Response{Log: "Parent domain doesn't exist, cannot delete sub domain!"}
	}

	if !parent.IsChangeable(ctx.Header.Height) {
		return false, action.Response{Log: "domain is not changeable"}
	}

	if !bytes.Equal(parent.Owner, del.Owner) {
		return false, action.Response{Log: "parent domain not owned"}
	}

	if isSub {

		err = ctx.Domains.DeleteASubdomain(del.Name)
		if err != nil {
			return false, action.Response{Log: err.Error()}
		}

	} else {

		err = ctx.Domains.DeleteAllSubdomains(del.Name)
		if err != nil {
			return false, action.Response{Log: err.Error()}
		}
	}

	parent.SetLastUpdatedHeight(ctx.Header.Height)

	return true, action.Response{Events: action.GetEvent(del.Tags(), "delete_subDomain")}
}
