package staking

import (
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/common"
	cmn "github.com/tendermint/tendermint/libs/common"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/identity"
)

var _ action.Msg = &Purge{}

type Purge struct {
	AdminAddress     action.Address
	ValidatorAddress action.Address
}

func (p Purge) Signers() []action.Address {
	return []action.Address{p.AdminAddress}
}

func (p Purge) Type() action.Type {
	return action.PURGE
}

func (p Purge) Tags() cmn.KVPairs {
	tags := make([]common.KVPair, 0)

	tag := common.KVPair{
		Key:   []byte("tx.type"),
		Value: []byte(p.Type().String()),
	}
	tag2 := common.KVPair{
		Key:   []byte("tx.owner"),
		Value: p.AdminAddress.Bytes(),
	}
	tags = append(tags, tag, tag2)
	return tags
}

func (p Purge) Marshal() ([]byte, error) {
	return json.Marshal(p)
}

func (p *Purge) Unmarshal(bytes []byte) error {
	return json.Unmarshal(bytes, p)
}

var _ action.Tx = purgeTx{}

type purgeTx struct {
}

func (p purgeTx) Validate(ctx *action.Context, signedTx action.SignedTx) (bool, error) {

	purge := &Purge{}
	err := purge.Unmarshal(signedTx.Data)
	if err != nil {
		return false, errors.Wrap(action.ErrWrongTxType, err.Error())
	}
	err = action.ValidateBasic(signedTx.RawBytes(), purge.Signers(), signedTx.Signatures)
	if err != nil {
		return false, err
	}

	err = action.ValidateFee(ctx.FeePool.GetOpt(), signedTx.Fee)
	if err != nil {
		return false, err
	}

	if purge.AdminAddress.Err() != nil || purge.ValidatorAddress.Err() != nil {
		return false, action.ErrInvalidAddress
	}
	return true, nil
}

func (p purgeTx) ProcessCheck(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	return runPurge(ctx, tx)
}

func (p purgeTx) ProcessDeliver(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	return runPurge(ctx, tx)
}

func (p purgeTx) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas) (bool, action.Response) {
	return action.BasicFeeHandling(ctx, signedTx, start, size, 1)
}

func runPurge(ctx *action.Context, tx action.RawTx) (bool, action.Response) {

	purge := &Purge{}
	err := purge.Unmarshal(tx.Data)
	if err != nil {
		return false, action.Response{Log: err.Error()}
	}

	validators := ctx.Validators

	validator, err := validators.Get(purge.ValidatorAddress)
	if err != nil {
		return false, action.Response{Log: errors.Wrap(err, "failed to get the specified validator at address: "+purge.ValidatorAddress.String()).Error()}
	}

	unstake := identity.Unstake{
		Address: validator.Address,
		Amount:  validator.Staking,
	}
	err = validators.HandleUnstake(unstake)
	if err != nil {
		return false, action.Response{Log: errors.Wrap(err, "failed to unstake specified validator at address: "+purge.ValidatorAddress.String()).Error()}
	}
	return true, action.Response{
		Tags: purge.Tags(),
	}
}
