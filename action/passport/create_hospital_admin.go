package passport

import (
	"encoding/json"
	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/action/helpers"
	"github.com/Oneledger/protocol/data/passport"
	"github.com/tendermint/tendermint/libs/kv"
)

type CreateHospitalAdmin struct {
	User             passport.UserID       `json:"userId"`
	TokenTypeID      passport.TokenTypeID  `json:"tokenTypeId"`
	TokenType        passport.TokenType    `json:"tokenType"`
	TokenSubType     passport.TokenSubType `json:"tokenSubType"`
	OwnerAddress     action.Address        `json:"ownerAddress"`
	SuperUserAddress action.Address        `json:"superUserAddress"`
	SuperUser        passport.UserID       `json:"superUser"`
	CreationTime     string                `json:"creationTime"`
}

var _ action.Msg = &CreateHospitalAdmin{}

func (chaMsg CreateHospitalAdmin) Signers() []action.Address {
	return []action.Address{chaMsg.SuperUserAddress}
}

func (chaMsg CreateHospitalAdmin) Type() action.Type {
	return action.PASSPORT_HOSP_ADMIN
}

func (chaMsg CreateHospitalAdmin) Tags() kv.Pairs {
	tags := make([]kv.Pair, 0)

	tag := kv.Pair{
		Key:   []byte("tx.type"),
		Value: []byte(chaMsg.Type().String()),
	}
	tag2 := kv.Pair{
		Key:   []byte("tx.owner"),
		Value: chaMsg.OwnerAddress,
	}

	tags = append(tags, tag, tag2)
	return tags
}

func (chaMsg CreateHospitalAdmin) Marshal() ([]byte, error) {
	return json.Marshal(chaMsg)
}

func (chaMsg *CreateHospitalAdmin) Unmarshal(bytes []byte) error {
	return json.Unmarshal(bytes, chaMsg)
}

type createHospitalAdminTx struct {
}

func (cha createHospitalAdminTx) Validate(ctx *action.Context, signedTx action.SignedTx) (bool, error) {
	createHospAdmin := CreateHospitalAdmin{}
	err := createHospAdmin.Unmarshal(signedTx.Data)
	if err != nil {
		return false, err
	}

	//Verify Message was signed by the Super User
	err = action.ValidateBasic(signedTx.RawBytes(), createHospAdmin.Signers(), signedTx.Signatures)
	if err != nil {
		return false, err
	}

	//Make sure the creation token type is not a Super User
	if createHospAdmin.TokenType == passport.TokenSuperAdmin {
		return false, ErrTypeIsSuperUsr
	}

	//Check if User name is valid
	if err = createHospAdmin.User.Err(); err != nil {
		return false, passport.ErrInvalidIdentifier
	}

	//Check if Super User name is valid
	if err = createHospAdmin.SuperUser.Err(); err != nil {
		return false, passport.ErrInvalidIdentifier
	}

	//Check if TokenType ID is valid
	if err = createHospAdmin.TokenTypeID.Err(); err != nil {
		return false, ErrInvalidTokenType
	}

	//Check creation time
	if !validateRFCTime(createHospAdmin.CreationTime) {
		return false, passport.ErrInvalidTimeStamp
	}

	return true, nil
}

func (cha createHospitalAdminTx) ProcessCheck(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	return runCreateHAdminAuthToken(ctx, tx)
}

func (cha createHospitalAdminTx) ProcessDeliver(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	return runCreateHAdminAuthToken(ctx, tx)
}

//No Fees Are charged for this transaction
func (cha createHospitalAdminTx) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas) (bool, action.Response) {
	return true, action.Response{}
}

var _ action.Tx = createHospitalAdminTx{}

func runCreateHAdminAuthToken(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	createHospAdmin := CreateHospitalAdmin{}
	err := createHospAdmin.Unmarshal(tx.Data)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, ErrSerialization, createHospAdmin.Tags(), err)
	}

	//Confirm Super User Address is part of the super user list
	if !ctx.AuthTokens.Exists(passport.TypeIDSuperAdmin, createHospAdmin.SuperUser) {
		return helpers.LogAndReturnFalse(ctx.Logger, ErrSuperUserNotFound, createHospAdmin.Tags(), nil)
	}

	//Check to see if Auth Token already exists
	if ctx.AuthTokens.Exists(createHospAdmin.TokenTypeID, createHospAdmin.User) {
		return helpers.LogAndReturnFalse(ctx.Logger, ErrTokenAlreadyExists, createHospAdmin.Tags(), nil)
	}

	//Create a new Auth token for a Hospital admin using input fields from the message
	authToken := passport.NewAuthToken("",
		passport.TokenHospital,
		createHospAdmin.TokenSubType,
		createHospAdmin.TokenTypeID,
		passport.RoleHospitalAdmin,
		passport.PermitUpload|passport.PermitScan|passport.PermitQueryTest|passport.PermitQueryRead|passport.PermitQueryTokens,
		createHospAdmin.User,
		createHospAdmin.OwnerAddress,
		createHospAdmin.CreationTime,
	)
	//Set the new Auth token in the passport auth token store
	err = ctx.AuthTokens.CreateAuthToken(authToken)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, ErrErrorCreatingToken, createHospAdmin.Tags(), err)
	}
	return helpers.LogAndReturnTrue(ctx.Logger, createHospAdmin.Tags(), "hospital_admin_token_created")
}
