package passport

import (
	codes "github.com/Oneledger/protocol/status_codes"
)

var (
	// Health passport errors
	ErrInvalidIdentifier   = codes.ProtocolError{codes.PsptErrInvalidIdentifier, "Invalid user identifier"}
	ErrDuplicateIdentifier = codes.ProtocolError{Code: codes.PsptErrIdentifierDuplicate, Msg: "Identifier must be unique"}
	ErrPermissionRequired  = codes.ProtocolError{Code: codes.PsptErrPermissionRequired, Msg: "This operation requires permission"}
	ErrAuthTokenNotFound   = codes.ProtocolError{Code: codes.PsptErrAuthTokenNotFound, Msg: "Auth token not found"}
	ErrAuthTokenRemove     = codes.ProtocolError{Code: codes.GeneralErr, Msg: "Auth token remove failure"}

	ErrInvalidTestType     = codes.ProtocolError{codes.PsptErrInvalidTestType, "Invalid test type"}
	ErrInvalidTestSubType  = codes.ProtocolError{codes.PsptErrInvalidTestSubType, "Invalid test sub type"}
	ErrInvalidTestResult   = codes.ProtocolError{codes.PsptErrInvalidTestResult, "Invalid test result"}
	ErrInvalidTimeStamp    = codes.ProtocolError{codes.PsptErrInvalidTimeStamp, "Invalid RFC3339 time stamp"}
	ErrUploadTestFailure   = codes.ProtocolError{codes.PsptErrUploadTestFailure, "Failed to upload test information"}
	ErrReadTestFailure     = codes.ProtocolError{codes.PsptErrReadTestFailure, "Failed to read test information"}
	ErrInvalidManufacturer = codes.ProtocolError{codes.PsptErrInvalidManufacturer, "Invalid Manufacturer"}
	ErrUpdateTestFailure   = codes.ProtocolError{codes.PsptErrUpdateTestFailure, "Failed to update test information"}
	ErrInvalidTestID       = codes.ProtocolError{codes.PsptErrInvalidTestID, "Invalid test ID"}
	ErrInitialResult       = codes.ProtocolError{codes.PsptErrInitialResult, "Test result must be initiated pending"}
)
