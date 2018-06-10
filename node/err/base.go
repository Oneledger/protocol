/*
	Copyright 2017-2018 OneLedger

	Common errors across the entire system
*/
package err

type Code = uint32 // Matches Tendermint status

const (
	SUCCESS         Code = 0
	INVALID         Code = 101
	PARSE_ERROR     Code = 201
	NOT_IMPLEMENTED Code = 301
	MISSING_VALUE   Code = 401
	EXPAND_ERROR    Code = 501
	DUPLICATE       Code = 601
)
