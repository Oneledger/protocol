package rpc

import (
	"github.com/powerman/rpc-codec/jsonrpc2"
)

type Code int

const (
	// Pre-defined errors for JSON-RPC 2.0 https://www.jsonrpc.org/specification#error_object
	CodeParseError     = -32700
	CodeInvalidRequest = -32600
	CodeMethodNotFound = -32601
	CodeInvalidParams  = -32602
	CodeInternalError  = -32603

	// -32000 to -32099 Are application level errors (i.e. define your rpc errors in this range)
	CodeNotAllowed = -32001

	CodeNotFound = -32000
)

type Error = jsonrpc2.Error

func NewError(code Code, msg string) *Error {
	return jsonrpc2.NewError(int(code), msg)
}

func InvalidParamsError(msg string) *Error {
	return NewError(CodeInvalidParams, msg)
}

func InternalError(msg string) *Error {
	return NewError(CodeInternalError, msg)
}

func ParseError(msg string) *Error {
	return NewError(CodeParseError, msg)
}

func InvalidRequestError(msg string) *Error {
	return NewError(CodeInvalidRequest, msg)
}

func MethodNotFoundError(msg string) *Error {
	return NewError(CodeMethodNotFound, msg)
}

/* Define application-specific errors here */

func NotAllowedError(msg string) *Error {
	return NewError(CodeNotAllowed, msg)
}
