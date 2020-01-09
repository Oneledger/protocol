package ons

import (
	"bytes"
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/common"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/data/ons"
)

type DeleteSub struct {
	Name  ons.Name       `json:"owner"`
	Owner action.Address `json:"address"`
}

var _ action.Msg = &DeleteSub{}

func (d DeleteSub) Signers() []action.Address {
	return []action.Address{d.Owner}
}

func (d DeleteSub) Type() action.Type {
	return action.DOMAIN_DELETE_SUB
}

func (d DeleteSub) Tags() common.KVPairs {
	tags := make([]common.KVPair, 0)

	tag := common.KVPair{
		Key:   []byte("tx.type"),
		Value: []byte(d.Type().String()),
	}
	tag2 := common.KVPair{
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

func (d deleteSubTx) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas) (bool, action.Response) {
	return action.BasicFeeHandling(ctx, signedTx, start, size, 1)
}

func runDeleteSub(ctx *action.Context, tx action.RawTx) (bool, action.Response) {

	del := &DeleteSub{}
	err := del.Unmarshal(tx.Data)
	if err != nil {
		return false, action.Response{Log: err.Error()}
	}

	parentName := del.Name
	isSub := del.Name.IsSub()
	if isSub {
		parentName, err = del.Name.GetParentName()
		if err != nil {
			return false, action.Response{Log: err.Error()}
		}
	}

	parent, err := ctx.Domains.Get(parentName)
	//Check if Parent exists in Domain Store
	if err != nil {
		return false, action.Response{Log: "Parent domain doesn't exist, cannot create sub domain!"}
	}
	if !bytes.Equal(parent.OwnerAddress, del.Owner) {
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

	return true, action.Response{Tags: del.Tags()}
}