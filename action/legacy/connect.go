package legacy

import (
	"encoding/json"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/action/helpers"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/kv"
)

type LegacyConnect struct {
	LegacyAddress action.Address `json:"legacyAddress"`
	NewAddress    action.Address `json:"newAddress"`
	Nonce         []byte         `json:"nonce"`
	Toggle        bool           `json:"toggle"`
}

func (legacy LegacyConnect) Marshal() ([]byte, error) {
	return json.Marshal(legacy)
}

func (legacy *LegacyConnect) Unmarshal(data []byte) error {
	return json.Unmarshal(data, legacy)
}

func (legacy LegacyConnect) Signers() []action.Address {
	if legacy.Toggle == true {
		return []action.Address{
			legacy.LegacyAddress.Bytes(),
			legacy.NewAddress.Bytes(),
		}
	}
	return []action.Address{legacy.LegacyAddress.Bytes()}
}

func (legacy LegacyConnect) Type() action.Type {
	return action.LEGACY_CONNECT
}

func (legacy LegacyConnect) Tags() kv.Pairs {
	tags := make([]kv.Pair, 0)

	tag := kv.Pair{
		Key:   []byte("tx.type"),
		Value: []byte(legacy.Type().String()),
	}
	tag2 := kv.Pair{
		Key:   []byte("tx.legacyAddress"),
		Value: legacy.LegacyAddress.Bytes(),
	}
	tag3 := kv.Pair{
		Key:   []byte("tx.newAddress"),
		Value: legacy.NewAddress.Bytes(),
	}
	tag4 := action.BoolTag("tx.toggle", legacy.Toggle)
	tags = append(tags, tag, tag2, tag3, tag4)
	return tags
}

var _ action.Tx = legacyConnectTx{}

type legacyConnectTx struct{}

func (ltx legacyConnectTx) Validate(ctx *action.Context, tx action.SignedTx) (bool, error) {
	legacy := &LegacyConnect{}
	err := legacy.Unmarshal(tx.Data)
	if err != nil {
		return false, errors.Wrap(action.ErrWrongTxType, err.Error())
	}

	if len(legacy.Nonce) != 32 {
		return false, action.ErrInvalidNonce
	}

	//validate basic signature
	err = action.ValidateBasic(legacy.Nonce, legacy.Signers(), tx.Signatures)
	if err != nil {
		return false, err
	}

	if tx.Signatures[0].Signer.KeyType != keys.ED25519 || (legacy.Toggle == true && tx.Signatures[1].Signer.KeyType != keys.ETHSECP) {
		return false, action.ErrInvalidSignature
	}

	err = action.ValidateFee(ctx.FeePool.GetOpt(), tx.Fee)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (legacyConnectTx) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas, gasUsed action.Gas) (bool, action.Response) {
	return action.BasicFeeHandling(ctx, signedTx, start, size, 1)
}

func (ltx legacyConnectTx) ProcessCheck(ctx *action.Context, tx action.RawTx) (ok bool, result action.Response) {
	ctx.Logger.Detail("Processing Legacy Connect Transaction for CheckTx", tx)
	ok, result = runConnect(ctx, tx)
	return
}

func (ltx legacyConnectTx) ProcessDeliver(ctx *action.Context, tx action.RawTx) (ok bool, result action.Response) {
	ctx.Logger.Detail("Processing Legacy Connect Transaction for DeliverTx", tx)
	ok, result = runConnect(ctx, tx)
	return
}

func runConnect(ctx *action.Context, tx action.RawTx) (ok bool, result action.Response) {
	legacy := &LegacyConnect{}
	err := legacy.Unmarshal(tx.Data)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, action.ErrUnserializable, legacy.Tags(), err)
	}

	// TODO: Add duplication check
	mapper := ctx.StateDB.GetAccountMapper()
	aj, err := mapper.Get(legacy.LegacyAddress, keys.ED25519)
	if err != nil {
		aj, err = mapper.GetOrCreateED25519ToETHECDSA(legacy.LegacyAddress, legacy.NewAddress)
		if err != nil {
			return helpers.LogAndReturnFalse(ctx.Logger, action.ErrUnserializable, legacy.Tags(), errors.Wrap(err, "failed to create mapping"))
		}
		return helpers.LogAndReturnTrue(ctx.Logger, legacy.Tags(), "legacy_connect")
	}
	if legacy.Toggle != aj.Enabled {
		aj.Enabled = legacy.Toggle
		err = mapper.Set(aj)
		if err != nil {
			return helpers.LogAndReturnFalse(ctx.Logger, action.ErrUnserializable, legacy.Tags(), errors.Wrap(err, "failed to update mapping"))
		}
		return helpers.LogAndReturnTrue(ctx.Logger, legacy.Tags(), "legacy_connect")
	}
	return helpers.LogAndReturnFalse(ctx.Logger, action.ErrUnserializable, legacy.Tags(), errors.Wrap(err, "no changes"))
}
