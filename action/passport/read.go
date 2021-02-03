package passport

import (
	"encoding/json"
	"github.com/Oneledger/protocol/action/helpers"

	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/kv"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/data/passport"
)

var _ action.Msg = &ReadTestInfo{}

type ReadTestInfo struct {
	Org          passport.TokenTypeID `json:"org"`
	Admin        passport.UserID      `json:"admin"`
	AdminAddress keys.Address         `json:"adminAddress"`
	Person       passport.UserID      `json:"person"`
	Address      keys.Address         `json:"address"`
	Test         passport.TestType    `json:"test"`
	ReadAt       string               `json:"readAt"`
}

func (rd ReadTestInfo) Signers() []action.Address {
	return []action.Address{rd.AdminAddress}
}

func (rd ReadTestInfo) Type() action.Type {
	return action.PASSPORT_READ_TEST
}

func (rd ReadTestInfo) Tags() kv.Pairs {
	tags := make([]kv.Pair, 0)

	tag := kv.Pair{
		Key:   []byte("tx.type"),
		Value: []byte(rd.Type().String()),
	}
	tag1 := kv.Pair{
		Key:   []byte("tx.org"),
		Value: []byte(rd.Org.String()),
	}
	tag2 := kv.Pair{
		Key:   []byte("tx.admin"),
		Value: []byte(rd.Admin.String()),
	}
	tag3 := kv.Pair{
		Key:   []byte("tx.person"),
		Value: []byte(rd.Person.String()),
	}
	tag4 := kv.Pair{
		Key:   []byte("tx.test"),
		Value: []byte(rd.Test.String()),
	}

	tags = append(tags, tag, tag1, tag2, tag3, tag4)
	return tags
}

func (rd ReadTestInfo) Marshal() ([]byte, error) {
	return json.Marshal(rd)
}

func (rd *ReadTestInfo) Unmarshal(bytes []byte) error {
	return json.Unmarshal(bytes, rd)
}

type readTestInfoTx struct {
}

var _ action.Tx = readTestInfoTx{}

func (readTestInfoTx) Validate(ctx *action.Context, signedTx action.SignedTx) (bool, error) {
	read := ReadTestInfo{}
	err := read.Unmarshal(signedTx.Data)
	if err != nil {
		return false, errors.Wrap(action.ErrWrongTxType, err.Error())
	}

	// Validate basic signature
	err = action.ValidateBasic(signedTx.RawBytes(), read.Signers(), signedTx.Signatures)
	if err != nil {
		return false, err
	}

	// Check if org identifier address is valid
	if err = read.Org.Err(); err != nil {
		return false, passport.ErrInvalidIdentifier.Wrap(err)
	}

	// Check if admin is valid
	if err = read.Admin.Err(); err != nil {
		return false, passport.ErrInvalidIdentifier.Wrap(err)
	}

	// Check if admin's address is valid OneLedger address
	if err = read.AdminAddress.Err(); err != nil {
		return false, action.ErrInvalidAddress.Wrap(err)
	}

	// Check if person is valid
	if err = read.Person.Err(); err != nil {
		return false, passport.ErrInvalidIdentifier.Wrap(err)
	}

	// If provided, check if person's address is valid OneLedger address
	if len(read.Address) != 0 && read.Address.Err() != nil {
		return false, action.ErrInvalidAddress.Wrap(err)
	}

	// Check if test type is valid
	if err = read.Test.Err(); err != nil {
		return false, passport.ErrInvalidTestType.Wrap(err)
	}

	// check time stamp
	if !validateRFCTime(read.ReadAt) {
		return false, passport.ErrInvalidTimeStamp
	}

	return true, nil
}

func (readTestInfoTx) ProcessCheck(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	ctx.Logger.Debug("Processing readTestInfoTx Transaction for CheckTx", tx)
	return runRead(ctx, tx)
}

func (readTestInfoTx) ProcessDeliver(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	ctx.Logger.Debug("Processing readTestInfoTx Transaction for DeliverTx", tx)
	return runRead(ctx, tx)
}

func (readTestInfoTx) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas) (bool, action.Response) {
	return true, action.Response{}
}

func runRead(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	read := &ReadTestInfo{}
	err := read.Unmarshal(tx.Data)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, ErrSerialization, read.Tags(), err)
	}

	// check permission
	permitted, err := ctx.AuthTokens.HasPermission(read.Org, read.Admin, passport.PermitScan)
	if !permitted {
		return helpers.LogAndReturnFalse(ctx.Logger, passport.ErrPermissionRequired, read.Tags(), err)
	}

	// Log this read Tx into store
	log := passport.NewReadLog(read.Org, read.Admin, read.AdminAddress, read.Person, read.Address, read.Test, read.ReadAt)
	err = ctx.Tests.LogRead(log)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, passport.ErrReadTestFailure, read.Tags(), err)
	}
	return helpers.LogAndReturnTrue(ctx.Logger, read.Tags(), "test_info_read")
}
