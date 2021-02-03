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

var _ action.Msg = &UploadTestInfo{}

type UploadTestInfo struct {
	TestID       string               `json:"testId"`
	Person       passport.UserID      `json:"person"`
	Test         passport.TestType    `json:"test"`
	SubTest      passport.TestSubType `json:"subTest"`
	Manufacturer string               `json:"manufacturer"`
	Result       passport.TestResult  `json:"testResult"`

	TestOrg  passport.TokenTypeID `json:"testOrg"`
	TestedAt string               `json:"testedAt"`
	TestedBy string               `json:"testedBy"`

	AnalysisOrg passport.TokenTypeID `json:"analysisOrg"`
	AnalyzedAt  string               `json:"analyzedAt"`
	AnalyzedBy  string               `json:"analyzedBy"`

	Admin        passport.UserID `json:"admin"`
	AdminAddress keys.Address    `json:"adminAddress"`
	UploadedAt   string          `json:"uploadedAt"`
	Notes        string          `json:"notes"`
}

func (up UploadTestInfo) Signers() []action.Address {
	return []action.Address{up.AdminAddress}
}

func (up UploadTestInfo) Type() action.Type {
	return action.PASSPORT_UPLOAD_TEST
}

func (up UploadTestInfo) Tags() kv.Pairs {
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
		Key:   []byte("tx.testOrg"),
		Value: []byte(up.TestOrg.String()),
	}
	tag3 := kv.Pair{
		Key:   []byte("tx.adminAddress"),
		Value: up.AdminAddress.Bytes(),
	}
	tag4 := kv.Pair{
		Key:   []byte("tx.person"),
		Value: []byte(up.Person.String()),
	}
	tag5 := kv.Pair{
		Key:   []byte("tx.testResult"),
		Value: []byte(up.Result.String()),
	}

	tags = append(tags, tag, tag1, tag2, tag3, tag4, tag5)
	return tags
}

func (up UploadTestInfo) Marshal() ([]byte, error) {
	return json.Marshal(up)
}

func (up *UploadTestInfo) Unmarshal(bytes []byte) error {
	return json.Unmarshal(bytes, up)
}

type uploadTestInfoTx struct {
}

var _ action.Tx = uploadTestInfoTx{}

func (uploadTestInfoTx) Validate(ctx *action.Context, signedTx action.SignedTx) (bool, error) {
	up := UploadTestInfo{}

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

	// Check if test sub type is valid
	if err = up.SubTest.Err(); err != nil {
		return false, passport.ErrInvalidTestSubType.Wrap(err)
	}

	// Check if manufacturer is valid
	if len(up.Manufacturer) == 0 {
		return false, passport.ErrInvalidManufacturer
	}

	// Check if test org identifier address is valid
	err = up.TestOrg.Err()
	if err != nil {
		return false, passport.ErrInvalidIdentifier.Wrap(err)
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
	if !validateRFCTime(up.TestedAt) {
		return false, passport.ErrInvalidTimeStamp
	}
	if up.AnalyzedAt != "" && !validateRFCTime(up.AnalyzedAt) {
		return false, passport.ErrInvalidTimeStamp
	}
	if !validateRFCTime(up.UploadedAt) {
		return false, passport.ErrInvalidTimeStamp
	}

	return true, nil
}

func (uploadTestInfoTx) ProcessCheck(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	ctx.Logger.Debug("Processing uploadTestInfoTx Transaction for CheckTx", tx)
	return runUpload(ctx, tx)
}

func (uploadTestInfoTx) ProcessDeliver(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	ctx.Logger.Debug("Processing uploadTestInfoTx Transaction for DeliverTx", tx)
	return runUpload(ctx, tx)
}

//No Fees Are charged for this transaction
func (uploadTestInfoTx) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas) (bool, action.Response) {
	return true, action.Response{}
}

func runUpload(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	up := &UploadTestInfo{}
	err := up.Unmarshal(tx.Data)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, ErrSerialization, up.Tags(), err)
	}

	// check permission
	permitted, err := ctx.AuthTokens.HasPermission(up.TestOrg, up.Admin, passport.PermitUpload)
	if !permitted {
		return helpers.LogAndReturnFalse(ctx.Logger, passport.ErrPermissionRequired, up.Tags(), err)
	}

	// create test record
	info := passport.NewTestInfo(up.TestID, up.Person, up.Test, up.SubTest, up.Manufacturer, up.Result,
		up.TestOrg, up.TestedAt, up.TestedBy, up.AnalysisOrg, up.AnalyzedAt, up.AnalyzedBy,
		up.UploadedAt, up.Admin, up.Notes)

	err = ctx.Tests.AddTestInfo(info)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, passport.ErrUploadTestFailure, up.Tags(), err)
	}
	return helpers.LogAndReturnTrue(ctx.Logger, up.Tags(), "test_info_upload")
}
