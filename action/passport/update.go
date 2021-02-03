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

var _ action.Msg = &UpdateTestInfo{}

type UpdateTestInfo struct {
	TestID  string               `json:"testId"`
	Person  passport.UserID      `json:"person"`
	Test    passport.TestType    `json:"test"`
	Result  passport.TestResult  `json:"testResult"`

	AnalysisOrg passport.TokenTypeID `json:"analysisOrg"`
	AnalyzedAt  string               `json:"analyzedAt"`
	AnalyzedBy  string               `json:"analyzedBy"`

	Admin        passport.UserID `json:"admin"`
	AdminAddress keys.Address    `json:"adminAddress"`
	Notes        string          `json:"notes"`
}

func (up UpdateTestInfo) Signers() []action.Address {
	return []action.Address{up.AdminAddress}
}

func (up UpdateTestInfo) Type() action.Type {
	return action.PASSPORT_UPDATE_TEST
}

func (up UpdateTestInfo) Tags() kv.Pairs {
	tags := make([]kv.Pair, 0)

	tag := kv.Pair{
		Key:   []byte("tx.type"),
		Value: []byte(up.Type().String()),
	}
	tag1 := kv.Pair{
		Key:   []byte("tx.testId"),
		Value: []byte(up.TestID),
	}
	tag2 := kv.Pair{
		Key:   []byte("tx.adminAddress"),
		Value: up.AdminAddress.Bytes(),
	}
	tag3 := kv.Pair{
		Key:   []byte("tx.person"),
		Value: []byte(up.Person.String()),
	}
	tag4 := kv.Pair{
		Key:   []byte("tx.testResult"),
		Value: []byte(up.Result.String()),
	}

	tags = append(tags, tag, tag1, tag2, tag3, tag4)
	return tags
}

func (up UpdateTestInfo) Marshal() ([]byte, error) {
	return json.Marshal(up)
}

func (up *UpdateTestInfo) Unmarshal(bytes []byte) error {
	return json.Unmarshal(bytes, up)
}

type updateTestInfoTx struct {
}

var _ action.Tx = updateTestInfoTx{}

func (updateTestInfoTx) Validate(ctx *action.Context, signedTx action.SignedTx) (bool, error) {
	up := UpdateTestInfo{}

	err := up.Unmarshal(signedTx.Data)
	if err != nil {
		return false, errors.Wrap(action.ErrWrongTxType, err.Error())
	}

	// Validate basic signature
	err = action.ValidateBasic(signedTx.RawBytes(), up.Signers(), signedTx.Signatures)
	if err != nil {
		return false, err
	}

	// Check if test ID is valid
	if len(up.TestID) != passport.TestIDLength {
		return false, passport.ErrInvalidTestID
	}

	// Check if test subject is valid
	if err = up.Person.Err(); err != nil {
		return false, passport.ErrInvalidIdentifier.Wrap(err)
	}

	// Check if test type is valid
	if err = up.Test.Err(); err != nil {
		return false, passport.ErrInvalidTestType.Wrap(err)
	}

	// Check if user admin is valid
	if err = up.Admin.Err(); err != nil {
		return false, passport.ErrInvalidIdentifier.Wrap(err)
	}

	// Check if admin address is valid OneLedger address
	err = up.AdminAddress.Err()
	if err != nil {
		return false, action.ErrInvalidAddress.Wrap(err)
	}

	// check time stamps
	if !validateRFCTime(up.AnalyzedAt) {
		return false, passport.ErrInvalidTimeStamp
	}

	return true, nil
}

func (updateTestInfoTx) ProcessCheck(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	ctx.Logger.Debug("Processing updateTestInfoTx Transaction for CheckTx", tx)
	return runUpdate(ctx, tx)
}

func (updateTestInfoTx) ProcessDeliver(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	ctx.Logger.Debug("Processing updateTestInfoTx Transaction for DeliverTx", tx)
	return runUpdate(ctx, tx)
}

//No Fees Are charged for this transaction
func (updateTestInfoTx) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas) (bool, action.Response) {
	return true, action.Response{}
}

func runUpdate(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	up := &UpdateTestInfo{}
	err := up.Unmarshal(tx.Data)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, ErrSerialization, up.Tags(), err)
	}

	// check permission
	permitted, err := ctx.AuthTokens.HasPermission(up.AnalysisOrg, up.Admin, passport.PermitUpload)
	if !permitted {
		return helpers.LogAndReturnFalse(ctx.Logger, passport.ErrPermissionRequired, up.Tags(), err)
	}

	// create update record
	info := passport.NewUpdateTestInfo(up.TestID, up.Person, up.Test, up.Result,
		up.AnalysisOrg, up.AnalyzedAt, up.AnalyzedBy, up.Notes)

	// update
	updated, err := ctx.Tests.UpdateTestInfo(info)
	if updated == false || err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, passport.ErrUpdateTestFailure, up.Tags(), err)
	}
	return helpers.LogAndReturnTrue(ctx.Logger, up.Tags(), "test_info_updated")
}
