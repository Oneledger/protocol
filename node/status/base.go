/*
	Copyright 2017-2018 OneLedger

	Common errors across the entire system
*/
package status

type Code = uint32 // Matches Tendermint status

const (
	SUCCESS           Code = 0
	INVALID           Code = 101
	INVALID_SIGNATURE Code = 102
	PARSE_ERROR       Code = 201
	NOT_IMPLEMENTED   Code = 301
	MISSING_VALUE     Code = 401
	EXPAND_ERROR      Code = 501
	DUPLICATE         Code = 601
	MISSING_DATA      Code = 701
	BAD_VALUE         Code = 801
	EXECUTE_ERROR     Code = 901
)
