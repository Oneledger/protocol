package passport

import (
	"encoding/json"
	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/action/helpers"
	"github.com/Oneledger/protocol/data/passport"
	"github.com/tendermint/tendermint/libs/kv"
)

type CreateScreenerAdmin struct {
	User             passport.UserID       `json:"userId"`
	TokenTypeID      passport.TokenTypeID  `json:"tokenTypeId"`
	TokenType        passport.TokenType    `json:"tokenType"`
	TokenSubType     passport.TokenSubType `json:"tokenSubType"`
	OwnerAddress     action.Address        `json:"ownerAddress"`
	SuperUserAddress action.Address        `json:"superUserAddress"`
	SuperUser        passport.UserID       `json:"superUser"`
	CreationTime     string                `json:"creationTime"`
}

var _ action.Msg = &CreateScreenerAdmin{}

func (csa CreateScreenerAdmin) Signers() []action.Address {
	return []action.Address{csa.SuperUserAddress}
}

func (csa CreateScreenerAdmin) Type() action.Type {
	return action.PASSPORT_SCR_ADMIN
}

func (csa CreateScreenerAdmin) Tags() kv.Pairs {
	tags := make([]kv.Pair, 0)

	tag := kv.Pair{
		Key:   []byte("tx.type"),
		Value: []byte(csa.Type().String()),
	}
	tag2 := kv.Pair{
		Key:   []byte("tx.owner"),
		Value: csa.OwnerAddress,
	}

	tags = append(tags, tag, tag2)
	return tags
}

func (csa CreateScreenerAdmin) Marshal() ([]byte, error) {
	return json.Marshal(csa)
}

func (csa *CreateScreenerAdmin) Unmarshal(bytes []byte) error {
	return json.Unmarshal(bytes, csa)
}

type createScreenerAdminTx struct {
}

var _ action.Tx = createScreenerAdminTx{}

func (c createScreenerAdminTx) Validate(ctx *action.Context, signedTx action.SignedTx) (bool, error) {
	createScreenerAdmin := CreateScreenerAdmin{}
	err := createScreenerAdmin.Unmarshal(signedTx.Data)
	if err != nil {
		return false, err
	}

	//Verify Message was signed by the Super User
	err = action.ValidateBasic(signedTx.RawBytes(), createScreenerAdmin.Signers(), signedTx.Signatures)
	if err != nil {
		return false, err
	}

	//Make sure the creation token type is not a Super User
	if createScreenerAdmin.TokenType == passport.TokenSuperAdmin {
		return false, ErrTypeIsSuperUsr
	}

	//Check if User name is valid
	if err = createScreenerAdmin.User.Err(); err != nil {
		return false, passport.ErrInvalidIdentifier
	}

	//Check if Super User name is valid
	if err = createScreenerAdmin.SuperUser.Err(); err != nil {
		return false, passport.ErrInvalidIdentifier
	}

	//Check if TokenType ID is valid
	if err = createScreenerAdmin.TokenTypeID.Err(); err != nil {
		return false, ErrInvalidTokenType
	}

	//Check creation time
	if !validateRFCTime(createScreenerAdmin.CreationTime) {
		return false, passport.ErrInvalidTimeStamp
	}

	return true, nil
}

func (c createScreenerAdminTx) ProcessCheck(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	return runCreateScrAdminAuthToken(ctx, tx)
}

func (c createScreenerAdminTx) ProcessDeliver(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	return runCreateScrAdminAuthToken(ctx, tx)
}

//No Fees Are charged for this transaction
func (c createScreenerAdminTx) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas) (bool, action.Response) {
	return true, action.Response{}
}

func runCreateScrAdminAuthToken(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	createScrAdmin := CreateScreenerAdmin{}
	err := createScrAdmin.Unmarshal(tx.Data)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, ErrSerialization, createScrAdmin.Tags(), err)
	}

	//Confirm Super User Address is part of the super user list
	if !ctx.AuthTokens.Exists(passport.TypeIDSuperAdmin, createScrAdmin.SuperUser) {
		return helpers.LogAndReturnFalse(ctx.Logger, ErrSuperUserNotFound, createScrAdmin.Tags(), nil)
	}

	//Check to see if Auth Token already exists
	if ctx.AuthTokens.Exists(createScrAdmin.TokenTypeID, createScrAdmin.User) {
		return helpers.LogAndReturnFalse(ctx.Logger, ErrTokenAlreadyExists, createScrAdmin.Tags(), nil)
	}

	//Create a new Auth token for a Screener admin using input fields from the message
	authToken := passport.NewAuthToken("",
		passport.TokenScreener,
		createScrAdmin.TokenSubType,
		createScrAdmin.TokenTypeID,
		passport.RoleScreenerAdmin,
		passport.PermitScan|passport.PermitQueryTest|passport.PermitQueryRead|passport.PermitQueryTokens,
		createScrAdmin.User,
		createScrAdmin.OwnerAddress,
		createScrAdmin.CreationTime,
	)

	//Set the new Auth token in the passport auth token store
	err = ctx.AuthTokens.CreateAuthToken(authToken)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, ErrErrorCreatingToken, createScrAdmin.Tags(), err)
	}
	return helpers.LogAndReturnTrue(ctx.Logger, createScrAdmin.Tags(), "screener_admin_token_created")
}
