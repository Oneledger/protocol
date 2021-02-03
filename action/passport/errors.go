package passport

import codes "github.com/Oneledger/protocol/status_codes"

var (
	ErrSuperUserNotFound  = codes.ProtocolError{Code: codes.PsptErrSuperUserNotFound, Msg: "super user doesn't exist"}
	ErrTokenAlreadyExists = codes.ProtocolError{Code: codes.PsptErrAuthTokenAlreadyExists, Msg: "auth token already exists"}
	ErrOwnerAddressExists = codes.ProtocolError{Code: codes.PsptErrOwnerAddressExists, Msg: "address already associated with token"}
	ErrAddressIsSuperUser = codes.ProtocolError{Code: codes.PsptErrAddressIsSuperUser, Msg: "owner address is a super user"}
	ErrSerialization      = codes.ProtocolError{Code: codes.PsptErrSerialization, Msg: "serialization error"}
	ErrErrorCreatingToken = codes.ProtocolError{Code: codes.PsptErrErrorCreatingToken, Msg: "error persisting token to db"}
	ErrTypeIsSuperUsr     = codes.ProtocolError{Code: codes.PsptErrTokenTypeIsSuperUser, Msg: "cannot create another super user"}
	ErrInvalidTokenType   = codes.ProtocolError{Code: codes.PsptErrInvalidTokenType, Msg: "token type id is invalid"}
)
